package dscommission

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"merryworld/surebank/internal/platform/auth"
	"merryworld/surebank/internal/platform/web/weberror"
	"merryworld/surebank/internal/postgres/models"

	"github.com/jinzhu/now"
	. "github.com/volatiletech/sqlboiler/v4/queries/qm"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func (repo *Repository) StartingNewCircle(ctx context.Context, accountID string, effectiveDate time.Time, dbTx *sql.Tx) (bool, error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "internal.commission.StartingNewCircle")
	defer span.Finish()

	lastCommission, err := repo.LattestCommission(ctx, accountID, dbTx)
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

func (repo *Repository) LattestCommission(ctx context.Context, accountID string, dbTx *sql.Tx) (*DsCommission, error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "internal.commission.LattestCommission")
	defer span.Finish()

	c, err := models.DSCommissions(
		models.DSCommissionWhere.AccountID.EQ(accountID),
		OrderBy(fmt.Sprintf("%s desc", models.DSCommissionColumns.EffectiveDate)),
		Limit(1),
	).One(ctx, dbTx)

	if err != nil {
		return nil, err
	}
	return FromModel(c), nil
}

// Find gets all the commissions from the database based on the request params.
func (repo *Repository) Find(ctx context.Context, claims auth.Claims, req FindRequest) (*PagedResponseList, error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "internal.commission.Find")
	defer span.Finish()

	var queries []QueryMod

	if req.Where != "" {
		queries = append(queries, Where(req.Where, req.Args...))
	}

	totalCount, err := models.DSCommissions(queries...).Count(ctx, repo.DbConn)
	if err != nil {
		return nil, weberror.WithMessage(ctx, err, "Cannot get commission count")
	}

	if req.IncludeAccount {
		queries = append(queries, Load(models.DSCommissionRels.Account))
	}

	if req.IncludeCustomer {
		queries = append(queries, Load(models.DSCommissionRels.Customer))
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

	slice, err := models.DSCommissions(queries...).All(ctx, repo.DbConn)
	if err != nil {
		if err.Error() == sql.ErrNoRows.Error() {
			return &PagedResponseList{}, nil
		}
		return nil, weberror.NewError(ctx, err, 500)
	}

	var result DsCommissions
	for _, rec := range slice {
		result = append(result, FromModel(rec))
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

func (repo *Repository) TotalAmountByWhere(ctx context.Context, where string, args []interface{}) (float64, error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "internal.commission.TotalAmountByWhere")
	defer span.Finish()

	statement := `select sum(amount) total from ds_commission `
	if len(where) > 0 {
		statement += fmt.Sprintf(" where %s ", where)
	}
	var result struct {
		Total sql.NullFloat64
	}
	err := models.NewQuery(SQL(statement, args...)).Bind(ctx, repo.DbConn, &result)
	return result.Total.Float64, err
}
