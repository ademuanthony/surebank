package branch

import (
	"context"
	"time"

	"merryworld/surebank/internal/platform/web"
	"merryworld/surebank/internal/postgres/models"

	"github.com/jmoiron/sqlx"
	"go.mongodb.org/mongo-driver/mongo"
)

// Repository defines the required dependencies for Branch.
type Repository struct {
	DbConn  *sqlx.DB
	mongoDb *mongo.Database
}

// NewRepository creates a new Repository that defines dependencies for Branch.
func NewRepository(db *sqlx.DB, mongoDb *mongo.Database) *Repository {
	return &Repository{
		DbConn:  db,
		mongoDb: mongoDb,
	}
}

// Branch represents a workflow.
type Branch struct {
	ID         string     `json:"id" validate:"required,uuid" example:"985f1746-1d9f-459f-a2d9-fc53ece5ae86"`
	Name       string     `json:"name"  validate:"required" example:"Rocket Launch"`
	CreatedAt  time.Time  `json:"created_at" truss:"api-read"`
	UpdatedAt  time.Time  `json:"updated_at" truss:"api-read"`
	ArchivedAt *time.Time `json:"archived_at,omitempty" truss:"api-hide"`
}

const CollectionName = "branch"

var Columns = struct {
	ID         string
	Name       string
	CreatedAt  string
	UpdatedAt  string
	ArchivedAt string
}{
	ID:         "id",
	Name:       "name",
	CreatedAt:  "created_at",
	UpdatedAt:  "updated_at",
	ArchivedAt: "archived_at",
}

func FromModel(rec *models.Branch) *Branch {
	b := &Branch{
		ID:        rec.ID,
		Name:      rec.Name,
		CreatedAt: time.Unix(rec.CreatedAt, 0),
		UpdatedAt: time.Unix(rec.UpdatedAt, 0),
	}
	if rec.ArchivedAt.Valid {
		archivedAt := time.Unix(rec.ArchivedAt.Int64, 0)
		b.ArchivedAt = &archivedAt
	}

	return b
}

// Response represents a workflow that is returned for display.
type Response struct {
	ID         string            `json:"id" validate:"required,uuid" example:"985f1746-1d9f-459f-a2d9-fc53ece5ae86"`
	Name       string            `json:"name"  validate:"required" example:"Rocket Launch"`
	CreatedAt  web.TimeResponse  `json:"created_at"`            // CreatedAt contains multiple format options for display.
	UpdatedAt  web.TimeResponse  `json:"updated_at"`            // UpdatedAt contains multiple format options for display.
	ArchivedAt *web.TimeResponse `json:"archived_at,omitempty"` // ArchivedAt contains multiple format options for display.
}

// Response transforms Branch to the Response that is used for display.
// Additional filtering by context values or translations could be applied.
func (m *Branch) Response(ctx context.Context) *Response {
	if m == nil {
		return nil
	}

	r := &Response{
		ID:        m.ID,
		Name:      m.Name,
		CreatedAt: web.NewTimeResponse(ctx, m.CreatedAt),
		UpdatedAt: web.NewTimeResponse(ctx, m.UpdatedAt),
	}

	if m.ArchivedAt != nil && !m.ArchivedAt.IsZero() {
		at := web.NewTimeResponse(ctx, *m.ArchivedAt)
		r.ArchivedAt = &at
	}

	return r
}

// Branches a list of Branches.
type Branches []*Branch

// Response transforms a list of Branches to a list of Responses.
func (m *Branches) Response(ctx context.Context) []*Response {
	var l []*Response
	if m != nil && len(*m) > 0 {
		for _, n := range *m {
			l = append(l, n.Response(ctx))
		}
	}

	return l
}

// CreateRequest contains information needed to create a new Branch.
type CreateRequest struct {
	Name string `json:"name" validate:"required"  example:"Rocket Launch"`
}

// ReadRequest defines the information needed to read a checklist.
type ReadRequest struct {
	ID              string `json:"id" validate:"required,uuid" example:"985f1746-1d9f-459f-a2d9-fc53ece5ae86"`
	IncludeArchived bool   `json:"include-archived" example:"false"`
}

// UpdateRequest defines what information may be provided to modify an existing
// Branch. All fields are optional so clients can send just the fields they want
// changed. It uses pointer fields so we can differentiate between a field that
// was not provided and a field that was provided as explicitly blank.
type UpdateRequest struct {
	ID   string  `json:"id" validate:"required,uuid" example:"985f1746-1d9f-459f-a2d9-fc53ece5ae86"`
	Name *string `json:"name,omitempty" validate:"omitempty,unique" example:"Rocket Launch to Moon"`
}

// ArchiveRequest defines the information needed to archive a checklist. This will archive (soft-delete) the
// existing database entry.
type ArchiveRequest struct {
	ID string `json:"id" validate:"required,uuid" example:"985f1746-1d9f-459f-a2d9-fc53ece5ae86"`
}

// DeleteRequest defines the information needed to delete a branch.
type DeleteRequest struct {
	ID string `json:"id" validate:"required,uuid" example:"985f1746-1d9f-459f-a2d9-fc53ece5ae86"`
}

// FindRequest defines the possible options to search for branches. By default
// archived checklist will be excluded from response.
type FindRequest struct {
	Where           string        `json:"where" example:"name = ? and status = ?"`
	Args            []interface{} `json:"args" swaggertype:"array,string" example:"Moon Launch,active"`
	Order           []string      `json:"order" example:"created_at desc"`
	Limit           *uint         `json:"limit" example:"10"`
	Offset          *uint         `json:"offset" example:"20"`
	IncludeArchived bool          `json:"include-archived" example:"false"`
}
