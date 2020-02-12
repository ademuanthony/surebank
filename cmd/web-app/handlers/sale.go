package handlers

import (
	"gopkg.in/DataDog/dd-trace-go.v1/contrib/go-redis/redis"
	"merryworld/surebank/internal/branch"
	"merryworld/surebank/internal/platform/web"
	"merryworld/surebank/internal/sale"
	"merryworld/surebank/internal/shop"
)

// Sales represents the sales API method handler set.
type Sales struct {
	Repo       *sale.Repository
	ShopRepo   *shop.Repository
	BranchRepo *branch.Repository
	Redis      *redis.Client
	Renderer   web.Renderer
}
