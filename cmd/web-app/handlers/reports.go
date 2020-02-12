package handlers

import (
	"gopkg.in/DataDog/dd-trace-go.v1/contrib/go-redis/redis"
	"merryworld/surebank/internal/account"
	"merryworld/surebank/internal/customer"
	"merryworld/surebank/internal/platform/web"
	"merryworld/surebank/internal/shop"
	"merryworld/surebank/internal/transaction"
)

// Customers represents the Customers API method handler set.
type Reports struct {
	CustomerRepo    *customer.Repository
	AccountRepo     *account.Repository
	TransactionRepo *transaction.Repository
	ShopRepo        *shop.Repository
	Renderer        web.Renderer
	Redis           *redis.Client
}
