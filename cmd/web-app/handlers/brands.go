package handlers

import (
	"context"
	"fmt"
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

// Brands represents the Brands API method handler set.
type Brands struct {
	ShopRepo *shop.Repository
	Redis    *redis.Client
	Renderer web.Renderer
}

func urlBrandsIndex() string {
	return fmt.Sprintf("/shop/brands")
}

func urlBrandsCreate() string {
	return fmt.Sprintf("/shop/brands/create")
}

func urlBrandsView(brandID string) string {
	return fmt.Sprintf("/shop/brands/%s", brandID)
}

func urlBrandsUpdate(brandID string) string {
	return fmt.Sprintf("/shop/brands/%s/update", brandID)
}

// Index handles listing all the brands.
func (h *Brands) Index(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {

	fields := []datatable.DisplayField{
		{Field: "id", Title: "ID", Visible: false, Searchable: true, Orderable: true, Filterable: false},
		{Field: "name", Title: "Brand", Visible: true, Searchable: true, Orderable: true, Filterable: true, FilterPlaceholder: "filter Name"},
	}

	mapFunc := func(q *shop.Brand, cols []datatable.DisplayField) (resp []datatable.ColumnValue, err error) {
		for i := 0; i < len(cols); i++ {
			col := cols[i]
			var v datatable.ColumnValue
			switch col.Field {
			case "id":
				v.Value = fmt.Sprintf("%s", q.ID)
			case "name":
				v.Value = q.Name
				v.Formatted = fmt.Sprintf("<a href='%s'>%s</a>", urlBrandsView(q.ID), v.Value)
			default:
				return resp, errors.Errorf("Failed to map value for %s.", col.Field)
			}
			resp = append(resp, v)
		}

		return resp, nil
	}

	loadFunc := func(ctx context.Context, sorting string, fields []datatable.DisplayField) (resp [][]datatable.ColumnValue, err error) {
		res, err := h.ShopRepo.FindBrand(ctx, shop.BrandFindRequest{
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
		"urlBrandsCreate": urlBrandsCreate(),
		"urlBrandsIndex": urlBrandsIndex(),
	}

	return h.Renderer.Render(ctx, w, r, TmplLayoutBase, "brands-index.gohtml", web.MIMETextHTMLCharsetUTF8, http.StatusOK, data)
}

// Create handles creating a new brand.
func (h *Brands) Create(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {

	ctxValues, err := webcontext.ContextValues(ctx)
	if err != nil {
		return err
	}

	claims, err := auth.ClaimsFromContext(ctx)
	if err != nil {
		return err
	}

	//
	req := new(shop.BrandCreateRequest)
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

			usr, err := h.ShopRepo.CreateBrand(ctx, claims, *req, ctxValues.Now)
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
				"Brand Created",
				"Brand successfully created.")

			return true, web.Redirect(ctx, w, r, urlBrandsView(usr.ID), http.StatusFound)
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
	data["urlBrandsIndex"] = urlBrandsIndex()

	if verr, ok := weberror.NewValidationError(ctx, webcontext.Validator().Struct(shop.BrandCreateRequest{})); ok {
		data["validationDefaults"] = verr.(*weberror.Error)
	}

	return h.Renderer.Render(ctx, w, r, TmplLayoutBase, "brands-create.gohtml", web.MIMETextHTMLCharsetUTF8, http.StatusOK, data)
}

// View handles displaying a brand.
func (h *Brands) View(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {

	brandID := params["brand_id"]

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
				err = h.ShopRepo.DeleteBrand(ctx, claims, shop.BrandDeleteRequest{
					ID: brandID,
				})
				if err != nil {
					return false, err
				}

				webcontext.SessionFlashSuccess(ctx,
					"Brand Deleted",
					"Brand successfully deleted.")

				return true, web.Redirect(ctx, w, r, urlBrandsIndex(), http.StatusFound)
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

	prj, err := h.ShopRepo.ReadBrandByID(ctx, claims, brandID)
	if err != nil {
		return err
	}
	data["brand"] = prj.Response(ctx)
	data["urlBrandsIndex"] = urlBrandsIndex()
	data["urlBrandsUpdate"] = urlBrandsUpdate(brandID)
	data["urlBrandsView"] = urlBrandsView(brandID)

	return h.Renderer.Render(ctx, w, r, TmplLayoutBase, "brands-view.gohtml", web.MIMETextHTMLCharsetUTF8, http.StatusOK, data)
}

// Update handles updating a brand.
func (h *Brands) Update(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {

	brandID := params["brand_id"]

	claims, err := auth.ClaimsFromContext(ctx)
	if err != nil {
		return err
	}

	//
	req := new(shop.BrandUpdateRequest)
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
			req.ID = brandID

			err = h.ShopRepo.UpdateBrand(ctx, claims, *req)
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
				"Brand Updated",
				"Brand successfully updated.")

			return true, web.Redirect(ctx, w, r, urlBrandsView(req.ID), http.StatusFound)
		}

		return false, nil
	}

	end, err := f()
	if err != nil {
		return web.RenderError(ctx, w, r, err, h.Renderer, TmplLayoutBase, TmplContentErrorGeneric, web.MIMETextHTMLCharsetUTF8)
	} else if end {
		return nil
	}

	prj, err := h.ShopRepo.ReadBrandByID(ctx, claims, brandID)
	if err != nil {
		return err
	}

	data["brand"] = prj.Response(ctx)

	data["urlBrandsIndex"] = urlBrandsIndex()
	data["urlBrandsView"] = urlBrandsView(brandID)

	if req.ID == "" {
		req.Name = &prj.Name
	}

	data["form"] = req

	if verr, ok := weberror.NewValidationError(ctx, webcontext.Validator().Struct(shop.BrandUpdateRequest{})); ok {
		data["validationDefaults"] = verr.(*weberror.Error)
	}

	return h.Renderer.Render(ctx, w, r, TmplLayoutBase, "brands-update.gohtml", web.MIMETextHTMLCharsetUTF8, http.StatusOK, data)
}
