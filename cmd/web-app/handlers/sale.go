package handlers

import (
	"context"
	"fmt"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"merryworld/surebank/internal/platform/datatable"
	"merryworld/surebank/internal/shop"
	"net/http"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/DataDog/dd-trace-go.v1/contrib/go-redis/redis"

	"merryworld/surebank/internal/platform/auth"
	"merryworld/surebank/internal/platform/web"
	"merryworld/surebank/internal/platform/web/webcontext"
	"merryworld/surebank/internal/platform/web/weberror"
	"merryworld/surebank/internal/sale"
)

// Sales represents the sales API method handler set.
type Sales struct {
	Repository *sale.Repository
	ShopRepo   *shop.Repository
	Redis      *redis.Client
	Renderer   web.Renderer 
}

func urlSalesIndex() string {
	return fmt.Sprintf("/sales")
}

func urlSalesView(saleID string) string {
	return fmt.Sprintf("/sales/%s", saleID)
}

// Index handles listing all the customers.
func (h *Sales) Index(ctx context.Context, w http.ResponseWriter, r *http.Request, _ map[string]string) error {

	claims, err := auth.ClaimsFromContext(ctx)
	if err != nil {
		return err
	}

	fields := []datatable.DisplayField{
		{Field: "id", Title: "ID", Visible: false, Searchable: false, Orderable: false, Filterable: false},
		{Field: "receipt_number", Title: "Receipt Number", Visible: true, Searchable: true, Orderable: false, Filterable: true, FilterPlaceholder: "filter Receipt"},
		{Field: "amount", Title: "Amount", Visible: true, Searchable: false, Orderable: true, Filterable: false},
		{Field: "customer_name", Title: "Customer Name", Visible: true, Searchable: true, Orderable: true, Filterable: true, FilterPlaceholder: "filter Customer name"},
		{Field: "phone_number", Title: "Phone Number", Visible: true, Searchable: true, Orderable: true, Filterable: true, FilterPlaceholder: "filter Phone Number"},
		{Field: "sales_rep", Title: "Sales Rep", Visible: true, Searchable: true, Orderable: false, Filterable: true, FilterPlaceholder: "filter Sales Rep"},
		{Field: "branch", Title: "Branch", Visible: true, Searchable: true, Orderable: false, Filterable: true, FilterPlaceholder: "filter Branch"},
		{Field: "created_at", Title: "Date", Visible: true, Searchable: true, Orderable: true, Filterable: true, FilterPlaceholder: "filter Date"},
	}

	mapFunc := func(q *sale.Response, cols []datatable.DisplayField) (resp []datatable.ColumnValue, err error) {
		for i := 0; i < len(cols); i++ {
			col := cols[i]
			var v datatable.ColumnValue
			switch col.Field {
			case "id":
				v.Value = fmt.Sprintf("%s", q.ID)
			case "receipt_number":
				v.Value = q.ReceiptNumber
				v.Formatted = fmt.Sprintf("<a href='%s'>%s</a>", urlSalesView(q.ID), v.Value)
			case "amount":
				p := message.NewPrinter(language.English)
				v.Value = p.Sprintf("%.2f", q.Amount)
				v.Formatted = v.Value
			case "customer_name":
				v.Value = q.CustomerName
				v.Formatted = v.Value
			case "phone_number":
				v.Value = q.PhoneNumber
				v.Formatted = v.Value
			case "sales_rep":
				if q.CreatedBy != nil {
					v.Value = *q.CreatedBy
					v.Formatted = fmt.Sprintf("<a href='%s'>%s</a>", urlUsersView(q.CreatedByID), v.Value)
				}
			case "branch":
				if q.Branch != nil {
					v.Value = *q.Branch
					v.Formatted = fmt.Sprintf("<a href='%s'>%s</a>", urlBranchesView(q.BranchID), v.Value)
				}
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
		res, err := h.Repository.Find(ctx, claims, sale.FindRequest{
			Order: order,
			IncludeBranch:true,
			IncludeCreatedBy:true,
		})
		if err != nil {
			return resp, err
		}

		for _, a := range res.Sales {
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

	products, err := h.ShopRepo.FindProduct(ctx, shop.ProductFindRequest{Order:[]string{"name DESC"}})
	if err != nil {
		return err
	}

	data := map[string]interface{}{
		"datatable":      dt.Response(),
		"urlCustomersCreate": urlCustomersCreate(),
		"urlCustomersIndex": urlCustomersIndex(),
		"products": products,
	}

	return h.Renderer.Render(ctx, w, r, TmplLayoutBase, "sales-index.gohtml", web.MIMETextHTMLCharsetUTF8, http.StatusOK, data)
}

// Sell handle request for making new sales
func (h *Sales) Sell(ctx context.Context, w http.ResponseWriter, r *http.Request, _ map[string]string) error {
	v, err := webcontext.ContextValues(ctx)
	if err != nil {
		return err
	}

	claims, err := auth.ClaimsFromContext(ctx)
	if err != nil {
		return err
	}

	var req sale.MakeSalesRequest
	if err := web.Decode(ctx, r, &req); err != nil {
		if _, ok := errors.Cause(err).(*weberror.Error); !ok {
			err = weberror.NewError(ctx, err, http.StatusBadRequest)
		}
		return web.RespondJsonError(ctx, w, err)
	}

	res, err := h.Repository.MakeSale(ctx, claims, req, v.Now)
	if err != nil {
		cause := errors.Cause(err)
		switch cause {
		case sale.ErrForbidden:
			return web.RespondJsonError(ctx, w, weberror.NewError(ctx, err, http.StatusForbidden))
		default:
			return web.RespondJsonError(ctx, w, weberror.NewError(ctx, err, http.StatusBadRequest))
			/*_, ok := cause.(validator.ValidationErrors)
			if ok {
				return web.RespondJsonError(ctx, w, weberror.NewError(ctx, err, http.StatusBadRequest))
			}
			return errors.Wrapf(err, "Customer: %+v", &req)*/
		}
	}

	result := res.Response(ctx)
	return web.RespondJson(ctx, w, result, http.StatusCreated)
}

// View handles displaying a sale.
func (h *Sales) View(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {

	ctxValues, err := webcontext.ContextValues(ctx)
	if err != nil {
		return err
	}

	saleID := params["sale_id"]

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
				err = h.Repository.Archive(ctx, claims, sale.ArchiveRequest{
					ID: saleID,
				}, ctxValues.Now)
				if err != nil {
					return false, err
				}

				webcontext.SessionFlashSuccess(ctx,
					"Sale Archived",
					"Sale successfully archived.")

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

	salesDetail, err := h.Repository.ReadByID(ctx, claims, saleID)
	if err != nil {
		return err
	}
	data["sale"] = salesDetail.Response(ctx)

	data["urlSalesIndex"] = urlSalesIndex()

	return h.Renderer.Render(ctx, w, r, TmplLayoutBase, "sales-view.gohtml", web.MIMETextHTMLCharsetUTF8, http.StatusOK, data)
}
