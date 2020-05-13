package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"merryworld/surebank/internal/accounting"
	"merryworld/surebank/internal/platform/auth"
	"merryworld/surebank/internal/platform/datatable"
	"merryworld/surebank/internal/platform/web"
	"merryworld/surebank/internal/postgres/models"
	"merryworld/surebank/internal/transaction"
	"merryworld/surebank/internal/platform/web/webcontext"
	"merryworld/surebank/internal/platform/web/weberror"

	"github.com/pkg/errors"
	"github.com/gofrs/uuid"
	"github.com/jinzhu/now"
	"github.com/jmoiron/sqlx"
	. "github.com/volatiletech/sqlboiler/v4/queries/qm"
	"gopkg.in/DataDog/dd-trace-go.v1/contrib/go-redis/redis"
)

// Branches represents the Branches API method handler set.
type Accounting struct {
	Redis    *redis.Client
	Renderer web.Renderer
	DbConn *sqlx.DB
}

// DailySummaries handles listing all the daily summaries.
func (h *Accounting) DailySummaries(ctx context.Context, w http.ResponseWriter,
	r *http.Request, params map[string]string) error {

   claims, err := auth.ClaimsFromContext(ctx)
   if err != nil {
	   return err
   }

   fields := []datatable.DisplayField{
	   {Field: "date", Title: "Date", Visible: true, Searchable: true, Orderable: true, Filterable: true, FilterPlaceholder: "filter Name"},
	   {Field: "income", Title: "Income", Visible: true, Searchable: true, Orderable: true, Filterable: true, FilterPlaceholder: "filter Name"},
	   {Field: "bank_deposit", Title: "Income", Visible: true, Searchable: true, Orderable: true, Filterable: true, FilterPlaceholder: "filter Name"},
	   {Field: "expenditure", Title: "Income", Visible: true, Searchable: true, Orderable: true, Filterable: true, FilterPlaceholder: "filter Name"},
	   {Field: "balance", Title: "Balance", Visible: true, Searchable: true, Orderable: true, Filterable: true, FilterPlaceholder: "filter Name"},

	}

   mapFunc := func(q *models.DailySummary, cols []datatable.DisplayField) (resp []datatable.ColumnValue, err error) {
	for i := 0; i < len(cols); i++ {
		   col := cols[i]
		   var v datatable.ColumnValue
		   switch col.Field {
		   case "date":
				dt := web.NewTimeResponse(ctx, time.Unix(q.Date, 0))
			   v.Value = dt.Local
			   v.Formatted = v.Value
			case "income":
				v.Value = fmt.Sprintf("%.2f", q.Income)
				v.Formatted = v.Value
			case "bank_deposit":
				v.Value = fmt.Sprintf("%.2f", q.BankDeposit)
				v.Formatted = fmt.Sprintf("<a href='/accounting/deposits?date=%d'>%s</a>", q.Date, v.Value)
			case "expenditure":
				v.Value = fmt.Sprintf("%.2f", q.Expenditure)
				v.Formatted = fmt.Sprintf("<a href='/accounting/expenditure?date=%d'>%s</a>", q.Date, v.Value)
			case "balance":
				v.Value = fmt.Sprintf("%.2f", q.Income - (q.BankDeposit + q.Expenditure))
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

	   var queries []QueryMod
	   for _, s := range order {
		   queries = append(queries, OrderBy(s))
	   }

	   res, err := models.Expenditures(queries...).All(ctx, h.DbConn)
	   if err != nil {
		   return resp, err
	   }

	   for _, a := range res {
		   l, err := mapFunc(a, fields)
		   if err != nil {
			   return resp, errors.Wrapf(err, "Failed to map branch for display.")
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
	   "datatable":           dt.Response(),
	   "urlBranchesCreate": urlBranchesCreate(),
   }

   return h.Renderer.Render(ctx, w, r, TmplLayoutBase, "accounting-deposits.gohtml", web.MIMETextHTMLCharsetUTF8, http.StatusOK, data)
}

// BankAccounts handles listing all the Bank Accounts.
func (h *Accounting) BankAccounts(ctx context.Context, w http.ResponseWriter,
	 r *http.Request, params map[string]string) error {

	claims, err := auth.ClaimsFromContext(ctx)
	if err != nil {
		return err
	}

	fields := []datatable.DisplayField{
		{Field: "id", Title: "ID", Visible: false, Searchable: true, Orderable: true, Filterable: false},
		{Field: "bank", Title: "Bank", Visible: true, Searchable: true, Orderable: true, Filterable: true, FilterPlaceholder: "filter Name"},
		{Field: "account_name", Title: "Account Name", Visible: true, Searchable: true, Orderable: true, Filterable: true, FilterPlaceholder: "filter Name"},
		{Field: "account_number", Title: "Account Number", Visible: true, Searchable: true, Orderable: true, Filterable: true, FilterPlaceholder: "filter Name"},
	}

	mapFunc := func(q *models.BankAccount, cols []datatable.DisplayField) (resp []datatable.ColumnValue, err error) {
		for i := 0; i < len(cols); i++ {
			col := cols[i]
			var v datatable.ColumnValue
			switch col.Field {
			case "id":
				v.Value = fmt.Sprintf("%s", q.ID)
			case "bank":
				v.Value = q.Name
				v.Formatted = v.Value
			case "account_name":
				v.Value = q.AccountName
				v.Formatted = v.Value
			case "account_number":
				v.Value = q.AccountNumber
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

		var queries []QueryMod
		for _, s := range order {
			queries = append(queries, OrderBy(s))
		}

		res, err := models.BankAccounts(queries...).All(ctx, h.DbConn)
		if err != nil {
			return resp, err
		}

		for _, a := range res {
			l, err := mapFunc(a, fields)
			if err != nil {
				return resp, errors.Wrapf(err, "Failed to map branch for display.")
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
		"datatable":           dt.Response(),
		"urlBranchesCreate": urlBranchesCreate(),
	}

	return h.Renderer.Render(ctx, w, r, TmplLayoutBase, "accounting-banks.gohtml", web.MIMETextHTMLCharsetUTF8, http.StatusOK, data)
}

// CreateBankAccount handles the json request from creating a new bank account
func (h *Accounting) CreateBankAccount(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {

	ctxValues, err := webcontext.ContextValues(ctx)
	if err != nil {
		return err
	}

	claims, err := auth.ClaimsFromContext(ctx)
	if err != nil {
		return err
	}

	//
	var req accounting.CreateBankAccount
	if err := web.Decode(ctx, r, &req); err != nil {
		if _, ok := errors.Cause(err).(*weberror.Error); !ok {
			err = weberror.NewError(ctx, err, http.StatusBadRequest)
		}
		return web.RespondJsonError(ctx, w, err)
	}

	
	id, _ := uuid.NewV4()
	bankAccount := models.BankAccount{
		ID: id.String(),
		AccountName: req.Name,
		AccountNumber: req.Number,
		Bank: req.Bank,
	}
	res, err := bankAccount.Insert(ctx, h.Repo.DbConn)
	if err != nil {
		cause := errors.Cause(err)
		switch cause {
		case branch.ErrForbidden:
			return web.RespondJsonError(ctx, w, weberror.NewError(ctx, err, http.StatusForbidden))
		default:
			return web.RespondJsonError(ctx, w, weberror.NewError(ctx, err, http.StatusBadRequest))
		}
	}

	result := res.Response(ctx)
	return web.RespondJson(ctx, w, result, http.StatusCreated)
}

// BankDeposits handles listing all the Bank Deposits.
func (h *Accounting) BankDeposits(ctx context.Context, w http.ResponseWriter,
	r *http.Request, params map[string]string) error {

   claims, err := auth.ClaimsFromContext(ctx)
   if err != nil {
	   return err
   }

   fields := []datatable.DisplayField{
	   {Field: "id", Title: "ID", Visible: false, Searchable: true, Orderable: true, Filterable: false},
	   {Field: "bank", Title: "Bank", Visible: true, Searchable: true, Orderable: true, Filterable: true, FilterPlaceholder: "filter Name"},
	   {Field: "amount", Title: "Amount", Visible: true, Searchable: true, Orderable: true, Filterable: true, FilterPlaceholder: "filter Name"},
	   {Field: "date", Title: "Date", Visible: true, Searchable: true, Orderable: true, Filterable: true, FilterPlaceholder: "filter Name"},
   }

   mapFunc := func(q *models.BankDeposit, cols []datatable.DisplayField) (resp []datatable.ColumnValue, err error) {
	bankInfo := fmt.Sprintf("%s (%s) - %s", q.R.BankAccount.AccountName,
	 q.R.BankAccount.AccountNumber, q.R.BankAccount.Bank)   
	for i := 0; i < len(cols); i++ {
		   col := cols[i]
		   var v datatable.ColumnValue
		   switch col.Field {
		   case "id":
			   v.Value = fmt.Sprintf("%s", q.ID)
		   case "bank":
			   v.Value = bankInfo
			   v.Formatted = v.Value
		   case "amount":
			   v.Value = fmt.Sprintf("%.2f", q.Amount)
			   v.Formatted = v.Value
		   case "date":
				dt := web.NewTimeResponse(ctx, time.Unix(q.Date, 0))
			   v.Value = dt.Local
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

	   var queries []QueryMod
	   for _, s := range order {
		   queries = append(queries, OrderBy(s))
	   }

	   r.ParseForm()
	   if date, err := strconv.ParseInt(r.FormValue("date"), 10, 64); err != nil {
		   queries = models.BankDepositWhere.Date.GTE(date)
	   }

	   res, err := models.BankDeposits(queries...).All(ctx, h.DbConn)
	   if err != nil {
		   return resp, err
	   }

	   for _, a := range res {
		   l, err := mapFunc(a, fields)
		   if err != nil {
			   return resp, errors.Wrapf(err, "Failed to map branch for display.")
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

   banks, err := models.BankAccounts(OrderBy("account_name")).All(ctx, h.DbConn)
   if err != nil {
	   return web.RespondError(ctx, w, err)
   }

   data := map[string]interface{}{
	   "datatable":           dt.Response(),
	   "urlBranchesCreate": urlBranchesCreate(),
	   "banks": banks,
   }

   return h.Renderer.Render(ctx, w, r, TmplLayoutBase, "accounting-deposits.gohtml", web.MIMETextHTMLCharsetUTF8, http.StatusOK, data)
}

// CreateBankDeposit handles the json request from creating a new bank deposit
func (h *Accounting) CreateBankDeposit(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {

   ctxValues, err := webcontext.ContextValues(ctx)
   if err != nil {
	   return err
   }

   claims, err := auth.ClaimsFromContext(ctx)
   if err != nil {
	   return err
   }

   //
   var req accounting.CreateBankDeposit
   if err := web.Decode(ctx, r, &req); err != nil {
	   if _, ok := errors.Cause(err).(*weberror.Error); !ok {
		   err = weberror.NewError(ctx, err, http.StatusBadRequest)
	   }
	   return web.RespondJsonError(ctx, w, err)
   }
   
   today := now.BeginningOfDay().UTC()
   id, _ := uuid.NewV4()
   model := models.BankDeposit{
	   ID: id.String(),
	   Amount	: req.Amount,
	   BankAccountID: req.BankID,
	   Date: ctxValues.Now,
   }

   tx, err := h.DbConn.Begin()
   if err != nil {
	   return err
   }
   res, err := model.Insert(ctx, tx)
   if err != nil {
	   cause := errors.Cause(err)
	   switch cause {
	   case branch.ErrForbidden:
		   return web.RespondJsonError(ctx, w, weberror.NewError(ctx, err, http.StatusForbidden))
	   default:
		   return web.RespondJsonError(ctx, w, weberror.NewError(ctx, err, http.StatusBadRequest))
	   }
   }

   
   if err = transaction.SaveDailySummary(ctx, 0, 0, req.Amount, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

   result := res.Response(ctx)
   return web.RespondJson(ctx, w, result, http.StatusCreated)
}

// Expenditures handles listing all the Expenditures.
func (h *Accounting) Expenditures(ctx context.Context, w http.ResponseWriter,
	r *http.Request, params map[string]string) error {

   claims, err := auth.ClaimsFromContext(ctx)
   if err != nil {
	   return err
   }

   fields := []datatable.DisplayField{
	   {Field: "id", Title: "ID", Visible: false, Searchable: true, Orderable: true, Filterable: false},
	   {Field: "amount", Title: "Amount", Visible: true, Searchable: true, Orderable: true, Filterable: true, FilterPlaceholder: "filter Name"},
	   {Field: "date", Title: "Date", Visible: true, Searchable: true, Orderable: true, Filterable: true, FilterPlaceholder: "filter Name"},
   }

   mapFunc := func(q *models.Expenditure, cols []datatable.DisplayField) (resp []datatable.ColumnValue, err error) {
	for i := 0; i < len(cols); i++ {
		   col := cols[i]
		   var v datatable.ColumnValue
		   switch col.Field {
		   case "id":
			   v.Value = fmt.Sprintf("%s", q.ID)
		   case "amount":
			   v.Value = fmt.Sprintf("%.2f", q.Amount)
			   v.Formatted = v.Value
		   case "date":
				dt := web.NewTimeResponse(ctx, time.Unix(q.Date, 0))
			   v.Value = dt.Local
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

	   var queries []QueryMod
	   for _, s := range order {
		   queries = append(queries, OrderBy(s))
	   }
	   
	   r.ParseForm()
	   if date, err := strconv.ParseInt(r.FormValue("date"), 10, 64); err != nil {
		   queries = models.ExpenditureWhere.Date.GTE(date)
	   }

	   res, err := models.Expenditures(queries...).All(ctx, h.DbConn)
	   if err != nil {
		   return resp, err
	   }

	   for _, a := range res {
		   l, err := mapFunc(a, fields)
		   if err != nil {
			   return resp, errors.Wrapf(err, "Failed to map branch for display.")
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
	   "datatable":           dt.Response(),
	   "urlBranchesCreate": urlBranchesCreate(),
   }

   return h.Renderer.Render(ctx, w, r, TmplLayoutBase, "accounting-deposits.gohtml", web.MIMETextHTMLCharsetUTF8, http.StatusOK, data)
}

// CreateExpenditure handles the json request from creating a new bank expenditure
func (h *Accounting) CreateExpenditure(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {

   ctxValues, err := webcontext.ContextValues(ctx)
   if err != nil {
	   return err
   }

   claims, err := auth.ClaimsFromContext(ctx)
   if err != nil {
	   return err
   }

   //
   var req accounting.CreateExpenditure
   if err := web.Decode(ctx, r, &req); err != nil {
	   if _, ok := errors.Cause(err).(*weberror.Error); !ok {
		   err = weberror.NewError(ctx, err, http.StatusBadRequest)
	   }
	   return web.RespondJsonError(ctx, w, err)
   }

   tx, err := h.DbConn.Begin()
   if err != nil {
	   return web.RespondJsonError(ctx, w, err)
   }
   id, _ := uuid.NewV4()
   model := models.Expenditure{
	   ID: id.String(),
	   Amount	: req.Amount,
	   Date: ctxValues.Now,
   }
   res, err := model.Insert(ctx, tx)
   if err != nil {
	   cause := errors.Cause(err)
	   switch cause {
	   case branch.ErrForbidden:
		   return web.RespondJsonError(ctx, w, weberror.NewError(ctx, err, http.StatusForbidden))
	   default:
		   return web.RespondJsonError(ctx, w, weberror.NewError(ctx, err, http.StatusBadRequest))
	   }
   }
   
   if err = transaction.SaveDailySummary(ctx, 0, req.Amount, 0, tx); err != nil {
	   tx.Rollback()
	   return err
   }

   if err = tx.Commit(); err != nil {
	   return nil, err
	}

   result := res.Response(ctx)
   return web.RespondJson(ctx, w, result, http.StatusCreated)
}
