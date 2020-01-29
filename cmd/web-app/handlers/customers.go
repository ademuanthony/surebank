package handlers

import (
	"context"
	"fmt"
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
	CustomerRepo *customer.Repository
	AccountRepo  *account.Repository
	Renderer     web.Renderer
	Redis        *redis.Client
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

	data["accountTypes"] = []string{"SB", "OM"}
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

	data["urlCustomersIndex"] = urlCustomersIndex()
	data["urlCustomersUpdate"] = urlCustomersUpdate(customerID)
	data["urlCustomersView"] = urlCustomersView(customerID)

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
