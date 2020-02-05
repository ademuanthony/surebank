package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/pkg/errors"

	"gopkg.in/DataDog/dd-trace-go.v1/contrib/go-redis/redis"
	"merryworld/surebank/internal/account"
	"merryworld/surebank/internal/customer"
	"merryworld/surebank/internal/platform/auth"
	"merryworld/surebank/internal/platform/datatable"
	"merryworld/surebank/internal/platform/web"
	"merryworld/surebank/internal/shop"
	"merryworld/surebank/internal/transaction"
)

// Customers represents the Customers API method handler set.
type Reports struct {
	CustomerRepo    *customer.Repository
	AccountRepo     *account.Repository
	TransactionRepo *transaction.Repository
	ShopRepo        *shop.Repository
	Renderer        web.Renderer
	Redis           *redis.Client
}

// Stocks handles listing all the stock info.
func (h *Reports) Stocks(ctx context.Context, w http.ResponseWriter, r *http.Request, _ map[string]string) error {

	claims, err := auth.ClaimsFromContext(ctx)
	if err != nil {
		return err
	}

	fields := []datatable.DisplayField{
		{Field: "id", Title: "ID", Visible: false, Searchable: true, Orderable: true, Filterable: false},
		{Field: "name", Title: "Name", Visible: true, Searchable: true, Orderable: true, Filterable: true, FilterPlaceholder: "filter Name"},
		{Field: "quantity", Title: "Quantity", Visible: true, Searchable: false, Orderable: true, Filterable: false},
	}

	mapFunc := func(q shop.StockInfo, cols []datatable.DisplayField) (resp []datatable.ColumnValue, err error) {
		for i := 0; i < len(cols); i++ {
			col := cols[i]
			var v datatable.ColumnValue
			switch col.Field {
			case "id":
				v.Value = fmt.Sprintf("%s", q.ProductID)
			case "name":
				v.Value = q.ProductName
				v.Formatted = fmt.Sprintf("<a href='%s'>%s</a>", urlCustomersView(q.ProductID), v.Value)
			case "quantity":
				v.Value = fmt.Sprintf("%d", q.Quantity)
				v.Formatted = v.Value
			default:
				return resp, errors.Errorf("Failed to map value for %s.", col.Field)
			}
			resp = append(resp, v)
		}

		return resp, nil
	}

	loadFunc := func(ctx context.Context, sorting string, fields []datatable.DisplayField) (resp [][]datatable.ColumnValue, err error) {
		res, err := h.ShopRepo.StockReport(ctx, claims, shop.StockReportRequest{
			Where:           "",
			Args:            nil,
			Order:           nil,
			Limit:           nil,
			Offset:          nil,
			IncludeArchived: false,
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
	}

	return h.Renderer.Render(ctx, w, r, TmplLayoutBase, "report-stocks.gohtml", web.MIMETextHTMLCharsetUTF8, http.StatusOK, data)
}

