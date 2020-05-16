package transaction

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/jinzhu/now"
	"github.com/pborman/uuid"
	"github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	. "github.com/volatiletech/sqlboiler/v4/queries/qm"
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
	var queries = []QueryMod {
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
		queries = append(queries, And(fmt.Sprintf("%s = '%s", models.TransactionColumns.SalesRepID, claims.Subject)))
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

// ReadByID gets the specified transaction by ID from the database.
func (repo *Repository) ReadByID(ctx context.Context, _ auth.Claims, id string) (*Transaction, error) {
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
	startDate := now.BeginningOfDay()
	return repo.TotalDepositAmount(ctx, claims, startDate.UTC().Unix(), time.Now().UTC().Unix())
}

func (repo *Repository) ThisWeekDepositAmount(ctx context.Context, claims auth.Claims) (float64, error) {
	startDate := now.BeginningOfWeek()
	return repo.TotalDepositAmount(ctx, claims, startDate.UTC().Unix(), time.Now().UTC().Unix())
}

func (repo *Repository) ThisMonthDepositAmount(ctx context.Context, claims auth.Claims) (float64, error) {
	startDate := now.BeginningOfMonth()
	return repo.TotalDepositAmount(ctx, claims, startDate.UTC().Unix(), time.Now().UTC().Unix())
}

func (repo *Repository) TotalDepositAmount(ctx context.Context, claims auth.Claims, startDate, endDate int64) (float64, error) {
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
	statement := `select sum(amount) total from transaction`
	if len(where) > 0 {
		statement += " where " + where
	}
	var result struct {
		Total sql.NullFloat64
	}
	err := models.NewQuery(SQL(statement, args...)).Bind(ctx, repo.DbConn, &result)
	return result.Total.Float64, err
}

// AccountBalance gets the balance of the specified account from the database.
func (repo *Repository) AccountBalance(ctx context.Context, _ auth.Claims, accountNumber string) (balance float64, err error) {
	account, err := models.Accounts(models.AccountWhere.Number.EQ(accountNumber)).One(ctx, repo.DbConn)
	if err != nil {
		return 0, weberror.WithMessage(ctx, err, "Invalid account number")
	}

	lastTx, err := repo.lastTransaction(ctx, account.ID)
	if err != nil {
		if err.Error() != sql.ErrNoRows.Error() {
			return 0, err
		}
	}

	return accountBalanceAtTx(lastTx), nil
}

// Create inserts a new transaction into the database.
func (repo *Repository) Create(ctx context.Context, claims auth.Claims, req CreateRequest, currentDate time.Time) (*Transaction, error) {
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
		Load(models.AccountRels.Customer)).One(ctx, repo.DbConn)
	if err != nil {
		return nil, weberror.NewErrorMessage(ctx, err, http.StatusBadRequest, "Invalid account number")
	}

	// If now empty set it to the current time.
	if currentDate.IsZero() {
		currentDate = time.Now()
	}

	// Always store the time as UTC.
	currentDate = currentDate.UTC()
	// Postgres truncates times to milliseconds when storing. We and do the same
	// here so the value we return is consistent with what we store.
	currentDate = currentDate.Truncate(time.Millisecond)

	repo.accNumMtx.Lock()
	defer repo.accNumMtx.Unlock()

	tx, err := repo.DbConn.Begin()
	if err != nil {
		return nil, err
	}

	lastTransaction, err := repo.lastTransaction(ctx, account.ID)
	if err != nil {
		if err.Error() != sql.ErrNoRows.Error() {
			_ = tx.Rollback()
			return nil, err
		}
	}

	m := models.Transaction{
		ID:             uuid.NewRandom().String(),
		AccountID:      account.ID,
		OpeningBalance: accountBalanceAtTx(lastTransaction),
		Amount:         req.Amount,
		Narration:      req.Narration,
		TXType:         req.Type.String(),
		SalesRepID:     claims.Subject,
		ReceiptNo: 		repo.generateReceiptNumber(ctx),
		CreatedAt:      currentDate.Unix(),
		UpdatedAt:      currentDate.Unix(),
	}

	effectiveDate := now.New(currentDate).BeginningOfDay()
	var isFirstContribution bool = true
	if account.AccountType == models.AccountTypeAJ {
		lastDeposit, err := repo.lastDeposit(ctx, account.ID)
		if err == nil {
			if effectiveDate.Year() == time.Unix(lastDeposit.EffectiveData, 0).Year() &&
			 effectiveDate.Month() == time.Unix(lastDeposit.EffectiveData, 0).Month() {
				effectiveDate = now.New(time.Unix(lastDeposit.EffectiveData, 0)).Time.Add(24 * time.Hour)
				isFirstContribution = false
			}
		}
	}
	m.EffectiveData = effectiveDate.UTC().Unix()

	if err := m.Insert(ctx, tx, boil.Infer()); err != nil {
		_ = tx.Rollback()
		return nil, errors.WithMessage(err, "Insert deposit failed")
	}

	if _, err := models.Accounts(models.AccountWhere.ID.EQ(account.ID)).UpdateAll(ctx, tx, models.M{
		models.AccountColumns.Balance: accountBalanceAtTx(&m),
	}); err != nil {

		_ = tx.Rollback()
		return nil, err
	}

	if err = SaveDailySummary(ctx, req.Amount, 0, 0, currentDate, tx); err != nil {
		tx.Rollback()
		return nil, err
	}

	if account.AccountType == models.AccountTypeAJ { 
		if err = repo.notifySMS.Send(ctx, account.R.Customer.PhoneNumber, "sms/ajor_received",
				map[string]interface{}{
					"Name":    account.R.Customer.Name,
					"EffectiveDate": web.NewTimeResponse(ctx, time.Unix(m.EffectiveData, 0)).LocalDate,
					"Amount":  req.Amount,
					"Balance": m.OpeningBalance + req.Amount,
				}); err != nil {
				// TODO: log critical error. Send message to monitoring account
				fmt.Println(err)
			}
	}else {
		// send SMS notification
		if req.Type == TransactionType_Deposit {
			if err = repo.notifySMS.Send(ctx, account.R.Customer.PhoneNumber, "sms/payment_received",
				map[string]interface{}{
					"Name":    account.R.Customer.Name,
					"Amount":  req.Amount,
					"Balance": m.OpeningBalance + req.Amount,
				}); err != nil {
				// TODO: log critical error. Send message to monitoring account
				fmt.Println(err)
			}
		} else {
			if err = repo.notifySMS.Send(ctx, account.R.Customer.PhoneNumber, "sms/payment_withdrawn",
				map[string]interface{}{
					"Name":    account.R.Customer.Name,
					"Amount":  req.Amount,
					"Balance": m.OpeningBalance + req.Amount,
				}); err != nil {
				// TODO: log critical error. Send message to monitoring account
				fmt.Println(err)
			}
		}
	}

	if account.AccountType == models.AccountTypeAJ && isFirstContribution {
		wm := models.Transaction{
			ID:             uuid.NewRandom().String(),
			AccountID:      account.ID,
			OpeningBalance: accountBalanceAtTx(&m),
			Amount:         req.Amount,
			Narration:      "Ajor fee deduction",
			TXType:         TransactionType_Withdrawal.String(),
			SalesRepID:     claims.Subject,
			ReceiptNo: 		repo.generateReceiptNumber(ctx),
			CreatedAt:      currentDate.Unix(),
			UpdatedAt:      currentDate.Unix(),
		}

		if err := wm.Insert(ctx, tx, boil.Infer()); err != nil {
			_ = tx.Rollback()
			return nil, errors.WithMessage(err, "Insert Ajor fee failed")
		}

		if _, err := models.Accounts(models.AccountWhere.ID.EQ(account.ID)).UpdateAll(ctx, tx, models.M{
			models.AccountColumns.Balance: accountBalanceAtTx(&wm),
		}); err != nil {
	
			_ = tx.Rollback()
			return nil, err
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, err
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
	}, nil
}

// Withdraw inserts a new withdrawal transaction into the database.
func (repo *Repository) Withdraw(ctx context.Context, claims auth.Claims, req WithdrawRequest, now time.Time) (*Transaction, error) {
	createReq := MakeDeductionRequest {
		AccountNumber: req.AccountNumber,
		Amount: req.Amount,
		Narration: fmt.Sprintf("%s - %s", req.PaymentMethod, req.Narration),
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

// lastTransaction returns the last transaction for the specified account
func (repo *Repository) lastTransaction(ctx context.Context, accountID string) (*models.Transaction, error) {
	return models.Transactions(
		models.TransactionWhere.AccountID.EQ(accountID),
		OrderBy(fmt.Sprintf("%s desc", models.TransactionColumns.CreatedAt)),
		Limit(1),
	).One(ctx, repo.DbConn)
}

func (repo *Repository) lastDeposit(ctx context.Context, accountID string) (*models.Transaction, error) {
	return models.Transactions(
		models.TransactionWhere.AccountID.EQ(accountID),
		models.TransactionWhere.TXType.EQ(TransactionType_Deposit.String()),
		OrderBy(fmt.Sprintf("%s desc", models.TransactionColumns.CreatedAt)),
		Limit(1),
	).One(ctx, repo.DbConn)
}

// accountBalanceAtTx returns the account balance as at the specified tx
func accountBalanceAtTx(tx *models.Transaction) float64 {
	var lastBalance float64
	if tx != nil {
		if tx.TXType == TransactionType_Deposit.String() {
			lastBalance = tx.OpeningBalance + tx.Amount
		} else {
			lastBalance = tx.OpeningBalance - tx.Amount
		}
	}

	return lastBalance
}

// Update replaces an exiting transaction in the database.
func (repo *Repository) Update(ctx context.Context, claims auth.Claims, req UpdateRequest, now time.Time) error {
	span, ctx := tracer.StartSpanFromContext(ctx, "internal.deposit.Update")
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

	_,err = models.Transactions(models.CustomerWhere.ID.EQ(req.ID)).UpdateAll(ctx, tx, cols)

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
	span, ctx := tracer.StartSpanFromContext(ctx, "internal.deposit.Archive")
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
		return err
	}

	_,err = models.Transactions(models.TransactionWhere.ID.EQ(req.ID)).UpdateAll(ctx, tx, models.M{models.TransactionColumns.ArchivedAt: now})

	var txAmount = tranx.Amount
	if tranx.TXType == TransactionType_Withdrawal.String() {
		txAmount *= -1
	}

	// updated all trans after it
	_, err = models.Transactions(
		models.TransactionWhere.CreatedAt.GT(tranx.CreatedAt)).
		UpdateAll(ctx, tx, models.M{models.TransactionColumns.Amount: fmt.Sprintf("amount + (%f)", txAmount)})

	if err != nil {
		_ = tx.Rollback()
		return err
	}

	// update all accounts balance
	_, err = models.Accounts(models.AccountWhere.ID.EQ(tranx.AccountID)).UpdateAll(ctx, tx,
		models.M{models.AccountColumns.Balance: fmt.Sprintf("balance + (%f)", txAmount)})

	if err != nil {
		_ = tx.Rollback()
		return err
	}

	return nil
}

// MakeDeduction inserts a new transaction of type withdrawal into the database.
func (repo *Repository) MakeDeduction(ctx context.Context, claims auth.Claims, req MakeDeductionRequest,
	now time.Time, tx *sql.Tx) (*Transaction, error) {

	span, ctx := tracer.StartSpanFromContext(ctx, "internal.transaction.makeDeduction")
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

	lastTransaction, err := models.Transactions(
		models.TransactionWhere.AccountID.EQ(account.ID),
		OrderBy(fmt.Sprintf("%s desc", models.TransactionColumns.CreatedAt)),
		Limit(1),
	).One(ctx, tx)
	if err != nil {
		if err.Error() != sql.ErrNoRows.Error() {
			return nil, err
		}
	}

	balance := accountBalanceAtTx(lastTransaction)
	if balance < req.Amount {
		return nil, fmt.Errorf("balance: %.2f, cannot deduct %.2f. Insufficient fund. Aborted", balance, req.Amount)
	}

	m := models.Transaction{
		ID:             uuid.NewRandom().String(),
		AccountID:      account.ID,
		TXType:         TransactionType_Withdrawal.String(),
		OpeningBalance: balance,
		Amount:         req.Amount,
		Narration:      req.Narration,
		SalesRepID:     claims.Subject,
		CreatedAt:      now.Unix(),
		UpdatedAt:      now.Unix(),
	}

	if err := m.Insert(ctx, tx, boil.Infer()); err != nil {
		return nil, errors.WithMessage(err, "Insert deposit failed")
	}

	if _, err := models.Accounts(models.AccountWhere.ID.EQ(account.ID)).UpdateAll(ctx, tx, models.M{
		models.AccountColumns.Balance: accountBalanceAtTx(&m),
	}); err != nil {

		_ = tx.Rollback()
		return nil, err
	}

	if err = repo.notifySMS.Send(ctx, account.R.Customer.PhoneNumber, "sms/payment_withdrawn",
			map[string]interface{}{
				"Name":    account.R.Customer.Name,
				"Amount":  req.Amount,
				"Balance": m.OpeningBalance + req.Amount,
			}); err != nil {
			// TODO: log critical error. Send message to monitoring account
			fmt.Println(err)
		}

	return FromModel(&m), nil
}

// SaveDailySummary saves the provided daily summary info to the db
func SaveDailySummary(ctx context.Context, income, expenditure, bankDeposit float64, date time.Time, tx *sql.Tx) error {
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
		Date: today,
		BankDeposit: bankDeposit,
		Income: income,
		Expenditure: expenditure,
	}

	return model.Insert(ctx, tx, boil.Infer())
}