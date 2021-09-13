package handlers

import (
	"context"
	"fmt"
	"log"
	"merryworld/surebank/internal/account"
	"merryworld/surebank/internal/customer"
	"merryworld/surebank/internal/dscommission"
	"merryworld/surebank/internal/expenditure"
	"merryworld/surebank/internal/inventory"
	"merryworld/surebank/internal/profit"
	"merryworld/surebank/internal/sale"
	"merryworld/surebank/internal/transaction"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"merryworld/surebank/internal/branch"
	"merryworld/surebank/internal/checklist"
	"merryworld/surebank/internal/geonames"
	"merryworld/surebank/internal/mid"
	"merryworld/surebank/internal/platform/auth"
	"merryworld/surebank/internal/platform/notify"
	"merryworld/surebank/internal/platform/web"
	"merryworld/surebank/internal/platform/web/webcontext"
	"merryworld/surebank/internal/platform/web/weberror"
	"merryworld/surebank/internal/shop"
	"merryworld/surebank/internal/signup"
	"merryworld/surebank/internal/tenant"
	"merryworld/surebank/internal/tenant/account_preference"
	"merryworld/surebank/internal/user"
	"merryworld/surebank/internal/user_account"
	"merryworld/surebank/internal/user_account/invite"
	"merryworld/surebank/internal/user_auth"
	"merryworld/surebank/internal/webroute"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/ikeikeikeike/go-sitemap-generator/v2/stm"
	"github.com/jmoiron/sqlx"
	"gopkg.in/DataDog/dd-trace-go.v1/contrib/go-redis/redis"
)

const (
	TmplLayoutBase          = "base.gohtml"
	tmplLayoutSite          = "site.gohtml"
	TmplContentErrorGeneric = "error-generic.gohtml"
)

type AppContext struct {
	Log               *log.Logger
	Env               webcontext.Env
	MasterDB          *sqlx.DB
	MasterDbHost      string
	Redis             *redis.Client
	UserRepo          *user.Repository
	UserAccountRepo   *user_account.Repository
	TenantRepo        *tenant.Repository
	AccountPrefRepo   *account_preference.Repository
	AuthRepo          *user_auth.Repository
	SignupRepo        *signup.Repository
	InviteRepo        *invite.Repository
	ChecklistRepo     *checklist.Repository
	GeoRepo           *geonames.Repository
	ProfitRepo        *profit.Repository
	ShopRepo          *shop.Repository
	InventoryRepo     *inventory.Repository
	BranchRepo        *branch.Repository
	CustomerRepo      *customer.Repository
	AccountRepo       *account.Repository
	CommissionRepo    *dscommission.Repository
	TransactionRepo   *transaction.Repository
	SaleRepo          *sale.Repository
	ExpendituresRepo  *expenditure.Repository
	NotifySMS         notify.SMS
	Authenticator     *auth.Authenticator
	StaticDir         string
	TemplateDir       string
	Renderer          web.Renderer
	WebRoute          webroute.WebRoute
	PreAppMiddleware  []web.Middleware
	PostAppMiddleware []web.Middleware
	AwsSession        *session.Session
}

// API returns a handler for a set of routes.
func APP(shutdown chan os.Signal, appCtx *AppContext, reopenDBFunc func() error) http.Handler {

	// Include the pre middlewares first.
	middlewares := appCtx.PreAppMiddleware

	// Define app middlewares applied to all requests.
	middlewares = append(middlewares,
		mid.Trace(),
		mid.Logger(appCtx.Log),
		mid.Errors(appCtx.Log, appCtx.Renderer),
		mid.Metrics(),
		mid.Panics())

	// Append any global middlewares that should be included after the app middlewares.
	if len(appCtx.PostAppMiddleware) > 0 {
		middlewares = append(middlewares, appCtx.PostAppMiddleware...)
	}

	// Construct the web.App which holds all routes as well as common Middleware.
	app := web.NewApp(shutdown, appCtx.Log, appCtx.Env, middlewares...)

	// Register serverless endpoint. This route is not authenticated.
	serverless := Serverless{
		MasterDB:     appCtx.MasterDB,
		MasterDbHost: appCtx.MasterDbHost,
		AwsSession:   appCtx.AwsSession,
		Renderer:     appCtx.Renderer,
	}
	app.Handle("GET", "/serverless/pending", serverless.Pending)

	// waitDbMid ensures the database is active before allowing the user to access the requested URI.
	waitDbMid := mid.WaitForDbResumed(mid.WaitForDbResumedConfig{
		// Database handle to be used to ensure its online.
		DB: appCtx.MasterDB,

		// WaitHandler defines the handler to render for the user to when the database is being resumed.
		WaitHandler:  serverless.Pending,
		ReopenDBFunc: reopenDBFunc,
	})

	// Build a sitemap.
	sm := stm.NewSitemap(1)
	sm.SetVerbose(false)
	sm.SetDefaultHost(appCtx.WebRoute.WebAppUrl(""))
	sm.Create()

	smLocAddModified := func(loc stm.URL, filename string) {
		contentPath := filepath.Join(appCtx.TemplateDir, "content", filename)

		file, err := os.Stat(contentPath)
		if err != nil {
			log.Fatalf("main : Add sitemap file modified for %s: %+v", filename, err)
		}

		lm := []interface{}{"lastmod", file.ModTime().Format(time.RFC3339)}
		loc = append(loc, lm)
		sm.Add(loc)
	}

	// Register checklist management pages.
	p := Checklists{
		ChecklistRepo: appCtx.ChecklistRepo,
		Redis:         appCtx.Redis,
		Renderer:      appCtx.Renderer,
	}
	app.Handle("POST", "/checklists/:checklist_id/update", p.Update, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasRole(auth.RoleAdmin))
	app.Handle("GET", "/checklists/:checklist_id/update", p.Update, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasRole(auth.RoleAdmin))
	app.Handle("POST", "/checklists/:checklist_id", p.View, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasRole(auth.RoleAdmin))
	app.Handle("GET", "/checklists/:checklist_id", p.View, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasAuth())
	app.Handle("POST", "/checklists/create", p.Create, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasRole(auth.RoleAdmin))
	app.Handle("GET", "/checklists/create", p.Create, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasRole(auth.RoleAdmin))
	app.Handle("GET", "/checklists", p.Index, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasAuth())

	// Brands
	branches := Branches{
		Repo:     appCtx.BranchRepo,
		Redis:    appCtx.Redis,
		Renderer: appCtx.Renderer,
	}
	app.Handle("POST", "/branches/:branch_id/update", branches.Update, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasRole(auth.RoleSuperAdmin))
	app.Handle("GET", "/branches/:branch_id/update", branches.Update, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasRole(auth.RoleSuperAdmin))
	app.Handle("POST", "/branches/:branch_id", branches.View, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasAuth(), mid.HasRole(auth.RoleSuperAdmin))
	app.Handle("GET", "/branches/:branch_id", branches.View, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasAuth(), mid.HasRole(auth.RoleSuperAdmin))
	app.Handle("POST", "/branches/create", branches.Create, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasRole(auth.RoleSuperAdmin))
	app.Handle("GET", "/branches/create", branches.Create, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasRole(auth.RoleSuperAdmin))
	app.Handle("GET", "/branches", branches.Index, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasAuth(), mid.HasRole(auth.RoleSuperAdmin))
	app.Handle("POST", "/api/v1/branches", branches.APICreate, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasRole(auth.RoleSuperAdmin))

	// Accounting
	accounting := Accounting{
		DbConn:    appCtx.MasterDB.DB,
		UserRepos: appCtx.UserRepo,
		Redis:     appCtx.Redis,
		Renderer:  appCtx.Renderer,
	}
	app.Handle("GET", "/accounting/banks", accounting.BankAccounts, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasAuth())
	app.Handle("POST", "/api/v1/accounting/banks", accounting.CreateBankAccount, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasRole(auth.RoleAdmin))
	app.Handle("GET", "/accounting/deposits", accounting.BankDeposits, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasAuth())
	app.Handle("POST", "/api/v1/accounting/deposits", accounting.CreateBankDeposit, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasRole(auth.RoleAdmin))
	app.Handle("GET", "/accounting/expenditures", accounting.Expenditures, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasAuth())
	app.Handle("POST", "/api/v1/accounting/expenditures", accounting.CreateExpenditure, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasRole(auth.RoleAdmin))
	app.Handle("GET", "/accounting/resp-summaries", accounting.RepsSummaries, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasAuth())
	app.Handle("GET", "/accounting", accounting.DailySummaries, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasAuth(), mid.HasRole(auth.RoleSuperAdmin))

	// /accounting/reps-expenditure
	repsExpenditure := Expenditures{
		ExpendituresRepo: appCtx.ExpendituresRepo,
		UserRepos:        appCtx.UserRepo,
		Redis:            appCtx.Redis,
		Renderer:         appCtx.Renderer,
	}
	app.Handle("GET", "/accounting/reps-expenditures", repsExpenditure.Index, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasAuth())
	app.Handle("DELETE", "/api/v1/accounting/reps-expenditures/:id", repsExpenditure.Delete, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasRole(auth.RoleAdmin))
	app.Handle("POST", "/api/v1/accounting/reps-expenditures", repsExpenditure.Create, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasRole(auth.RoleAdmin))

	// Register shop management pages
	// Brands
	brands := Brands{
		ShopRepo: appCtx.ShopRepo,
		Redis:    appCtx.Redis,
		Renderer: appCtx.Renderer,
	}
	app.Handle("POST", "/shop/brands/:brand_id/update", brands.Update, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasRole(auth.RoleAdmin))
	app.Handle("GET", "/shop/brands/:brand_id/update", brands.Update, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasRole(auth.RoleAdmin))
	app.Handle("POST", "/shop/brands/:brand_id", brands.View, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasAuth())
	app.Handle("GET", "/shop/brands/:brand_id", brands.View, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasAuth())
	app.Handle("POST", "/shop/brands/create", brands.Create, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasRole(auth.RoleAdmin))
	app.Handle("GET", "/shop/brands/create", brands.Create, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasRole(auth.RoleAdmin))
	app.Handle("GET", "/shop/brands", brands.Index, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasAuth())

	// Category
	cat := Categories{
		ShopRepo: appCtx.ShopRepo,
		Redis:    appCtx.Redis,
		Renderer: appCtx.Renderer,
	}
	app.Handle("POST", "/shop/categories/:category_id/update", cat.Update, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasRole(auth.RoleAdmin))
	app.Handle("GET", "/shop/categories/:category_id/update", cat.Update, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasRole(auth.RoleAdmin))
	app.Handle("POST", "/shop/categories/:category_id", cat.View, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasAuth())
	app.Handle("GET", "/shop/categories/:category_id", cat.View, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasAuth())
	app.Handle("POST", "/shop/categories/create", cat.Create, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasRole(auth.RoleAdmin))
	app.Handle("GET", "/shop/categories/create", cat.Create, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasRole(auth.RoleAdmin))
	app.Handle("GET", "/shop/categories", cat.Index, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasAuth())

	// Products
	prod := Products{
		ShopRepo: appCtx.ShopRepo,
		Redis:    appCtx.Redis,
		Renderer: appCtx.Renderer,
	}
	app.Handle("POST", "/shop/products/:product_id/update", prod.Update, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasRole(auth.RoleAdmin))
	app.Handle("GET", "/shop/products/:product_id/update", prod.Update, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasRole(auth.RoleAdmin))
	app.Handle("POST", "/shop/products/:product_id", prod.View, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasAuth())
	app.Handle("GET", "/shop/products/:product_id", prod.View, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasAuth())
	app.Handle("POST", "/shop/products/create", prod.Create, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasRole(auth.RoleAdmin))
	app.Handle("GET", "/shop/products/create", prod.Create, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasRole(auth.RoleAdmin))
	app.Handle("GET", "/shop/products", prod.Index, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasAuth())

	// Products
	prof := Profits{
		ProfRepo: appCtx.ProfitRepo,
		Redis:    appCtx.Redis,
		Renderer: appCtx.Renderer,
	}
	app.Handle("POST", "/profits/:profit_id/update", prof.Update, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasRole(auth.RoleAdmin))
	app.Handle("GET", "/profits/:profit_id/update", prof.Update, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasRole(auth.RoleAdmin))
	app.Handle("POST", "/profits/:profit_id", prof.View, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasAuth())
	app.Handle("GET", "/profits/:profit_id", prof.View, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasAuth())
	app.Handle("POST", "/profits/create", prof.Create, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasRole(auth.RoleAdmin))
	app.Handle("GET", "/profits/create", prof.Create, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasRole(auth.RoleAdmin))
	app.Handle("GET", "/profits", prof.Index, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasAuth())

	// Stocks
	stock := Stocks{
		Repo:       appCtx.InventoryRepo,
		ShopRepo:   appCtx.ShopRepo,
		BranchRepo: appCtx.BranchRepo,
		Redis:      appCtx.Redis,
		Renderer:   appCtx.Renderer,
	}
	// app.Handle("POST", "/shop/inventory/:stock_id/update", stock.Update, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasRole(auth.RoleAdmin))
	// app.Handle("GET", "/shop/inventory/:stock_id/update", stock.Update, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasRole(auth.RoleAdmin))
	app.Handle("POST", "/shop/inventory/:stock_id", stock.View, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasAuth())
	app.Handle("GET", "/shop/inventory/:stock_id", stock.View, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasAuth())
	app.Handle("POST", "/shop/inventory/create", stock.Create, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasRole(auth.RoleAdmin))
	app.Handle("GET", "/shop/inventory/create", stock.Create, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasRole(auth.RoleAdmin))
	app.Handle("POST", "/shop/inventory/remove", stock.Remove, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasRole(auth.RoleAdmin))
	app.Handle("GET", "/shop/inventory/remove", stock.Remove, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasRole(auth.RoleAdmin))
	app.Handle("GET", "/shop/inventory/report", stock.Report, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasAuth())
	app.Handle("GET", "/shop/inventory", stock.Index, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasAuth())

	// Customers
	custs := Customers{
		CustomerRepo:    appCtx.CustomerRepo,
		AccountRepo:     appCtx.AccountRepo,
		NotifySMS:       appCtx.NotifySMS,
		TransactionRepo: appCtx.TransactionRepo,
		Redis:           appCtx.Redis,
		Renderer:        appCtx.Renderer,
	}
	app.Handle("POST", "/customers/:customer_id/update", custs.Update, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasRole(auth.RoleAdmin))
	app.Handle("GET", "/customers/:customer_id/update", custs.Update, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasRole(auth.RoleAdmin))
	app.Handle("GET", "/customers/:customer_id/add-account", custs.AddAccount, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasAuth())
	app.Handle("POST", "/customers/:customer_id/add-account", custs.AddAccount, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasAuth())

	app.Handle("GET", "/customers/:customer_id/accounts/:account_id/transactions/deposit", custs.Deposit, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasAuth())
	app.Handle("POST", "/customers/:customer_id/accounts/:account_id/transactions/deposit", custs.Deposit, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasAuth())
	app.Handle("GET", "/deposit", custs.DirectDeposit, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasAuth())
	app.Handle("POST", "/deposit", custs.DirectDeposit, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasAuth())
	app.Handle("GET", "/api/v1/customers/account-name", custs.AccountName, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasAuth())
	app.Handle("GET", "/api/db-stat", custs.DBStat)
	app.Handle("GET", "/customers/:customer_id/accounts/:account_id/transactions/withdraw",
		custs.Withraw, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasAuth())
	app.Handle("POST", "/customers/:customer_id/accounts/:account_id/transactions/withdraw",
		custs.Withraw, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasAuth())
	app.Handle("GET", "/customers/:customer_id/accounts/:account_id/transactions", custs.AccountTransactions, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasAuth())
	app.Handle("POST", "/customers/:customer_id/accounts/:account_id/transactions/:transaction_id", custs.Transaction, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasAuth())
	app.Handle("GET", "/customers/:customer_id/accounts/:account_id/transactions/:transaction_id", custs.Transaction, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasAuth())
	app.Handle("GET", "/customers/:customer_id/accounts/:account_id", custs.Account, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasAuth())
	app.Handle("POST", "/customers/:customer_id/accounts/:account_id/update", custs.UpdateAccount, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasAuth())
	app.Handle("GET", "/customers/:customer_id/accounts/:account_id/update", custs.UpdateAccount, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasAuth())
	app.Handle("GET", "/customers/:customer_id/transactions", custs.Transactions, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasAuth())
	app.Handle("POST", "/customers/:customer_id", custs.View, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasAuth())
	app.Handle("GET", "/customers/:customer_id", custs.View, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasAuth())
	app.Handle("POST", "/customers/create", custs.Create, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasAuth())
	app.Handle("GET", "/customers/create", custs.Create, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasAuth())
	app.Handle("GET", "/customers", custs.Index, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasAuth())

	// Customers
	sms := BulkSMS{
		CustomerRepo: appCtx.CustomerRepo,
		AccountRepo:  appCtx.AccountRepo,
		NotifySMS:    appCtx.NotifySMS,
		Redis:        appCtx.Redis,
		Renderer:     appCtx.Renderer,
		DbConn:       appCtx.MasterDB.DB,
	}
	app.Handle("POST", "/sms", sms.Index, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasRole(auth.RoleAdmin))
	app.Handle("GET", "/sms", sms.Index, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasRole(auth.RoleAdmin))

	reports := Reports{
		CustomerRepo:    appCtx.CustomerRepo,
		AccountRepo:     appCtx.AccountRepo,
		CommissionRepo:  appCtx.CommissionRepo,
		TransactionRepo: appCtx.TransactionRepo,
		ShopRepo:        appCtx.ShopRepo,
		UserRepos:       appCtx.UserRepo,
		Renderer:        appCtx.Renderer,
		Redis:           appCtx.Redis,
	}
	app.Handle("GET", "/reports/withdrawals", reports.Withdrawals, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasAuth())
	app.Handle("GET", "/reports/collections", reports.Transactions, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasAuth())
	app.Handle("GET", "/reports/ds/commissions", reports.DsCommissions, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasAuth())
	app.Handle("GET", "/reports/ds", reports.Ds, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasAuth())
	app.Handle("GET", "/reports/debtors", reports.Debtors, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasAuth())

	// Register sales endpoint
	sales := Sales{
		Repository: appCtx.SaleRepo,
		ShopRepo:   appCtx.ShopRepo,
		Redis:      appCtx.Redis,
		Renderer:   appCtx.Renderer,
	}
	app.Handle("POST", "/api/v1/sales/sell", sales.Sell, mid.AuthenticateSessionRequired(appCtx.Authenticator))
	app.Handle("GET", "/sales/:sale_id", sales.View, mid.AuthenticateSessionRequired(appCtx.Authenticator))
	app.Handle("GET", "/sales", sales.Index, mid.AuthenticateSessionRequired(appCtx.Authenticator))

	// Register user management pages.
	us := Users{
		UserRepo:        appCtx.UserRepo,
		UserAccountRepo: appCtx.UserAccountRepo,
		AuthRepo:        appCtx.AuthRepo,
		InviteRepo:      appCtx.InviteRepo,
		GeoRepo:         appCtx.GeoRepo,
		BranchRepo:      appCtx.BranchRepo,
		Redis:           appCtx.Redis,
		Renderer:        appCtx.Renderer,
	}
	app.Handle("POST", "/users/:user_id/update", us.Update, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasRole(auth.RoleSuperAdmin))
	app.Handle("GET", "/users/:user_id/update", us.Update, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasRole(auth.RoleSuperAdmin))
	app.Handle("POST", "/users/:user_id", us.View, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasRole(auth.RoleSuperAdmin))
	app.Handle("GET", "/users/:user_id", us.View, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasAuth())
	app.Handle("POST", "/users/invite/:hash", us.InviteAccept)
	app.Handle("GET", "/users/invite/:hash", us.InviteAccept)
	app.Handle("POST", "/users/invite", us.Invite, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasRole(auth.RoleSuperAdmin))
	app.Handle("GET", "/users/invite", us.Invite, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasRole(auth.RoleSuperAdmin))
	app.Handle("POST", "/users/create", us.Create, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasRole(auth.RoleSuperAdmin))
	app.Handle("GET", "/users/create", us.Create, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasRole(auth.RoleSuperAdmin))
	app.Handle("GET", "/users", us.Index, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasAuth())

	// Register user management and authentication endpoints.
	u := UserRepos{
		UserRepo:        appCtx.UserRepo,
		UserAccountRepo: appCtx.UserAccountRepo,
		AccountRepo:     appCtx.TenantRepo,
		AuthRepo:        appCtx.AuthRepo,
		GeoRepo:         appCtx.GeoRepo,
		Renderer:        appCtx.Renderer,
	}
	app.Handle("POST", "/user/login", u.Login)
	app.Handle("GET", "/user/login", u.Login, waitDbMid)
	app.Handle("GET", "/user/logout", u.Logout)
	app.Handle("POST", "/user/reset-password/:hash", u.ResetConfirm)
	app.Handle("GET", "/user/reset-password/:hash", u.ResetConfirm)
	app.Handle("POST", "/user/reset-password", u.ResetPassword)
	app.Handle("GET", "/user/reset-password", u.ResetPassword)
	app.Handle("POST", "/user/update", u.Update, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasAuth())
	app.Handle("GET", "/user/update", u.Update, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasAuth())
	app.Handle("GET", "/user/account", u.Account, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasAuth())
	app.Handle("GET", "/user/virtual-login/:user_id", u.VirtualLogin, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasRole(auth.RoleAdmin))
	app.Handle("POST", "/user/virtual-login", u.VirtualLogin, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasRole(auth.RoleAdmin))
	app.Handle("GET", "/user/virtual-login", u.VirtualLogin, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasRole(auth.RoleAdmin))
	app.Handle("GET", "/user/virtual-logout", u.VirtualLogout, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasAuth())
	app.Handle("GET", "/user/switch-account/:account_id", u.SwitchAccount, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasAuth())
	app.Handle("POST", "/user/switch-account", u.SwitchAccount, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasAuth())
	app.Handle("GET", "/user/switch-account", u.SwitchAccount, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasAuth())
	app.Handle("POST", "/user", u.View, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasAuth())
	app.Handle("GET", "/user", u.View, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasAuth())

	// Register account management endpoints.
	acc := Account{
		AccountRepo:     appCtx.TenantRepo,
		AccountPrefRepo: appCtx.AccountPrefRepo,
		AuthRepo:        appCtx.AuthRepo,
		Authenticator:   appCtx.Authenticator,
		GeoRepo:         appCtx.GeoRepo,
		Renderer:        appCtx.Renderer,
	}
	app.Handle("POST", "/account/update", acc.Update, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasRole(auth.RoleAdmin))
	app.Handle("GET", "/account/update", acc.Update, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasRole(auth.RoleAdmin))
	app.Handle("POST", "/account", acc.View, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasRole(auth.RoleAdmin))
	app.Handle("GET", "/account", acc.View, mid.AuthenticateSessionRequired(appCtx.Authenticator), mid.HasRole(auth.RoleAdmin))

	// Register signup endpoints.
	s := Signup{
		SignupRepo: appCtx.SignupRepo,
		AuthRepo:   appCtx.AuthRepo,
		GeoRepo:    appCtx.GeoRepo,
		Renderer:   appCtx.Renderer,
	}
	// This route is not authenticated
	app.Handle("POST", "/signup", s.Step1)
	app.Handle("GET", "/signup", s.Step1, waitDbMid)

	// Register example endpoints.
	ex := Examples{
		Renderer: appCtx.Renderer,
	}
	app.Handle("POST", "/examples/flash-messages", ex.FlashMessages, mid.AuthenticateSessionOptional(appCtx.Authenticator))
	app.Handle("GET", "/examples/flash-messages", ex.FlashMessages, mid.AuthenticateSessionOptional(appCtx.Authenticator))
	app.Handle("GET", "/examples/images", ex.Images, mid.AuthenticateSessionOptional(appCtx.Authenticator))

	// Register geo
	g := Geo{
		GeoRepo: appCtx.GeoRepo,
		Redis:   appCtx.Redis,
	}
	app.Handle("GET", "/geo/regions/autocomplete", g.RegionsAutocomplete)
	app.Handle("GET", "/geo/postal_codes/autocomplete", g.PostalCodesAutocomplete)
	app.Handle("GET", "/geo/geonames/postal_code/:postalCode", g.GeonameByPostalCode)
	app.Handle("GET", "/geo/country/:countryCode/timezones", g.CountryTimezones)

	// Register root
	r := Root{
		ShopRepo:        appCtx.ShopRepo,
		CustomerRepo:    appCtx.CustomerRepo,
		AccountRepo:     appCtx.AccountRepo,
		TransactionRepo: appCtx.TransactionRepo,
		Renderer:        appCtx.Renderer,
		Sitemap:         sm,
		WebRoute:        appCtx.WebRoute,
	}
	app.Handle("GET", "/api", r.SitePage)
	app.Handle("GET", "/pricing", r.SitePage)
	app.Handle("GET", "/support", r.SitePage)
	app.Handle("GET", "/legal/privacy", r.SitePage)
	app.Handle("GET", "/legal/terms", r.SitePage)
	// app.Handle("GET", "/", r.Index, mid.AuthenticateSessionOptional(appCtx.Authenticator))
	app.Handle("GET", "/", r.Index, mid.AuthenticateSessionRequired(appCtx.Authenticator))
	app.Handle("GET", "/index.html", r.IndexHtml)
	app.Handle("GET", "/robots.txt", r.RobotTxt)
	app.Handle("GET", "/sitemap.xml", r.SitemapXml)

	// Register health check endpoint. This route is not authenticated.
	check := Check{
		MasterDB: appCtx.MasterDB,
		Redis:    appCtx.Redis,
	}

	app.Handle("GET", "/v1/health", check.Health)
	app.Handle("GET", "/ping", check.Ping)

	// Add sitemap entries for Root.
	smLocAddModified(stm.URL{{"loc", "/"}, {"changefreq", "weekly"}, {"mobile", true}, {"priority", 0.9}}, "site-index.gohtml")
	smLocAddModified(stm.URL{{"loc", "/pricing"}, {"changefreq", "monthly"}, {"mobile", true}, {"priority", 0.8}}, "site-pricing.gohtml")
	smLocAddModified(stm.URL{{"loc", "/support"}, {"changefreq", "monthly"}, {"mobile", true}, {"priority", 0.8}}, "site-support.gohtml")
	smLocAddModified(stm.URL{{"loc", "/api"}, {"changefreq", "monthly"}, {"mobile", true}, {"priority", 0.7}}, "site-api.gohtml")
	smLocAddModified(stm.URL{{"loc", "/legal/privacy"}, {"changefreq", "monthly"}, {"mobile", true}, {"priority", 0.5}}, "legal-privacy.gohtml")
	smLocAddModified(stm.URL{{"loc", "/legal/terms"}, {"changefreq", "monthly"}, {"mobile", true}, {"priority", 0.5}}, "legal-terms.gohtml")

	// Handle static files/pages. Render a custom 404 page when file not found.
	static := func(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {
		err := web.StaticHandler(ctx, w, r, params, appCtx.StaticDir, "")
		if err != nil {
			if os.IsNotExist(err) {
				rmsg := fmt.Sprintf("%s %s not found", r.Method, r.RequestURI)
				err = weberror.NewErrorMessage(ctx, err, http.StatusNotFound, rmsg)
			} else {
				err = weberror.NewError(ctx, err, http.StatusInternalServerError)
			}

			return web.RenderError(ctx, w, r, err, appCtx.Renderer, TmplLayoutBase, TmplContentErrorGeneric, web.MIMETextHTMLCharsetUTF8)
		}

		return nil
	}

	// Static file server
	app.Handle("GET", "/*", static)

	return app
}
