package dscommission

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"merryworld/surebank/internal/dal"
	"merryworld/surebank/internal/platform/auth"
	"merryworld/surebank/internal/platform/web/weberror"
	"merryworld/surebank/internal/postgres/models"

	"github.com/jinzhu/now"
	. "github.com/volatiletech/sqlboiler/queries/qm"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func (repo *Repository) StartingNewCircle(ctx context.Context, accountID string, effectiveDate time.Time) (bool, error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "internal.commission.StartingNewCircle")
	defer span.Finish()

	lastCommission, err := repo.LattestCommission(ctx, accountID)
	if err != nil {
		if err.Error() == sql.ErrNoRows.Error() {
			return true, nil
		}
		return false, err
	}
	lastDate := now.New(time.Unix(lastCommission.EffectiveDate, 0)).BeginningOfDay()
	effectiveDate = now.New(effectiveDate).BeginningOfDay()
	duration := effectiveDate.Sub(lastDate)
	r := duration.Hours() >= float64(31*24)
	return r, nil
}

func (repo *Repository) LattestCommission(ctx context.Context, accountID string) (*DsCommission, error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "internal.commission.LattestCommission")
	defer span.Finish()

	var comm DsCommission
	q := bson.M{dal.DSCommissionColumns.AccountID: accountID}
	findOption := options.FindOne()
	findOption.SetSort(bson.M{dal.DSCommissionColumns.EffectiveDate: -1})
	err := repo.mongoDb.Collection(dal.C.DSCommission).FindOne(ctx, q, findOption).Decode(&comm)

	if err != nil && err != mongo.ErrNoDocuments {
		return nil, err
	}
	return &comm, nil
}

// Find gets all the commissions from the database based on the request params.
func (repo *Repository) Find(ctx context.Context, claims auth.Claims, req FindRequest) (*PagedResponseList, error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "internal.commission.Find")
	defer span.Finish()

	var queries bson.M
	if req.StateDate > 0 {
		queries[dal.DSCommissionColumns.Date] = bson.M{"$gte": req.StateDate}
	}
	if req.EndDate > 0 {
		queries[dal.DSCommissionColumns.Date] = bson.M{"$lte": req.EndDate}
	}

	collection := repo.mongoDb.Collection(dal.C.DSCommission)
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
			s, _ := strconv.Atoi(sortInfo[1])
			sort = append(sort, primitive.E{Key: sortInfo[0], Value: s})
		}
	}
	findOptions.SetSort(sort)

	if req.Limit != nil {
		findOptions.SetLimit(int64(*req.Limit))
	}

	if req.Offset != nil {
		findOptions.Skip = req.Offset
	}

	cursor, err := collection.Find(ctx, queries, findOptions)
	if err != nil {
		return nil, weberror.WithMessage(ctx, err, "Cannot get customer list")
	}
	defer cursor.Close(ctx)
	var result DsCommissions
	for cursor.Next(ctx) {
		var c DsCommission
		if err = cursor.Decode(&c); err != nil {
			return nil, fmt.Errorf("Decode, %v", err)
		}
		result = append(result, c)
	}
	if err := cursor.Err(); err != nil {
		return nil, weberror.WithMessage(ctx, err, "Cannot get customer list")
	}

	if len(result) == 0 {
		return &PagedResponseList{}, nil
	}

	return &PagedResponseList{
		Items:      result.Response(ctx),
		TotalCount: totalCount,
	}, nil
}

// ReadByID gets the specified commission by ID from the database.
func (repo *Repository) ReadByID(ctx context.Context, _ auth.Claims, id string) (*DsCommission, error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "internal.commission.ReadByID")
	defer span.Finish()

	queries := []QueryMod{
		models.DSCommissionWhere.ID.EQ(id),
		Load(models.DSCommissionRels.Account),
		Load(models.DSCommissionRels.Customer),
	}
	model, err := models.DSCommissions(queries...).One(ctx, repo.DbConn)
	if err != nil {
		return nil, err
	}

	return FromModel(model), nil
}

func (repo *Repository) TotalAmountByWhere(ctx context.Context, req FindRequest) (float64, error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "internal.commission.TotalAmountByWhere")
	defer span.Finish()

	var queries bson.M
	if req.StateDate > 0 {
		queries[dal.DSCommissionColumns.Date] = bson.M{"$gte": req.StateDate}
	}
	if req.EndDate > 0 {
		queries[dal.DSCommissionColumns.Date] = bson.M{"$lte": req.EndDate}
	}

	var result []struct {
		Total float64
	}

	pipeline := []bson.M{
		{
			"$match": queries,
		},
		{
			"$group": bson.M{
				"_id":   "",
				"total": bson.M{"$sum": "$" + dal.DSCommissionColumns.Amount},
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
		return 0, fmt.Errorf("TotalAmountByWhere -> Aggregate, %s", err.Error())
	}

	err = cursor.All(ctx, &result)
	if err != nil {
		return 0, fmt.Errorf("TotalAmountByWhere -> All, %s", err.Error())
	}
	if len(result) > 0 {
		return result[0].Total, err
	}
	return 0, err
}
