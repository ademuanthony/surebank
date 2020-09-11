package customer

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/pborman/uuid"
	"github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"

	"merryworld/surebank/internal/dal"
	"merryworld/surebank/internal/platform/auth"
	"merryworld/surebank/internal/platform/web/webcontext"
	"merryworld/surebank/internal/platform/web/weberror"
	"merryworld/surebank/internal/postgres/models"
	"merryworld/surebank/internal/user"
)

// Find gets all the accounts from the database based on the request params.
func (repo *AccountRepository) Find(ctx context.Context, claims auth.Claims, req FindAccountRequest) (*AccountPagedResponseList, error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "internal.account.Find")
	defer span.Finish()

	return repo.find(ctx, claims, req, primitive.M{})
}

// FindDs gets all the accounts from the database that are of DS type and have > 0 balance.
func (repo *AccountRepository) FindDs(ctx context.Context, claims auth.Claims, req FindAccountRequest) (*AccountPagedResponseList, error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "internal.account.FindDs")
	defer span.Finish()

	var queries = bson.M{
		dal.AccountColumns.AccountType: models.AccountTypeDS,
		dal.AccountColumns.Balance:     bson.M{"$gt": 0},
	}

	return repo.find(ctx, claims, req, queries)
}

// FindDs gets all the accounts from the database that are of DS type and have > 0 balance.
func (repo *AccountRepository) Debtors(ctx context.Context, claims auth.Claims, req FindAccountRequest, currentDate time.Time) (*AccountPagedResponseList, error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "internal.account.Debtors")
	defer span.Finish()

	threeDaysAgo := currentDate.Add(-3 * 24 * time.Hour).Unix()
	thirtyDaysAgo := currentDate.Add(-30 * 24 * time.Hour).Unix()

	var queries = bson.M{
		dal.AccountColumns.LastPaymentDate: bson.M{"$gte": thirtyDaysAgo},
		dal.AccountColumns.LastPaymentDate: bson.M{"$lte": threeDaysAgo},
	}

	req.Order = append(req.Order, fmt.Sprintf("%s -1", dal.AccountColumns.LastPaymentDate))

	return repo.find(ctx, claims, req, queries)
}

func (repo *AccountRepository) find(ctx context.Context, claims auth.Claims, req FindAccountRequest,
	queries primitive.M) (*AccountPagedResponseList, error) {

	collection := repo.mongoDb.Collection(dal.C.Account)

	if !req.IncludeArchived {
		queries[dal.AccountColumns.ArchivedAt] = nil
	}

	if !claims.HasRole(auth.RoleAdmin) {
		queries[dal.AccountColumns.SalesRepID] = claims.Subject
	}

	if req.CustomerID != "" {
		queries[dal.AccountColumns.CustomerID] = req.CustomerID
	}

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
	var result Accounts
	for cursor.Next(ctx) {
		var c Account
		cursor.Decode(&c)
		result = append(result, c)
	}
	if err := cursor.Err(); err != nil {
		return nil, weberror.WithMessage(ctx, err, "Cannot get customer list")
	}

	if len(result) == 0 {
		return &AccountPagedResponseList{}, nil
	}

	return &AccountPagedResponseList{
		Accounts:   result.Response(ctx),
		TotalCount: totalCount,
	}, nil
}

// ReadByID gets the specified branch by ID from the database.
func (repo *AccountRepository) ReadByID(ctx context.Context, id string) (*Account, error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "internal.account.ReadByID")
	defer span.Finish()

	var rec Account
	collection := repo.mongoDb.Collection(dal.C.Account)
	err := collection.FindOne(ctx, bson.M{dal.AccountColumns.ID: id}).Decode(&rec)
	return &rec, err
}

func (repo *AccountRepository) AccountsCount(ctx context.Context, claims auth.Claims) (int64, error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "internal.account.AccountsCount")
	defer span.Finish()

	query := bson.M{}
	if !claims.HasRole(auth.RoleAdmin) {
		query[dal.AccountColumns.SalesRepID] = claims.Subject
	}

	collection := repo.mongoDb.Collection(dal.C.Account)
	return collection.CountDocuments(ctx, query)
}

// Create inserts a new account into the database.
func (repo *AccountRepository) Create(ctx context.Context, claims auth.Claims, req CreateAccountRequest, now time.Time) (*Account, error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "internal.account.Create")
	defer span.Finish()
	if claims.Audience == "" {
		return nil, errors.WithStack(ErrForbidden)
	}

	salesRep, err := models.Users(models.UserWhere.ID.EQ(claims.Subject)).One(ctx, repo.DbConn)
	if err != nil {
		return nil, weberror.NewErrorMessage(ctx, err, 400, "Something went wrong. Are you logged in?")
	}
	req.BranchID = salesRep.BranchID
	branch, err := repo.branchRepo.ReadByID(ctx, salesRep.BranchID)
	if err != nil {
		return nil, weberror.NewErrorMessage(ctx, err, 400, "Something went wrong. Cannot get your branch")
	}

	spew.Dump("\n Get customer \n")
	cust, err := repo.customerRepo.ReadByID(ctx, req.CustomerID)
	if err != nil {
		return nil, weberror.NewErrorMessage(ctx, err, 500, "Cannot read customer info")
	}

	spew.Dump("\n Got customer \n")
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

	m := Account{
		ID:         uuid.NewRandom().String(),
		Number:     repo.generateAccountNumber(ctx, req.Type),
		CustomerID: req.CustomerID,
		Type:       req.Type,
		Target:     req.Target,
		TargetInfo: req.TargetInfo,
		SalesRepID: claims.Subject,
		CreatedAt:  now,
		BranchID:   req.BranchID,
		UpdatedAt:  now,

		SalesRep:    user.FromModel(salesRep),
		Branch:      branch,
		Customer:    cust.Name,
		PhoneNumber: cust.PhoneNumber,
	}

	if _, err := repo.mongoDb.Collection(dal.C.Account).InsertOne(ctx, m); err != nil {
		return nil, weberror.WithMessage(ctx, err, "Insert account failed")
	}
	// add the account to the customer
	filter := bson.M{dal.CustomerColumns.ID: cust.ID}
	update := bson.M{"$set": bson.M{dal.CustomerColumns.Accounts: cust.Accounts}}
	if _, err := repo.mongoDb.Collection(dal.C.Customer).UpdateOne(ctx, filter, update); err != nil {
		_, _ = repo.mongoDb.Collection(dal.C.Customer).DeleteOne(ctx, bson.M{dal.AccountColumns.ID: m.ID})
		return nil, err
	}

	return &m, nil
}

func (repo *AccountRepository) generateAccountNumber(ctx context.Context, accountType string) string {
	var accountNumber string
	for accountNumber == "" || repo.accountNumberExists(ctx, accountNumber) {
		accountNumber = accountType
		rand.Seed(time.Now().UTC().UnixNano())
		for i := 0; i < 5; i++ {
			accountNumber += strconv.Itoa(rand.Intn(10))
		}
	}
	return accountNumber
}

func (repo *AccountRepository) accountNumberExists(ctx context.Context, number string) bool {
	var rec Account
	collection := repo.mongoDb.Collection(dal.C.Account)
	_ = collection.FindOne(ctx, bson.M{dal.AccountColumns.Number: number}).Decode(&rec)
	return rec.ID != ""
}

// Update replaces an account in the database.
func (repo *AccountRepository) Update(ctx context.Context, claims auth.Claims, req UpdateAccountRequest, now time.Time) error {
	span, ctx := tracer.StartSpanFromContext(ctx, "internal.account.Update")
	defer span.Finish()

	if claims.Audience == "" {
		return errors.WithStack(ErrForbidden)
	}
	// Admin users can update branches they have access to.
	if !claims.HasRole(auth.RoleAdmin) {
		return errors.WithStack(ErrForbidden)
	}

	// Validate the request.
	v := webcontext.Validator()
	err := v.Struct(req)
	if err != nil {
		return err
	}

	cols := bson.M{}
	if req.TargetInfo != nil {
		cols[dal.AccountColumns.TargetInfo] = *req.TargetInfo
	}

	if req.Target != nil {
		cols[dal.AccountColumns.Target] = *req.Target
	}

	if req.Type != nil {
		cols[dal.AccountColumns.AccountType] = *req.Type
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

	cols[dal.AccountColumns.UpdatedAt] = now.Unix()

	collection := repo.mongoDb.Collection(dal.C.Account)
	_, err = collection.UpdateOne(ctx, bson.M{dal.AccountColumns.ID: req.ID}, bson.M{"$set": cols})

	if err != nil {
		return weberror.NewError(ctx, err, 500)
	}

	return nil
}

// Archive soft deleted the account from the database.
func (repo *AccountRepository) Archive(ctx context.Context, claims auth.Claims, req ArchiveAccountRequest) error {
	span, ctx := tracer.StartSpanFromContext(ctx, "internal.account.Archive")
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

	collection := repo.mongoDb.Collection(dal.C.Account)
	_, err = collection.DeleteOne(ctx, bson.M{dal.AccountColumns.ID: req.ID})
	// TODO: delete transactions

	return err
}

func (repo *AccountRepository) DbConnCount(ctx context.Context) (int, error) {
	result := struct {
		Number int `json:"number"`
	}{}
	statement := "SELECT sum(numbackends) as number FROM pg_stat_database;"
	if err := models.NewQuery(qm.SQL(statement)).Bind(ctx, repo.DbConn, &result); err != nil {
		return 0, err
	}
	return result.Number, nil
}
