package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"merryworld/surebank/internal/accounting"
	"merryworld/surebank/internal/platform/auth"
	"merryworld/surebank/internal/platform/datatable"
	"merryworld/surebank/internal/platform/web"
	"merryworld/surebank/internal/platform/web/webcontext"
	"merryworld/surebank/internal/platform/web/weberror"
	"merryworld/surebank/internal/postgres/models"
	"merryworld/surebank/internal/transaction"
	"merryworld/surebank/internal/user"

	"github.com/gofrs/uuid"
	"github.com/jinzhu/now"
	"github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"gopkg.in/DataDog/dd-trace-go.v1/contrib/go-redis/redis"
)

// Accounting represents the Accounting API method handler set.
type Accounting struct {
	Redis    *redis.Client
	Renderer web.Renderer
	DbConn   *sql.DB
	UserRepos *user.Repository
}

// DailySummaries handles listing all the daily summaries.
func (h *Accounting) DailySummaries(ctx context.Context, w http.ResponseWriter,
	r *http.Request, params map[string]string) error {

	fields := []datatable.DisplayField{
		{Field: "date", Title: "Date", Visible: true, Searchable: true, Orderable: true, Filterable: true, FilterPlaceholder: "filter Name"},
		{Field: "income", Title: "Income", Visible: true, Searchable: true, Orderable: true, Filterable: true, FilterPlaceholder: "filter Name"},
		{Field: "bank_deposit", Title: "Bank Deposit", Visible: true, Searchable: true, Orderable: true, Filterable: true, FilterPlaceholder: "filter Name"},
		{Field: "expenditure", Title: "Expenditure", Visible: true, Searchable: true, Orderable: true, Filterable: true, FilterPlaceholder: "filter Name"},
		{Field: "balance", Title: "Balance", Visible: true, Searchable: true, Orderable: true, Filterable: true, FilterPlaceholder: "filter Name"},
	}

	mapFunc := func(q *models.DailySummary, cols []datatable.DisplayField) (resp []datatable.ColumnValue, err error) {
		for i := 0; i < len(cols); i++ {
			col := cols[i]
			var v datatable.ColumnValue
			switch col.Field {
			case "date":
				dt := web.NewTimeResponse(ctx, time.Unix(q.Date, 0))
				v.Value = dt.LocalDate
				v.Formatted = v.Value
			case "income":
				v.Value = fmt.Sprintf("%.2f", q.Income)
				v.Formatted = v.Value
			case "bank_deposit":
				v.Value = fmt.Sprintf("%.2f", q.BankDeposit)
				v.Formatted = fmt.Sprintf("<a href='/accounting/deposits?date=%d'>%s</a>", q.Date, v.Value)
			case "expenditure":
				v.Value = fmt.Sprintf("%.2f", q.Expenditure)
				v.Formatted = fmt.Sprintf("<a href='/accounting/expenditures?date=%d'>%s</a>", q.Date, v.Value)
			case "balance":
				v.Value = fmt.Sprintf("%.2f", q.Income-(q.BankDeposit+q.Expenditure))
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

		var queries []qm.QueryMod
		for _, s := range order {
			queries = append(queries, qm.OrderBy(s))
		}

		res, err := models.DailySummaries(queries...).All(ctx, h.DbConn)
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
		"datatable":         dt.Response(),
		"urlBranchesCreate": urlBranchesCreate(),
	}

	return h.Renderer.Render(ctx, w, r, TmplLayoutBase, "accounting-summary.gohtml", web.MIMETextHTMLCharsetUTF8, http.StatusOK, data)
}

// RepsSummaries handles listing all the expenditures.
func (h *Accounting) RepsSummaries(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {

	claims, err := auth.ClaimsFromContext(ctx)
	if err != nil {
		return err
	}
	
	data := map[string]interface{}{}
	var salesRepID, startDate, endDate string
	// todo sales rep filtering
	if v := r.URL.Query().Get("sales_rep_id"); v != "" {
		data["salesRepID"] = v
		salesRepID = v
	}

	if v := r.URL.Query().Get("start_date"); v != "" {
		startDate = v
	} else {
		startDate = time.Now().Format("01/02/2006")
	}
	data["startDate"] = startDate

	if v := r.URL.Query().Get("end_date"); v != "" {
		endDate = v
	} else {
		endDate = time.Now().Format("01/02/2006")
	}
	data["endDate"] = endDate

	fields := []datatable.DisplayField{
		{Field: "id", Title: "ID", Visible: false, Searchable: true, Orderable: true, Filterable: false},
		{Field: "sales_rep", Title: "Sales Rep", Visible: true, Searchable: true, Orderable: true, Filterable: true, FilterPlaceholder: "filter Sales Rep"},
		{Field: "income", Title: "Income", Visible: true, Searchable: true, Orderable: true, Filterable: true, FilterPlaceholder: "filter Sales Rep"},
		{Field: "expenditure", Title: "Expenditure", Visible: true, Searchable: true, Orderable: true, Filterable: true, FilterPlaceholder: "filter Sales Rep"},
		{Field: "balance", Title: "Balance", Visible: true, Searchable: true, Orderable: true, Filterable: true, FilterPlaceholder: "filter Sales Rep"},
	}

	mapFunc := func(q accounting.RepsSummary, cols []datatable.DisplayField) (resp []datatable.ColumnValue, err error) {
		incomeUrlTemplate := `/reports/collections?sales_rep_id=%s&start_date=%s&end_date=%s`
		expenditureUrlTemplate := `/accounting/reps-expenditures?sales_rep_id=%s&start_date=%s&end_date=%s`
		for i := 0; i < len(cols); i++ {
			col := cols[i]
			var v datatable.ColumnValue
			switch col.Field {
			case "id":
				v.Value = fmt.Sprintf("%s", q.SalesRepID)
			case "sales_rep":
				v.Value = q.SalesRep
				v.Formatted = fmt.Sprintf("<a href='%s'>%s</a>", urlUsersView(q.SalesRepID), v.Value)
			case "income":
				var amount float64
				if q.Income.Valid {
					amount = q.Income.Float64
				}
				v.Value = fmt.Sprintf("%f", amount)
				p := message.NewPrinter(language.English)
				url := fmt.Sprintf(incomeUrlTemplate, q.SalesRepID, startDate, endDate)
				v.Formatted = fmt.Sprintf("<a href='%s'>%s</a>", url, p.Sprintf("%.2f", amount))
			case "expenditure":
				var amount float64
				if q.Expenditure.Valid {
					amount = q.Expenditure.Float64
				}
				v.Value = fmt.Sprintf("%f", amount)
				p := message.NewPrinter(language.English)
				url := fmt.Sprintf(expenditureUrlTemplate, q.SalesRepID, startDate, endDate)
				v.Formatted = fmt.Sprintf("<a href='%s'>%s</a>", url, p.Sprintf("%.2f", amount))
			case "balance":
				balance := q.Income.Float64 - q.Expenditure.Float64
				v.Value = fmt.Sprintf("%f", balance)
				p := message.NewPrinter(language.English)
				v.Formatted = p.Sprintf("%.2f", balance)
			default:
				return resp, errors.Errorf("Failed to map value for %s.", col.Field)
			}
			resp = append(resp, v)
		}

		return resp, nil
	}

	loadFunc := func(ctx context.Context, sorting string, fields []datatable.DisplayField) (resp [][]datatable.ColumnValue, err error) {

		var userWhere string = "WHERE archived_at IS NULL"
		if salesRepID != "" {
			userWhere += fmt.Sprintf(" AND u.id = '%s'", salesRepID)
		}
		var incomeWhere, expenditureWhere string
		if startDate != "" {
			date, err := time.Parse("01/02/2006", startDate)
			if err != nil {
				date = time.Now()
			}
			date = date.Truncate(time.Millisecond)
			date = now.New(date).BeginningOfDay().Add(-1 * time.Hour)
			incomeWhere = fmt.Sprintf(" AND created_at >= %d ", date.UTC().Unix())
			expenditureWhere = fmt.Sprintf(" AND date >= %d ", date.UTC().Unix())
		}
		if endDate != "" {
			date, err := time.Parse("01/02/2006", endDate)
			if err != nil {
				date = time.Now()
			}
			date = date.Truncate(time.Millisecond)
			date = now.New(date).EndOfDay().Add(-1 * time.Hour)
			incomeWhere += fmt.Sprintf(" AND created_at <= %d ", date.UTC().Unix())
			expenditureWhere += fmt.Sprintf(" AND date <= %d ", date.UTC().Unix())
		}

		statement := fmt.Sprintf(`select 
		u.id as sales_rep_id, 
		concat(u.first_name, ' ', u.last_name) as sales_rep,
		(SELECT SUM(amount) FROM transaction 
			 WHERE tx_type = 'deposit' AND (payment_method = 'cash' or payment_method = 'null') AND sales_rep_id = u.id %s ) AS income,
		(SELECT SUM(amount) FROM reps_expense 
			 WHERE sales_rep_id = u.id %s) AS expenditure
		from users u %s
		order by u.first_name, u.last_name;`, incomeWhere, expenditureWhere, userWhere)

		var summaries []accounting.RepsSummary
		err = models.NewQuery(qm.SQL(statement)).Bind(ctx, h.DbConn, &summaries)

		for _, a := range summaries {
			l, err := mapFunc(a, fields)
			if err != nil {
				return resp, errors.Wrapf(err, "Failed to map reps summary for display.")
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

	users, err := h.UserRepos.Find(ctx, claims, user.UserFindRequest{
		Order: []string{"first_name", "last_name"},
	})
	if err != nil {
		return err
	}
	
	data["datatable"] = dt.Response()
	data["users"] = users

	return h.Renderer.Render(ctx, w, r, TmplLayoutBase, "accounting-reps-summaries.gohtml", web.MIMETextHTMLCharsetUTF8, http.StatusOK, data)
}

// BankAccounts handles listing all the Bank Accounts.
func (h *Accounting) BankAccounts(ctx context.Context, w http.ResponseWriter,
	r *http.Request, params map[string]string) error {

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
				v.Value = q.Bank
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

		var queries []qm.QueryMod
		for _, s := range order {
			queries = append(queries, qm.OrderBy(s))
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
		"datatable":         dt.Response(),
		"urlBranchesCreate": urlBranchesCreate(),
	}

	return h.Renderer.Render(ctx, w, r, TmplLayoutBase, "accounting-banks.gohtml", web.MIMETextHTMLCharsetUTF8, http.StatusOK, data)
}

// CreateBankAccount handles the json request from creating a new bank account
func (h *Accounting) CreateBankAccount(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {

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
		ID:            id.String(),
		AccountName:   req.Name,
		AccountNumber: req.Number,
		Bank:          req.Bank,
	}
	err := bankAccount.Insert(ctx, h.DbConn, boil.Infer())
	if err != nil {
		return web.RespondJsonError(ctx, w, weberror.NewError(ctx, err, http.StatusBadRequest))
	}

	return web.RespondJson(ctx, w, bankAccount, http.StatusCreated)
}

// BankDeposits handles listing all the Bank Deposits.
func (h *Accounting) BankDeposits(ctx context.Context, w http.ResponseWriter,
	r *http.Request, params map[string]string) error {

	fields := []datatable.DisplayField{
		{Field: "id", Title: "ID", Visible: false, Searchable: true, Orderable: true, Filterable: false},
		{Field: "bank", Title: "Bank", Visible: true, Searchable: true, Orderable: true, Filterable: true, FilterPlaceholder: "filter Name"},
		{Field: "amount", Title: "Amount", Visible: true, Searchable: true, Orderable: true, Filterable: true, FilterPlaceholder: "filter Name"},
		{Field: "info", Title: "Customer Information", Visible: true, Searchable: true, Orderable: true, Filterable: true, FilterPlaceholder: "filter Name"},
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
			case "info":
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

		var queries = []qm.QueryMod{
			qm.Load(models.BankDepositRels.BankAccount),
		}
		for _, s := range order {
			queries = append(queries, qm.OrderBy(s))
		}

		r.ParseForm()
		if date, err := strconv.ParseInt(r.FormValue("date"), 10, 64); err != nil {
			queries = append(queries, models.BankDepositWhere.Date.GTE(date))
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

	banks, err := models.BankAccounts(qm.OrderBy("account_name")).All(ctx, h.DbConn)
	if err != nil {
		return web.RespondError(ctx, w, err)
	}

	data := map[string]interface{}{
		"datatable":         dt.Response(),
		"urlBranchesCreate": urlBranchesCreate(),
		"banks":             banks,
	}

	return h.Renderer.Render(ctx, w, r, TmplLayoutBase, "accounting-deposits.gohtml", web.MIMETextHTMLCharsetUTF8, http.StatusOK, data)
}

// CreateBankDeposit handles the json request from creating a new bank deposit
func (h *Accounting) CreateBankDeposit(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {

	ctxValues, err := webcontext.ContextValues(ctx)
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

	today := now.New(ctxValues.Now).BeginningOfDay().UTC()
	id, _ := uuid.NewV4()
	model := models.BankDeposit{
		ID:            id.String(),
		Amount:        req.Amount,
		BankAccountID: req.BankID,
		Date:          today.Unix(),
	}

	tx, err := h.DbConn.Begin()
	if err != nil {
		return err
	}
	err = model.Insert(ctx, tx, boil.Infer())
	if err != nil {
		return web.RespondJsonError(ctx, w, weberror.NewError(ctx, err, http.StatusBadRequest))
	}

	if err = transaction.SaveDailySummary(ctx, 0, 0, req.Amount, ctxValues.Now, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return web.RespondJson(ctx, w, model, http.StatusCreated)
}

// Expenditures handles listing all the Expenditures.
func (h *Accounting) Expenditures(ctx context.Context, w http.ResponseWriter,
	r *http.Request, params map[string]string) error {

	fields := []datatable.DisplayField{
		{Field: "id", Title: "ID", Visible: false, Searchable: true, Orderable: true, Filterable: false},
		{Field: "amount", Title: "Amount", Visible: true, Searchable: true, Orderable: true, Filterable: true, FilterPlaceholder: "filter Name"},
		{Field: "memo", Title: "Memo", Visible: true, Searchable: true, Orderable: true, Filterable: true, FilterPlaceholder: "filter Name"},
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
			case "memo":
				v.Value = q.Reason
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

		var queries []qm.QueryMod
		for _, s := range order {
			queries = append(queries, qm.OrderBy(s))
		}

		r.ParseForm()
		if date, err := strconv.ParseInt(r.FormValue("date"), 10, 64); err != nil {
			queries = append(queries, models.ExpenditureWhere.Date.GTE(date))
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
		"datatable":         dt.Response(),
		"urlBranchesCreate": urlBranchesCreate(),
	}

	return h.Renderer.Render(ctx, w, r, TmplLayoutBase, "accounting-expenditures.gohtml", web.MIMETextHTMLCharsetUTF8, http.StatusOK, data)
}

// CreateExpenditure handles the json request from creating a new bank expenditure
func (h *Accounting) CreateExpenditure(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {

	ctxValues, err := webcontext.ContextValues(ctx)
	if err != nil {
		return err
	}

	// claims, err := auth.ClaimsFromContext(ctx)
	// if err != nil {
	// 	return err
	// }

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
		ID:     id.String(),
		Amount: req.Amount,
		Reason: req.Memo,
		Date:   ctxValues.Now.Unix(),
	}
	err = model.Insert(ctx, tx, boil.Infer())
	if err != nil {
		return web.RespondJsonError(ctx, w, weberror.NewError(ctx, err, http.StatusBadRequest))
	}

	if err = transaction.SaveDailySummary(ctx, 0, req.Amount, 0, ctxValues.Now, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return web.RespondJson(ctx, w, model, http.StatusCreated)
}
