package handlers

import (
	"context"
	"fmt"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"merryworld/surebank/internal/shop"
	"merryworld/surebank/internal/branch"
	"net/http"
	"strings"

	"merryworld/surebank/internal/platform/auth"
	"merryworld/surebank/internal/platform/datatable"
	"merryworld/surebank/internal/platform/web"
	"merryworld/surebank/internal/platform/web/webcontext"
	"merryworld/surebank/internal/platform/web/weberror"

	"github.com/gorilla/schema"
	"github.com/pkg/errors"
	"gopkg.in/DataDog/dd-trace-go.v1/contrib/go-redis/redis"
)

// Stocks represents the Stock API method handler set.
type Stocks struct {
	ShopRepo *shop.Repository
	BranchRepo *branch.Repository
	Redis    *redis.Client
	Renderer web.Renderer
}

func urlStocksIndex() string {
	return fmt.Sprintf("/shop/inventory")
}

func urlStocksCreate() string {
	return fmt.Sprintf("/shop/inventory/create")
}

func urlStocksView(id string) string {
	return fmt.Sprintf("/shop/inventory/%s", id)
}

func urlStocksUpdate(id string) string {
	return fmt.Sprintf("/shop/inventory/%s/update", id)
}

// Index handles listing all stocks.
func (h *Stocks) Index(ctx context.Context, w http.ResponseWriter, r *http.Request, _ map[string]string) error {

	fields := []datatable.DisplayField{
		{Field: "id", Title: "ID", Visible: false, Searchable: true, Orderable: true, Filterable: false},
		{Field: "product_id", Title: "Product ID", Visible: false, Searchable: true, Orderable: true, Filterable: false},
		{Field: "product_name", Title: "Product", Visible: true, Searchable: true, Orderable: true, Filterable: true, FilterPlaceholder: "filter Name"},
		{Field: "batch_number", Title: "Batch", Visible: true, Searchable: true, Orderable: false, Filterable: true, FilterPlaceholder: "filter Batch"},
		{Field: "quantity", Title: "Quantity", Visible: true, Searchable: false, Orderable: true, Filterable: false, },
		{Field: "unit_cost_price", Title: "Unit Cost Price", Visible: true, Searchable: false, Orderable: true, Filterable: false, },
		{Field: "manufacture_date", Title: "Manufacture Date", Visible: true, Searchable: false, Orderable: true, Filterable: false, },
		{Field: "expiry_date", Title: "Expiry Date", Visible: true, Searchable: false, Orderable: true, Filterable: false, },
	}

	mapFunc := func(q *shop.Stock, cols []datatable.DisplayField) (resp []datatable.ColumnValue, err error) {
		for i := 0; i < len(cols); i++ {
			col := cols[i]
			var v datatable.ColumnValue
			switch col.Field {
			case "id":
				v.Value = fmt.Sprintf("%s", q.ID)
			case "product_id":
				v.Value = fmt.Sprintf("%s", q.ProductID)
			case "product_name":
				v.Value = q.ProductName
				v.Formatted = fmt.Sprintf("<a href='%s'>%s</a>", urlStocksView(q.ID), v.Value)
			case "batch_number":
				v.Value = q.BatchNumber
				v.Formatted = q.BatchNumber
			case "quantity":
				v.Value = fmt.Sprintf("%d", q.Quantity)
				p := message.NewPrinter(language.English)
				v.Formatted = p.Sprintf("%.2d", q.Quantity)
			case "unit_cost_price":
				v.Value = fmt.Sprintf("%f", q.UnitCostPrice)
				p := message.NewPrinter(language.English)
				v.Formatted = p.Sprintf("%.2f", q.UnitCostPrice)
			case "manufacture_date":
				if q.ManufactureDate != nil {
					dt := web.NewTimeResponse(ctx, *q.ManufactureDate)
					v.Value = dt.Local
					v.Formatted = fmt.Sprintf("<span class='cell-font-date'>%s</span>", v.Value)
				} else {
					v.Value = "N/A"
					v.Formatted = "N/A"
				}
			case "expiry_date":
				if q.ExpiryDate != nil {
					dt := web.NewTimeResponse(ctx, *q.ExpiryDate)
					v.Value = dt.Local
					v.Formatted = fmt.Sprintf("<span class='cell-font-date'>%s</span>", v.Value)
				} else {
					v.Value = "N/A"
					v.Formatted = "N/A"
				}

			default:
				return resp, errors.Errorf("Failed to map value for %s.", col.Field)
			}
			resp = append(resp, v)
		}

		return resp, nil
	}

	loadFunc := func(ctx context.Context, sorting string, fields []datatable.DisplayField) (resp [][]datatable.ColumnValue, err error) {
		res, err := h.ShopRepo.FindStock(ctx, shop.StockFindRequest{
			Order: strings.Split(sorting, ","), IncludeProducts: true,
		})
		
		if err != nil {
			return resp, err
		}

		for _, a := range res {
			l, err := mapFunc(a, fields)
			if err != nil {
				return resp, errors.Wrapf(err, "Failed to map category for display.")
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
		"urlStocksCreate": urlStocksCreate(),
		"urlStocksIndex": urlStocksIndex(),
	}

	return h.Renderer.Render(ctx, w, r, TmplLayoutBase, "stocks-index.gohtml", web.MIMETextHTMLCharsetUTF8, http.StatusOK, data)
}

// Create handles creating a new stock.
func (h *Stocks) Create(ctx context.Context, w http.ResponseWriter, r *http.Request, _ map[string]string) error {

	ctxValues, err := webcontext.ContextValues(ctx)
	if err != nil {
		return err
	}

	claims, err := auth.ClaimsFromContext(ctx)
	if err != nil {
		return err
	} 

	//
	req := new(shop.StockCreateRequest)
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
				return false, errors.WithMessage(err, "Something wrong")
			}

			resp, err := h.ShopRepo.CreateStock(ctx, claims, *req, ctxValues.Now)
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

			// Display a success message to the product.
			webcontext.SessionFlashSuccess(ctx,
				"Stock Created",
				"Stock successfully created.")

			return true, web.Redirect(ctx, w, r, urlStocksView(resp.ID), http.StatusFound)
		}

		return false, nil
	}

	end, err := f()
	if err != nil {
		return web.RenderError(ctx, w, r, err, h.Renderer, TmplLayoutBase, TmplContentErrorGeneric, web.MIMETextHTMLCharsetUTF8)
	} else if end {
		return nil
	}

	data["products"], err = h.ShopRepo.FindProduct(ctx, shop.ProductFindRequest{ Order: []string{"name"} })
	if err != nil {
		return err
	}

	data["branches"], err = h.BranchRepo.Find(ctx, claims, branch.FindRequest{
		Order: []string{"name asc"},
	})
	if err != nil {
		return err 
	}

	data["form"] = req
	data["urlStocksIndex"] = urlStocksIndex()

	if verr, ok := weberror.NewValidationError(ctx, webcontext.Validator().Struct(shop.StockCreateRequest{})); ok {
		data["validationDefaults"] = verr.(*weberror.Error)
	}

	return h.Renderer.Render(ctx, w, r, TmplLayoutBase, "stocks-create.gohtml", web.MIMETextHTMLCharsetUTF8, http.StatusOK, data)
}

// View handles displaying a stock.
func (h *Stocks) View(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {

	id := params["stock_id"]

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
				err = h.ShopRepo.DeleteStock(ctx, claims, shop.StockDeleteRequest{
					ID: id,
				})
				if err != nil {
					return false, err
				}

				webcontext.SessionFlashSuccess(ctx,
					"Stock Deleted",
					"Stock successfully deleted.")

				return true, web.Redirect(ctx, w, r, urlStocksIndex(), http.StatusFound)
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

	prj, err := h.ShopRepo.ReadStockByID(ctx, claims, id)
	if err != nil {
		return err
	}
	data["stock"] = prj.Response(ctx)
	data["urlStocksIndex"] = urlStocksIndex()
	data["urlStocksView"] = urlStocksView(id)
	data["urlStocksUpdate"] = urlStocksUpdate(id)

	return h.Renderer.Render(ctx, w, r, TmplLayoutBase, "stocks-view.gohtml", web.MIMETextHTMLCharsetUTF8, http.StatusOK, data)
}

// Update handles updating a stock.
func (h *Stocks) Update(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {

	ctxValues, err := webcontext.ContextValues(ctx)
	if err != nil {
		return err
	}

	id := params["stock_id"]

	claims, err := auth.ClaimsFromContext(ctx)
	if err != nil {
		return err
	}

	//
	req := new(shop.StockUpdateRequest)
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
			req.ID = id

			err = h.ShopRepo.UpdateStock(ctx, claims, *req, ctxValues.Now)
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
				"Stock Updated",
				"Stock successfully updated.")

			return true, web.Redirect(ctx, w, r, urlStocksView(req.ID), http.StatusFound)
		}

		return false, nil
	}

	end, err := f()
	if err != nil {
		return web.RenderError(ctx, w, r, err, h.Renderer, TmplLayoutBase, TmplContentErrorGeneric, web.MIMETextHTMLCharsetUTF8)
	} else if end {
		return nil
	}

	prj, err := h.ShopRepo.ReadStockByID(ctx, claims, id)
	if err != nil {
		return err
	}

	data["stock"] = prj.Response(ctx)

	data["products"], err = h.ShopRepo.FindProduct(ctx, shop.ProductFindRequest{ Order: []string{"name"} })
	if err != nil {
		return err
	}

	data["branches"], err = h.BranchRepo.Find(ctx, claims, branch.FindRequest{
		Order: []string{"name asc"},
	})
	if err != nil {
		return err 
	}

	data["urlStocksIndex"] = urlStocksIndex()
	data["urlStocksView"] = urlStocksView(id)

	if req.ID == "" {
		req.UnitCostPrice = &prj.UnitCostPrice
		req.ExpiryDate = prj.ExpiryDate
		req.ManufactureDate = prj.ManufactureDate
		req.Quantity = &prj.Quantity
		req.BatchNumber = &prj.BatchNumber
	}

	data["form"] = req

	if verr, ok := weberror.NewValidationError(ctx, webcontext.Validator().Struct(shop.StockUpdateRequest{})); ok {
		data["validationDefaults"] = verr.(*weberror.Error)
	}

	return h.Renderer.Render(ctx, w, r, TmplLayoutBase, "stocks-update.gohtml", web.MIMETextHTMLCharsetUTF8, http.StatusOK, data)
}

// Index handles listing all stocks.
func (h *Stocks) Report(ctx context.Context, w http.ResponseWriter, r *http.Request, _ map[string]string) error {

	claims, err := auth.ClaimsFromContext(ctx)
	if err != nil { 
		return err
	}

	fields := []datatable.DisplayField{
		{Field: "product_id", Title: "Product ID", Visible: false, Searchable: true, Orderable: true, Filterable: false},
		{Field: "product_name", Title: "Product", Visible: true, Searchable: true, Orderable: true, Filterable: true, FilterPlaceholder: "filter Name"},
		{Field: "quantity", Title: "Quantity", Visible: true, Searchable: false, Orderable: true, Filterable: false, },
	}

	mapFunc := func(q shop.StockInfo, cols []datatable.DisplayField) (resp []datatable.ColumnValue, err error) {
		for i := 0; i < len(cols); i++ {
			col := cols[i]
			var v datatable.ColumnValue
			switch col.Field {
			case "product_id":
				v.Value = fmt.Sprintf("%s", q.ProductID)
			case "product_name":
				v.Value = q.ProductName
				v.Formatted = q.ProductName
			case "quantity":
				v.Value = fmt.Sprintf("%d", q.Quantity)
				p := message.NewPrinter(language.English)
				v.Formatted = p.Sprintf("%.2d", q.Quantity)
			default:
				return resp, errors.Errorf("Failed to map value for %s.", col.Field)
			}
			resp = append(resp, v)
		}

		return resp, nil
	}

	loadFunc := func(ctx context.Context, sorting string, fields []datatable.DisplayField) (resp [][]datatable.ColumnValue, err error) {
		res, err := h.ShopRepo.StockReport(ctx, claims, shop.StockReportRequest{
			Order: strings.Split(sorting, ","),
		})
		
		if err != nil {
			return resp, err
		}

		for _, a := range res {
			l, err := mapFunc(a, fields)
			if err != nil {
				return resp, errors.Wrapf(err, "Failed to map category for display.")
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
		"urlStocksCreate": urlStocksCreate(),
		"urlStocksIndex": urlStocksIndex(),
	}

	return h.Renderer.Render(ctx, w, r, TmplLayoutBase, "stocks-report.gohtml", web.MIMETextHTMLCharsetUTF8, http.StatusOK, data)
}