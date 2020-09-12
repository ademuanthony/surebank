package transaction

import (
	"context"
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
	"github.com/volatiletech/sqlboiler/queries/qm"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"

	"merryworld/surebank/internal/dal"
	"merryworld/surebank/internal/dscommission"
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

func (repo *Repository) BuildQuery(claims auth.Claims, req FindRequest) primitive.M {
	var queries primitive.M
	if !req.IncludeArchived {
		queries[dal.TransactionColumns.ArchivedAt] = bson.M{"$ne": nil}
	}

	if !claims.HasRole(auth.RoleAdmin) {
		queries[dal.TransactionColumns.SalesRepID] = claims.Subject
	}
	if req.CustomerID != "" {
		queries[dal.TransactionColumns.CustomerID] = req.CustomerID
	}
	if req.AccountID != "" {
		queries[dal.TransactionColumns.CustomerID] = req.CustomerID
	}
	if req.AccountNumber != "" {
		queries[dal.TransactionColumns.AccountNumber] = req.AccountNumber
	}
	if req.SalesRepID != "" {
		queries[dal.TransactionColumns.SalesRepID] = req.SalesRepID
	}

	if req.PaymentMethod != "" {
		queries[dal.TransactionColumns.PaymentMethod] = req.PaymentMethod
	}

	if req.StartDate > 0 {
		queries[dal.TransactionColumns.CreatedAt] = bson.M{"$gte": req.StartDate}
	}

	if req.EndDate > 0 {
		queries[dal.TransactionColumns.CreatedAt] = bson.M{"$lte": req.EndDate}
	}
	return queries
}

// Find gets all the transaction from the database based on the request params.
func (repo *Repository) Find(ctx context.Context, claims auth.Claims, req FindRequest) (*PagedResponseList, error) {

	collection := repo.mongoDb.Collection(dal.C.Transaction)

	queries := repo.BuildQuery(claims, req)

	totalCount, err := collection.CountDocuments(ctx, queries)
	if err != nil {
		return nil, weberror.WithMessage(ctx, err, "Cannot get customer count")
	}

	findOptions := options.Find()

	sort := bson.D{}
	if len(req.Order) > 0 {
		for _, s := range req.Order {
			sortInfo := strings.Split(s, " ")
			if len(sortInfo) != 2 {
				continue
			}
			sort = append(sort, primitive.E{Key: sortInfo[0], Value: sortInfo[1]})
		}
	}
	findOptions.SetSort(sort)

	if req.Limit != nil {
		findOptions.SetLimit(int64(*req.Limit))
	}

	if req.Offset != nil {
		findOptions.Skip = req.Offset
	}

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, weberror.WithMessage(ctx, err, "Cannot get customer list")
	}
	defer cursor.Close(ctx)
	var result Transactions
	for cursor.Next(ctx) {
		var c Transaction
		cursor.Decode(&c)
		result = append(result, c)
	}
	if err := cursor.Err(); err != nil {
		return nil, weberror.WithMessage(ctx, err, "Cannot get customer list")
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
func (repo *Repository) ReadByID(ctx context.Context, id string) (*Transaction, error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "internal.transaction.ReadByID")
	defer span.Finish()

	var rec Transaction
	collection := repo.mongoDb.Collection(dal.C.Transaction)
	err := collection.FindOne(ctx, bson.M{dal.TransactionColumns.ID: id}).Decode(&rec)
	return &rec, err
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

	queries := bson.M{
		dal.TransactionColumns.Type:      string(TransactionType_Deposit),
		dal.TransactionColumns.CreatedAt: bson.M{"$gt": startDate},
		dal.TransactionColumns.CreatedAt: bson.M{"$lt": endDate},
	}

	if !claims.HasRole(auth.RoleAdmin) {
		queries[dal.TransactionColumns.SalesRepID] = claims.Subject
	}

	return repo.DepositAmountByWhere(ctx, queries)
}

func (repo *Repository) DepositAmountByWhere(ctx context.Context, queries primitive.M) (float64, error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "internal.transaction.DepositAmountByWhere")
	defer span.Finish()

	var result struct {
		Total float64
	}

	pipeline := []bson.M{
		{
			"$match": queries,
		},
		{
			"$group": bson.M{
				"_id":   "",
				"total": bson.M{"$sum": "$" + dal.TransactionColumns.Amount},
			},
		},
		{
			"$project": bson.M{
				"_id":   0,
				"total": 1,
			},
		},
	}

	cursor, err := repo.mongoDb.Collection(dal.C.Transaction).Aggregate(ctx, pipeline)
	if err != nil {
		return 0, fmt.Errorf("TotalDepositAmount -> Aggregate, %s", err.Error())
	}

	err = cursor.All(ctx, &result)
	return result.Total, err
}

// AccountBalance gets the balance of the specified account from the database.
func (repo *Repository) AccountBalance(ctx context.Context, accountID string) (float64, error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "internal.transaction.AccountBalance")
	defer span.Finish()

	account, err := repo.AccountRepo.ReadByID(ctx, accountID)
	if err != nil {
		return 0, err
	}
	return account.Balance, nil
}

func (repo *Repository) Deposit(ctx context.Context, claims auth.Claims, req CreateRequest, currentDate time.Time) (*Transaction, error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "internal.transaction.Deposit")
	defer span.Finish()

	account, err := repo.AccountRepo.ReadByNumber(ctx, req.AccountNumber)
	if err != nil {
		return nil, weberror.NewErrorMessage(ctx, err, http.StatusBadRequest, "Invalid account number")
	}

	if account.Type != models.AccountTypeDS {
		m, err := repo.create(ctx, claims, req, currentDate)
		return m, err
	}

	if math.Mod(req.Amount, account.Target) != 0 {
		return nil, weberror.NewError(ctx, fmt.Errorf("Amount must be a multiple of %f", account.Target), 400)
	}

	if req.Amount/account.Target > 50 {
		return nil, weberror.NewError(ctx, fmt.Errorf("Please pay for max of 50 days at a time, one day is %.2f", account.Target), 400)
	}

	if req.PaymentMethod != "bank_deposit" {
		req.PaymentMethod = "cash"
	}

	var tx *Transaction
	amount, reqAmount := req.Amount, req.Amount
	req.Amount = account.Target
	for amount > 0 {
		tx, err = repo.create(ctx, claims, req, currentDate)
		if err != nil {
			return nil, err
		}
		amount -= account.Target
		currentDate = currentDate.Add(4 * time.Second)
	}

	if serr := repo.notifySMS.Send(ctx, account.PhoneNumber, "sms/ds_received",
		map[string]interface{}{
			"Name":          account.Customer,
			"EffectiveDate": web.NewTimeResponse(ctx, time.Unix(tx.EffectiveDate, 0)).LocalDate,
			"Amount":        reqAmount,
			"Balance":       tx.OpeningBalance + tx.Amount,
		}); err != nil {
		// TODO: log critical error. Send message to monitoring account
		fmt.Println(serr)
	}
	return tx, err
}

// create inserts a new transaction into the database.
func (repo *Repository) create(ctx context.Context, claims auth.Claims, req CreateRequest, currentDate time.Time) (*Transaction, error) {
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

	account, err := repo.AccountRepo.ReadByNumber(ctx, req.AccountNumber)
	if err != nil {
		return nil, weberror.NewErrorMessage(ctx, err, http.StatusBadRequest, "Invalid account number")
	}

	// If now empty set it to the current time.
	if currentDate.IsZero() {
		currentDate = time.Now()
	}

	effectiveDate := now.New(currentDate).BeginningOfDay()
	if account.Type == models.AccountTypeDS {
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

	accountBalance := account.Balance

	m := Transaction{
		ID:             uuid.NewRandom().String(),
		AccountID:      account.ID,
		AccountNumber:  account.Number,
		CustomerID:     account.CustomerID,
		Customer:       account.Customer,
		SalesRep:       account.SalesRep,
		OpeningBalance: accountBalance,
		Amount:         req.Amount,
		Narration:      req.Narration,
		PaymentMethod:  req.PaymentMethod,
		Type:           req.Type,
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

	if _, err := repo.mongoDb.Collection(dal.C.Transaction).InsertOne(ctx, m); err != nil {
		return nil, errors.WithMessage(err, "Insert deposit failed")
	}

	var lastDepositDate = account.LastPaymentDate
	if req.Type == TransactionType_Deposit {
		lastDepositDate = m.EffectiveDate
		accountBalance += m.Amount
	} else {
		accountBalance -= m.Amount
	}

	if err = repo.SaveDailySummary(ctx, req.Amount, 0, 0, currentDate); err != nil {
		return nil, err
	}

	// send SMS notification
	if req.Type == TransactionType_Deposit {
		if account.Type == models.AccountTypeSB {
			if err = repo.notifySMS.Send(ctx, account.PhoneNumber, "sms/payment_received",
				map[string]interface{}{
					"Name":    account.Customer,
					"Amount":  req.Amount,
					"Balance": accountBalance,
				}); err != nil {
				// TODO: log critical error. Send message to monitoring account
				fmt.Println(err)
			}
		}
	} else {
		if err = repo.notifySMS.Send(ctx, account.PhoneNumber, "sms/payment_withdrawn",
			map[string]interface{}{
				"Name":    account.Customer,
				"Amount":  req.Amount,
				"Balance": accountBalance,
			}); err != nil {
			// TODO: log critical error. Send message to monitoring account
			fmt.Println(err)
		}
	}

	if req.Type == TransactionType_Deposit && account.Type == models.AccountTypeDS && isFirstContribution {
		wm := Transaction{
			ID:             uuid.NewRandom().String(),
			AccountID:      account.ID,
			OpeningBalance: accountBalance,
			Amount:         req.Amount,
			Narration:      "DS fee deduction",
			Type:           TransactionType_Withdrawal,
			SalesRepID:     claims.Subject,
			ReceiptNo:      repo.generateReceiptNumber(ctx),
			CreatedAt:      currentDate.Add(2 * time.Second).Unix(),
			UpdatedAt:      currentDate.Unix(),
			AccountNumber:  account.Number,
			Customer:       account.Customer,
			CustomerID:     account.CustomerID,
			SalesRep:       m.SalesRep,
		}

		if _, err := repo.mongoDb.Collection(dal.C.Transaction).InsertOne(ctx, wm); err != nil {
			return nil, errors.WithMessage(err, "Insert DS fee failed")
		}

		accountBalance -= req.Amount

		commission := dscommission.DsCommission{
			ID:            uuid.NewRandom().String(),
			AccountID:     account.ID,
			AccountNumber: account.Number,
			CustomerID:    account.CustomerID,
			Customer:      account.Customer,
			Amount:        req.Amount,
			Date:          currentDate.Unix(),
			EffectiveDate: effectiveDate.Unix(),
		}
		if _, err := repo.mongoDb.Collection(dal.C.DSCommission).InsertOne(ctx, commission); err != nil {
			return nil, err
		}
	}

	if _, err := repo.mongoDb.Collection(dal.C.Account).UpdateOne(ctx, bson.M{dal.AccountColumns.ID: account.ID}, bson.M{
		dal.AccountColumns.Balance:         accountBalance,
		dal.AccountColumns.LastPaymentDate: lastDepositDate,
	}); err != nil {
		return nil, err
	}

	return &m, nil
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
	if req.PaymentMethod == "Transfer" {
		if len(req.Narration) > 0 {
			createReq.Narration += " -"
		}
		if len(req.Bank) > 0 && len(req.BankAccountNumber) > 0 {
			createReq.Narration += fmt.Sprintf("%s - %s", req.Bank, req.BankAccountNumber)
		}

	}

	txn, err := repo.MakeDeduction(ctx, claims, createReq, now)
	if err != nil {
		return nil, err
	}
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
	var rec Transaction
	collection := repo.mongoDb.Collection(dal.C.Transaction)
	_ = collection.FindOne(ctx, bson.M{dal.TransactionColumns.ReceiptNo: receipt}).Decode(&rec)
	return rec.ID != ""
}

func (repo *Repository) lastDeposit(ctx context.Context, accountID string) (*Transaction, error) {
	queries := bson.M{
		dal.TransactionColumns.AccountID: accountID,
		dal.TransactionColumns.Type:      TransactionType_Deposit,
	}
	findOptions := options.Find()
	var limit int64 = 1
	findOptions.Limit = &limit
	findOptions.Sort = bson.M{dal.TransactionColumns.CreatedAt: "-1"}
	cursor, err := repo.mongoDb.Collection(dal.C.Transaction).Find(ctx, queries, findOptions)
	if err != nil {
		return nil, weberror.WithMessage(ctx, err, "Cannot get customer list")
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var c Transaction
		cursor.Decode(&c)
		return &c, nil
	}
	return nil, errors.New("Not found")
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

	tranx, err := repo.ReadByID(ctx, req.ID)
	if err != nil {
		return err
	}

	if tranx.ArchivedAt != nil {
		return errors.New("This transaction has been archived")
	}

	_, err = repo.mongoDb.Collection(dal.C.Transaction).
		UpdateOne(ctx, bson.M{dal.TransactionColumns.ID: req.ID}, bson.M{models.TransactionColumns.ArchivedAt: now})
	if err != nil {
		return err
	}

	var txAmount = tranx.Amount
	if tranx.Type == TransactionType_Deposit {
		txAmount *= -1
	}

	filter := bson.M{dal.AccountColumns.ID: tranx.AccountID}
	update := bson.M{dal.AccountColumns.Balance: bson.M{"$inc": txAmount}}
	_, err = repo.mongoDb.Collection(dal.C.Account).UpdateOne(ctx, filter, update)

	if err != nil {
		return err
	}

	return nil
}

// MakeDeduction inserts a new transaction of type withdrawal into the database.
func (repo *Repository) MakeDeduction(ctx context.Context, claims auth.Claims, req MakeDeductionRequest,
	now time.Time) (*Transaction, error) {

	span, ctx := tracer.StartSpanFromContext(ctx, "internal.transaction.MakeDeduction")
	defer span.Finish()
	if claims.Subject == "" {
		return nil, errors.WithStack(ErrForbidden)
	}

	account, err := repo.AccountRepo.ReadByNumber(ctx, req.AccountNumber)
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

	accountBalance := account.Balance
	if err != nil {
		return nil, err
	}

	if accountBalance < req.Amount {
		return nil, weberror.NewError(ctx, errors.New("insufficient fund"), 400)
	}

	salesRep, err := models.FindUser(ctx, repo.DbConn, claims.Subject)
	if err != nil {
		return nil, fmt.Errorf("FindUser, %v", err)
	}

	m := Transaction{
		ID:             uuid.NewRandom().String(),
		AccountID:      account.ID,
		AccountNumber:  account.Number,
		Customer:       account.Customer,
		CustomerID:     account.CustomerID,
		Type:           TransactionType_Withdrawal,
		OpeningBalance: accountBalance,
		Amount:         req.Amount,
		Narration:      req.Narration,
		SalesRepID:     claims.Subject,
		SalesRep:       salesRep.FirstName + " " + salesRep.LastName,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	if _, err := repo.mongoDb.Collection(dal.C.Transaction).InsertOne(ctx, m); err != nil {
		return nil, errors.WithMessage(err, "Insert deduction failed")
	}

	accountBalance -= req.Amount
	if _, err := repo.mongoDb.Collection(dal.C.Account).UpdateOne(ctx, bson.M{dal.AccountColumns.ID: m.AccountID}, bson.M{
		models.AccountColumns.Balance: bson.M{"$inc": accountBalance},
	}); err != nil {
		return nil, err
	}

	if err = repo.notifySMS.Send(ctx, account.PhoneNumber, "sms/payment_withdrawn",
		map[string]interface{}{
			"Name":    account.Customer,
			"Amount":  req.Amount,
			"Balance": accountBalance,
		}); err != nil {
		// TODO: log critical error. Send message to monitoring account
		fmt.Println(err)
	}

	return &m, nil
}

// SaveDailySummary saves the provided daily summary info to the db
func (repo *Repository) SaveDailySummary(ctx context.Context, income, expenditure, bankDeposit float64, date time.Time) error {
	span, ctx := tracer.StartSpanFromContext(ctx, "internal.transaction.SaveDailySummary")
	defer span.Finish()

	collection := repo.mongoDb.Collection(dal.C.DailySummary)

	today := now.New(date).BeginningOfDay().Unix()
	existingSummary, err := repo.FindDailySummary(ctx, today)
	if err == nil {
		existingSummary.BankDeposit += bankDeposit
		existingSummary.Income += income
		existingSummary.Expenditure += expenditure
		cols := bson.M{
			dal.DailySummaryColumns.BankDeposit: bson.M{"$inc": bankDeposit},
			dal.DailySummaryColumns.Income:      bson.M{"$inc": income},
			dal.DailySummaryColumns.Expenditure: bson.M{"$inc": expenditure},
		}

		_, err = collection.UpdateOne(ctx, bson.M{dal.DailySummaryColumns.Date: today}, cols)
		return err
	}

	model := DailySummary{
		Date:        today,
		BankDeposit: bankDeposit,
		Income:      income,
		Expenditure: expenditure,
	}

	_, err = collection.InsertOne(ctx, model)
	return err
}

func (repo *Repository) FindDailySummary(ctx context.Context, date int64) (*DailySummary, error) {
	var rec DailySummary
	err := repo.mongoDb.Collection(dal.C.DailySummary).FindOne(ctx, bson.M{dal.DailySummaryColumns.Date: date}).Decode(&rec)
	return &rec, err
}

func (repo *Repository) Migrate(ctx context.Context) error {
	if c, _ := repo.mongoDb.Collection(dal.C.Account).CountDocuments(ctx, bson.M{}); c > 0 {
		return nil
	}
	records, err := models.Transactions(
		qm.Load(models.TransactionRels.Account),
		qm.Load(models.TransactionRels.SalesRep),
	).All(ctx, repo.DbConn)
	if err != nil {
		return err
	}
	for _, b := range records {
		tx := FromModel(b)
		if _, err := repo.mongoDb.Collection(dal.C.Transaction).InsertOne(ctx, tx); err != nil {
			return err
		}
	}
	return nil
}
