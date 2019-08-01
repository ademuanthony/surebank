package handlers

import (
	"context"
	"net/http"
	"time"

	"geeks-accelerator/oss/saas-starter-kit/internal/account"
	"geeks-accelerator/oss/saas-starter-kit/internal/platform/auth"
	"geeks-accelerator/oss/saas-starter-kit/internal/platform/web"
	"geeks-accelerator/oss/saas-starter-kit/internal/platform/web/webcontext"
	"geeks-accelerator/oss/saas-starter-kit/internal/platform/web/weberror"
	"geeks-accelerator/oss/saas-starter-kit/internal/signup"
	"geeks-accelerator/oss/saas-starter-kit/internal/user"
	"github.com/gorilla/schema"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// Signup represents the Signup API method handler set.
type Signup struct {
	MasterDB      *sqlx.DB
	Renderer      web.Renderer
	Authenticator *auth.Authenticator
}

// Step1 handles collecting the first detailed needed to create a new account.
func (h *Signup) Step1(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {

	ctxValues, err := webcontext.ContextValues(ctx)
	if err != nil {
		return err
	}

	//
	req := new(signup.SignupRequest)
	data := make(map[string]interface{})
	f := func() error {
		claims, _ := auth.ClaimsFromContext(ctx)

		if r.Method == http.MethodPost {

			err := r.ParseForm()
			if err != nil {
				return err
			}

			decoder := schema.NewDecoder()
			if err := decoder.Decode(req, r.PostForm); err != nil {
				return err
			}

			// Execute the account / user signup.
			_, err = signup.Signup(ctx, claims, h.MasterDB, *req, ctxValues.Now)
			if err != nil {
				switch errors.Cause(err) {
				case account.ErrForbidden:
					return web.RespondError(ctx, w, weberror.NewError(ctx, err, http.StatusForbidden))
				default:
					if verr, ok := weberror.NewValidationError(ctx, err); ok {
						data["validationErrors"] = verr.(*weberror.Error)
						return nil
					} else {
						return err
					}
				}
			}

			// Authenticated the new user.
			token, err := user.Authenticate(ctx, h.MasterDB, h.Authenticator, req.User.Email, req.User.Password, time.Hour, ctxValues.Now)
			if err != nil {
				return err
			}

			// Add the token to the users session.
			err = handleSessionToken(ctx, w, r, token)
			if err != nil {
				return err
			}

			// Redirect the user to the dashboard.
			http.Redirect(w, r, "/", http.StatusFound)
		}

		return nil
	}

	if err := f(); err != nil {
		return web.RenderError(ctx, w, r, err, h.Renderer, tmplLayoutBase, tmplContentErrorGeneric, web.MIMETextHTMLCharsetUTF8)
	}

	data["form"] = req

	if verr, ok := weberror.NewValidationError(ctx, webcontext.Validator().Struct(signup.SignupRequest{})); ok {
		data["validationDefaults"] = verr.(*weberror.Error)
	}

	return h.Renderer.Render(ctx, w, r, tmplLayoutBase, "signup-step1.tmpl", web.MIMETextHTMLCharsetUTF8, http.StatusOK, data)
}