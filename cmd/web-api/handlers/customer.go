package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"merryworld/surebank/internal/account"
	"merryworld/surebank/internal/customer"
	"merryworld/surebank/internal/platform/auth"
	"merryworld/surebank/internal/platform/notify"
	"merryworld/surebank/internal/platform/web"
	"merryworld/surebank/internal/platform/web/webcontext"
	"merryworld/surebank/internal/platform/web/weberror"

	"github.com/pkg/errors"
	"gopkg.in/go-playground/validator.v9"
)

// Customers represents the Customer API method handler set.
type Customers struct {
	Repository  *customer.Repository
	AccountRepo *account.Repository
	notifySMS notify.SMS

	// ADD OTHER STATE LIKE THE LOGGER IF NEEDED.
}

// Find godoc
// @Summary List customers
// @Description Find returns the existing customers in the system.
// @Tags customer
// @Accept  json
// @Produce  json
// @Security OAuth2Password
// @Param where				query string 	false	"Filter string, example: name = 'Oluwafe Dami'"
// @Param order				query string   	false 	"Order columns separated by comma, example: created_at desc"
// @Param limit				query integer  	false 	"Limit, example: 10"
// @Param offset			query integer  	false 	"Offset, example: 20"
// @Param include-archived query boolean 	false 	"Included Archived, example: false"
// @Success 200 {object} customer.PagedResponseList
// @Failure 400 {object} weberror.ErrorResponse
// @Failure 403 {object} weberror.ErrorResponse
// @Failure 500 {object} weberror.ErrorResponse
// @Router /customers [get]
func (h *Customers) Find(ctx context.Context, w http.ResponseWriter, r *http.Request, _ map[string]string) error {
	claims, ok := ctx.Value(auth.Key).(auth.Claims)
	if !ok {
		return errors.New("claims missing from context")
	}

	var req customer.FindRequest

	// Handle where query value if set.
	if v := r.URL.Query().Get("where"); v != "" {
		where, args, err := web.ExtractWhereArgs(v)
		if err != nil {
			return web.RespondJsonError(ctx, w, weberror.NewError(ctx, err, http.StatusBadRequest))
		}
		req.Where = where
		req.Args = args
	}

	// Handle order query value if set.
	req.Order = []string{"created_at desc"}
	if v := r.URL.Query().Get("order"); v != "" {
		for _, o := range strings.Split(v, ",") {
			o = strings.TrimSpace(o)
			if o != "" {
				req.Order = append(req.Order, o)
			}
		}
	}

	// Handle limit query value if set.
	var l uint = 20
	req.Limit = &l
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

	// Handle include-account query value if set.
	if v := r.URL.Query().Get("include-account"); v != "" {
		b, err := strconv.ParseBool(v)
		if err != nil {
			err = errors.WithMessagef(err, "unable to parse %s as boolean for include-account param", v)
			return web.RespondJsonError(ctx, w, weberror.NewError(ctx, err, http.StatusBadRequest))
		}
		req.IncludeArchived = b
	}

	req.IncludeAccountNo = true

	res, err := h.Repository.Find(ctx, claims, req)
	if err != nil {
		return err
	}

	return web.RespondJson(ctx, w, res, http.StatusOK)
}

// Read godoc
// @Summary Get customer by ID.
// @Description Read returns the specified customer from the system.
// @Tags customer
// @Accept  json
// @Produce  json
// @Security OAuth2Password
// @Param id path string true "Customer ID"
// @Success 200 {object} customer.Response
// @Failure 400 {object} weberror.ErrorResponse
// @Failure 404 {object} weberror.ErrorResponse
// @Failure 500 {object} weberror.ErrorResponse
// @Router /customers/{id} [get]
func (h *Customers) Read(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {
	claims, ok := ctx.Value(auth.Key).(auth.Claims)
	if !ok {
		return errors.New("claims missing from context")
	}

	res, err := h.Repository.ReadByID(ctx, claims, params["id"])
	if err != nil {
		cause := errors.Cause(err)
		switch cause {
		case customer.ErrNotFound:
			return web.RespondJsonError(ctx, w, weberror.NewError(ctx, err, http.StatusNotFound))
		default:
			return errors.Wrapf(err, "ID: %s", params["id"])
		}
	}

	return web.RespondJson(ctx, w, res.Response(ctx), http.StatusOK)
}

// Create godoc
// @Summary Create new customer.
// @Description Create inserts a new customer into the system.
// @Tags customer
// @Accept  json
// @Produce  json
// @Security OAuth2Password
// @Param data body customer.CreateRequest true "Customer details"
// @Success 201 {object} account.Response
// @Failure 400 {object} weberror.ErrorResponse
// @Failure 403 {object} weberror.ErrorResponse
// @Failure 404 {object} weberror.ErrorResponse
// @Failure 500 {object} weberror.ErrorResponse
// @Router /customers [post]
func (h *Customers) Create(ctx context.Context, w http.ResponseWriter, r *http.Request, _ map[string]string) error {
	v, err := webcontext.ContextValues(ctx)
	if err != nil {
		return err
	}

	claims, err := auth.ClaimsFromContext(ctx)
	if err != nil {
		return err
	}

	var req customer.CreateRequest
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
		case customer.ErrForbidden:
			return web.RespondJsonError(ctx, w, weberror.NewError(ctx, err, http.StatusForbidden))
		default:
			_, ok := cause.(validator.ValidationErrors)
			if ok {
				return web.RespondJsonError(ctx, w, weberror.NewError(ctx, err, http.StatusBadRequest))
			}
			return errors.Wrapf(err, "Customer: %+v", &req)
		}
	}

	accReq := account.CreateRequest{
		CustomerID: res.ID,
		Type:       req.Type,
		Target:     req.Target,
		TargetInfo: req.TargetInfo,
		BranchID:   req.BranchID,
	}
	accRes, err := h.AccountRepo.Create(ctx, claims, accReq, v.Now)
	if err != nil {
		// delete the created customer account
		_ = h.Repository.Archive(ctx, claims, customer.ArchiveRequest{ID: res.ID}) // TODO: log delete error for debug
		cause := errors.Cause(err)
		switch cause {
		case customer.ErrForbidden:
			return web.RespondJsonError(ctx, w, weberror.NewError(ctx, err, http.StatusForbidden))
		default:
			_, ok := cause.(validator.ValidationErrors)
			if ok {
				return web.RespondJsonError(ctx, w, weberror.NewError(ctx, err, http.StatusBadRequest))
			}
			return errors.Wrapf(err, "Customer: %+v", &req)
		}
	}

	if err = h.notifySMS.Send(ctx, req.PhoneNumber, "sms/welcome_message",
		map[string]interface{}{
			"Name":          req.Name,
			"AccountNumber": accRes.Number,
			"Target":        req.Target,
		}); err != nil {
		// TODO: log critical error. Send message to monitoring account
		fmt.Println(err)
	}

	result := accRes.Response(ctx)
	result.Customer = res.Response(ctx)
	return web.RespondJson(ctx, w, result, http.StatusCreated)
}

// Read godoc
// @Summary Update customer by ID
// @Description Update updates the specified customer in the system.
// @Tags customer
// @Accept  json
// @Produce  json
// @Security OAuth2Password
// @Param data body customer.UpdateRequest true "Update fields"
// @Success 204
// @Failure 400 {object} weberror.ErrorResponse
// @Failure 403 {object} weberror.ErrorResponse
// @Failure 500 {object} weberror.ErrorResponse
// @Router /customers [patch]
func (h *Customers) Update(ctx context.Context, w http.ResponseWriter, r *http.Request, _ map[string]string) error {
	v, err := webcontext.ContextValues(ctx)
	if err != nil {
		return err
	}

	claims, err := auth.ClaimsFromContext(ctx)
	if err != nil {
		return err
	}

	var req customer.UpdateRequest
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
		case customer.ErrForbidden:
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
// @Summary Archive customer by ID
// @Description Archive soft-deletes the specified customer from the system.
// @Tags customer
// @Accept  json
// @Produce  json
// @Security OAuth2Password
// @Param data body customer.ArchiveRequest true "Update fields"
// @Success 204
// @Failure 400 {object} weberror.ErrorResponse
// @Failure 403 {object} weberror.ErrorResponse
// @Failure 500 {object} weberror.ErrorResponse
// @Router /customers/archive [patch]
func (h *Customers) Archive(ctx context.Context, w http.ResponseWriter, r *http.Request, _ map[string]string) error {

	claims, err := auth.ClaimsFromContext(ctx)
	if err != nil {
		return err 
	}

	var req customer.ArchiveRequest
	if err := web.Decode(ctx, r, &req); err != nil {
		if _, ok := errors.Cause(err).(*weberror.Error); !ok {
			err = weberror.NewError(ctx, err, http.StatusBadRequest)
		}
		return web.RespondJsonError(ctx, w, err)
	}

	err = h.Repository.Archive(ctx, claims, req)
	if err != nil {
		cause := errors.Cause(err)
		switch cause {
		case customer.ErrForbidden:
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

// Delete godoc
// @Summary Delete customer by ID
// @Description Delete removes the specified customer from the system.
// @Tags customer
// @Accept  json
// @Produce  json
// @Security OAuth2Password
// @Param id path string true "Customer ID"
// @Success 204
// @Failure 400 {object} weberror.ErrorResponse
// @Failure 403 {object} weberror.ErrorResponse
// @Failure 500 {object} weberror.ErrorResponse
// @Router /customers/{id} [delete]
func (h *Customers) Delete(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {
	claims, err := auth.ClaimsFromContext(ctx)
	if err != nil {
		return err
	}

	err = h.Repository.Archive(ctx, claims,
		customer.ArchiveRequest{ID: params["id"]})
	if err != nil {
		cause := errors.Cause(err)
		switch cause {
		case customer.ErrForbidden:
			return web.RespondJsonError(ctx, w, weberror.NewError(ctx, err, http.StatusForbidden))
		default:
			_, ok := cause.(validator.ValidationErrors)
			if ok {
				return web.RespondJsonError(ctx, w, weberror.NewError(ctx, err, http.StatusBadRequest))
			}

			return errors.Wrapf(err, "Id: %s", params["id"])
		}
	}

	return web.RespondJson(ctx, w, nil, http.StatusNoContent)
}
