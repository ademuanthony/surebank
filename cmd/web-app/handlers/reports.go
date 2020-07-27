package handlers

import (
	"context"
	"fmt"
	"merryworld/surebank/internal/account"
	"merryworld/surebank/internal/customer"
	"merryworld/surebank/internal/dscommission"
	"merryworld/surebank/internal/platform/auth"
	"merryworld/surebank/internal/platform/datatable"
	"merryworld/surebank/internal/platform/web"
	"merryworld/surebank/internal/platform/web/webcontext"
	"merryworld/surebank/internal/shop"
	"merryworld/surebank/internal/transaction"
	"merryworld/surebank/internal/user"
	"net/http"
	"strings"
	"time"

	"github.com/jinzhu/now"
	"github.com/pkg/errors"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"gopkg.in/DataDog/dd-trace-go.v1/contrib/go-redis/redis"
)

// Customers represents the Customers API method handler set.
type Reports struct {
	CustomerRepo    *customer.Repository
	AccountRepo     *account.Repository
	TransactionRepo *transaction.Repository
	ShopRepo        *shop.Repository
	UserRepos		*user.Repository
	CommissionRepo	*dscommission.Repository
	Renderer        web.Renderer
	Redis           *redis.Client
}

// Transactions handles listing all the customers transactions across all his accounts.
func (h *Reports) Transactions(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {

	var data = make(map[string]interface{})
	var total float64

	claims, err := auth.ClaimsFromContext(ctx)
	if err != nil {
		return err
	}

	fields := []datatable.DisplayField{
		{Field: "id", Title: "ID", Visible: false, Searchable: true, Orderable: true, Filterable: false},
		{Field: "amount", Title: "Quantity", Visible: true, Searchable: false, Orderable: true, Filterable: true, FilterPlaceholder: "filter Quantity"},
		{Field: "created_at", Title: "Date", Visible: true, Searchable: true, Orderable: true, Filterable: true, FilterPlaceholder: "filter Date"},
		{Field: "customer_name", Title: "Customer", Visible: true, Searchable: true, Orderable: true, Filterable: true, FilterPlaceholder: "filter Account"},
		{Field: "account", Title: "Account Number", Visible: true, Searchable: true, Orderable: true, Filterable: true, FilterPlaceholder: "filter Account"},
		{Field: "sales_rep_id", Title: "Recorded By", Visible: true, Searchable: true, Orderable: false, Filterable: true, FilterPlaceholder: "filter Recorder"},
	}

	mapFunc := func(q transaction.TxReportResponse, cols []datatable.DisplayField) (resp []datatable.ColumnValue, err error) {
		for i := 0; i < len(cols); i++ {
			col := cols[i]
			var v datatable.ColumnValue
			switch col.Field { 
			case "id":
				v.Value = fmt.Sprintf("%s", q.ID)
			case "amount":
				v.Value = fmt.Sprintf("%f", q.Amount)
				p := message.NewPrinter(language.English)
				v.Formatted = p.Sprintf("<a href='%s'>%.2f</a>", urlCustomersTransactionsView(q.CustomerID, q.AccountID, q.ID), q.Amount)
			case "created_at":
				date := web.NewTimeResponse(ctx, time.Unix(q.CreatedAt, 0))
				v.Value = date.LocalDate
				v.Formatted = date.LocalDate
			case "narration":
				values := strings.Split(q.Narration, ":")
				if len(values) > 1 {
					if values[0] == "sale" {
						v.Value = values[1]
						v.Formatted = fmt.Sprintf("<a href='%s'>%s</a>", urlSalesView(values[2]), v.Value)
					}
				} else {
					v.Value = q.Narration
					v.Formatted = q.Narration
				}
			case "payment_method": 
				v.Value = q.PaymentMethod
				v.Formatted = q.PaymentMethod
			case "customer_name":
				v.Value = q.CustomerName
				v.Formatted = fmt.Sprintf("<a href='%s'>%s</a>", urlCustomersView(q.CustomerID), v.Value)
			case "account":
				v.Value = q.AccountNumber
				v.Formatted = fmt.Sprintf("<a href='%s'>%s</a>", urlCustomersAccountsView(q.CustomerID, q.AccountID), v.Value)
			case "sales_rep_id":
				v.Value = q.SalesRepID
				v.Formatted = fmt.Sprintf("<a href='%s'>%s</a>", urlUsersView(q.SalesRepID), q.SalesRep)
			default:
				return resp, errors.Errorf("Failed to map value for %s.", col.Field)
			}
			resp = append(resp, v)
		}

		return resp, nil
	}

	var txWhere = []string { "tx_type = 'deposit'"}
	var txArgs []interface{}
	// todo sales rep filtering
	if v := r.URL.Query().Get("sales_rep_id"); v != "" {
		txWhere = append(txWhere, "sales_rep_id = $1")
		txArgs = append(txArgs, v)
		data["salesRepID"] = v
	}

	if v := r.URL.Query().Get("payment_method"); v != "" {
		txWhere = append(txWhere, fmt.Sprintf("payment_method = $%d", len(txArgs) + 1))
		txArgs = append(txArgs, v)
		data["paymentMethod"] = v
	}

	if v := r.URL.Query().Get("start_date"); v != "" {
		date, err := time.Parse("01/02/2006", v)
		if err != nil {
			return err
		}
		date = date.Truncate(time.Millisecond)
		date = now.New(date).BeginningOfDay().Add(-1 * time.Hour)
		txWhere = append(txWhere, fmt.Sprintf("created_at >= $%d", len(txArgs) + 1))
		txArgs = append(txArgs, date.UTC().Unix())
		data["startDate"] = v
		// 1581897600
		// 1581897323 
	}

	if v := r.URL.Query().Get("end_date"); v != "" {
		date, err := time.Parse("01/02/2006", v)
		if err != nil {
			return err
		}
		date = date.Truncate(time.Millisecond)
		date = now.New(date).EndOfDay().Add(-1 * time.Hour)
		txWhere = append(txWhere, fmt.Sprintf("created_at <= $%d", len(txArgs) + 1))
		txArgs = append(txArgs, date.Unix())
		data["endDate"] = v
	}

	loadFunc := func(ctx context.Context, sorting string, fields []datatable.DisplayField) (resp [][]datatable.ColumnValue, err error) {

		var order []string
		if len(sorting) > 0 {
			order = strings.Split(sorting, ",")
		}

		for i := range txWhere {
			txWhere[i] = "tx." + txWhere[i]
		}
		res, err := h.TransactionRepo.TxReport(ctx, claims, transaction.FindRequest{
			Order: order, Where: strings.Join(txWhere, " AND "), Args: txArgs,
		})
		if err != nil {
			return resp, err
		}

		for _, a := range res {
			l, err := mapFunc(a, fields)
			if err != nil {
				return resp, errors.Wrapf(err, "Failed to map brand for display.")
			}

			resp = append(resp, l)
		}

		return resp, nil
	}

	dt, err := datatable.New(ctx, w, r, h.Redis, fields, loadFunc)
	if err != nil {
		return err
	}

	if dt.HasCache() {
		return nil
	}

	if ok, err := dt.Render(); ok {
		if err != nil {
			return err
		}
		return nil
	}

	users, err := h.UserRepos.Find(ctx, claims, user.UserFindRequest{
		Order: []string{"first_name", "last_name"},
	})
	if err != nil {
		return err
	}

	total, err = h.TransactionRepo.DepositAmountByWhere(ctx, strings.Join(txWhere, " and "), txArgs)
	if err != nil {
		return err
	}

	data["paymentMethods"] = transaction.PaymentMethods
	data["users"] = users
	data["total"] = total
	data["datatable"] = dt.Response()

	return h.Renderer.Render(ctx, w, r, TmplLayoutBase, "report-transactions.gohtml", web.MIMETextHTMLCharsetUTF8, http.StatusOK, data)
}

// Ds handles listing of all the Ds account types
func (h *Reports) Ds(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {

	var data = make(map[string]interface{})

	claims, err := auth.ClaimsFromContext(ctx)
	if err != nil {
		return err 
	}

	fields := []datatable.DisplayField{
		{Field: "id", Title: "ID", Visible: false, Searchable: true, Orderable: true, Filterable: false},
		{Field: "customer", Title: "Name", Visible: true, Searchable: false, Orderable: true, Filterable: true, FilterPlaceholder: "filter Quantity"},
		{Field: "number", Title: "Account Number", Visible: true, Searchable: true, Orderable: true, Filterable: true, FilterPlaceholder: "filter Date"},
		{Field: "target", Title: "Daily Contribution", Visible: true, Searchable: true, Orderable: true, Filterable: true, FilterPlaceholder: "filter Narration"},
		{Field: "balance", Title: "Account Balance", Visible: true, Searchable: false, Orderable: true, Filterable: false},
		{Field: "sales_rep_id", Title: "Account Manager", Visible: true, Searchable: true, Orderable: false, Filterable: true, FilterPlaceholder: "filter Recorder"},
		{Field: "created_at", Title: "Registration Date", Visible: true, Searchable: false, Orderable: true, Filterable: false},
	}

	mapFunc := func(q *account.Response, cols []datatable.DisplayField) (resp []datatable.ColumnValue, err error) {
		for i := 0; i < len(cols); i++ {
			col := cols[i]
			var v datatable.ColumnValue
			switch col.Field {
			case "id":
				v.Value = fmt.Sprintf("%s", q.ID)
			case "customer":
				v.Value = q.Customer.Name
				v.Formatted = fmt.Sprintf("<a href='%s'>%s</a>", urlCustomersView(q.CustomerID), q.Customer.Name)
			case "number":
				v.Value = q.Number
				v.Formatted = fmt.Sprintf("<a href='%s'>%s</a>", urlCustomersAccountsView(q.CustomerID, q.ID), q.Number)
			case "target":
				v.Value = fmt.Sprintf("%f", q.Balance)
				p := message.NewPrinter(language.English)
				v.Formatted = p.Sprintf("%.2f", q.Balance)
			case "balance":
				v.Value = fmt.Sprintf("%f", q.Balance)
				p := message.NewPrinter(language.English)
				v.Formatted = p.Sprintf("%.2f", q.Balance)
			case "sales_rep_id":
				v.Value = q.SalesRepID
				v.Formatted = fmt.Sprintf("<a href='%s'>%s</a>", urlUsersView(q.SalesRepID), q.SalesRep)
			case "created_at":
				v.Value = q.CreatedAt.Local
				v.Formatted = v.Value
			default:
				return resp, errors.Errorf("Failed to map value for %s.", col.Field)
			}
			resp = append(resp, v)
		}

		return resp, nil
	}

	loadFunc := func(ctx context.Context, sorting string, fields []datatable.DisplayField) (resp [][]datatable.ColumnValue, err error) {

		var order []string
		if len(sorting) > 0 {
			order = strings.Split(sorting, ",")
		}

		res, err := h.AccountRepo.FindDs(ctx, claims, account.FindRequest{
			Order: order,
			IncludeBranch: true,
			IncludeCustomer: true,
			IncludeSalesRep: true,
		})
		if err != nil {
			return resp, err
		}

		for _, a := range res.Accounts {
			l, err := mapFunc(a, fields)
			if err != nil {
				return resp, errors.Wrapf(err, "Failed to map DS accounts for display.")
			}

			resp = append(resp, l)
		}

		return resp, nil
	}

	dt, err := datatable.New(ctx, w, r, h.Redis, fields, loadFunc)
	if err != nil {
		return err
	}

	if dt.HasCache() {
		return nil
	}

	if ok, err := dt.Render(); ok {
		if err != nil {
			return err
		}
		return nil
	}

	data["datatable"] = dt.Response()

	return h.Renderer.Render(ctx, w, r, TmplLayoutBase, "report-ds.gohtml", web.MIMETextHTMLCharsetUTF8, http.StatusOK, data)
}

// Ds handles listing of all the Ds account types
func (h *Reports) Debtors(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {

	ctxValue, err := webcontext.ContextValues(ctx)
	if err != nil {
		return err	
	}

	var data = make(map[string]interface{})

	claims, err := auth.ClaimsFromContext(ctx)
	if err != nil {
		return err 
	}
 
	fields := []datatable.DisplayField{
		{Field: "id", Title: "ID", Visible: false, Searchable: true, Orderable: true, Filterable: false},
		{Field: "customer", Title: "Name", Visible: true, Searchable: false, Orderable: true, Filterable: true, FilterPlaceholder: "filter Quantity"},
		{Field: "phone_number", Title: "Phone Number", Visible: true, Searchable: false, Orderable: true, Filterable: true, FilterPlaceholder: "filter Quantity"},
		{Field: "number", Title: "Account Number", Visible: true, Searchable: true, Orderable: true, Filterable: true, FilterPlaceholder: "filter Date"},
		{Field: "last_payment_date", Title: "Last Payment Date", Visible: true, Searchable: true, Orderable: true, Filterable: true, FilterPlaceholder: "filter Narration"},
		{Field: "target", Title: "Daily Contribution", Visible: true, Searchable: true, Orderable: true, Filterable: true, FilterPlaceholder: "filter Narration"},
		{Field: "balance", Title: "Account Balance", Visible: true, Searchable: false, Orderable: true, Filterable: false},
		{Field: "sales_rep_id", Title: "Account Manager", Visible: true, Searchable: true, Orderable: false, Filterable: true, FilterPlaceholder: "filter Recorder"},
		{Field: "created_at", Title: "Registration Date", Visible: true, Searchable: false, Orderable: true, Filterable: false},
	}

	mapFunc := func(q *account.Response, cols []datatable.DisplayField) (resp []datatable.ColumnValue, err error) {
		for i := 0; i < len(cols); i++ {
			col := cols[i]
			var v datatable.ColumnValue
			switch col.Field {
			case "id":
				v.Value = fmt.Sprintf("%s", q.ID)
			case "customer":
				v.Value = q.Customer.ShortName
				v.Formatted = fmt.Sprintf("<a href='%s'>%s</a>", urlCustomersView(q.CustomerID), q.Customer.ShortName)
			case "phone_number":
				v.Value = q.Customer.PhoneNumber
				v.Formatted = v.Value
			case "number":
				v.Value = q.Number
				v.Formatted = fmt.Sprintf("<a href='%s'>%s</a>", urlCustomersAccountsView(q.CustomerID, q.ID), q.Number)
			case "last_payment_date":
				v.Value = q.LastPaymentDate.LocalDate
				v.Formatted = v.Value
			case "target":
				v.Value = fmt.Sprintf("%f", q.Balance)
				p := message.NewPrinter(language.English)
				v.Formatted = p.Sprintf("%.2f", q.Balance)
			case "balance":
				v.Value = fmt.Sprintf("%f", q.Balance)
				p := message.NewPrinter(language.English)
				v.Formatted = p.Sprintf("%.2f", q.Balance)
			case "sales_rep_id":
				v.Value = q.SalesRepID
				v.Formatted = fmt.Sprintf("<a href='%s'>%s</a>", urlUsersView(q.SalesRepID), q.SalesRep)
			case "created_at":
				v.Value = q.CreatedAt.LocalDate
				v.Formatted = v.Value
			default:
				return resp, errors.Errorf("Failed to map value for %s.", col.Field)
			}
			resp = append(resp, v)
		}

		return resp, nil
	}

	loadFunc := func(ctx context.Context, sorting string, fields []datatable.DisplayField) (resp [][]datatable.ColumnValue, err error) {

		var order []string
		if len(sorting) > 0 {
			order = strings.Split(sorting, ",")
		}

		res, err := h.AccountRepo.Debtors(ctx, claims, account.FindRequest{
			Order: order,
			IncludeBranch: true,
			IncludeCustomer: true,
			IncludeSalesRep: true,
		}, ctxValue.Now)
		if err != nil {
			return resp, err
		}

		for _, a := range res.Accounts {
			l, err := mapFunc(a, fields)
			if err != nil {
				return resp, errors.Wrapf(err, "Failed to map DS accounts for display.")
			}

			resp = append(resp, l)
		}

		return resp, nil
	}

	dt, err := datatable.New(ctx, w, r, h.Redis, fields, loadFunc)
	if err != nil {
		return err
	}

	if dt.HasCache() {
		return nil
	}

	if ok, err := dt.Render(); ok {
		if err != nil {
			return err
		}
		return nil
	}

	data["datatable"] = dt.Response()

	return h.Renderer.Render(ctx, w, r, TmplLayoutBase, "report-debtors.gohtml", web.MIMETextHTMLCharsetUTF8, http.StatusOK, data)
}

// DsCommissions handles listing of all the Ds Commissions
func (h *Reports) DsCommissions(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {

	var data = make(map[string]interface{})

	claims, err := auth.ClaimsFromContext(ctx)
	if err != nil {
		return err
	}

	fields := []datatable.DisplayField{
		{Field: "id", Title: "ID", Visible: false, Searchable: true, Orderable: true, Filterable: false},
		{Field: "customer", Title: "Name", Visible: true, Searchable: false, Orderable: true, Filterable: true, FilterPlaceholder: "filter Quantity"},
		{Field: "number", Title: "Account Number", Visible: true, Searchable: true, Orderable: true, Filterable: true, FilterPlaceholder: "filter Date"},
		{Field: "amount", Title: "Daily Contribution", Visible: true, Searchable: true, Orderable: true, Filterable: true, FilterPlaceholder: "filter Narration"},
		{Field: "effective_date", Title: "Effective Date", Visible: true, Searchable: true, Orderable: true, Filterable: true, FilterPlaceholder: "filter Date"},
		{Field: "date", Title: "Created At", Visible: true, Searchable: true, Orderable: true, Filterable: false},
	}

	mapFunc := func(q *dscommission.Response, cols []datatable.DisplayField) (resp []datatable.ColumnValue, err error) {
		for i := 0; i < len(cols); i++ {
			col := cols[i]
			var v datatable.ColumnValue
			switch col.Field {
			case "id":
				v.Value = fmt.Sprintf("%s", q.ID)
			case "customer": 
				v.Value = q.CustomerName
				v.Formatted = fmt.Sprintf("<a href='%s'>%s</a>", urlCustomersView(q.CustomerID), q.CustomerName)
			case "number":
				v.Value = q.AccountNumber
				v.Formatted = fmt.Sprintf("<a href='%s'>%s</a>", urlCustomersAccountsView(q.CustomerID, q.AccountID), q.AccountNumber)
			case "amount":
				v.Value = fmt.Sprintf("%f", q.Amount)
				p := message.NewPrinter(language.English)
				v.Formatted = p.Sprintf("%.2f", q.Amount)
			case "effective_date":
				v.Value = q.EffectiveDate.Local
				v.Formatted = v.Value
			case "date":
				v.Value = q.Date.Local
				v.Formatted = v.Value
			default:
				return resp, errors.Errorf("Failed to map value for %s.", col.Field)
			}
			resp = append(resp, v)
		}

		return resp, nil
	}

	var txWhere []string
	var txArgs []interface{}

	if v := r.URL.Query().Get("start_date"); v != "" {
		date, err := time.Parse("01/02/2006", v)
		if err != nil {
			return err
		}
		date = date.Truncate(time.Millisecond)
		date = now.New(date).BeginningOfDay().Add(-1 * time.Hour)
		txWhere = append(txWhere, fmt.Sprintf("effective_date >= $%d", len(txArgs) + 1))
		txArgs = append(txArgs, date.UTC().Unix())
		data["startDate"] = v
	}

	if v := r.URL.Query().Get("end_date"); v != "" {
		date, err := time.Parse("01/02/2006", v)
		if err != nil {
			return err
		}
		date = date.Truncate(time.Millisecond)
		date = now.New(date).EndOfDay().Add(-1 * time.Hour)
		txWhere = append(txWhere, fmt.Sprintf("effective_date <= $%d", len(txArgs) + 1))
		txArgs = append(txArgs, date.Unix())
		data["endDate"] = v
	}

	loadFunc := func(ctx context.Context, sorting string, fields []datatable.DisplayField) (resp [][]datatable.ColumnValue, err error) {

		var order []string
		if len(sorting) > 0 {
			order = strings.Split(sorting, ",")
		}

		res, err := h.CommissionRepo.Find(ctx, claims, dscommission.FindRequest{
			Order: order, Where: strings.Join(txWhere, " AND "), Args: txArgs,
			IncludeCustomer: true,
			IncludeAccount: true,
		})
		if err != nil {
			return resp, err
		} 

		for _, a := range res.Items {
			l, err := mapFunc(a, fields)
			if err != nil {
				return resp, errors.Wrapf(err, "Failed to map DS commission for display.")
			}
 
			resp = append(resp, l)
		}

		return resp, nil
	}

	dt, err := datatable.New(ctx, w, r, h.Redis, fields, loadFunc)
	if err != nil {
		return err
	}

	if dt.HasCache() {
		return nil
	}

	if ok, err := dt.Render(); ok {
		if err != nil {
			return err
		}
		return nil
	}

	data["datatable"] = dt.Response()

	return h.Renderer.Render(ctx, w, r, TmplLayoutBase, "report-ds-commissions.gohtml", web.MIMETextHTMLCharsetUTF8, http.StatusOK, data)
}
