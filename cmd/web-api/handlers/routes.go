package handlers

import (
	"log"
	"net/http"
	"os"

	"merryworld/surebank/internal/account"
	"merryworld/surebank/internal/checklist"
	"merryworld/surebank/internal/customer"
	"merryworld/surebank/internal/deposit"
	"merryworld/surebank/internal/mid"
	saasSwagger "merryworld/surebank/internal/mid/saas-swagger"
	"merryworld/surebank/internal/platform/auth"
	"merryworld/surebank/internal/platform/web"
	"merryworld/surebank/internal/platform/web/webcontext"
	_ "merryworld/surebank/internal/platform/web/weberror"
	"merryworld/surebank/internal/signup"
	"merryworld/surebank/internal/tenant"
	"merryworld/surebank/internal/tenant/account_preference"
	"merryworld/surebank/internal/user"
	"merryworld/surebank/internal/user_account"
	"merryworld/surebank/internal/user_account/invite"
	"merryworld/surebank/internal/user_auth"

	"github.com/jmoiron/sqlx"
	"gopkg.in/DataDog/dd-trace-go.v1/contrib/go-redis/redis"
)

type AppContext struct {
	Log               *log.Logger
	Env               webcontext.Env
	MasterDB          *sqlx.DB
	Redis             *redis.Client
	UserRepo          *user.Repository
	UserAccountRepo   *user_account.Repository
	RenantRepo        *tenant.Repository
	AccountPrefRepo   *account_preference.Repository
	AuthRepo          *user_auth.Repository
	SignupRepo        *signup.Repository
	InviteRepo        *invite.Repository
	ChecklistRepo     *checklist.Repository
	CustomerRepo      *customer.Repository
	AccountRepo 	  *account.Repository
	DepositRepo		  *deposit.Repository
	Authenticator     *auth.Authenticator
	PreAppMiddleware  []web.Middleware
	PostAppMiddleware []web.Middleware
}

// API returns a handler for a set of routes.
func API(shutdown chan os.Signal, appCtx *AppContext) http.Handler {

	// Include the pre middlewares first.
	middlewares := appCtx.PreAppMiddleware

	// Define app middlewares applied to all requests.
	middlewares = append(middlewares,
		mid.Trace(),
		mid.Logger(appCtx.Log),
		mid.Errors(appCtx.Log, nil),
		mid.Metrics(),
		mid.Panics())

	// Append any global middlewares that should be included after the app middlewares.
	if len(appCtx.PostAppMiddleware) > 0 {
		middlewares = append(middlewares, appCtx.PostAppMiddleware...)
	}

	// Construct the web.App which holds all routes as well as common Middleware.
	app := web.NewApp(shutdown, appCtx.Log, appCtx.Env, middlewares...)

	// Register health check endpoint. This route is not authenticated.
	check := Check{
		MasterDB: appCtx.MasterDB,
		Redis:    appCtx.Redis,
	}
	app.Handle("GET", "/v1/health", check.Health)
	app.Handle("GET", "/ping", check.Ping)

	// Register example endpoints.
	ex := Example{
		Checklist: appCtx.ChecklistRepo,
	}
	app.Handle("GET", "/v1/examples/error-response", ex.ErrorResponse)

	// Register user management and authentication endpoints.
	u := Users{
		UserRepo: appCtx.UserRepo,
		AuthRepo: appCtx.AuthRepo,
	}
	app.Handle("GET", "/v1/users", u.Find, mid.AuthenticateHeader(appCtx.Authenticator))
	app.Handle("POST", "/v1/users", u.Create, mid.AuthenticateHeader(appCtx.Authenticator), mid.HasRole(auth.RoleAdmin))
	app.Handle("GET", "/v1/users/:id", u.Read, mid.AuthenticateHeader(appCtx.Authenticator))
	app.Handle("PATCH", "/v1/users", u.Update, mid.AuthenticateHeader(appCtx.Authenticator))
	app.Handle("PATCH", "/v1/users/password", u.UpdatePassword, mid.AuthenticateHeader(appCtx.Authenticator))
	app.Handle("PATCH", "/v1/users/archive", u.Archive, mid.AuthenticateHeader(appCtx.Authenticator), mid.HasRole(auth.RoleAdmin))
	app.Handle("DELETE", "/v1/users/:id", u.Delete, mid.AuthenticateHeader(appCtx.Authenticator), mid.HasRole(auth.RoleAdmin))
	app.Handle("PATCH", "/v1/users/switch-account/:account_id", u.SwitchAccount, mid.AuthenticateHeader(appCtx.Authenticator))

	// This route is not authenticated
	app.Handle("POST", "/v1/oauth/token", u.Token)

	// Register user account management endpoints.
	ua := UserAccount{
		Repository: appCtx.UserAccountRepo,
	}
	app.Handle("GET", "/v1/user_accounts", ua.Find, mid.AuthenticateHeader(appCtx.Authenticator))
	app.Handle("POST", "/v1/user_accounts", ua.Create, mid.AuthenticateHeader(appCtx.Authenticator), mid.HasRole(auth.RoleAdmin))
	app.Handle("GET", "/v1/user_accounts/:user_id/:account_id", ua.Read, mid.AuthenticateHeader(appCtx.Authenticator))
	app.Handle("PATCH", "/v1/user_accounts", ua.Update, mid.AuthenticateHeader(appCtx.Authenticator))
	app.Handle("PATCH", "/v1/user_accounts/archive", ua.Archive, mid.AuthenticateHeader(appCtx.Authenticator), mid.HasRole(auth.RoleAdmin))
	app.Handle("DELETE", "/v1/user_accounts", ua.Delete, mid.AuthenticateHeader(appCtx.Authenticator), mid.HasRole(auth.RoleAdmin))

	// Register account endpoints.
	a := Tenants{
		Repository: appCtx.RenantRepo,
	}
	app.Handle("GET", "/v1/tenants/:id", a.Read, mid.AuthenticateHeader(appCtx.Authenticator))
	app.Handle("PATCH", "/v1/tenants", a.Update, mid.AuthenticateHeader(appCtx.Authenticator), mid.HasRole(auth.RoleAdmin))

	// Register customer.
	cus := Customers{
		Repository: appCtx.CustomerRepo,
	}
	app.Handle("GET", "/v1/customers", cus.Find, mid.AuthenticateHeader(appCtx.Authenticator))
	app.Handle("POST", "/v1/customers", cus.Create, mid.AuthenticateHeader(appCtx.Authenticator), mid.HasRole(auth.RoleAdmin))
	app.Handle("GET", "/v1/customers/:id", cus.Read, mid.AuthenticateHeader(appCtx.Authenticator))
	app.Handle("PATCH", "/v1/customers", cus.Update, mid.AuthenticateHeader(appCtx.Authenticator), mid.HasRole(auth.RoleAdmin))
	app.Handle("PATCH", "/v1/customers/archive", cus.Archive, mid.AuthenticateHeader(appCtx.Authenticator), mid.HasRole(auth.RoleAdmin))
	app.Handle("DELETE", "/v1/customers/:id", cus.Delete, mid.AuthenticateHeader(appCtx.Authenticator), mid.HasRole(auth.RoleAdmin))

	// Register customer.
	acc := Accounts{
		Repository: appCtx.AccountRepo,
	}
	app.Handle("GET", "/v1/accounts", acc.Find, mid.AuthenticateHeader(appCtx.Authenticator))
	app.Handle("POST", "/v1/accounts", acc.Create, mid.AuthenticateHeader(appCtx.Authenticator), mid.HasRole(auth.RoleAdmin))
	app.Handle("GET", "/v1/accounts/:id", acc.Read, mid.AuthenticateHeader(appCtx.Authenticator))
	app.Handle("PATCH", "/v1/accounts", acc.Update, mid.AuthenticateHeader(appCtx.Authenticator), mid.HasRole(auth.RoleAdmin))
	app.Handle("PATCH", "/v1/accounts/archive", acc.Archive, mid.AuthenticateHeader(appCtx.Authenticator), mid.HasRole(auth.RoleAdmin))

	// Register deposit.
	dep := Deposits{
		Repository: appCtx.DepositRepo,
	}
	app.Handle("GET", "/v1/deposits", dep.Find, mid.AuthenticateHeader(appCtx.Authenticator))
	app.Handle("POST", "/v1/deposits", dep.Create, mid.AuthenticateHeader(appCtx.Authenticator), mid.HasRole(auth.RoleAdmin))
	app.Handle("GET", "/v1/deposits/:id", dep.Read, mid.AuthenticateHeader(appCtx.Authenticator))
	app.Handle("PATCH", "/v1/deposits", dep.Update, mid.AuthenticateHeader(appCtx.Authenticator), mid.HasRole(auth.RoleAdmin))
	app.Handle("PATCH", "/v1/deposits/archive", dep.Archive, mid.AuthenticateHeader(appCtx.Authenticator), mid.HasRole(auth.RoleAdmin))

	// Register swagger documentation.
	// TODO: Add authentication. Current authenticator requires an Authorization header
	// 		 which breaks the browser experience.
	app.Handle("GET", "/docs/", saasSwagger.WrapHandler)
	app.Handle("GET", "/docs/*", saasSwagger.WrapHandler)

	return app
}

// Types godoc
// @Summary List of types.
// @Param data body weberror.FieldError false "Field Error"
// @Param data body web.TimeResponse false "Time Response"
// @Param data body web.EnumResponse false "Enum Response"
// @Param data body web.EnumMultiResponse false "Enum Multi Response"
// @Param data body web.EnumOption false "Enum Option"
// @Param data body signup.SignupAccount false "SignupAccount"
// @Param data body signup.SignupUser false "SignupUser"
// To support nested types not parsed by swag.
func Types() {}
