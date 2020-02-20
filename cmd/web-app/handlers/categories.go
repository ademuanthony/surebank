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

// Categories represents the Category API method handler set.
type Categories struct {
	ShopRepo *shop.Repository
	Redis    *redis.Client
	Renderer web.Renderer
}

func urlCategoriesIndex() string {
	return fmt.Sprintf("/shop/categories")
}

func urlCategoriesCreate() string {
	return fmt.Sprintf("/shop/categories/create")
}

func urlCategoriesView(categoryID string) string {
	return fmt.Sprintf("/shop/categories/%s", categoryID)
}

func urlCategoriesUpdate(categoryID string) string {
	return fmt.Sprintf("/shop/categories/%s/update", categoryID)
}

// Index handles listing all the categories.
func (h *Categories) Index(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {

	fields := []datatable.DisplayField{
		{Field: "id", Title: "ID", Visible: false, Searchable: true, Orderable: true, Filterable: false},
		{Field: "name", Title: "Category", Visible: true, Searchable: true, Orderable: true, Filterable: true, FilterPlaceholder: "search Name"},
	}

	mapFunc := func(q *shop.Category, cols []datatable.DisplayField) (resp []datatable.ColumnValue, err error) {
		for i := 0; i < len(cols); i++ {
			col := cols[i]
			var v datatable.ColumnValue
			switch col.Field {
			case "id":
				v.Value = fmt.Sprintf("%s", q.ID)
			case "name":
				v.Value = q.Name
				v.Formatted = fmt.Sprintf("<a href='%s'>%s</a>", urlCategoriesView(q.ID), v.Value)
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

		res, err := h.ShopRepo.FindCategory(ctx, shop.CategoryFindRequest{
			Order: order,
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
		"urlCategoriesCreate": urlCategoriesCreate(),
		"urlCategoriesIndex": urlCategoriesIndex(),
	}

	return h.Renderer.Render(ctx, w, r, TmplLayoutBase, "categories-index.gohtml", web.MIMETextHTMLCharsetUTF8, http.StatusOK, data)
}

// Create handles creating a new category.
func (h *Categories) Create(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {

	claims, err := auth.ClaimsFromContext(ctx)
	if err != nil {
		return err
	}

	//
	req := new(shop.CategoryCreateRequest)
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

			usr, err := h.ShopRepo.CreateCategory(ctx, claims, *req)
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
				"Category Created",
				"Category successfully created.")

			return true, web.Redirect(ctx, w, r, urlCategoriesView(usr.ID), http.StatusFound)
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
	data["urlCategoriesIndex"] = urlCategoriesIndex()

	if verr, ok := weberror.NewValidationError(ctx, webcontext.Validator().Struct(shop.CategoryCreateRequest{})); ok {
		data["validationDefaults"] = verr.(*weberror.Error)
	}

	return h.Renderer.Render(ctx, w, r, TmplLayoutBase, "categories-create.gohtml", web.MIMETextHTMLCharsetUTF8, http.StatusOK, data)
}

// View handles displaying a brand.
func (h *Categories) View(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {

	categoryID := params["category_id"]

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
				err = h.ShopRepo.DeleteCategory(ctx, claims, shop.CategoryDeleteRequest{
					ID: categoryID,
				})
				if err != nil {
					return false, err
				}

				webcontext.SessionFlashSuccess(ctx,
					"Category Deleted",
					"Category successfully deleted.")

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

	prj, err := h.ShopRepo.ReadCategoryByID(ctx, claims, categoryID)
	if err != nil {
		return err
	}

	where := "category_id = $1"
	args := []interface{}{categoryID}
	dt, ok, err := productDatatable(ctx, h.ShopRepo, h.Redis, w, r, where, args)
	if ok {
		if err != nil {
			return err
		}
		return nil
	}

	if err != nil {
		return err
	}

	data["category"] = prj.Response(ctx)
	data["datatable"] = dt.Response()
	data["urlCategoriesCreate"] = urlCategoriesCreate()
	data["urlCategoriesIndex"] = urlCategoriesIndex()
	data["urlCategoriesView"] = urlCategoriesView(categoryID)
	data["urlCategoriesUpdate"] = urlCategoriesUpdate(categoryID)

	return h.Renderer.Render(ctx, w, r, TmplLayoutBase, "categories-view.gohtml", web.MIMETextHTMLCharsetUTF8, http.StatusOK, data)
}

// Update handles updating a brand.
func (h *Categories) Update(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {

	categoryID := params["category_id"]

	claims, err := auth.ClaimsFromContext(ctx)
	if err != nil {
		return err
	}

	//
	req := new(shop.CategoryUpdateRequest)
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
			req.ID = categoryID

			err = h.ShopRepo.UpdateCategory(ctx, claims, *req)
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
				"Category Updated",
				"Category successfully updated.")

			return true, web.Redirect(ctx, w, r, urlCategoriesView(req.ID), http.StatusFound)
		}

		return false, nil
	}

	end, err := f()
	if err != nil {
		return web.RenderError(ctx, w, r, err, h.Renderer, TmplLayoutBase, TmplContentErrorGeneric, web.MIMETextHTMLCharsetUTF8)
	} else if end {
		return nil
	}

	prj, err := h.ShopRepo.ReadBrandByID(ctx, claims, categoryID)
	if err != nil {
		return err
	}

	data["brand"] = prj.Response(ctx)

	data["urlCategoriesIndex"] = urlCategoriesIndex()
	data["urlCategoriesView"] = urlCategoriesView(categoryID)

	if req.ID == "" {
		req.Name = &prj.Name
	}

	data["form"] = req

	if verr, ok := weberror.NewValidationError(ctx, webcontext.Validator().Struct(shop.CategoryUpdateRequest{})); ok {
		data["validationDefaults"] = verr.(*weberror.Error)
	}

	return h.Renderer.Render(ctx, w, r, TmplLayoutBase, "categories-update.gohtml", web.MIMETextHTMLCharsetUTF8, http.StatusOK, data)
}
