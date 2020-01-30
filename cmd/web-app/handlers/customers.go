package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"merryworld/surebank/internal/postgres/models"
	"merryworld/surebank/internal/transaction"
	"net/http"
	"strings"

	"merryworld/surebank/internal/account"
	"merryworld/surebank/internal/customer"
	"merryworld/surebank/internal/platform/auth"
	"merryworld/surebank/internal/platform/datatable"
	"merryworld/surebank/internal/platform/web"
	"merryworld/surebank/internal/platform/web/webcontext"
	"merryworld/surebank/internal/platform/web/weberror"

	"github.com/gorilla/schema"
	"github.com/pkg/errors"
	"gopkg.in/DataDog/dd-trace-go.v1/contrib/go-redis/redis"
)

// Customers represents the Customers API method handler set.
type Customers struct {
	CustomerRepo    *customer.Repository
	AccountRepo     *account.Repository
	TransactionRepo *transaction.Repository
	Renderer        web.Renderer
	Redis           *redis.Client
}

var accountTypes = []string {
	models.AccountTypeSB,
	models.AccountTypeDS,
}

func urlCustomersIndex() string {
	return fmt.Sprintf("/customers")
}

func urlCustomersCreate() string {
	return fmt.Sprintf("/customers/create")
}

func urlCustomersView(customerID string) string {
	return fmt.Sprintf("/customers/%s", customerID)
}

func urlCustomersUpdate(customerID string) string {
	return fmt.Sprintf("/customers/%s/update", customerID)
}

func urlCustomersTransactions(customerID string) string {
	return fmt.Sprintf("/customers/%s/transactions", customerID)
}

func urlCustomersAddAccount(customerID string) string {
	return fmt.Sprintf("/customers/%s/add-account", customerID)
}

func urlCustomersAccountsView(customerID, accountID string) string {
	return fmt.Sprintf("/customers/%s/accounts/%s", customerID, accountID)
}

func urlCustomersAccountTransactions(customerID, accountID string) string {
	return fmt.Sprintf("/customers/%s/accounts/%s/transactions", customerID, accountID)
}

func urlCustomersTransactionsCreate(customerID, accountID string) string {
	return fmt.Sprintf("/customers/%s/accounts/%s/transactions/create", customerID, accountID)
}

// Index handles listing all the customers.
func (h *Customers) Index(ctx context.Context, w http.ResponseWriter, r *http.Request, _ map[string]string) error {

	claims, err := auth.ClaimsFromContext(ctx)
	if err != nil {
		return err
	}

	fields := []datatable.DisplayField{
		{Field: "id", Title: "ID", Visible: false, Searchable: true, Orderable: true, Filterable: false},
		{Field: "name", Title: "Name", Visible: true, Searchable: true, Orderable: true, Filterable: true, FilterPlaceholder: "filter Name"},
		{Field: "email", Title: "Email", Visible: true, Searchable: true, Orderable: true, Filterable: true, FilterPlaceholder: "filter Email"},
		{Field: "phone_number", Title: "Phone Number", Visible: true, Searchable: true, Orderable: true, Filterable: true, FilterPlaceholder: "filter Phone Number"},
		{Field: "sales_rep", Title: "Manager", Visible: true, Searchable: true, Orderable: true, Filterable: true, FilterPlaceholder: "filter Manager"},
		{Field: "branch", Title: "Branch", Visible: true, Searchable: true, Orderable: true, Filterable: true, FilterPlaceholder: "filter Branch"},
	}

	mapFunc := func(q *customer.Response, cols []datatable.DisplayField) (resp []datatable.ColumnValue, err error) {
		for i := 0; i < len(cols); i++ {
			col := cols[i]
			var v datatable.ColumnValue
			switch col.Field {
			case "id":
				v.Value = fmt.Sprintf("%s", q.ID)
			case "name":
				v.Value = q.Name
				v.Formatted = fmt.Sprintf("<a href='%s'>%s</a>", urlCustomersView(q.ID), v.Value)
			case "email":
				v.Value = q.Email
				v.Formatted = q.Email
			case "phone_number":
				v.Value = q.PhoneNumber
				v.Formatted = q.PhoneNumber
			case "sales_rep":
				v.Value = q.SalesRep
				v.Formatted = fmt.Sprintf("<a href='%s'>%s</a>", urlUsersView(q.SalesRepID), v.Value)
			case "branch":
				v.Value = q.Branch
				v.Formatted = fmt.Sprintf("<a href='%s'>%s</a>", urlBranchesView(q.BranchID), v.Value)
			default:
				return resp, errors.Errorf("Failed to map value for %s.", col.Field)
			}
			resp = append(resp, v)
		}

		return resp, nil
	}

	loadFunc := func(ctx context.Context, sorting string, fields []datatable.DisplayField) (resp [][]datatable.ColumnValue, err error) {
		res, err := h.CustomerRepo.Find(ctx, claims, customer.FindRequest{
			Order: strings.Split(sorting, ","),
		})
		if err != nil {
			return resp, err
		}

		for _, a := range res.Customers {
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

	data := map[string]interface{}{
		"datatable":      dt.Response(),
		"urlCustomersCreate": urlCustomersCreate(),
		"urlCustomersIndex": urlCustomersIndex(),
	}

	return h.Renderer.Render(ctx, w, r, TmplLayoutBase, "customers-index.gohtml", web.MIMETextHTMLCharsetUTF8, http.StatusOK, data)
}

// Create handles creating a new customer.
func (h *Customers) Create(ctx context.Context, w http.ResponseWriter, r *http.Request, _ map[string]string) error {

	ctxValues, err := webcontext.ContextValues(ctx)
	if err != nil {
		return err
	}

	claims, err := auth.ClaimsFromContext(ctx)
	if err != nil {
		return err
	}

	//
	req := new(customer.CreateRequest)
	data := make(map[string]interface{})
	f := func() (bool, error) {
		if r.Method == http.MethodPost {
			err := r.ParseForm()
			if err != nil {
				return false, err
			}

			decoder := schema.NewDecoder()
			decoder.IgnoreUnknownKeys(true)

			if err := decoder.Decode(req, r.PostForm); err != nil {
				return false, err
			}

			res, err := h.CustomerRepo.Create(ctx, claims, *req, ctxValues.Now)
			if err != nil {
				switch errors.Cause(err) {
				default:
					if verr, ok := weberror.NewValidationError(ctx, err); ok {
						data["validationErrors"] = verr.(*weberror.Error)
						return false, nil
					} else {
						return false, err
					}
				}
			}

			accReq := account.CreateRequest{
				CustomerID: res.ID,
				Type:       req.Type,
				Target:     req.Target,
				TargetInfo: req.TargetInfo,
				BranchID:   req.BranchID,
			}

			_, err = h.AccountRepo.Create(ctx, claims, accReq, ctxValues.Now)
			if err != nil {
				// delete the created customer account
				_ = h.CustomerRepo.Delete(ctx, claims, customer.DeleteRequest{ID: res.ID}) // TODO: log delete error for debug
				cause := errors.Cause(err)
				switch cause {
				default:
					if verr, ok := weberror.NewValidationError(ctx, err); ok {
						data["validationErrors"] = verr.(*weberror.Error)
						return false, nil
					} else {
						return false, err
					}
				}
			}

			// Display a success message to the checklist.
			webcontext.SessionFlashSuccess(ctx,
				"Customer Created",
				"Customer successfully created.")

			return true, web.Redirect(ctx, w, r, urlCustomersView(res.ID), http.StatusFound)
		}

		return false, nil
	}

	end, err := f()
	if err != nil {
		return web.RenderError(ctx, w, r, err, h.Renderer, TmplLayoutBase, TmplContentErrorGeneric, web.MIMETextHTMLCharsetUTF8)
	} else if end {
		return nil
	}

	data["accountTypes"] = accountTypes
	data["form"] = req
	data["urlCustomersIndex"] = urlCustomersIndex()

	if verr, ok := weberror.NewValidationError(ctx, webcontext.Validator().Struct(customer.CreateRequest{})); ok {
		data["validationDefaults"] = verr.(*weberror.Error)
	}

	return h.Renderer.Render(ctx, w, r, TmplLayoutBase, "customers-create.gohtml", web.MIMETextHTMLCharsetUTF8, http.StatusOK, data)
}

// View handles displaying a customer.
func (h *Customers) View(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {

	ctxValues, err := webcontext.ContextValues(ctx)
	if err != nil {
		return err
	}

	customerID := params["customer_id"]

	claims, err := auth.ClaimsFromContext(ctx)
	if err != nil {
		return err
	}

	data := make(map[string]interface{})
	f := func() (bool, error) {
		if r.Method == http.MethodPost {
			err := r.ParseForm()
			if err != nil {
				return false, err
			}

			switch r.PostForm.Get("action") {
			case "archive":
				err = h.CustomerRepo.Archive(ctx, claims, customer.ArchiveRequest{
					ID: customerID,
				}, ctxValues.Now)
				if err != nil {
					return false, err
				}

				webcontext.SessionFlashSuccess(ctx,
					"Customer Archived",
					"Customer successfully archived.")

				return true, web.Redirect(ctx, w, r, urlCustomersIndex(), http.StatusFound)
			}
		}

		return false, nil
	}

	end, err := f()
	if err != nil {
		return web.RenderError(ctx, w, r, err, h.Renderer, TmplLayoutBase, TmplContentErrorGeneric, web.MIMETextHTMLCharsetUTF8)
	} else if end {
		return nil
	}

	cust, err := h.CustomerRepo.ReadByID(ctx, claims, customerID)
	if err != nil {
		return err
	}
	data["customer"] = cust.Response(ctx)

	accountsResp, err := h.AccountRepo.Find(ctx, claims, account.FindRequest{
		Where: "customer_id = ?", Args: []interface{}{customerID}, IncludeSalesRep: true, IncludeBranch: true,
	})
	if err != nil {
		return err
	}

	var accountBalance float64
	var txWhere []string
	var txArgs []interface{}
	for _, acc := range accountsResp.Accounts {
		accountBalance += acc.Balance
		txWhere = append(txWhere, "account_id = ?")
		txArgs = append(txArgs, acc.ID)
	}

	data["accounts"] = accountsResp.Accounts
	data["accountBalance"] = accountBalance

	var limit uint = 5
	var offset uint = 0
	var tranxListResp = &transaction.PagedResponseList{}
	// 0 length implies that the customer has no associated account
	if len(txWhere) > 0 {
		tranxListResp, err = h.TransactionRepo.Find(ctx, claims, transaction.FindRequest{
			Where:           strings.Join(txWhere, " OR "),
			Args:            txArgs,
			Order:           []string{"created_at desc"},
			Limit:           &limit,
			Offset:          &offset,
			IncludeAccount:  true,
			IncludeSalesRep: true,
		})
		if err != nil && err.Error() != sql.ErrNoRows.Error() {
			return err
		}
	}

	data["transactions"] = tranxListResp.Transactions

	data["urlCustomersIndex"] = urlCustomersIndex()
	data["urlCustomersUpdate"] = urlCustomersUpdate(customerID)
	data["urlCustomersView"] = urlCustomersView(customerID)
	data["urlCustomersAddAccount"] = urlCustomersAddAccount(customerID)
	data["urlCustomersTransactions"] = urlCustomersTransactions(customerID)
	var accountID string
	if len(accountsResp.Accounts) > 0 {
		accountID = accountsResp.Accounts[0].ID
	}
	data["urlCustomersTransactionsCreate"] = urlCustomersTransactionsCreate(customerID, accountID)

	return h.Renderer.Render(ctx, w, r, TmplLayoutBase, "customers-view.gohtml", web.MIMETextHTMLCharsetUTF8, http.StatusOK, data)
}

// Update handles updating a customer.
func (h *Customers) Update(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {

	ctxValues, err := webcontext.ContextValues(ctx)
	if err != nil {
		return err
	}

	customerID := params["customer_id"]

	claims, err := auth.ClaimsFromContext(ctx)
	if err != nil {
		return err
	}

	//
	req := new(customer.UpdateRequest)
	data := make(map[string]interface{})
	f := func() (bool, error) {
		if r.Method == http.MethodPost {
			err := r.ParseForm()
			if err != nil {
				return false, err
			}

			decoder := schema.NewDecoder()
			decoder.IgnoreUnknownKeys(true)

			if err := decoder.Decode(req, r.PostForm); err != nil {
				return false, err
			}
			req.ID = customerID

			err = h.CustomerRepo.Update(ctx, claims, *req, ctxValues.Now)
			if err != nil {
				switch errors.Cause(err) {
				default:
					if verr, ok := weberror.NewValidationError(ctx, err); ok {
						data["validationErrors"] = verr.(*weberror.Error)
						return false, nil
					} else {
						return false, err
					}
				}
			}

			webcontext.SessionFlashSuccess(ctx,
				"Customer Updated",
				"Customer successfully updated.")

			return true, web.Redirect(ctx, w, r, urlCustomersView(req.ID), http.StatusFound)
		}

		return false, nil
	}

	end, err := f()
	if err != nil {
		return web.RenderError(ctx, w, r, err, h.Renderer, TmplLayoutBase, TmplContentErrorGeneric, web.MIMETextHTMLCharsetUTF8)
	} else if end {
		return nil
	}

	prj, err := h.CustomerRepo.ReadByID(ctx, claims, customerID)
	if err != nil {
		return err
	}

	data["customer"] = prj.Response(ctx)

	data["urlCustomersIndex"] = urlCustomersIndex()
	data["urlCustomersView"] = urlCustomersView(customerID)

	if req.ID == "" {
		req.Name = &prj.Name
	}

	data["form"] = req

	if verr, ok := weberror.NewValidationError(ctx, webcontext.Validator().Struct(customer.UpdateRequest{})); ok {
		data["validationDefaults"] = verr.(*weberror.Error)
	}

	return h.Renderer.Render(ctx, w, r, TmplLayoutBase, "customers-update.gohtml", web.MIMETextHTMLCharsetUTF8, http.StatusOK, data)
}

// Transactions handles listing all the customers transactions across all his accounts.
func (h *Customers) Transactions(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {

	claims, err := auth.ClaimsFromContext(ctx)
	if err != nil {
		return err
	}

	customerID := params["customer_id"]

	fields := []datatable.DisplayField{
		{Field: "id", Title: "ID", Visible: false, Searchable: true, Orderable: true, Filterable: false},
		{Field: "amount", Title: "Amount", Visible: true, Searchable: false, Orderable: true, Filterable: true, FilterPlaceholder: "filter Amount"},
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
				v.Formatted = p.Sprintf("<a href='%s'>%.2f</a>", urlTransactionsView(q.ID), q.Amount)
			case "created_at":
				v.Value = q.CreatedAt.Local
				v.Formatted = q.CreatedAt.Local
			case "narration":
				v.Value = q.Narration
				v.Formatted = q.Narration
			case "account":
				v.Value = q.AccountNumber
				v.Formatted = fmt.Sprintf("<a href='%s'>%s</a>", urlCustomerAccountsView(q.AccountID), v.Value)
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

	customer, err := h.CustomerRepo.ReadByID(ctx, claims, customerID)
	if err != nil {
		return  err
	}

	accountsResp, err := h.AccountRepo.Find(ctx, claims, account.FindRequest{
		Where: "customer_id = ?", Args: []interface{}{customerID},
	})
	if err != nil {
		return err
	}

	var txWhere []string
	var txArgs []interface{}
	for _, acc := range accountsResp.Accounts {
		txWhere = append(txWhere, "account_id = ?")
		txArgs = append(txArgs, acc.ID)
	}

	loadFunc := func(ctx context.Context, sorting string, fields []datatable.DisplayField) (resp [][]datatable.ColumnValue, err error) {

		var res = &transaction.PagedResponseList{}
		// 0 where means this customer has no associated account
		if len(txWhere) > 0 {
			res, err = h.TransactionRepo.Find(ctx, claims, transaction.FindRequest{
				Order: strings.Split(sorting, ","), Where: strings.Join(txWhere, " OR "), Args: txArgs,
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

	var accountID string
	if len(accountsResp.Accounts) > 0 {
		accountID = accountsResp.Accounts[0].ID
	}

	data := map[string]interface{}{
		"customer":              customer,
		"datatable":             dt.Response(),
		"urlCustomersTransactionsCreate": urlCustomersTransactionsCreate(customerID, accountID),
		"urlCustomersIndex":     urlCustomersIndex(),
		"urlCustomersView":      urlCustomersView(customerID),
	}

	return h.Renderer.Render(ctx, w, r, TmplLayoutBase, "customers-transactions.gohtml", web.MIMETextHTMLCharsetUTF8, http.StatusOK, data)
}

// AddAccount handles add a new customer account.
func (h *Customers) AddAccount(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {

	ctxValues, err := webcontext.ContextValues(ctx)
	if err != nil {
		return err
	}

	customerID := params["customer_id"]

	claims, err := auth.ClaimsFromContext(ctx)
	if err != nil {
		return err
	}

	//
	req := new(account.CreateRequest)
	data := make(map[string]interface{})
	f := func() (bool, error) {
		if r.Method == http.MethodPost {
			err := r.ParseForm()
			if err != nil {
				return false, err
			}

			decoder := schema.NewDecoder()
			decoder.IgnoreUnknownKeys(true)

			if err := decoder.Decode(req, r.PostForm); err != nil {
				return false, err
			}
			req.CustomerID = customerID

			_, err = h.AccountRepo.Create(ctx, claims, *req, ctxValues.Now)
			if err != nil {
				switch errors.Cause(err) {
				default:
					if verr, ok := weberror.NewValidationError(ctx, err); ok {
						data["validationErrors"] = verr.(*weberror.Error)
						return false, nil
					} else {
						return false, err
					}
				}
			}

			// Display a success message to the checklist.
			webcontext.SessionFlashSuccess(ctx,
				"Account Added",
				"Account successfully Added.")

			return true, web.Redirect(ctx, w, r, urlCustomersView(customerID), http.StatusFound)
		}

		return false, nil
	}

	end, err := f()
	if err != nil {
		return web.RenderError(ctx, w, r, err, h.Renderer, TmplLayoutBase, TmplContentErrorGeneric, web.MIMETextHTMLCharsetUTF8)
	} else if end {
		return nil
	}

	customerRes, err := h.CustomerRepo.ReadByID(ctx, claims, customerID)
	if err != nil {
		return err
	}

	data["form"] = req
	data["accountTypes"] = accountTypes
	data["customer"] = customerRes
	data["urlCustomersIndex"] = urlCustomersIndex()
	data["urlCustomersView"] = urlCustomersView(customerID)

	if verr, ok := weberror.NewValidationError(ctx, webcontext.Validator().Struct(account.CreateRequest{})); ok {
		data["validationDefaults"] = verr.(*weberror.Error)
	}

	return h.Renderer.Render(ctx, w, r, TmplLayoutBase, "customers-add-account.gohtml", web.MIMETextHTMLCharsetUTF8, http.StatusOK, data)
}

// Account handles displaying an account for a customer
func (h *Customers) Account(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {

	ctxValues, err := webcontext.ContextValues(ctx)
	if err != nil {
		return err
	}

	customerID := params["customer_id"]
	accountID := params["account_id"]

	claims, err := auth.ClaimsFromContext(ctx)
	if err != nil {
		return err
	}

	data := make(map[string]interface{})
	f := func() (bool, error) {
		if r.Method == http.MethodPost {
			err := r.ParseForm()
			if err != nil {
				return false, err
			}

			switch r.PostForm.Get("action") {
			case "archive":
				err = h.CustomerRepo.Archive(ctx, claims, customer.ArchiveRequest{
					ID: customerID,
				}, ctxValues.Now)
				if err != nil {
					return false, err
				}

				webcontext.SessionFlashSuccess(ctx,
					"Customer Archived",
					"Customer successfully archived.")

				return true, web.Redirect(ctx, w, r, urlCustomersIndex(), http.StatusFound)
			}
		}

		return false, nil
	}

	end, err := f()
	if err != nil {
		return web.RenderError(ctx, w, r, err, h.Renderer, TmplLayoutBase, TmplContentErrorGeneric, web.MIMETextHTMLCharsetUTF8)
	} else if end {
		return nil
	}

	cust, err := h.CustomerRepo.ReadByID(ctx, claims, customerID)
	if err != nil {
		return err
	}
	data["customer"] = cust.Response(ctx)

	acc, err := h.AccountRepo.ReadByID(ctx, claims, accountID)
	if err != nil {
		return  weberror.NewError(ctx, err, 404)
	}
	data["account"] = acc.Response(ctx)

	var limit uint = 5
	var offset uint = 0
	tranxListResp, err := h.TransactionRepo.Find(ctx, claims, transaction.FindRequest{
		Where:           "account_id = ?",
		Args:            []interface{}{accountID},
		Order:           []string{"created_at desc"},
		Limit:           &limit,
		Offset:          &offset,
		IncludeAccount:  true,
		IncludeSalesRep: true,
	})
	if err != nil && err.Error() != sql.ErrNoRows.Error() {
		return err
	}

	data["transactions"] = tranxListResp.Transactions

	data["urlCustomersIndex"] = urlCustomersIndex()
	data["urlCustomersView"] = urlCustomersView(customerID)
	data["urlCustomersAccountTransactions"] = urlCustomersAccountTransactions(customerID, accountID)
	data["urlCustomersTransactionsCreate"] = urlCustomersTransactionsCreate(customerID, accountID)

	return h.Renderer.Render(ctx, w, r, TmplLayoutBase, "customers-account.gohtml", web.MIMETextHTMLCharsetUTF8, http.StatusOK, data)
}

// AccountTransactions handles listing all the transactions for the selected account.
func (h *Customers) AccountTransactions(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {

	claims, err := auth.ClaimsFromContext(ctx)
	if err != nil {
		return err
	}

	accountID := params["account_id"]

	acc, err := h.AccountRepo.ReadByID(ctx, claims, accountID)
	if err != nil {
		return  err
	}

	cust, err := h.CustomerRepo.ReadByID(ctx, claims, acc.CustomerID)
	if err != nil {
		return err
	}

	fields := []datatable.DisplayField{
		{Field: "id", Title: "ID", Visible: false, Searchable: true, Orderable: true, Filterable: false},
		{Field: "amount", Title: "Amount", Visible: true, Searchable: false, Orderable: true, Filterable: true, FilterPlaceholder: "filter Amount"},
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
				v.Formatted = p.Sprintf("<a href='%s'>%.2f</a>", urlTransactionsView(q.ID), q.Amount)
			case "created_at":
				v.Value = q.CreatedAt.Local
				v.Formatted = q.CreatedAt.Local
			case "narration":
				v.Value = q.Narration
				v.Formatted = q.Narration
			case "account":
				v.Value = q.AccountNumber
				v.Formatted = fmt.Sprintf("<a href='%s'>%s</a>", urlCustomerAccountsView(q.AccountID), v.Value)
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

	loadFunc := func(ctx context.Context, sorting string, fields []datatable.DisplayField) (resp [][]datatable.ColumnValue, err error) {

		var res = &transaction.PagedResponseList{}
		res, err = h.TransactionRepo.Find(ctx, claims, transaction.FindRequest{
			Order: strings.Split(sorting, ","),
			Where: "account_id = ?",
			Args: []interface{}{accountID},
		})
		if err != nil {
			return resp, err
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

	data := map[string]interface{}{
		"customer":                       cust.Response(ctx),
		"account":                        acc.Response(ctx),
		"datatable":                      dt.Response(),
		"urlCustomersTransactionsCreate": urlCustomersTransactionsCreate(cust.ID, accountID),
		"urlCustomersAccountsView":       urlCustomersAccountsView(cust.ID, accountID),
		"urlCustomersIndex":              urlCustomersIndex(),
		"urlCustomersView":               urlCustomersView(cust.ID),
	}

	return h.Renderer.Render(ctx, w, r, TmplLayoutBase, "customers-account-transactions.gohtml", web.MIMETextHTMLCharsetUTF8, http.StatusOK, data)
}
