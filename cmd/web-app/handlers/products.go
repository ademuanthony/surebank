package handlers

import (
	"context"
	"fmt"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"merryworld/surebank/internal/shop"
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

// Products represents the Product API method handler set.
type Products struct {
	ShopRepo *shop.Repository
	Redis    *redis.Client
	Renderer web.Renderer
}

func urlProductsIndex() string {
	return fmt.Sprintf("/shop/products")
}

func urlProductsCreate() string {
	return fmt.Sprintf("/shop/products/create")
}

func urlProductsView(productID string) string {
	return fmt.Sprintf("/shop/products/%s", productID)
}

func urlProductsUpdate(categoryID string) string {
	return fmt.Sprintf("/shop/products/%s/update", categoryID)
}

// Index handles listing all products.
func (h *Products) Index(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {

	dt, ok, err := productDatatable(ctx, h.ShopRepo, h.Redis, w, r, "", nil)
	if ok {
		if err != nil {
			return err
		}
		return nil
	}

	if err != nil {
		return err
	}

	data := map[string]interface{}{
		"datatable":      dt.Response(),
		"urlProductsCreate": urlProductsCreate(),
		"urlProductsIndex": urlProductsIndex(),
	}

	return h.Renderer.Render(ctx, w, r, TmplLayoutBase, "products-index.gohtml", web.MIMETextHTMLCharsetUTF8, http.StatusOK, data)
}

func productDatatable(ctx context.Context, shopRepo *shop.Repository, redisClient *redis.Client, w http.ResponseWriter, r *http.Request,
	where string, args []interface{}) (*datatable.Datatable, bool, error) {

	fields := []datatable.DisplayField{
		{Field: "id", Title: "ID", Visible: false, Searchable: true, Orderable: true, Filterable: false},
		{Field: "name", Title: "Product", Visible: true, Searchable: true, Orderable: true, Filterable: true, FilterPlaceholder: "filter Name"},
		{Field: "sku", Title: "SKU", Visible: true, Searchable: true, Orderable: false, Filterable: true, FilterPlaceholder: "filter SKU"},
		{Field: "price", Title: "Price", Visible: true, Searchable: false, Orderable: true, Filterable: true, FilterPlaceholder: "filter Price"},
	}

	mapFunc := func(q *shop.Product, cols []datatable.DisplayField) (resp []datatable.ColumnValue, err error) {

		for i := 0; i < len(cols); i++ {
			col := cols[i]
			var v datatable.ColumnValue
			switch col.Field {
			case "id":
				v.Value = fmt.Sprintf("%s", q.ID)
			case "name":
				v.Value = q.Name
				v.Formatted = fmt.Sprintf("<a href='%s'>%s</a>", urlProductsView(q.ID), v.Value)
			case "sku":
				v.Value = q.Sku
				v.Formatted = q.Sku
			case "price":
				v.Value = fmt.Sprintf("%f", q.Price)
				p := message.NewPrinter(language.English)
				v.Formatted = p.Sprintf("%.2f", q.Price)
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

		res, err := shopRepo.FindProduct(ctx, shop.ProductFindRequest{
			Order: order,
			Where: where,
			Args: args,
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

	dt, err := datatable.New(ctx, w, r, redisClient, fields, loadFunc)
	if err != nil {
		return nil, false, err
	}

	if dt.HasCache() {
		return nil, false, nil
	}

	ok, err := dt.Render()

	return dt, ok, err
}

// Create handles creating a new product.
func (h *Products) Create(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {

	ctxValues, err := webcontext.ContextValues(ctx)
	if err != nil {
		return err
	}

	claims, err := auth.ClaimsFromContext(ctx)
	if err != nil {
		return err
	}

	//
	req := new(shop.ProductCreateRequest)
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

			resp, err := h.ShopRepo.CreateProduct(ctx, claims, *req, ctxValues.Now)
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
				"Product Created",
				"Product successfully created.")

			return true, web.Redirect(ctx, w, r, urlProductsView(resp.ID), http.StatusFound)
		}

		return false, nil
	}

	end, err := f()
	if err != nil {
		return web.RenderError(ctx, w, r, err, h.Renderer, TmplLayoutBase, TmplContentErrorGeneric, web.MIMETextHTMLCharsetUTF8)
	} else if end {
		return nil
	}

	data["categories"], err = h.ShopRepo.FindCategory(ctx, shop.CategoryFindRequest{ Order: []string{"name"} })
	if err != nil {
		return err
	}

	data["brands"], err = h.ShopRepo.FindBrand(ctx, shop.BrandFindRequest{ Order:[]string{"name"} })
	if err != nil {
		return err
	}

	data["form"] = req
	data["urlProductsIndex"] = urlProductsIndex()

	if verr, ok := weberror.NewValidationError(ctx, webcontext.Validator().Struct(shop.ProductCreateRequest{})); ok {
		data["validationDefaults"] = verr.(*weberror.Error)
	}

	return h.Renderer.Render(ctx, w, r, TmplLayoutBase, "products-create.gohtml", web.MIMETextHTMLCharsetUTF8, http.StatusOK, data)
}

// View handles displaying a product.
func (h *Products) View(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {

	productID := params["product_id"]

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
				err = h.ShopRepo.DeleteProduct(ctx, claims, shop.ProductDeleteRequest{
					ID: productID,
				})
				if err != nil {
					return false, err
				}

				webcontext.SessionFlashSuccess(ctx,
					"Product Archived",
					"Product successfully archived.")

				return true, web.Redirect(ctx, w, r, urlProductsIndex(), http.StatusFound)
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

	prj, err := h.ShopRepo.ReadProductByID(ctx, claims, productID)
	if err != nil {
		return err
	}
	data["product"] = prj.Response(ctx)
	data["urlProductsIndex"] = urlProductsIndex()
	data["urlProductsView"] = urlProductsView(productID)
	data["urlProductsUpdate"] = urlProductsUpdate(productID)
	data["urlProductsCreate"] = urlProductsCreate()

	return h.Renderer.Render(ctx, w, r, TmplLayoutBase, "products-view.gohtml", web.MIMETextHTMLCharsetUTF8, http.StatusOK, data)
}

// Update handles updating a product.
func (h *Products) Update(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {

	ctxValues, err := webcontext.ContextValues(ctx)
	if err != nil {
		return err
	}

	productID := params["product_id"]

	claims, err := auth.ClaimsFromContext(ctx)
	if err != nil {
		return err
	}

	//
	req := new(shop.ProductUpdateRequest)
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
			req.ID = productID

			err = h.ShopRepo.UpdateProduct(ctx, claims, *req, ctxValues.Now)
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
				"Product Updated",
				"Product successfully updated.")

			return true, web.Redirect(ctx, w, r, urlProductsView(req.ID), http.StatusFound)
		}

		return false, nil
	}

	end, err := f()
	if err != nil {
		return web.RenderError(ctx, w, r, err, h.Renderer, TmplLayoutBase, TmplContentErrorGeneric, web.MIMETextHTMLCharsetUTF8)
	} else if end {
		return nil
	}

	prj, err := h.ShopRepo.ReadProductByID(ctx, claims, productID)
	if err != nil {
		return err
	}

	data["product"] = prj.Response(ctx)

	data["categories"], err = h.ShopRepo.FindCategory(ctx, shop.CategoryFindRequest{ Order: []string{"name"} })
	if err != nil {
		return err
	}

	data["brands"], err = h.ShopRepo.FindBrand(ctx, shop.BrandFindRequest{ Order:[]string{"name"} })
	if err != nil {
		return err
	}

	data["urlProductsIndex"] = urlProductsIndex()
	data["urlProductsView"] = urlProductsView(productID)

	if req.ID == "" {
		req.Name = &prj.Name
		req.CategoryID = &prj.CategoryID
		req.Price = &prj.Price
		req.Sku = &prj.Sku
		req.BrandID = &prj.BrandID
		req.Barcode = &prj.Barcode
	}

	data["form"] = req

	if verr, ok := weberror.NewValidationError(ctx, webcontext.Validator().Struct(shop.ProductUpdateRequest{})); ok {
		data["validationDefaults"] = verr.(*weberror.Error)
	}

	return h.Renderer.Render(ctx, w, r, TmplLayoutBase, "products-update.gohtml", web.MIMETextHTMLCharsetUTF8, http.StatusOK, data)
}
