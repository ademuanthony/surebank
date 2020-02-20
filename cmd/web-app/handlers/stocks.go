package handlers

import (
	"context"
	"fmt"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"merryworld/surebank/internal/branch"
	"merryworld/surebank/internal/inventory"
	"merryworld/surebank/internal/shop"
	"merryworld/surebank/internal/transaction"
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

// Stocks represents the Inventory API method handler set.
type Stocks struct {
	Repo       *inventory.Repository
	ShopRepo   *shop.Repository
	BranchRepo *branch.Repository
	Redis      *redis.Client
	Renderer   web.Renderer
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

// Index handles listing all stock transactions.
func (h *Stocks) Index(ctx context.Context, w http.ResponseWriter, r *http.Request, _ map[string]string) error {

	claims, err := auth.ClaimsFromContext(ctx)
	if err != nil {
		return err
	}

	fields := []datatable.DisplayField{
		{Field: "id", Title: "ID", Visible: false, Searchable: true, Orderable: true, Filterable: false},
		{Field: "product_id", Title: "Product ID", Visible: false, Searchable: false, Orderable: false, Filterable: false},
		{Field: "product_name", Title: "Product", Visible: true, Searchable: true, Orderable: true, Filterable: true, FilterPlaceholder: "filter Name"},
		{Field: "branch_id", Title: "Branch ID", Visible: false, Searchable: false, Orderable: false, Filterable: false},
		{Field: "branch_name", Title: "Branch", Visible: true, Searchable: true, Orderable: true, Filterable: true, FilterPlaceholder: "filter Branch"},
		{Field: "opening_balance", Title: "Opening Balance", Visible: true, Searchable: false, Orderable: false, Filterable: false,},
		{Field: "quantity", Title: "Quantity", Visible: true, Searchable: false, Orderable: true, Filterable: false, },
		{Field: "type", Title: "Type", Visible: true, Searchable: false, Orderable: true, Filterable: true, FilterPlaceholder: "filter Transaction Type"},
		{Field: "created_at", Title: "Date", Visible: true, Searchable: false, Orderable: true, Filterable: true, },
		{Field: "sales_rep", Title: "Added By", Visible: true, Searchable: false, Orderable: true, Filterable: true, },
	}

	mapFunc := func(q *inventory.Response, cols []datatable.DisplayField) (resp []datatable.ColumnValue, err error) {
		for i := 0; i < len(cols); i++ {
			col := cols[i]
			var v datatable.ColumnValue
			switch col.Field {
			case "id":
				v.Value = fmt.Sprintf("%s", q.ID)
			case "product_id":
				v.Value = fmt.Sprintf("%s", q.ProductID)
			case "product_name":
				v.Value = q.Product
				v.Formatted = fmt.Sprintf("<a href='%s'>%s</a>", urlProductsView(q.ProductID), v.Value)
			case "branch_id":
				v.Value = fmt.Sprintf("%s", q.BranchID)
			case "branch_name":
				v.Value = q.Branch
				v.Formatted = fmt.Sprintf("<a href='%s'>%s</a>", urlBranchesView(q.BranchID), v.Value)
			case "opening_balance":
				v.Value = fmt.Sprintf("%d", q.OpeningBalance)
				p := message.NewPrinter(language.English)
				v.Formatted = p.Sprintf("%d", q.OpeningBalance)
				v.Formatted = v.Value
			case "quantity":
				v.Value = fmt.Sprintf("%d", q.Quantity)
				v.Formatted = fmt.Sprintf("<a href='%s'>%s</a>", urlStocksView(q.ID), v.Value)
			case "type":
				v.Value = q.TXType
				v.Formatted = v.Value
			case "created_at":
				v.Value = q.CreatedAt.Local
				v.Formatted = fmt.Sprintf("<a href='%s'>%s</a>", urlStocksView(q.ID), v.Value)
			case "sales_rep":
				v.Value = q.SalesRep
				v.Formatted = fmt.Sprintf("<a href='%s'>%s</a>", urlUsersView(q.SalesRepID), v.Value)

			default:
				return resp, errors.Errorf("Failed to map value for %s.", col.Field)
			}
			resp = append(resp, v)
		}

		return resp, nil
	}

	loadFunc := func(ctx context.Context, sorting string, fields []datatable.DisplayField) (resp [][]datatable.ColumnValue, err error) {
		if len(sorting) == 0 {
			sorting = "created_at desc"
		}
		res, err := h.Repo.Find(ctx, claims, inventory.FindRequest{
			Order: strings.Split(sorting, ","), IncludeProduct: true, IncludeBranch: true, IncludeSalesRep: true,
		})

		if err != nil {
			return resp, err
		}

		for _, a := range res.Transactions {
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

// Create handles creating a new stock transaction.
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
	req := new(inventory.AddStockRequest)
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

			resp, err := h.Repo.AddStock(ctx, claims, *req, ctxValues.Now)
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
				"Inventory Created",
				"Inventory successfully created.")

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

	data["form"] = req
	data["urlStocksIndex"] = urlStocksIndex()

	if verr, ok := weberror.NewValidationError(ctx, webcontext.Validator().Struct(inventory.AddStockRequest{})); ok {
		data["validationDefaults"] = verr.(*weberror.Error)
	}

	return h.Renderer.Render(ctx, w, r, TmplLayoutBase, "stocks-create.gohtml", web.MIMETextHTMLCharsetUTF8, http.StatusOK, data)
}

// View handles displaying a stock.
func (h *Stocks) View(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {

	id := params["stock_id"]

	ctxValues, err := webcontext.ContextValues(ctx)
	if err != nil {
		return err
	}

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
				err = h.Repo.Archive(ctx, claims, inventory.ArchiveRequest{
					ID: id,
				}, ctxValues.Now)
				if err != nil {
					return false, err
				}

				webcontext.SessionFlashSuccess(ctx,
					"Inventory Deleted",
					"Inventory successfully deleted.")

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

	prj, err := h.Repo.ReadByID(ctx, claims, id)
	if err != nil {
		return err
	}
	data["stock"] = prj.Response(ctx)
	data["urlStocksCreate"] = urlStocksCreate()
	data["urlStocksIndex"] = urlStocksIndex()
	data["urlStocksView"] = urlStocksView(id)
	data["urlStocksUpdate"] = urlStocksUpdate(id)

	return h.Renderer.Render(ctx, w, r, TmplLayoutBase, "stocks-view.gohtml", web.MIMETextHTMLCharsetUTF8, http.StatusOK, data)
}

// Report handles listing all stock balance.
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

	mapFunc := func(q inventory.StockInfo, cols []datatable.DisplayField) (resp []datatable.ColumnValue, err error) {
		for i := 0; i < len(cols); i++ {
			col := cols[i]
			var v datatable.ColumnValue
			switch col.Field {
			case "product_id":
				v.Value = fmt.Sprintf("%s", q.ProductID)
			case "product_name":
				v.Value = q.ProductName
				v.Formatted = fmt.Sprintf("<a href='%s'>%s</a>", urlProductsView(q.ProductID), v.Value)
			case "quantity":
				var qnt int64
				if q.TxType == transaction.TransactionType_Deposit.String() {
					qnt = q.OpeningBalance + q.Quantity
				} else {
					qnt = q.OpeningBalance - q.Quantity
				}
				v.Value = fmt.Sprintf("%d", qnt)
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
			order = order
		} else {
			order = append(order, "p.name")
		}

		res, err := h.Repo.Report(ctx, claims, inventory.ReportRequest{
			Order: order,
		})

		if err != nil {
			return resp, err
		}

		for _, a := range res.Inventories {
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
