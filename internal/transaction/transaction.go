package transaction

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/jinzhu/now"
	"github.com/pborman/uuid"
	"github.com/pkg/errors"
	"github.com/volatiletech/null"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"
	. "github.com/volatiletech/sqlboiler/queries/qm"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"

	"merryworld/surebank/internal/platform/auth"
	"merryworld/surebank/internal/platform/web"
	"merryworld/surebank/internal/platform/web/webcontext"
	"merryworld/surebank/internal/platform/web/weberror"
	"merryworld/surebank/internal/postgres/models"
)

var (
	// ErrNotFound abstracts the postgres not found error.
	ErrNotFound = errors.New("Entity not found")

	// ErrForbidden occurs when a user tries to do something that is forbidden to them according to our access control policies.
	ErrForbidden = errors.New("Attempted action is not allowed")
)

// Find gets all the transaction from the database based on the request params.
func (repo *Repository) Find(ctx context.Context, claims auth.Claims, req FindRequest) (*PagedResponseList, error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "internal.transaction.Find")
	defer span.Finish()
	var queries = []QueryMod{
		Load(models.TransactionRels.SalesRep),
		Load(models.TransactionRels.Account),
	}

	if req.Where != "" {
		queries = append(queries, Where(req.Where, req.Args...))
	}

	if !req.IncludeArchived {
		queries = append(queries, And("archived_at is null"))
	}

	// if the current sales resp is not an admin, show only his transactions
	if !claims.HasRole(auth.RoleAdmin) {
		queries = append(queries, And(fmt.Sprintf("%s = '%s'", models.TransactionColumns.SalesRepID, claims.Subject)))
	}

	totalCount, err := models.Transactions(queries...).Count(ctx, repo.DbConn)
	if err != nil {
		return nil, weberror.WithMessage(ctx, err, "Cannot get transaction count")
	}

	if req.IncludeAccount {
		queries = append(queries, Load(models.TransactionRels.Account))
	}

	if req.IncludeSalesRep {
		queries = append(queries, Load(models.TransactionRels.SalesRep))
	}

	if len(req.Order) > 0 {
		for _, s := range req.Order {
			queries = append(queries, OrderBy(s))
		}
	}

	if req.Limit != nil {
		queries = append(queries, Limit(int(*req.Limit)))
	}

	if req.Offset != nil {
		queries = append(queries, Offset(int(*req.Offset)))
	}

	slice, err := models.Transactions(queries...).All(ctx, repo.DbConn)
	if err != nil {
		if err.Error() == sql.ErrNoRows.Error() {
			return &PagedResponseList{}, nil
		}
		return nil, weberror.NewError(ctx, err, 500)
	}

	var result Transactions
	for _, rec := range slice {
		result = append(result, FromModel(rec))
	}

	if len(result) == 0 {
		return &PagedResponseList{}, nil
	}

	return &PagedResponseList{
		Transactions: result.Response(ctx),
		TotalCount:   totalCount,
	}, nil
}

func (repo *Repository) TxReport(ctx context.Context, claims auth.Claims, req FindRequest) ([]TxReportResponse, error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "internal.transaction.TxReport")
	defer span.Finish()
	statement := `select 
		tx.id, 
		tx.sales_rep_id,
		tx.account_id,
		tx.tx_type,
		tx.narration,
		tx.receipt_no,
		c.id as customer_id,
		concat(u.first_name, ' ', u.last_name) as sales_rep, 
		c.name as customer_name, 
		ac.number as account_number, 
		tx.amount, 
		tx.created_at from 
	transaction tx
		inner join account ac on ac.id = tx.account_id
		inner join customer c on c.id = ac.customer_id
		inner join users u on u.id = tx.sales_rep_id `

	var wheres []string
	var args []interface{}

	if req.Where != "" {
		wheres = append(wheres, req.Where)
		args = append(args, req.Args...)
	}

	if !req.IncludeArchived {
		wheres = append(wheres, "tx.archived_at is null")
	}

	if !claims.HasRole(auth.RoleAdmin) {
		wheres = append(wheres, fmt.Sprintf("tx.%s = $%d", models.TransactionColumns.SalesRepID, len(args)+1))
		args = append(args, claims.Subject)
	}

	if len(wheres) > 0 {
		statement = fmt.Sprintf("%s where %s", statement, strings.Join(wheres, " and "))
	}

	order := "order by tx.created_at"
	if len(req.Order) > 0 {
		order = fmt.Sprintf("order by %s", strings.Join(req.Order, ","))
	}
	statement = fmt.Sprintf("%s %s", statement, order)

	if req.Limit != nil {
		statement = fmt.Sprintf("%s limit $%d", statement, len(args)+1)
		args = append(args, req.Limit)
	}

	if req.Offset != nil {
		statement = fmt.Sprintf("%s offset $%d", statement, len(args)+1)
		args = append(args, req.Offset)
	}

	fmt.Println(statement)
	var result []TxReportResponse
	err := models.NewQuery(qm.SQL(statement, args...)).Bind(ctx, repo.DbConn, &result)
	return result, err
}

// ReadByID gets the specified transaction by ID from the database.
func (repo *Repository) ReadByID(ctx context.Context, _ auth.Claims, id string) (*Transaction, error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "internal.transaction.ReadByID")
	defer span.Finish()
	queries := []QueryMod{
		models.TransactionWhere.ID.EQ(id),
		Load(models.TransactionRels.Account),
		Load(models.TransactionRels.SalesRep),
	}
	model, err := models.Transactions(queries...).One(ctx, repo.DbConn)
	if err != nil {
		return nil, err
	}

	return FromModel(model), nil
}

func (repo *Repository) TodayDepositAmount(ctx context.Context, claims auth.Claims) (float64, error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "internal.transaction.TodayDepositAmount")
	defer span.Finish()
	startDate := now.BeginningOfDay()
	return repo.TotalDepositAmount(ctx, claims, startDate.UTC().Unix(), time.Now().UTC().Unix())
}

func (repo *Repository) ThisWeekDepositAmount(ctx context.Context, claims auth.Claims) (float64, error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "internal.transaction.ThisWeekDepositAmount")
	defer span.Finish()
	startDate := now.BeginningOfWeek()
	return repo.TotalDepositAmount(ctx, claims, startDate.UTC().Unix(), time.Now().UTC().Unix())
}

func (repo *Repository) ThisMonthDepositAmount(ctx context.Context, claims auth.Claims) (float64, error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "internal.transaction.ThisMonthDepositAmount")
	defer span.Finish()
	startDate := now.BeginningOfMonth()
	return repo.TotalDepositAmount(ctx, claims, startDate.UTC().Unix(), time.Now().UTC().Unix())
}

func (repo *Repository) TotalDepositAmount(ctx context.Context, claims auth.Claims, startDate, endDate int64) (float64, error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "internal.transaction.TotalDepositAmount")
	defer span.Finish()
	statement := `select sum(amount) total from transaction where tx_type = 'deposit' and created_at > $1 and created_at < $2`
	args := []interface{}{startDate, endDate}
	if !claims.HasRole(auth.RoleAdmin) {
		statement += " and sales_rep_id = $3"
		args = append(args, claims.Subject)
	}

	var result struct {
		Total sql.NullFloat64
	}
	err := models.NewQuery(SQL(statement, args...)).Bind(ctx, repo.DbConn, &result)
	return result.Total.Float64, err
}

func (repo *Repository) DepositAmountByWhere(ctx context.Context, where string, args []interface{}) (float64, error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "internal.transaction.DepositAmountByWhere")
	defer span.Finish()
	statement := `select sum(amount) total from transaction where `
	if len(where) > 0 {
		statement += where + fmt.Sprintf(" AND %s = '%s' ", models.TransactionColumns.TXType, TransactionType_Deposit)
	} else {
		statement += fmt.Sprintf(" %s = '%s' ", models.TransactionColumns.TXType, TransactionType_Deposit)
	}
	var result struct {
		Total sql.NullFloat64
	}
	err := models.NewQuery(SQL(statement, args...)).Bind(ctx, repo.DbConn, &result)
	return result.Total.Float64, err
}

const accountBalanceStatement = `SELECT 
	SUM(amount) AS balance FROM (
		SELECT
			CASE WHEN tx.tx_type = 'deposit' THEN tx.amount ELSE -1 * tx.amount END AS amount 
		FROM transaction tx
		WHERE tx.account_id = $1
	) res`

// AccountBalance gets the balance of the specified account from the database.
func (repo *Repository) AccountBalance(ctx context.Context, accountID string) (float64, error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "internal.transaction.AccountBalance")
	defer span.Finish()
	var result null.Float64
	err := repo.DbConn.QueryRow(accountBalanceStatement, accountID).Scan(&result)
	if err != nil && err.Error() == sql.ErrNoRows.Error() {
		return 0, nil
	}
	return result.Float64, err
}

// AccountBalanceTx gets the balance of the specified account from the database within a DB tx.
func (repo *Repository) AccountBalanceTx(ctx context.Context, accountID string, tx *sql.Tx) (float64, error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "internal.transaction.AccountBalanceTx")
	defer span.Finish()
	var result null.Float64
	err := tx.QueryRow(accountBalanceStatement, accountID).Scan(&result)
	return result.Float64, err
}

func (repo *Repository) Deposit(ctx context.Context, claims auth.Claims, req CreateRequest, currentDate time.Time) (*Transaction, error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "internal.transaction.Deposit")
	defer span.Finish()
	//open a new db
	db, err := repo.creatDB()
	if err != nil {
		return nil, weberror.NewErrorMessage(ctx, err, http.StatusBadRequest, "Cannot create DB connection")
	}
	defer db.Close()
	dbTx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	account, err := models.Accounts(models.AccountWhere.Number.EQ(req.AccountNumber),
		Load(models.AccountRels.Customer)).One(ctx, dbTx)
	if err != nil {
		dbTx.Rollback()
		return nil, weberror.NewErrorMessage(ctx, err, http.StatusBadRequest, "Invalid account number")
	}

	if account.AccountType != models.AccountTypeDS {
		m, err := repo.create(ctx, claims, req, currentDate, dbTx)
		if err != nil {
			dbTx.Rollback()
			return nil, err
		}
		if err := dbTx.Commit(); err != nil {
			return nil, err
		}
		return m, nil
	}

	if math.Mod(req.Amount, account.Target) != 0 {
		dbTx.Rollback()
		return nil, weberror.NewError(ctx, fmt.Errorf("Amount must be a multiple of %f", account.Target), 400)
	}

	if req.Amount/account.Target > 50 {
		dbTx.Rollback()
		return nil, weberror.NewError(ctx, fmt.Errorf("Please pay for max of 50 days at a time, one day is %.2f", account.Target), 400)
	}

	if req.PaymentMethod != "bank_deposit" {
		req.PaymentMethod = "cash"
	}

	var tx *Transaction
	amount, reqAmount := req.Amount, req.Amount
	req.Amount = account.Target
	for amount > 0 {
		tx, err = repo.create(ctx, claims, req, currentDate, dbTx)
		if err != nil {
			dbTx.Rollback()
			return nil, err
		}
		amount -= account.Target
		currentDate = currentDate.Add(4 * time.Second)
	}

	if err = dbTx.Commit(); err != nil {
		return nil, err
	}
	if serr := repo.notifySMS.Send(ctx, account.R.Customer.PhoneNumber, "sms/ds_received",
		map[string]interface{}{
			"Name":          account.R.Customer.Name,
			"EffectiveDate": web.NewTimeResponse(ctx, tx.EffectiveDate).LocalDate,
			"Amount":        reqAmount,
			"Balance":       tx.OpeningBalance + tx.Amount,
		}); err != nil {
		// TODO: log critical error. Send message to monitoring account
		fmt.Println(serr)
	}
	return tx, err
}

// create inserts a new transaction into the database.
func (repo *Repository) create(ctx context.Context, claims auth.Claims, req CreateRequest, currentDate time.Time, dbTx *sql.Tx) (*Transaction, error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "internal.transaction.Create")
	defer span.Finish()
	if claims.Audience == "" {
		return nil, errors.WithStack(ErrForbidden)
	}

	// Validate the request.
	v := webcontext.Validator()
	err := v.Struct(req)
	if err != nil {
		return nil, err
	}

	account, err := models.Accounts(models.AccountWhere.Number.EQ(req.AccountNumber),
		Load(models.AccountRels.Customer)).One(ctx, dbTx)
	if err != nil {
		return nil, weberror.NewErrorMessage(ctx, err, http.StatusBadRequest, "Invalid account number")
	}

	// If now empty set it to the current time.
	if currentDate.IsZero() {
		currentDate = time.Now()
	}

	effectiveDate := now.New(currentDate).BeginningOfDay()
	if account.AccountType == models.AccountTypeDS {
		lastDeposit, err := repo.lastDeposit(ctx, account.ID)
		if err == nil {
			effectiveDate = now.New(time.Unix(lastDeposit.EffectiveDate, 0)).Time.Add(24 * time.Hour)
		}
	}

	effectiveDate = effectiveDate.UTC()

	// Always store the time as UTC.
	currentDate = currentDate.UTC()
	// Postgres truncates times to milliseconds when storing. We and do the same
	// here so the value we return is consistent with what we store.
	currentDate = currentDate.Truncate(time.Millisecond)

	repo.accNumMtx.Lock()
	defer repo.accNumMtx.Unlock()

	accountBalance, err := repo.AccountBalanceTx(ctx, account.ID, dbTx)
	if err != nil {
		if err.Error() != sql.ErrNoRows.Error() {
			dbTx.Rollback()
			return nil, err
		}
	}

	m := models.Transaction{
		ID:             uuid.NewRandom().String(),
		AccountID:      account.ID,
		OpeningBalance: accountBalance,
		Amount:         req.Amount,
		Narration:      req.Narration,
		PaymentMethod:  req.PaymentMethod,
		TXType:         req.Type.String(),
		SalesRepID:     claims.Subject,
		ReceiptNo:      repo.generateReceiptNumber(ctx),
		CreatedAt:      currentDate.Unix(),
		UpdatedAt:      currentDate.Unix(),
		EffectiveDate:  effectiveDate.Unix(),
	}

	isFirstContribution, err := repo.CommissionRepo.StartingNewCircle(ctx, account.ID, effectiveDate)
	if err != nil {
		return nil, err
	}

	if err := m.Insert(ctx, dbTx, boil.Infer()); err != nil {
		return nil, errors.WithMessage(err, "Insert deposit failed")
	}

	var lastDepositDate int64
	if req.Type == TransactionType_Deposit {
		lastDepositDate = m.EffectiveDate
		accountBalance += m.Amount
	} else {
		accountBalance -= m.Amount
	}

	if _, err := models.Accounts(models.AccountWhere.ID.EQ(account.ID)).UpdateAll(ctx, dbTx, models.M{
		models.AccountColumns.Balance:         accountBalance,
		models.AccountColumns.LastPaymentDate: lastDepositDate,
	}); err != nil {
		return nil, err
	}

	if err = SaveDailySummary(ctx, req.Amount, 0, 0, currentDate, dbTx); err != nil {
		return nil, err
	}

	// send SMS notification
	if req.Type == TransactionType_Deposit {
		if account.AccountType == models.AccountTypeSB {
			if err = repo.notifySMS.Send(ctx, account.R.Customer.PhoneNumber, "sms/payment_received",
				map[string]interface{}{
					"Name":    account.R.Customer.Name,
					"Amount":  req.Amount,
					"Balance": accountBalance,
				}); err != nil {
				// TODO: log critical error. Send message to monitoring account
				fmt.Println(err)
			}
		}
	} else {
		if err = repo.notifySMS.Send(ctx, account.R.Customer.PhoneNumber, "sms/payment_withdrawn",
			map[string]interface{}{
				"Name":    account.R.Customer.Name,
				"Amount":  req.Amount,
				"Balance": accountBalance,
			}); err != nil {
			// TODO: log critical error. Send message to monitoring account
			fmt.Println(err)
		}
	}

	if req.Type == TransactionType_Deposit && account.AccountType == models.AccountTypeDS && isFirstContribution {
		wm := models.Transaction{
			ID:             uuid.NewRandom().String(),
			AccountID:      account.ID,
			OpeningBalance: accountBalance,
			Amount:         req.Amount,
			Narration:      "DS fee deduction",
			TXType:         TransactionType_Withdrawal.String(),
			SalesRepID:     claims.Subject,
			ReceiptNo:      repo.generateReceiptNumber(ctx),
			CreatedAt:      currentDate.Add(2 * time.Second).Unix(),
			UpdatedAt:      currentDate.Unix(),
		}

		if err := wm.Insert(ctx, dbTx, boil.Infer()); err != nil {
			return nil, errors.WithMessage(err, "Insert DS fee failed")
		}

		accountBalance -= req.Amount
		if _, err := models.Accounts(models.AccountWhere.ID.EQ(account.ID)).UpdateAll(ctx, dbTx, models.M{
			models.AccountColumns.Balance: accountBalance,
		}); err != nil {
			return nil, err
		}

		commission := models.DSCommission{
			ID:            uuid.NewRandom().String(),
			AccountID:     account.ID,
			CustomerID:    account.CustomerID,
			Amount:        req.Amount,
			Date:          currentDate.Unix(),
			EffectiveDate: effectiveDate.Unix(),
		}
		if err := commission.Insert(ctx, repo.DbConn, boil.Infer()); err != nil {
			return nil, err
		}
	}

	return &Transaction{
		ID:             m.ID,
		AccountID:      m.AccountID,
		Amount:         m.Amount,
		OpeningBalance: m.OpeningBalance,
		Narration:      m.Narration,
		Type:           TransactionType(m.TXType),
		SalesRepID:     m.SalesRepID,
		CreatedAt:      time.Unix(m.CreatedAt, 0),
		UpdatedAt:      time.Unix(m.UpdatedAt, 0),
		EffectiveDate:  time.Unix(m.EffectiveDate, 0),
	}, nil
}

// Withdraw inserts a new withdrawal transaction into the database.
func (repo *Repository) Withdraw(ctx context.Context, claims auth.Claims, req WithdrawRequest, now time.Time) (*Transaction, error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "internal.transaction.Withdraw")
	defer span.Finish()
	createReq := MakeDeductionRequest{
		AccountNumber: req.AccountNumber,
		Amount:        req.Amount,
		Narration:     fmt.Sprintf("%s - %s", req.PaymentMethod, req.Narration),
	}
	tx, err := repo.DbConn.Begin()
	if err != nil {
		return nil, err
	}
	if req.PaymentMethod == "Transfer" {
		if len(req.Narration) > 0 {
			createReq.Narration += " -"
		}
		if len(req.Bank) > 0 && len(req.BankAccountNumber) > 0 {
			createReq.Narration += fmt.Sprintf("%s - %s", req.Bank, req.BankAccountNumber)
		}

	}

	txn, err := repo.MakeDeduction(ctx, claims, createReq, now, tx)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}
	_ = tx.Commit()
	return txn, nil
}

func (repo *Repository) generateReceiptNumber(ctx context.Context) string {
	var receipt string
	for receipt == "" || repo.receiptExists(ctx, receipt) {
		receipt = "TX"
		rand.Seed(time.Now().UTC().UnixNano())
		for i := 0; i < 6; i++ {
			receipt += strconv.Itoa(rand.Intn(10))
		}
	}
	return receipt
}

func (repo *Repository) receiptExists(ctx context.Context, receipt string) bool {
	exists, _ := models.Transactions(models.TransactionWhere.ReceiptNo.EQ(receipt)).Exists(ctx, repo.DbConn)
	return exists
}

func (repo *Repository) lastDeposit(ctx context.Context, accountID string) (*models.Transaction, error) {
	return models.Transactions(
		models.TransactionWhere.AccountID.EQ(accountID),
		models.TransactionWhere.TXType.EQ(TransactionType_Deposit.String()),
		OrderBy(fmt.Sprintf("%s desc", models.TransactionColumns.CreatedAt)),
		Limit(1),
	).One(ctx, repo.DbConn)
}

// Update replaces an exiting transaction in the database.
func (repo *Repository) Update(ctx context.Context, claims auth.Claims, req UpdateRequest, now time.Time) error {
	span, ctx := tracer.StartSpanFromContext(ctx, "internal.transaction.Update")
	defer span.Finish()

	if claims.Audience == "" {
		return errors.WithStack(ErrForbidden)
	}
	// Admin users can update transactions they have access to.
	if !claims.HasRole(auth.RoleAdmin) {
		return errors.WithStack(ErrForbidden)
	}

	// Validate the request.
	v := webcontext.Validator()
	err := v.Struct(req)
	if err != nil {
		return err
	}

	cols := models.M{}
	if req.Narration != nil {
		cols[models.TransactionColumns.Narration] = req.Narration
	}

	if req.Amount != nil {
		cols[models.TransactionColumns.Amount] = req.Amount
	}

	if len(cols) == 0 {
		return nil
	}

	// If now empty set it to the current time.
	if now.IsZero() {
		now = time.Now()
	}

	// Always store the time as UTC.
	now = now.UTC()
	// Postgres truncates times to milliseconds when storing. We and do the same
	// here so the value we return is consistent with what we store.
	now = now.Truncate(time.Millisecond)

	cols[models.CustomerColumns.UpdatedAt] = now

	tx, err := repo.DbConn.Begin()
	if err != nil {
		return err
	}

	_, err = models.Transactions(models.CustomerWhere.ID.EQ(req.ID)).UpdateAll(ctx, tx, cols)
	if err != nil {
		tx.Rollback()
		return err
	}

	if req.Amount != nil {
		tranx, err := models.FindTransaction(ctx, repo.DbConn, req.ID)
		if err != nil {
			_ = tx.Rollback()
			return err
		}

		diff := *req.Amount - tranx.Amount
		if tranx.TXType == TransactionType_Withdrawal.String() {
			diff *= -1
		}

		if diff != 0 {
			if _, err := models.Accounts(models.AccountWhere.ID.EQ(tranx.AccountID)).UpdateAll(ctx, tx, models.M{
				models.AccountColumns.Balance: fmt.Sprintf("%s + (%f)", models.AccountColumns.Balance, diff),
			}); err != nil {

				_ = tx.Rollback()
				return err
			}

			_, err = models.Transactions(
				models.TransactionWhere.CreatedAt.GT(tranx.CreatedAt)).
				UpdateAll(ctx, tx, models.M{models.TransactionColumns.Amount: fmt.Sprintf("amount + (%f)", diff)})

			if err != nil {
				_ = tx.Rollback()
				return err
			}
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

// Archive soft deleted the transaction from the database.
func (repo *Repository) Archive(ctx context.Context, claims auth.Claims, req ArchiveRequest, now time.Time) error {
	span, ctx := tracer.StartSpanFromContext(ctx, "internal.transaction.Archive")
	defer span.Finish()

	if claims.Audience == "" {
		return errors.WithStack(ErrForbidden)
	}
	// Admin users can update customer they have access to.
	if !claims.HasRole(auth.RoleAdmin) {
		return errors.WithStack(ErrForbidden)
	}
	// Validate the request.
	v := webcontext.Validator()
	err := v.Struct(req)
	if err != nil {
		return err
	}

	// If now empty set it to the current time.
	if now.IsZero() {
		now = time.Now()
	}

	// Always store the time as UTC.
	now = now.UTC()
	// Postgres truncates times to milliseconds when storing. We and do the same
	// here so the value we return is consistent with what we store.
	now = now.Truncate(time.Millisecond)

	tx, err := repo.DbConn.Begin()
	if err != nil {
		return err
	}

	tranx, err := models.Transactions(models.TransactionWhere.ID.EQ(req.ID)).One(ctx, tx)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	if tranx.ArchivedAt.Valid {
		_ = tx.Rollback()
		return errors.New("This transaction has been archived")
	}

	_, err = models.Transactions(models.TransactionWhere.ID.EQ(req.ID)).
		UpdateAll(ctx, tx, models.M{models.TransactionColumns.ArchivedAt: now.Unix()})
	if err != nil {
		tx.Rollback()
		return err
	}

	var txAmount = tranx.Amount
	if tranx.TXType == TransactionType_Deposit.String() {
		txAmount *= -1
	}

	// updated all trans after it
	statement := fmt.Sprintf("UPDATE %s SET %s = %s + $1 WHERE %s = $2 AND %s > $3",
		models.TableNames.Transaction,
		models.TransactionColumns.OpeningBalance, models.TransactionColumns.OpeningBalance,
		models.TransactionColumns.AccountID, models.TransactionColumns.CreatedAt)
	_, err = models.NewQuery(qm.SQL(statement, txAmount, tranx.AccountID, tranx.CreatedAt)).Exec(tx)

	if err != nil {
		_ = tx.Rollback()
		return err
	}

	// update the account balance
	statement = fmt.Sprintf("UPDATE %s SET %s = %s + $1 WHERE %s = $2",
		models.TableNames.Account,
		models.AccountColumns.Balance,
		models.AccountColumns.Balance,
		models.AccountColumns.ID,
	)
	_, err = models.NewQuery(qm.SQL(statement, txAmount, tranx.AccountID)).Exec(tx)

	if err != nil {
		_ = tx.Rollback()
		return err
	}

	if err = tx.Commit(); err != nil {
		return errors.Wrap(err, "commintin transaction")
	}

	return nil
}

// MakeDeduction inserts a new transaction of type withdrawal into the database.
func (repo *Repository) MakeDeduction(ctx context.Context, claims auth.Claims, req MakeDeductionRequest,
	now time.Time, tx *sql.Tx) (*Transaction, error) {

	span, ctx := tracer.StartSpanFromContext(ctx, "internal.transaction.MakeDeduction")
	defer span.Finish()
	if claims.Subject == "" {
		return nil, errors.WithStack(ErrForbidden)
	}

	account, err := models.Accounts(
		qm.Load(models.AccountRels.Customer),
		models.AccountWhere.Number.EQ(req.AccountNumber)).One(ctx, tx)
	if err != nil {
		return nil, weberror.NewErrorMessage(ctx, err, 400, "invalid account number")
	}

	// Validate the request.
	v := webcontext.Validator()
	err = v.Struct(req)
	if err != nil {
		return nil, err
	}

	// If now empty set it to the current time.
	if now.IsZero() {
		now = time.Now()
	}

	// Always store the time as UTC.
	now = now.UTC()
	// Postgres truncates times to milliseconds when storing. We and do the same
	// here so the value we return is consistent with what we store.
	now = now.Truncate(time.Millisecond)

	repo.accNumMtx.Lock()
	defer repo.accNumMtx.Unlock()

	accountBalance, err := repo.AccountBalanceTx(ctx, account.ID, tx)
	if err != nil {
		return nil, err
	}

	if accountBalance < req.Amount {
		return nil, weberror.NewError(ctx, errors.New("insufficient fund"), 400)
	}

	m := models.Transaction{
		ID:             uuid.NewRandom().String(),
		AccountID:      account.ID,
		TXType:         TransactionType_Withdrawal.String(),
		OpeningBalance: accountBalance,
		Amount:         req.Amount,
		Narration:      req.Narration,
		SalesRepID:     claims.Subject,
		CreatedAt:      now.Unix(),
		UpdatedAt:      now.Unix(),
	}

	if err := m.Insert(ctx, tx, boil.Infer()); err != nil {
		return nil, errors.WithMessage(err, "Insert deduction failed")
	}

	accountBalance -= req.Amount
	if _, err := models.Accounts(models.AccountWhere.ID.EQ(account.ID)).UpdateAll(ctx, tx, models.M{
		models.AccountColumns.Balance: accountBalance,
	}); err != nil {

		_ = tx.Rollback()
		return nil, err
	}

	if err = repo.notifySMS.Send(ctx, account.R.Customer.PhoneNumber, "sms/payment_withdrawn",
		map[string]interface{}{
			"Name":    account.R.Customer.Name,
			"Amount":  req.Amount,
			"Balance": accountBalance,
		}); err != nil {
		// TODO: log critical error. Send message to monitoring account
		fmt.Println(err)
	}

	return FromModel(&m), nil
}

// SaveDailySummary saves the provided daily summary info to the db
func SaveDailySummary(ctx context.Context, income, expenditure, bankDeposit float64, date time.Time, tx *sql.Tx) error {
	span, ctx := tracer.StartSpanFromContext(ctx, "internal.transaction.SaveDailySummary")
	defer span.Finish()

	today := now.New(date).BeginningOfDay().Unix()
	existingSummary, err := models.FindDailySummary(ctx, tx, today)
	if err == nil {
		existingSummary.BankDeposit += bankDeposit
		existingSummary.Income += income
		existingSummary.Expenditure += expenditure

		_, err = existingSummary.Update(ctx, tx, boil.Infer())
		return err
	}

	model := models.DailySummary{
		Date:        today,
		BankDeposit: bankDeposit,
		Income:      income,
		Expenditure: expenditure,
	}

	return model.Insert(ctx, tx, boil.Infer())
}
