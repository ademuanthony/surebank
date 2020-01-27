package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"merryworld/surebank/internal/account"
	"merryworld/surebank/internal/platform/auth"
	"merryworld/surebank/internal/platform/datatable"
	"merryworld/surebank/internal/platform/web"
	"merryworld/surebank/internal/platform/web/webcontext"
	"merryworld/surebank/internal/platform/web/weberror"

	"github.com/gorilla/schema"
	"github.com/pkg/errors"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"gopkg.in/DataDog/dd-trace-go.v1/contrib/go-redis/redis"
)

// CustomerAccounts represents the Customer Accounts API method handler set.
type CustomerAccounts struct {
	Repository *account.Repository
	Redis    *redis.Client
	Renderer web.Renderer
}

func urlCustomerAccountsIndex() string {
	return fmt.Sprintf("/customer-accounts")
}

func urlCustomerAccountsCreate() string {
	return fmt.Sprintf("/customer-accounts/create")
}

func urlCustomerAccountsView(brandID string) string {
	return fmt.Sprintf("/customer-accounts/%s", brandID)
}

// Index handles listing all the customer accounts.
func (h *CustomerAccounts) Index(ctx context.Context, w http.ResponseWriter, r *http.Request, _ map[string]string) error {

	claims, err := auth.ClaimsFromContext(ctx)
	if err != nil {
		return err
	}

	fields := []datatable.DisplayField{
		{Field: "id", Title: "ID", Visible: false, Searchable: true, Orderable: true, Filterable: false},
		{Field: "customer_name", Title: "Customer Name", Visible: true, Searchable: true, Orderable: true, Filterable: true, FilterPlaceholder: "filter Customer Name"},
		{Field: "number", Title: "Account Number", Visible: true, Searchable: true, Orderable: true, Filterable: true, FilterPlaceholder: "filter Account Number"},
		{Field: "type", Title: "Account Type", Visible: true, Searchable: true, Orderable: true, Filterable: true, FilterPlaceholder: "filter Account Type"},
		{Field: "account_manager", Title: "Account Manager", Visible: true, Searchable: true, Orderable: true, Filterable: true, FilterPlaceholder: "filter Account Manager"},
		{Field: "target_amount", Title: "Target Amount", Visible: true, Searchable: false, Orderable: true, Filterable: false},
		{Field: "target_info", Title: "Target Info", Visible: true, Searchable: true, Orderable: true, Filterable: true, FilterPlaceholder: "filter Target"},
		{Field: "balance", Title: "Account Balance", Visible: true, Searchable: false, Orderable: true, Filterable: false},
	}

	mapFunc := func(q *account.Account, cols []datatable.DisplayField) (resp []datatable.ColumnValue, err error) {
		for i := 0; i < len(cols); i++ {
			col := cols[i]
			var v datatable.ColumnValue
			switch col.Field {
			case "id":
				v.Value = fmt.Sprintf("%s", q.ID)
			case "customer_name":
				v.Value = "N/A"
				if q.Customer != nil {
					v.Value = q.Customer.Name
				}
				v.Formatted = fmt.Sprintf("<a href='%s'>%s</a>", urlCustomersView(q.CustomerID), v.Value)
			case "number":
				v.Value = q.Number
				v.Formatted = fmt.Sprintf("<a href='%s'>%s</a>", urlCustomerAccountsView(q.ID), v.Value)
			case "type":
				v.Value = q.Type
				v.Formatted = q.Type
			case "account_manager":
				v.Value = "N/A"
				if q.Customer != nil {
					v.Value = q.SalesRep.LastName + " " + q.SalesRep.FirstName
				}
				v.Formatted = fmt.Sprintf("<a href='%s'>%s</a>", urlUsersView(q.SalesRepID), v.Value)
			case "target_amount":
				v.Value = fmt.Sprintf("%f", q.Target)
				p := message.NewPrinter(language.English)
				v.Formatted = p.Sprintf("%.2f", q.Target)
			case "target_info":
				v.Value = q.TargetInfo
				v.Formatted = q.TargetInfo
			case "balance":
				v.Value = fmt.Sprintf("%f", q.Balance)
				p := message.NewPrinter(language.English)
				v.Formatted = p.Sprintf("%.2f", q.Balance)
			default:
				return resp, errors.Errorf("Failed to map value for %s.", col.Field)
			}
			resp = append(resp, v)
		}

		return resp, nil
	}

	loadFunc := func(ctx context.Context, sorting string, fields []datatable.DisplayField) (resp [][]datatable.ColumnValue, err error) {
		res, err := h.Repository.Find(ctx, claims, account.FindRequest{
			Order: strings.Split(sorting, ","),
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

	data := map[string]interface{}{
		"datatable":      dt.Response(),
		"urlCustomerAccountsCreate": urlCustomersCreate(),
		"urlCustomerAccountsIndex": urlCustomerAccountsIndex(),
	}

	return h.Renderer.Render(ctx, w, r, TmplLayoutBase, "customer-accounts-index.gohtml", web.MIMETextHTMLCharsetUTF8, http.StatusOK, data)
}

// Create handles creating a new customer account.
func (h *CustomerAccounts) Create(ctx context.Context, w http.ResponseWriter, r *http.Request, _ map[string]string) error {

	ctxValues, err := webcontext.ContextValues(ctx)
	if err != nil {
		return err
	}

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

			usr, err := h.Repository.Create(ctx, claims, *req, ctxValues.Now)
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
				"Account Created",
				"Account successfully created.")

			return true, web.Redirect(ctx, w, r, urlCustomerAccountsView(usr.ID), http.StatusFound)
		}

		return false, nil
	}

	end, err := f()
	if err != nil {
		return web.RenderError(ctx, w, r, err, h.Renderer, TmplLayoutBase, TmplContentErrorGeneric, web.MIMETextHTMLCharsetUTF8)
	} else if end {
		return nil
	}

	data["form"] = req
	data["urlCustomerAccountsIndex"] = urlCustomerAccountsIndex()

	if verr, ok := weberror.NewValidationError(ctx, webcontext.Validator().Struct(account.CreateRequest{})); ok {
		data["validationDefaults"] = verr.(*weberror.Error)
	}

	return h.Renderer.Render(ctx, w, r, TmplLayoutBase, "customer-accounts-create.gohtml", web.MIMETextHTMLCharsetUTF8, http.StatusOK, data)
}

// View handles displaying a customer account.
func (h *CustomerAccounts) View(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {

	ctxValues, err := webcontext.ContextValues(ctx)
	if err != nil {
		return  err
	}

	brandID := params["account_id"]

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
				err = h.Repository.Archive(ctx, claims, account.ArchiveRequest{
					ID: brandID,
				}, ctxValues.Now)
				if err != nil {
					return false, err
				}

				webcontext.SessionFlashSuccess(ctx,
					"Account Archived",
					"Account successfully archived.")

				return true, web.Redirect(ctx, w, r, urlCustomerAccountsIndex(), http.StatusFound)
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

	prj, err := h.Repository.ReadByID(ctx, claims, brandID)
	if err != nil {
		return err
	}
	data["brand"] = prj.Response(ctx)
	data["urlCustomerAccountsIndex"] = urlCustomerAccountsIndex()
	data["urlCustomerAccountsView"] = urlCustomerAccountsView(brandID)

	return h.Renderer.Render(ctx, w, r, TmplLayoutBase, "customer-accounts-view.gohtml", web.MIMETextHTMLCharsetUTF8, http.StatusOK, data)
}
