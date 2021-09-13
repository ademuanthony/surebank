package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"merryworld/surebank/internal/platform/auth"
	"merryworld/surebank/internal/platform/datatable"
	"merryworld/surebank/internal/platform/web"
	"merryworld/surebank/internal/platform/web/webcontext"
	"merryworld/surebank/internal/platform/web/weberror"
	"merryworld/surebank/internal/profit"

	"github.com/gorilla/schema"
	"github.com/pkg/errors"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"gopkg.in/DataDog/dd-trace-go.v1/contrib/go-redis/redis"
)

// Profits represents the Profit API method handler set.
type Profits struct {
	ProfRepo *profit.Repository
	Redis    *redis.Client
	Renderer web.Renderer
}

func urlProfitsIndex() string {
	return "/profits"
}

func urlProfitsCreate() string {
	return "/profits/create"
}

func urlProfitsView(profitID string) string {
	return fmt.Sprintf("/profits/%s", profitID)
}

func urlProfitsUpdate(categoryID string) string {
	return fmt.Sprintf("/profits/%s/update", categoryID)
}

// Index handles listing all profits.
func (h *Profits) Index(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {

	dt, ok, err := profitDatatable(ctx, h.ProfRepo, h.Redis, w, r)
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
		"datatable":        dt.Response(),
		"urlProfitsCreate": urlProfitsCreate(),
		"urlProfitsIndex":  urlProfitsIndex(),
	}

	return h.Renderer.Render(ctx, w, r, TmplLayoutBase, "profits-index.gohtml", web.MIMETextHTMLCharsetUTF8, http.StatusOK, data)
}

func profitDatatable(ctx context.Context, repo *profit.Repository, redisClient *redis.Client, w http.ResponseWriter, r *http.Request,
	) (*datatable.Datatable, bool, error) {

	fields := []datatable.DisplayField{
		{Field: "id", Title: "ID", Visible: false, Searchable: true, Orderable: true, Filterable: false},
		{Field: "narration", Title: "Profit", Visible: true, Searchable: true, Orderable: true, Filterable: true, FilterPlaceholder: "filter Name"},
		{Field: "amount", Title: "Amount", Visible: true, Searchable: false, Orderable: true, Filterable: true, FilterPlaceholder: "filter Price"},
	}

	mapFunc := func(q *profit.Profit, cols []datatable.DisplayField) (resp []datatable.ColumnValue, err error) {

		for i := 0; i < len(cols); i++ {
			col := cols[i]
			var v datatable.ColumnValue
			switch col.Field {
			case "id":
				v.Value = q.ID
			case "narration":
				v.Value = q.Narration
				v.Formatted = fmt.Sprintf("<a href='%s'>%s</a>", urlProfitsView(q.ID), v.Value)
			case "amount":
				v.Value = fmt.Sprintf("%f", q.Amount)
				p := message.NewPrinter(language.English)
				v.Formatted = p.Sprintf("%.2f", q.Amount)
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

		res, err := repo.FindProfit(ctx, profit.ProfitFindRequest{
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

// Create handles creating a new profit.
func (h *Profits) Create(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {

	ctxValues, err := webcontext.ContextValues(ctx)
	if err != nil {
		return err
	}

	claims, err := auth.ClaimsFromContext(ctx)
	if err != nil {
		return err
	}

	//
	req := new(profit.ProfitCreateRequest)
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

			resp, err := h.ProfRepo.CreateProfit(ctx, claims, *req, ctxValues.Now)
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

			// Display a success message to the profit.
			webcontext.SessionFlashSuccess(ctx,
				"Profit Created",
				"Profit successfully created.")

			return true, web.Redirect(ctx, w, r, urlProfitsView(resp.ID), http.StatusFound)
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
	data["urlProfitsIndex"] = urlProfitsIndex()

	if verr, ok := weberror.NewValidationError(ctx, webcontext.Validator().Struct(profit.ProfitCreateRequest{})); ok {
		data["validationDefaults"] = verr.(*weberror.Error)
	}

	return h.Renderer.Render(ctx, w, r, TmplLayoutBase, "profits-create.gohtml", web.MIMETextHTMLCharsetUTF8, http.StatusOK, data)
}

// View handles displaying a profit.
func (h *Profits) View(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {

	profitID := params["profit_id"]

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
				err = h.ProfRepo.DeleteProfit(ctx, claims, profit.ProfitDeleteRequest{
					ID: profitID,
				})
				if err != nil {
					return false, err
				}

				webcontext.SessionFlashSuccess(ctx,
					"Profit Archived",
					"Profit successfully archived.")

				return true, web.Redirect(ctx, w, r, urlProfitsIndex(), http.StatusFound)
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

	prj, err := h.ProfRepo.ReadProfitByID(ctx, claims, profitID)
	if err != nil {
		return err
	}
	data["profit"] = prj.Response(ctx)
	data["urlProfitsIndex"] = urlProfitsIndex()
	data["urlProfitsView"] = urlProfitsView(profitID)
	data["urlProfitsUpdate"] = urlProfitsUpdate(profitID)
	data["urlProfitsCreate"] = urlProfitsCreate()

	return h.Renderer.Render(ctx, w, r, TmplLayoutBase, "profits-view.gohtml", web.MIMETextHTMLCharsetUTF8, http.StatusOK, data)
}

// Update handles updating a profit.
func (h *Profits) Update(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {

	ctxValues, err := webcontext.ContextValues(ctx)
	if err != nil {
		return err
	}

	profitID := params["profit_id"]

	claims, err := auth.ClaimsFromContext(ctx)
	if err != nil {
		return err
	}

	//
	req := new(profit.ProfitUpdateRequest)
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
			req.ID = profitID

			err = h.ProfRepo.UpdateProfit(ctx, claims, *req, ctxValues.Now)
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
				"Profit Updated",
				"Profit successfully updated.")

			return true, web.Redirect(ctx, w, r, urlProfitsView(req.ID), http.StatusFound)
		}

		return false, nil
	}

	end, err := f()
	if err != nil {
		return web.RenderError(ctx, w, r, err, h.Renderer, TmplLayoutBase, TmplContentErrorGeneric, web.MIMETextHTMLCharsetUTF8)
	} else if end {
		return nil
	}

	prj, err := h.ProfRepo.ReadProfitByID(ctx, claims, profitID)
	if err != nil {
		return err
	}

	data["profit"] = prj.Response(ctx)

	data["urlProfitsIndex"] = urlProfitsIndex()
	data["urlProfitsView"] = urlProfitsView(profitID)

	if req.ID == "" {
		req.Narration = &prj.Narration
		req.Amount = &prj.Amount
	}

	data["form"] = req

	if verr, ok := weberror.NewValidationError(ctx, webcontext.Validator().Struct(profit.ProfitUpdateRequest{})); ok {
		data["validationDefaults"] = verr.(*weberror.Error)
	}

	return h.Renderer.Render(ctx, w, r, TmplLayoutBase, "profits-update.gohtml", web.MIMETextHTMLCharsetUTF8, http.StatusOK, data)
}
