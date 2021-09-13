package profit

import (
	"context"
	"database/sql"
	"time"

	"merryworld/surebank/internal/platform/auth"
	"merryworld/surebank/internal/platform/web"
	"merryworld/surebank/internal/platform/web/webcontext"
	"merryworld/surebank/internal/postgres/models"

	"github.com/jmoiron/sqlx"
	"github.com/pborman/uuid"
	"github.com/pkg/errors"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	. "github.com/volatiletech/sqlboiler/v4/queries/qm"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

var (
	// ErrNotFound abstracts the postgres not found error.
	ErrNotFound = errors.New("Entity not found")

	// ErrForbidden occurs when a user tries to do something that is forbidden to them according to our access control policies.
	ErrForbidden = errors.New("Attempted action is not allowed")
)

// Repository defines the required dependencies for Project.
type Repository struct {
	DbConn *sqlx.DB
}

// NewRepository creates a new Repository that defines dependencies for Project.
func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		DbConn: db,
	}
}

// Profit
type Profit struct {
	ID         string    `json:"id"`
	Amount     float64   `json:"price"`
	Narration  string    `json:"image"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	ArchivedAt null.Time `json:"archived_at"`
}

func (m Profit) ToModel() models.Profit {
	return models.Profit{
		ID:         m.ID,
		Narration:  m.Narration,
		Amount:     m.Amount,
		CreatedAt:  m.CreatedAt.Unix(),
		UpdatedAt:  m.UpdatedAt.Unix(),
		ArchivedAt: null.Int64From(m.ArchivedAt.Time.Unix()),
	}
}

func ProfitFromModel(profit *models.Profit) *Profit {
	p := &Profit{
		ID:        profit.ID,
		Amount:    profit.Amount,
		Narration: profit.Narration,
		CreatedAt: time.Unix(profit.CreatedAt, 0),
		UpdatedAt: time.Unix(profit.UpdatedAt, 0),
	}
	if profit.ArchivedAt.Valid {
		p.ArchivedAt = null.TimeFrom(time.Unix(profit.ArchivedAt.Int64, 0))
	}

	return p
}

// ProfitResponse represent a Profit that is returned for display
type ProfitResponse struct {
	ID        string  `json:"id"`
	Narration string  `json:"narration"`
	Amount    float64 `json:"amount"`

	CreatedAt web.TimeResponse `json:"created_at"`
	UpdatedAt web.TimeResponse `json:"updated_at"`
}

// Response transforms Profit to ProfitResponse that is used for display.
func (m *Profit) Response(ctx context.Context) *ProfitResponse {
	if m == nil {
		return nil
	}

	r := &ProfitResponse{
		ID:        m.ID,
		Narration: m.Narration,
		Amount:    m.Amount,
		CreatedAt: web.NewTimeResponse(ctx, m.CreatedAt),
		UpdatedAt: web.NewTimeResponse(ctx, m.UpdatedAt),
	}

	return r
}

// Profits a list of Profits.
type Profits []*Profit

// Response transforms a list of Profits to a list of ProfitResponse.
func (m *Profits) Response(ctx context.Context) []*ProfitResponse {
	var l []*ProfitResponse
	if m != nil && len(*m) > 0 {
		for _, n := range *m {
			l = append(l, n.Response(ctx))
		}
	}

	return l
}

// ProfitCreateRequest contains information needed to create a new Profit.
type ProfitCreateRequest struct {
	Narration string  `json:"narration" validate:"required,unique"  example:"Bread"`
	Amount    float64 `json:"amount" validate:"required"`
}

type ProfitUpdateRequest struct {
	ID        string   `json:"id"`
	Amount    *float64 `json:"amount"`
	Narration *string  `json:"narration"`
}

// ProfitReadRequest defines the information needed to read a Profit.
type ProfitReadRequest struct {
	ID              string `json:"id" validate:"required,uuid" example:"985f1746-1d9f-459f-a2d9-fc53ece5ae86"`
	IncludeArchived bool   `json:"include-archived" example:"false"`
}

// ProfitArchiveRequest defines the information needed to archive a Profit. This will archive (soft-delete) the
// existing database entry.
type ProfitArchiveRequest struct {
	ID string `json:"id" validate:"required,uuid" example:"985f1746-1d9f-459f-a2d9-fc53ece5ae86"`
}

// ProfitDeleteRequest defines the information needed to delete a Profit.
type ProfitDeleteRequest struct {
	ID string `json:"id" validate:"required,uuid" example:"985f1746-1d9f-459f-a2d9-fc53ece5ae86"`
}

// ServiceTypeFindRequest defines the possible options to search for ServiceTypes. By default
// archived ServiceType will be excluded from response.
type ProfitFindRequest struct {
	Where           string        `json:"where" example:"name = ? and status = ?"`
	Args            []interface{} `json:"args" swaggertype:"array,string" example:"Moon Launch,active"`
	Order           []string      `json:"order" example:"created_at desc"`
	Limit           *uint         `json:"limit" example:"10"`
	Offset          *uint         `json:"offset" example:"20"`
	IncludeArchived bool          `json:"include-archived" example:"false"`
}

// FindProfit gets all the profits from the database based on the request params.
func (repo Repository) FindProfit(ctx context.Context, req ProfitFindRequest) (Profits, error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "internal.shop.FindProfit")
	defer span.Finish()
	var queries []QueryMod

	if req.Where != "" {
		queries = append(queries, Where(req.Where, req.Args...))
	}

	if req.Limit != nil {
		queries = append(queries, Limit(int(*req.Limit)))
	}

	if req.Offset != nil {
		queries = append(queries, Offset(int(*req.Offset)))
	}

	ProfitSlice, err := models.Profits(queries...).All(ctx, repo.DbConn)
	if err != nil {
		if err.Error() == sql.ErrNoRows.Error() {
			return Profits{}, nil
		}
		return nil, err
	}

	var result Profits
	for _, rec := range ProfitSlice {
		result = append(result, ProfitFromModel(rec))
	}

	return result, nil
}

// ReadProfitByID gets the specified Profit by ID from the database.
func (repo *Repository) ReadProfitByID(ctx context.Context, _ auth.Claims, id string) (*Profit, error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "internal.shop.ReadProfitByID")
	defer span.Finish()
	queries := []QueryMod{
		models.ProfitWhere.ID.EQ(id),
	}
	profitModel, err := models.Profits(queries...).One(ctx, repo.DbConn)
	if err != nil {
		return nil, err
	}

	return ProfitFromModel(profitModel), nil
}

func (repo *Repository) ReadProfitByIDTx(ctx context.Context, _ auth.Claims, id string, tx *sql.Tx) (*Profit, error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "internal.shop.ReadProfitByIDTx")
	defer span.Finish()
	queries := []QueryMod{
		models.ProfitWhere.ID.EQ(id),
	}
	profitModel, err := models.Profits(queries...).One(ctx, tx)
	if err != nil {
		return nil, err
	}

	return ProfitFromModel(profitModel), nil
}

func (repo *Repository) CreateProfit(ctx context.Context, claims auth.Claims, req ProfitCreateRequest, now time.Time) (*Profit, error) {

	tx, err := repo.DbConn.Begin()
	if err != nil {
		return nil, errors.WithStack(errors.WithMessage(err, "create profit failed, cannot start db transaction"))
	}

	s, err := repo.CreateProfitTx(ctx, tx, claims, req, now)
	if err != nil {
		_ = tx.Rollback()
		return nil, errors.WithStack(errors.WithMessage(err, "create profit failed, cannot commit db transaction"))
	}

	if err = tx.Commit(); err != nil {
		return nil, errors.WithStack(errors.WithMessage(err, "create profit failed, cannot commit db transaction"))
	}

	return s, nil
}

func (repo *Repository) CreateProfitTx(ctx context.Context, tx *sql.Tx, claims auth.Claims, req ProfitCreateRequest, now time.Time) (*Profit, error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "internal.shop.CreateProfit")
	defer span.Finish()
	if claims.Audience == "" {
		return nil, errors.WithStack(ErrForbidden)
	}

	if !claims.HasRole(auth.RoleAdmin) {
		return nil, errors.WithStack(ErrForbidden)
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

	s := Profit{
		ID:        uuid.NewRandom().String(),
		Amount:    req.Amount,
		Narration: req.Narration,
		CreatedAt: now,
		UpdatedAt: now,
	}

	prodModel := s.ToModel()
	if err := prodModel.Insert(ctx, tx, boil.Infer()); err != nil {
		_ = tx.Rollback()
		return nil, errors.WithMessage(err, "create Profit failed")
	}

	return &s, nil
}

// UpdateProfit replaces an project in the database.
func (repo *Repository) UpdateProfit(ctx context.Context, claims auth.Claims, req ProfitUpdateRequest, now time.Time) error {
	span, ctx := tracer.StartSpanFromContext(ctx, "internal.shop.UpdateProfit")
	defer span.Finish()

	if claims.Audience == "" {
		return errors.WithStack(ErrForbidden)
	}
	// Admin users can update Profits they have access to.
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
		cols[models.ProfitColumns.Narration] = *req.Narration
	}

	if req.Amount != nil {
		cols[models.ProfitColumns.Amount] = *req.Amount
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

	cols[models.ProfitColumns.UpdatedAt] = now

	if len(cols) == 0 {
		return nil
	}

	_, err = models.Profits(models.ProfitWhere.ID.EQ(req.ID)).UpdateAll(ctx, repo.DbConn, cols)

	return err
}

// Archive soft deletes the profit from the database.
func (repo *Repository) ArchiveProfit(ctx context.Context, claims auth.Claims, req ProfitArchiveRequest, now time.Time) error {
	span, ctx := tracer.StartSpanFromContext(ctx, "internal.shop.ArchiveProfit")
	defer span.Finish()

	// Validate the request.
	v := webcontext.Validator()
	err := v.Struct(req)
	if err != nil {
		return err
	}

	if claims.Audience == "" {
		return errors.WithStack(ErrForbidden)
	}
	// Admin users can update Profits they have access to.
	if !claims.HasRole(auth.RoleAdmin) {
		return errors.WithStack(ErrForbidden)
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

	cols := models.M{models.ProfitColumns.ArchivedAt: now}
	if _, err := models.Profits(models.ProfitWhere.ID.EQ(req.ID)).UpdateAll(ctx, repo.DbConn, cols); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// Delete removes an Profit from the database.
func (repo *Repository) DeleteProfit(ctx context.Context, claims auth.Claims, req ProfitDeleteRequest) error {
	span, ctx := tracer.StartSpanFromContext(ctx, "internal.shop.DeleteProfit")
	defer span.Finish()

	// Validate the request.
	v := webcontext.Validator()
	err := v.Struct(req)
	if err != nil {
		return err
	}

	// Ensure the claims can modify the project specified in the request.
	if claims.Audience == "" {
		return errors.WithStack(ErrForbidden)
	}
	// Admin users can update Profits they have access to.
	if !claims.HasRole(auth.RoleAdmin) {
		return errors.WithStack(ErrForbidden)
	}

	_, err = models.Profits(models.ProfitWhere.ID.EQ(req.ID)).DeleteAll(ctx, repo.DbConn)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
