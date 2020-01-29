package handlers

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"merryworld/surebank/internal/account"
	"merryworld/surebank/internal/checklist"
	"merryworld/surebank/internal/platform/auth"
	"merryworld/surebank/internal/platform/web"
	"merryworld/surebank/internal/platform/web/webcontext"
	"merryworld/surebank/internal/platform/web/weberror"

	"github.com/pkg/errors"
	"gopkg.in/go-playground/validator.v9"
)

// Accounts represents the Account API method handler set.
type Accounts struct {
	Repository *account.Repository

	// ADD OTHER STATE LIKE THE LOGGER IF NEEDED.
}

// Find godoc
// @Summary List accounts
// @Description Find returns the existing accounts in the system.
// @Tags account
// @Accept  json
// @Produce  json
// @Security OAuth2Password
// @Param where				query string 	false	"Filter string, example: number = 'SB10000001'"
// @Param order				query string   	false 	"Order columns separated by comma, example: created_at desc"
// @Param limit				query integer  	false 	"Limit, example: 10"
// @Param offset			query integer  	false 	"Offset, example: 20"
// @Param include-archived query boolean 	false 	"Included Archived, example: false"
// @Param include-customer query boolean 	false 	"Included Customer info, example: false"
// @Param include-branch query boolean 	false 	"Included Branch info, example: false"
// @Param include-sales-rep query boolean 	false 	"Included Sale rep info, example: false"
// @Success 200 {object} account.PagedResponseList
// @Failure 400 {object} weberror.ErrorResponse
// @Failure 403 {object} weberror.ErrorResponse
// @Failure 500 {object} weberror.ErrorResponse
// @Router /accounts [get]
func (h *Accounts) Find(ctx context.Context, w http.ResponseWriter, r *http.Request, _ map[string]string) error {
	claims, ok := ctx.Value(auth.Key).(auth.Claims)
	if !ok {
		return errors.New("claims missing from context")
	}

	var req account.FindRequest

	// Handle where query value if set.
	if v := r.URL.Query().Get("where"); v != "" {
		where, args, err := web.ExtractWhereArgs(v)
		if err != nil {
			return web.RespondJsonError(ctx, w, weberror.NewError(ctx, err, http.StatusBadRequest))
		}
		req.Where = web.SqlBoilderWhere(where)
		req.Args = args
	}

	// Handle order query value if set.
	if v := r.URL.Query().Get("order"); v != "" {
		for _, o := range strings.Split(v, ",") {
			o = strings.TrimSpace(o)
			if o != "" {
				req.Order = append(req.Order, o)
			}
		}
	}

	// Handle limit query value if set.
	if v := r.URL.Query().Get("limit"); v != "" {
		l, err := strconv.Atoi(v)
		if err != nil {
			err = errors.WithMessagef(err, "unable to parse %s as int for limit param", v)
			return web.RespondJsonError(ctx, w, weberror.NewError(ctx, err, http.StatusBadRequest))
		}
		ul := uint(l)
		req.Limit = &ul
	}

	// Handle offset query value if set.
	if v := r.URL.Query().Get("offset"); v != "" {
		l, err := strconv.Atoi(v)
		if err != nil {
			err = errors.WithMessagef(err, "unable to parse %s as int for offset param", v)
			return web.RespondJsonError(ctx, w, weberror.NewError(ctx, err, http.StatusBadRequest))
		}
		ul := uint(l)
		req.Limit = &ul
	}

	// Handle include-archive query value if set.
	if v := r.URL.Query().Get("include-archived"); v != "" {
		b, err := strconv.ParseBool(v)
		if err != nil {
			err = errors.WithMessagef(err, "unable to parse %s as boolean for include-archived param", v)
			return web.RespondJsonError(ctx, w, weberror.NewError(ctx, err, http.StatusBadRequest))
		}
		req.IncludeArchived = b
	}

	// Handle include-branch query value if set.
	if v := r.URL.Query().Get("include-branch"); v != "" {
		b, err := strconv.ParseBool(v)
		if err != nil {
			err = errors.WithMessagef(err, "unable to parse %s as boolean for include-archived param", v)
			return web.RespondJsonError(ctx, w, weberror.NewError(ctx, err, http.StatusBadRequest))
		}
		req.IncludeBranch = b
	}

	// Handle include-customer query value if set.
	if v := r.URL.Query().Get("include-customer"); v != "" {
		b, err := strconv.ParseBool(v)
		if err != nil {
			err = errors.WithMessagef(err, "unable to parse %s as boolean for include-archived param", v)
			return web.RespondJsonError(ctx, w, weberror.NewError(ctx, err, http.StatusBadRequest))
		}
		req.IncludeCustomer = b
	}

	// Handle include-sales-rep query value if set.
	if v := r.URL.Query().Get("include-sales-rep"); v != "" {
		b, err := strconv.ParseBool(v)
		if err != nil {
			err = errors.WithMessagef(err, "unable to parse %s as boolean for include-archived param", v)
			return web.RespondJsonError(ctx, w, weberror.NewError(ctx, err, http.StatusBadRequest))
		}
		req.IncludeSalesRep = b
	}

	res, err := h.Repository.Find(ctx, claims, req)
	if err != nil {
		return err
	}

	return web.RespondJson(ctx, w, res, http.StatusOK)
}

// Read godoc
// @Summary Get account by ID.
// @Description Read returns the specified account from the system.
// @Tags account
// @Accept  json
// @Produce  json
// @Security OAuth2Password
// @Param id path string true "Account ID"
// @Success 200 {object} account.Response
// @Failure 400 {object} weberror.ErrorResponse
// @Failure 404 {object} weberror.ErrorResponse
// @Failure 500 {object} weberror.ErrorResponse
// @Router /accounts/{id} [get]
func (h *Accounts) Read(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {
	claims, ok := ctx.Value(auth.Key).(auth.Claims)
	if !ok {
		return errors.New("claims missing from context")
	}

	res, err := h.Repository.ReadByID(ctx, claims, params["id"])
	if err != nil {
		cause := errors.Cause(err)
		switch cause {
		case checklist.ErrNotFound:
			return web.RespondJsonError(ctx, w, weberror.NewError(ctx, err, http.StatusNotFound))
		default:
			return errors.Wrapf(err, "ID: %s", params["id"])
		}
	}

	return web.RespondJson(ctx, w, res.Response(ctx), http.StatusOK)
}

// Create godoc
// @Summary Create new account.
// @Description Create inserts a new account into the system.
// @Tags account
// @Accept  json
// @Produce  json
// @Security OAuth2Password
// @Param data body account.CreateRequest true "Account details"
// @Success 201 {object} account.Response
// @Failure 400 {object} weberror.ErrorResponse
// @Failure 403 {object} weberror.ErrorResponse
// @Failure 404 {object} weberror.ErrorResponse
// @Failure 500 {object} weberror.ErrorResponse
// @Router /accounts [post]
func (h *Accounts) Create(ctx context.Context, w http.ResponseWriter, r *http.Request, _ map[string]string) error {
	v, err := webcontext.ContextValues(ctx)
	if err != nil {
		return err
	}

	claims, err := auth.ClaimsFromContext(ctx)
	if err != nil {
		return err
	}

	var req account.CreateRequest
	if err := web.Decode(ctx, r, &req); err != nil {
		if _, ok := errors.Cause(err).(*weberror.Error); !ok {
			err = weberror.NewError(ctx, err, http.StatusBadRequest)
		}
		return web.RespondJsonError(ctx, w, err)
	}

	res, err := h.Repository.Create(ctx, claims, req, v.Now)
	if err != nil {
		cause := errors.Cause(err)
		switch cause {
		case checklist.ErrForbidden:
			return web.RespondJsonError(ctx, w, weberror.NewError(ctx, err, http.StatusForbidden))
		default:
			_, ok := cause.(validator.ValidationErrors)
			if ok {
				return web.RespondJsonError(ctx, w, weberror.NewError(ctx, err, http.StatusBadRequest))
			}
			return errors.Wrapf(err, "Customer: %+v", &req)
		}
	}

	return web.RespondJson(ctx, w, res.Response(ctx), http.StatusCreated)
}

// Read godoc
// @Summary Update account by ID
// @Description Update updates the specified account in the system.
// @Tags account
// @Accept  json
// @Produce  json
// @Security OAuth2Password
// @Param data body account.UpdateRequest true "Update fields"
// @Success 204
// @Failure 400 {object} weberror.ErrorResponse
// @Failure 403 {object} weberror.ErrorResponse
// @Failure 500 {object} weberror.ErrorResponse
// @Router /accounts [patch]
func (h *Accounts) Update(ctx context.Context, w http.ResponseWriter, r *http.Request, _ map[string]string) error {
	v, err := webcontext.ContextValues(ctx)
	if err != nil {
		return err
	}

	claims, err := auth.ClaimsFromContext(ctx)
	if err != nil {
		return err
	}

	var req account.UpdateRequest
	if err := web.Decode(ctx, r, &req); err != nil {
		if _, ok := errors.Cause(err).(*weberror.Error); !ok {
			err = weberror.NewError(ctx, err, http.StatusBadRequest)
		}
		return web.RespondJsonError(ctx, w, err)
	}

	err = h.Repository.Update(ctx, claims, req, v.Now)
	if err != nil {
		cause := errors.Cause(err)
		switch cause {
		case checklist.ErrForbidden:
			return web.RespondJsonError(ctx, w, weberror.NewError(ctx, err, http.StatusForbidden))
		default:
			_, ok := cause.(validator.ValidationErrors)
			if ok {
				return web.RespondJsonError(ctx, w, weberror.NewError(ctx, err, http.StatusBadRequest))
			}

			return errors.Wrapf(err, "ID: %s Update: %+v", req.ID, req)
		}
	}

	return web.RespondJson(ctx, w, nil, http.StatusNoContent)
}

// Read godoc
// @Summary Archive account by ID
// @Description Archive soft-deletes the specified account from the system.
// @Tags account
// @Accept  json
// @Produce  json
// @Security OAuth2Password
// @Param data body account.ArchiveRequest true "Update fields"
// @Success 204
// @Failure 400 {object} weberror.ErrorResponse
// @Failure 403 {object} weberror.ErrorResponse
// @Failure 500 {object} weberror.ErrorResponse
// @Router /accounts/archive [patch]
func (h *Accounts) Archive(ctx context.Context, w http.ResponseWriter, r *http.Request, _ map[string]string) error {
	v, err := webcontext.ContextValues(ctx)
	if err != nil {
		return err
	}

	claims, err := auth.ClaimsFromContext(ctx)
	if err != nil {
		return err
	}

	var req account.ArchiveRequest
	if err := web.Decode(ctx, r, &req); err != nil {
		if _, ok := errors.Cause(err).(*weberror.Error); !ok {
			err = weberror.NewError(ctx, err, http.StatusBadRequest)
		}
		return web.RespondJsonError(ctx, w, err)
	}

	err = h.Repository.Archive(ctx, claims, req, v.Now)
	if err != nil {
		cause := errors.Cause(err)
		switch cause {
		case checklist.ErrForbidden:
			return web.RespondJsonError(ctx, w, weberror.NewError(ctx, err, http.StatusForbidden))
		default:
			_, ok := cause.(validator.ValidationErrors)
			if ok {
				return web.RespondJsonError(ctx, w, weberror.NewError(ctx, err, http.StatusBadRequest))
			}

			return errors.Wrapf(err, "Id: %s", req.ID)
		}
	}

	return web.RespondJson(ctx, w, nil, http.StatusNoContent)
}
