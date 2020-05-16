package handlers

import (
	"context"
	"fmt"
	"github.com/jinzhu/now"
	"github.com/pkg/errors"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"gopkg.in/DataDog/dd-trace-go.v1/contrib/go-redis/redis"
	"merryworld/surebank/internal/account"
	"merryworld/surebank/internal/customer"
	"merryworld/surebank/internal/platform/auth"
	"merryworld/surebank/internal/platform/datatable"
	"merryworld/surebank/internal/platform/web"
	"merryworld/surebank/internal/shop"
	"merryworld/surebank/internal/transaction"
	"merryworld/surebank/internal/user"
	"net/http"
	"strings"
	"time"
)

// Customers represents the Customers API method handler set.
type Reports struct {
	CustomerRepo    *customer.Repository
	AccountRepo     *account.Repository
	TransactionRepo *transaction.Repository
	ShopRepo        *shop.Repository
	UserRepos		*user.Repository
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
		{Field: "narration", Title: "Narration", Visible: true, Searchable: true, Orderable: true, Filterable: true, FilterPlaceholder: "filter Narration"},
		{Field: "account", Title: "Account", Visible: true, Searchable: true, Orderable: true, Filterable: true, FilterPlaceholder: "filter Account"},
		{Field: "sales_rep_id", Title: "Recorded By", Visible: true, Searchable: true, Orderable: false, Filterable: true, FilterPlaceholder: "filter Recorder"},
		{Field: "opening_balance", Title: "Opening Balance", Visible: true, Searchable: false, Orderable: true, Filterable: false},
	}

	mapFunc := func(q *transaction.Response, cols []datatable.DisplayField) (resp []datatable.ColumnValue, err error) {
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
				v.Value = q.CreatedAt.Local
				v.Formatted = q.CreatedAt.Local
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
			case "account":
				v.Value = q.AccountNumber
				v.Formatted = fmt.Sprintf("<a href='%s'>%s</a>", urlCustomersAccountsView(q.CustomerID, q.AccountID), v.Value)
			case "sales_rep_id":
				v.Value = q.SalesRepID
				v.Formatted = fmt.Sprintf("<a href='%s'>%s</a>", urlUsersView(q.SalesRepID), q.SalesRep)
			case "opening_balance":
				v.Value = fmt.Sprintf("%f", q.OpeningBalance)
				p := message.NewPrinter(language.English)
				v.Formatted = p.Sprintf("%.2f", q.OpeningBalance)
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

		var res = &transaction.PagedResponseList{}
		// 0 where means this customer has no associated account
		if len(txWhere) > 0 {
			res, err = h.TransactionRepo.Find(ctx, claims, transaction.FindRequest{
				Order: order, Where: strings.Join(txWhere, " AND "), Args: txArgs,
			})
			if err != nil {
				return resp, err
			}
		}

		for _, a := range res.Transactions {
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

	data["users"] = users
	data["total"] = total
	data["datatable"] = dt.Response()

	return h.Renderer.Render(ctx, w, r, TmplLayoutBase, "report-transactions.gohtml", web.MIMETextHTMLCharsetUTF8, http.StatusOK, data)
}

// Ajor handles listing all the Ajor account types
func (h *Reports) Ajor(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {

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

		res, err := h.AccountRepo.FindAjor(ctx, claims, account.FindRequest{
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
				return resp, errors.Wrapf(err, "Failed to map Ajor accounts for display.")
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

	return h.Renderer.Render(ctx, w, r, TmplLayoutBase, "report-ajor.gohtml", web.MIMETextHTMLCharsetUTF8, http.StatusOK, data)
}
