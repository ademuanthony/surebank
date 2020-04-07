package handlers

import (
	"gopkg.in/DataDog/dd-trace-go.v1/contrib/go-redis/redis"

	"merryworld/surebank/internal/account"
	"merryworld/surebank/internal/platform/web"
	"merryworld/surebank/internal/transaction"
)

// Transactions represents the Transaction HTTP method handler set.
type Transactions struct {
	AccountRepo *account.Repository
	Repository  *transaction.Repository
	Redis       *redis.Client
	Renderer    web.Renderer
}
 