package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"merryworld/surebank/internal/expenditure"
	"merryworld/surebank/internal/platform/auth"
	"merryworld/surebank/internal/platform/datatable"
	"merryworld/surebank/internal/platform/web"
	"merryworld/surebank/internal/platform/web/webcontext"
	"merryworld/surebank/internal/platform/web/weberror"

	"github.com/pkg/errors"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"gopkg.in/DataDog/dd-trace-go.v1/contrib/go-redis/redis"
)

// Expenditures represents the Expenditures API method handler set.
type Expenditures struct {
	ExpendituresRepo *expenditure.Repository
	Redis    *redis.Client
	Renderer web.Renderer
}

func urlExpendituresIndex() string {
	return fmt.Sprintf("/accounting/reps-expenditure")
}

// Index handles listing all the expenditures.
func (h *Expenditures) Index(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {

	claims, err := auth.ClaimsFromContext(ctx)
	if err != nil {
		return err
	}

	fields := []datatable.DisplayField{
		{Field: "id", Title: "ID", Visible: false, Searchable: true, Orderable: true, Filterable: false},
		{Field: "sales_rep", Title: "Sales Rep", Visible: true, Searchable: true, Orderable: true, Filterable: true, FilterPlaceholder: "filter Sales Rep"},
		{Field: "amount", Title: "Amount", Visible: true, Searchable: true, Orderable: true, Filterable: true, FilterPlaceholder: "filter Sales Rep"},
		{Field: "date", Title: "Date", Visible: true, Searchable: true, Orderable: true, Filterable: true, FilterPlaceholder: "filter Sales Rep"},
		{Field: "reason", Title: "Reason", Visible: true, Searchable: true, Orderable: true, Filterable: true, FilterPlaceholder: "filter Sales Rep"},
		{Field: "action", Title: "", Visible: true},
	}

	mapFunc := func(q *expenditure.Response, cols []datatable.DisplayField) (resp []datatable.ColumnValue, err error) {
		for i := 0; i < len(cols); i++ {
			col := cols[i]
			var v datatable.ColumnValue
			switch col.Field {
			case "id":
				v.Value = fmt.Sprintf("%s", q.ID)
			case "sales_rep":
				v.Value = q.SalesRep
				v.Formatted = fmt.Sprintf("<a href='%s'>%s</a>", urlUsersView(q.SalesRepID), v.Value)
			case "amount":
				v.Value = fmt.Sprintf("%f", q.Amount)
				p := message.NewPrinter(language.English)
				v.Formatted = p.Sprintf("%.2f", q.Amount)
			case "date":
				v.Value = q.Date.Local
				v.Formatted = v.Value
			case "reason":
				v.Value = q.Reason
				v.Formatted = v.Value
			case "action":
				v.Value = "action"
				v.Formatted = fmt.Sprintf("<button type='button' data-action='click->repsExpenditure#remove' data-id='%s' class='btn btn-danger'>Delete</button>", q.ID)
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

		res, err := h.ExpendituresRepo.Find(ctx, claims, expenditure.FindRequest{
			Order: order,
			IncludeSalesRep: true,
		})
		if err != nil {
			return resp, err
		}

		for _, a := range res.Expenditures {
			l, err := mapFunc(a, fields)
			if err != nil {
				return resp, errors.Wrapf(err, "Failed to map expenditure for display.")
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

	return h.Renderer.Render(ctx, w, r, TmplLayoutBase, "accounting-reps-expenditures.gohtml", web.MIMETextHTMLCharsetUTF8, http.StatusOK, data)
}

// Create handles creating a new expenditure.
func (h *Expenditures) Create(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {

	ctxValues, err := webcontext.ContextValues(ctx)
	if err != nil {
		return err
	}

	claims, err := auth.ClaimsFromContext(ctx)
	if err != nil {
		return err
	}

	//
	var req expenditure.CreateRequest
	if err := web.Decode(ctx, r, &req); err != nil {
		if _, ok := errors.Cause(err).(*weberror.Error); !ok {
			err = weberror.NewError(ctx, err, http.StatusBadRequest)
		}
		return web.RespondJsonError(ctx, w, err)
	}

	res, err := h.ExpendituresRepo.Create(ctx, claims, req, ctxValues.Now)
	if err != nil {
		return web.RespondJsonError(ctx, w, weberror.NewError(ctx, err, http.StatusBadRequest))
	}
	return web.RespondJson(ctx, w, res, http.StatusCreated)
}

// Create handles creating a new expenditure.
func (h *Expenditures) Delete(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {

	claims, err := auth.ClaimsFromContext(ctx)
	if err != nil {
		return err
	}

	id := params["id"] 

	err = h.ExpendituresRepo.Delete(ctx, claims, expenditure.DeleteRequest{ ID: id })
	if err != nil {
		return web.RespondJsonError(ctx, w, weberror.NewError(ctx, err, http.StatusBadRequest))
	}
	return web.RespondJson(ctx, w, true, http.StatusOK)
}

