// Code generated by SQLBoiler 3.6.0 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package models

import "testing"

// This test suite runs each operation test in parallel.
// Example, if your database has 3 tables, the suite will run:
// table1, table2 and table3 Delete in parallel
// table1, table2 and table3 Insert in parallel, and so forth.
// It does NOT run each operation group in parallel.
// Separating the tests thusly grants avoidance of Postgres deadlocks.
func TestParent(t *testing.T) {
	t.Run("Accounts", testAccounts)
	t.Run("Branches", testBranches)
	t.Run("Brands", testBrands)
	t.Run("Categories", testCategories)
	t.Run("Customers", testCustomers)
	t.Run("Deposits", testDeposits)
	t.Run("Payments", testPayments)
	t.Run("Products", testProducts)
	t.Run("ProductCategories", testProductCategories)
	t.Run("Sales", testSales)
	t.Run("SaleItems", testSaleItems)
	t.Run("Stocks", testStocks)
	t.Run("Users", testUsers)
	t.Run("Withdrawals", testWithdrawals)
}

func TestDelete(t *testing.T) {
	t.Run("Accounts", testAccountsDelete)
	t.Run("Branches", testBranchesDelete)
	t.Run("Brands", testBrandsDelete)
	t.Run("Categories", testCategoriesDelete)
	t.Run("Customers", testCustomersDelete)
	t.Run("Deposits", testDepositsDelete)
	t.Run("Payments", testPaymentsDelete)
	t.Run("Products", testProductsDelete)
	t.Run("ProductCategories", testProductCategoriesDelete)
	t.Run("Sales", testSalesDelete)
	t.Run("SaleItems", testSaleItemsDelete)
	t.Run("Stocks", testStocksDelete)
	t.Run("Users", testUsersDelete)
	t.Run("Withdrawals", testWithdrawalsDelete)
}

func TestQueryDeleteAll(t *testing.T) {
	t.Run("Accounts", testAccountsQueryDeleteAll)
	t.Run("Branches", testBranchesQueryDeleteAll)
	t.Run("Brands", testBrandsQueryDeleteAll)
	t.Run("Categories", testCategoriesQueryDeleteAll)
	t.Run("Customers", testCustomersQueryDeleteAll)
	t.Run("Deposits", testDepositsQueryDeleteAll)
	t.Run("Payments", testPaymentsQueryDeleteAll)
	t.Run("Products", testProductsQueryDeleteAll)
	t.Run("ProductCategories", testProductCategoriesQueryDeleteAll)
	t.Run("Sales", testSalesQueryDeleteAll)
	t.Run("SaleItems", testSaleItemsQueryDeleteAll)
	t.Run("Stocks", testStocksQueryDeleteAll)
	t.Run("Users", testUsersQueryDeleteAll)
	t.Run("Withdrawals", testWithdrawalsQueryDeleteAll)
}

func TestSliceDeleteAll(t *testing.T) {
	t.Run("Accounts", testAccountsSliceDeleteAll)
	t.Run("Branches", testBranchesSliceDeleteAll)
	t.Run("Brands", testBrandsSliceDeleteAll)
	t.Run("Categories", testCategoriesSliceDeleteAll)
	t.Run("Customers", testCustomersSliceDeleteAll)
	t.Run("Deposits", testDepositsSliceDeleteAll)
	t.Run("Payments", testPaymentsSliceDeleteAll)
	t.Run("Products", testProductsSliceDeleteAll)
	t.Run("ProductCategories", testProductCategoriesSliceDeleteAll)
	t.Run("Sales", testSalesSliceDeleteAll)
	t.Run("SaleItems", testSaleItemsSliceDeleteAll)
	t.Run("Stocks", testStocksSliceDeleteAll)
	t.Run("Users", testUsersSliceDeleteAll)
	t.Run("Withdrawals", testWithdrawalsSliceDeleteAll)
}

func TestExists(t *testing.T) {
	t.Run("Accounts", testAccountsExists)
	t.Run("Branches", testBranchesExists)
	t.Run("Brands", testBrandsExists)
	t.Run("Categories", testCategoriesExists)
	t.Run("Customers", testCustomersExists)
	t.Run("Deposits", testDepositsExists)
	t.Run("Payments", testPaymentsExists)
	t.Run("Products", testProductsExists)
	t.Run("ProductCategories", testProductCategoriesExists)
	t.Run("Sales", testSalesExists)
	t.Run("SaleItems", testSaleItemsExists)
	t.Run("Stocks", testStocksExists)
	t.Run("Users", testUsersExists)
	t.Run("Withdrawals", testWithdrawalsExists)
}

func TestFind(t *testing.T) {
	t.Run("Accounts", testAccountsFind)
	t.Run("Branches", testBranchesFind)
	t.Run("Brands", testBrandsFind)
	t.Run("Categories", testCategoriesFind)
	t.Run("Customers", testCustomersFind)
	t.Run("Deposits", testDepositsFind)
	t.Run("Payments", testPaymentsFind)
	t.Run("Products", testProductsFind)
	t.Run("ProductCategories", testProductCategoriesFind)
	t.Run("Sales", testSalesFind)
	t.Run("SaleItems", testSaleItemsFind)
	t.Run("Stocks", testStocksFind)
	t.Run("Users", testUsersFind)
	t.Run("Withdrawals", testWithdrawalsFind)
}

func TestBind(t *testing.T) {
	t.Run("Accounts", testAccountsBind)
	t.Run("Branches", testBranchesBind)
	t.Run("Brands", testBrandsBind)
	t.Run("Categories", testCategoriesBind)
	t.Run("Customers", testCustomersBind)
	t.Run("Deposits", testDepositsBind)
	t.Run("Payments", testPaymentsBind)
	t.Run("Products", testProductsBind)
	t.Run("ProductCategories", testProductCategoriesBind)
	t.Run("Sales", testSalesBind)
	t.Run("SaleItems", testSaleItemsBind)
	t.Run("Stocks", testStocksBind)
	t.Run("Users", testUsersBind)
	t.Run("Withdrawals", testWithdrawalsBind)
}

func TestOne(t *testing.T) {
	t.Run("Accounts", testAccountsOne)
	t.Run("Branches", testBranchesOne)
	t.Run("Brands", testBrandsOne)
	t.Run("Categories", testCategoriesOne)
	t.Run("Customers", testCustomersOne)
	t.Run("Deposits", testDepositsOne)
	t.Run("Payments", testPaymentsOne)
	t.Run("Products", testProductsOne)
	t.Run("ProductCategories", testProductCategoriesOne)
	t.Run("Sales", testSalesOne)
	t.Run("SaleItems", testSaleItemsOne)
	t.Run("Stocks", testStocksOne)
	t.Run("Users", testUsersOne)
	t.Run("Withdrawals", testWithdrawalsOne)
}

func TestAll(t *testing.T) {
	t.Run("Accounts", testAccountsAll)
	t.Run("Branches", testBranchesAll)
	t.Run("Brands", testBrandsAll)
	t.Run("Categories", testCategoriesAll)
	t.Run("Customers", testCustomersAll)
	t.Run("Deposits", testDepositsAll)
	t.Run("Payments", testPaymentsAll)
	t.Run("Products", testProductsAll)
	t.Run("ProductCategories", testProductCategoriesAll)
	t.Run("Sales", testSalesAll)
	t.Run("SaleItems", testSaleItemsAll)
	t.Run("Stocks", testStocksAll)
	t.Run("Users", testUsersAll)
	t.Run("Withdrawals", testWithdrawalsAll)
}

func TestCount(t *testing.T) {
	t.Run("Accounts", testAccountsCount)
	t.Run("Branches", testBranchesCount)
	t.Run("Brands", testBrandsCount)
	t.Run("Categories", testCategoriesCount)
	t.Run("Customers", testCustomersCount)
	t.Run("Deposits", testDepositsCount)
	t.Run("Payments", testPaymentsCount)
	t.Run("Products", testProductsCount)
	t.Run("ProductCategories", testProductCategoriesCount)
	t.Run("Sales", testSalesCount)
	t.Run("SaleItems", testSaleItemsCount)
	t.Run("Stocks", testStocksCount)
	t.Run("Users", testUsersCount)
	t.Run("Withdrawals", testWithdrawalsCount)
}

func TestInsert(t *testing.T) {
	t.Run("Accounts", testAccountsInsert)
	t.Run("Accounts", testAccountsInsertWhitelist)
	t.Run("Branches", testBranchesInsert)
	t.Run("Branches", testBranchesInsertWhitelist)
	t.Run("Brands", testBrandsInsert)
	t.Run("Brands", testBrandsInsertWhitelist)
	t.Run("Categories", testCategoriesInsert)
	t.Run("Categories", testCategoriesInsertWhitelist)
	t.Run("Customers", testCustomersInsert)
	t.Run("Customers", testCustomersInsertWhitelist)
	t.Run("Deposits", testDepositsInsert)
	t.Run("Deposits", testDepositsInsertWhitelist)
	t.Run("Payments", testPaymentsInsert)
	t.Run("Payments", testPaymentsInsertWhitelist)
	t.Run("Products", testProductsInsert)
	t.Run("Products", testProductsInsertWhitelist)
	t.Run("ProductCategories", testProductCategoriesInsert)
	t.Run("ProductCategories", testProductCategoriesInsertWhitelist)
	t.Run("Sales", testSalesInsert)
	t.Run("Sales", testSalesInsertWhitelist)
	t.Run("SaleItems", testSaleItemsInsert)
	t.Run("SaleItems", testSaleItemsInsertWhitelist)
	t.Run("Stocks", testStocksInsert)
	t.Run("Stocks", testStocksInsertWhitelist)
	t.Run("Users", testUsersInsert)
	t.Run("Users", testUsersInsertWhitelist)
	t.Run("Withdrawals", testWithdrawalsInsert)
	t.Run("Withdrawals", testWithdrawalsInsertWhitelist)
}

// TestToOne tests cannot be run in parallel
// or deadlocks can occur.
func TestToOne(t *testing.T) {
	t.Run("AccountToBranchUsingBranch", testAccountToOneBranchUsingBranch)
	t.Run("AccountToCustomerUsingCustomer", testAccountToOneCustomerUsingCustomer)
	t.Run("AccountToUserUsingSalesRep", testAccountToOneUserUsingSalesRep)
	t.Run("CustomerToBranchUsingBranch", testCustomerToOneBranchUsingBranch)
	t.Run("CustomerToUserUsingSalesRep", testCustomerToOneUserUsingSalesRep)
	t.Run("DepositToAccountUsingAccount", testDepositToOneAccountUsingAccount)
	t.Run("DepositToUserUsingSalesRep", testDepositToOneUserUsingSalesRep)
	t.Run("PaymentToSaleUsingSale", testPaymentToOneSaleUsingSale)
	t.Run("PaymentToUserUsingSalesRep", testPaymentToOneUserUsingSalesRep)
	t.Run("ProductToUserUsingArchivedBy", testProductToOneUserUsingArchivedBy)
	t.Run("ProductToBrandUsingBrand", testProductToOneBrandUsingBrand)
	t.Run("ProductToCategoryUsingCategory", testProductToOneCategoryUsingCategory)
	t.Run("ProductToUserUsingCreatedBy", testProductToOneUserUsingCreatedBy)
	t.Run("ProductToUserUsingUpdatedBy", testProductToOneUserUsingUpdatedBy)
	t.Run("ProductCategoryToCategoryUsingCategory", testProductCategoryToOneCategoryUsingCategory)
	t.Run("ProductCategoryToProductUsingProduct", testProductCategoryToOneProductUsingProduct)
	t.Run("SaleToUserUsingArchivedBy", testSaleToOneUserUsingArchivedBy)
	t.Run("SaleToBranchUsingBranch", testSaleToOneBranchUsingBranch)
	t.Run("SaleToUserUsingCreatedBy", testSaleToOneUserUsingCreatedBy)
	t.Run("SaleToUserUsingUpdatedBy", testSaleToOneUserUsingUpdatedBy)
	t.Run("SaleItemToProductUsingProduct", testSaleItemToOneProductUsingProduct)
	t.Run("SaleItemToSaleUsingSale", testSaleItemToOneSaleUsingSale)
	t.Run("StockToUserUsingArchivedBy", testStockToOneUserUsingArchivedBy)
	t.Run("StockToBranchUsingBranch", testStockToOneBranchUsingBranch)
	t.Run("StockToUserUsingCreatedBy", testStockToOneUserUsingCreatedBy)
	t.Run("StockToProductUsingProduct", testStockToOneProductUsingProduct)
	t.Run("StockToUserUsingUpdatedBy", testStockToOneUserUsingUpdatedBy)
	t.Run("UserToBranchUsingBranch", testUserToOneBranchUsingBranch)
	t.Run("WithdrawalToAccountUsingAccount", testWithdrawalToOneAccountUsingAccount)
	t.Run("WithdrawalToUserUsingSalesRep", testWithdrawalToOneUserUsingSalesRep)
}

// TestOneToOne tests cannot be run in parallel
// or deadlocks can occur.
func TestOneToOne(t *testing.T) {}

// TestToMany tests cannot be run in parallel
// or deadlocks can occur.
func TestToMany(t *testing.T) {
	t.Run("AccountToDeposits", testAccountToManyDeposits)
	t.Run("AccountToWithdrawals", testAccountToManyWithdrawals)
	t.Run("BranchToAccounts", testBranchToManyAccounts)
	t.Run("BranchToCustomers", testBranchToManyCustomers)
	t.Run("BranchToSales", testBranchToManySales)
	t.Run("BranchToStocks", testBranchToManyStocks)
	t.Run("BranchToUsers", testBranchToManyUsers)
	t.Run("BrandToProducts", testBrandToManyProducts)
	t.Run("CategoryToProducts", testCategoryToManyProducts)
	t.Run("CategoryToProductCategories", testCategoryToManyProductCategories)
	t.Run("CustomerToAccounts", testCustomerToManyAccounts)
	t.Run("ProductToProductCategories", testProductToManyProductCategories)
	t.Run("ProductToSaleItems", testProductToManySaleItems)
	t.Run("ProductToStocks", testProductToManyStocks)
	t.Run("SaleToPayments", testSaleToManyPayments)
	t.Run("SaleToSaleItems", testSaleToManySaleItems)
	t.Run("UserToSalesRepAccounts", testUserToManySalesRepAccounts)
	t.Run("UserToSalesRepCustomers", testUserToManySalesRepCustomers)
	t.Run("UserToSalesRepDeposits", testUserToManySalesRepDeposits)
	t.Run("UserToSalesRepPayments", testUserToManySalesRepPayments)
	t.Run("UserToArchivedByProducts", testUserToManyArchivedByProducts)
	t.Run("UserToCreatedByProducts", testUserToManyCreatedByProducts)
	t.Run("UserToUpdatedByProducts", testUserToManyUpdatedByProducts)
	t.Run("UserToArchivedBySales", testUserToManyArchivedBySales)
	t.Run("UserToCreatedBySales", testUserToManyCreatedBySales)
	t.Run("UserToUpdatedBySales", testUserToManyUpdatedBySales)
	t.Run("UserToArchivedByStocks", testUserToManyArchivedByStocks)
	t.Run("UserToCreatedByStocks", testUserToManyCreatedByStocks)
	t.Run("UserToUpdatedByStocks", testUserToManyUpdatedByStocks)
	t.Run("UserToSalesRepWithdrawals", testUserToManySalesRepWithdrawals)
}

// TestToOneSet tests cannot be run in parallel
// or deadlocks can occur.
func TestToOneSet(t *testing.T) {
	t.Run("AccountToBranchUsingAccounts", testAccountToOneSetOpBranchUsingBranch)
	t.Run("AccountToCustomerUsingAccounts", testAccountToOneSetOpCustomerUsingCustomer)
	t.Run("AccountToUserUsingSalesRepAccounts", testAccountToOneSetOpUserUsingSalesRep)
	t.Run("CustomerToBranchUsingCustomers", testCustomerToOneSetOpBranchUsingBranch)
	t.Run("CustomerToUserUsingSalesRepCustomers", testCustomerToOneSetOpUserUsingSalesRep)
	t.Run("DepositToAccountUsingDeposits", testDepositToOneSetOpAccountUsingAccount)
	t.Run("DepositToUserUsingSalesRepDeposits", testDepositToOneSetOpUserUsingSalesRep)
	t.Run("PaymentToSaleUsingPayments", testPaymentToOneSetOpSaleUsingSale)
	t.Run("PaymentToUserUsingSalesRepPayments", testPaymentToOneSetOpUserUsingSalesRep)
	t.Run("ProductToUserUsingArchivedByProducts", testProductToOneSetOpUserUsingArchivedBy)
	t.Run("ProductToBrandUsingProducts", testProductToOneSetOpBrandUsingBrand)
	t.Run("ProductToCategoryUsingProducts", testProductToOneSetOpCategoryUsingCategory)
	t.Run("ProductToUserUsingCreatedByProducts", testProductToOneSetOpUserUsingCreatedBy)
	t.Run("ProductToUserUsingUpdatedByProducts", testProductToOneSetOpUserUsingUpdatedBy)
	t.Run("ProductCategoryToCategoryUsingProductCategories", testProductCategoryToOneSetOpCategoryUsingCategory)
	t.Run("ProductCategoryToProductUsingProductCategories", testProductCategoryToOneSetOpProductUsingProduct)
	t.Run("SaleToUserUsingArchivedBySales", testSaleToOneSetOpUserUsingArchivedBy)
	t.Run("SaleToBranchUsingSales", testSaleToOneSetOpBranchUsingBranch)
	t.Run("SaleToUserUsingCreatedBySales", testSaleToOneSetOpUserUsingCreatedBy)
	t.Run("SaleToUserUsingUpdatedBySales", testSaleToOneSetOpUserUsingUpdatedBy)
	t.Run("SaleItemToProductUsingSaleItems", testSaleItemToOneSetOpProductUsingProduct)
	t.Run("SaleItemToSaleUsingSaleItems", testSaleItemToOneSetOpSaleUsingSale)
	t.Run("StockToUserUsingArchivedByStocks", testStockToOneSetOpUserUsingArchivedBy)
	t.Run("StockToBranchUsingStocks", testStockToOneSetOpBranchUsingBranch)
	t.Run("StockToUserUsingCreatedByStocks", testStockToOneSetOpUserUsingCreatedBy)
	t.Run("StockToProductUsingStocks", testStockToOneSetOpProductUsingProduct)
	t.Run("StockToUserUsingUpdatedByStocks", testStockToOneSetOpUserUsingUpdatedBy)
	t.Run("UserToBranchUsingUsers", testUserToOneSetOpBranchUsingBranch)
	t.Run("WithdrawalToAccountUsingWithdrawals", testWithdrawalToOneSetOpAccountUsingAccount)
	t.Run("WithdrawalToUserUsingSalesRepWithdrawals", testWithdrawalToOneSetOpUserUsingSalesRep)
}

// TestToOneRemove tests cannot be run in parallel
// or deadlocks can occur.
func TestToOneRemove(t *testing.T) {
	t.Run("ProductToUserUsingArchivedByProducts", testProductToOneRemoveOpUserUsingArchivedBy)
	t.Run("ProductToBrandUsingProducts", testProductToOneRemoveOpBrandUsingBrand)
	t.Run("SaleToUserUsingArchivedBySales", testSaleToOneRemoveOpUserUsingArchivedBy)
	t.Run("SaleToUserUsingUpdatedBySales", testSaleToOneRemoveOpUserUsingUpdatedBy)
	t.Run("StockToUserUsingArchivedByStocks", testStockToOneRemoveOpUserUsingArchivedBy)
}

// TestOneToOneSet tests cannot be run in parallel
// or deadlocks can occur.
func TestOneToOneSet(t *testing.T) {}

// TestOneToOneRemove tests cannot be run in parallel
// or deadlocks can occur.
func TestOneToOneRemove(t *testing.T) {}

// TestToManyAdd tests cannot be run in parallel
// or deadlocks can occur.
func TestToManyAdd(t *testing.T) {
	t.Run("AccountToDeposits", testAccountToManyAddOpDeposits)
	t.Run("AccountToWithdrawals", testAccountToManyAddOpWithdrawals)
	t.Run("BranchToAccounts", testBranchToManyAddOpAccounts)
	t.Run("BranchToCustomers", testBranchToManyAddOpCustomers)
	t.Run("BranchToSales", testBranchToManyAddOpSales)
	t.Run("BranchToStocks", testBranchToManyAddOpStocks)
	t.Run("BranchToUsers", testBranchToManyAddOpUsers)
	t.Run("BrandToProducts", testBrandToManyAddOpProducts)
	t.Run("CategoryToProducts", testCategoryToManyAddOpProducts)
	t.Run("CategoryToProductCategories", testCategoryToManyAddOpProductCategories)
	t.Run("CustomerToAccounts", testCustomerToManyAddOpAccounts)
	t.Run("ProductToProductCategories", testProductToManyAddOpProductCategories)
	t.Run("ProductToSaleItems", testProductToManyAddOpSaleItems)
	t.Run("ProductToStocks", testProductToManyAddOpStocks)
	t.Run("SaleToPayments", testSaleToManyAddOpPayments)
	t.Run("SaleToSaleItems", testSaleToManyAddOpSaleItems)
	t.Run("UserToSalesRepAccounts", testUserToManyAddOpSalesRepAccounts)
	t.Run("UserToSalesRepCustomers", testUserToManyAddOpSalesRepCustomers)
	t.Run("UserToSalesRepDeposits", testUserToManyAddOpSalesRepDeposits)
	t.Run("UserToSalesRepPayments", testUserToManyAddOpSalesRepPayments)
	t.Run("UserToArchivedByProducts", testUserToManyAddOpArchivedByProducts)
	t.Run("UserToCreatedByProducts", testUserToManyAddOpCreatedByProducts)
	t.Run("UserToUpdatedByProducts", testUserToManyAddOpUpdatedByProducts)
	t.Run("UserToArchivedBySales", testUserToManyAddOpArchivedBySales)
	t.Run("UserToCreatedBySales", testUserToManyAddOpCreatedBySales)
	t.Run("UserToUpdatedBySales", testUserToManyAddOpUpdatedBySales)
	t.Run("UserToArchivedByStocks", testUserToManyAddOpArchivedByStocks)
	t.Run("UserToCreatedByStocks", testUserToManyAddOpCreatedByStocks)
	t.Run("UserToUpdatedByStocks", testUserToManyAddOpUpdatedByStocks)
	t.Run("UserToSalesRepWithdrawals", testUserToManyAddOpSalesRepWithdrawals)
}

// TestToManySet tests cannot be run in parallel
// or deadlocks can occur.
func TestToManySet(t *testing.T) {
	t.Run("BrandToProducts", testBrandToManySetOpProducts)
	t.Run("UserToArchivedByProducts", testUserToManySetOpArchivedByProducts)
	t.Run("UserToArchivedBySales", testUserToManySetOpArchivedBySales)
	t.Run("UserToUpdatedBySales", testUserToManySetOpUpdatedBySales)
	t.Run("UserToArchivedByStocks", testUserToManySetOpArchivedByStocks)
}

// TestToManyRemove tests cannot be run in parallel
// or deadlocks can occur.
func TestToManyRemove(t *testing.T) {
	t.Run("BrandToProducts", testBrandToManyRemoveOpProducts)
	t.Run("UserToArchivedByProducts", testUserToManyRemoveOpArchivedByProducts)
	t.Run("UserToArchivedBySales", testUserToManyRemoveOpArchivedBySales)
	t.Run("UserToUpdatedBySales", testUserToManyRemoveOpUpdatedBySales)
	t.Run("UserToArchivedByStocks", testUserToManyRemoveOpArchivedByStocks)
}

func TestReload(t *testing.T) {
	t.Run("Accounts", testAccountsReload)
	t.Run("Branches", testBranchesReload)
	t.Run("Brands", testBrandsReload)
	t.Run("Categories", testCategoriesReload)
	t.Run("Customers", testCustomersReload)
	t.Run("Deposits", testDepositsReload)
	t.Run("Payments", testPaymentsReload)
	t.Run("Products", testProductsReload)
	t.Run("ProductCategories", testProductCategoriesReload)
	t.Run("Sales", testSalesReload)
	t.Run("SaleItems", testSaleItemsReload)
	t.Run("Stocks", testStocksReload)
	t.Run("Users", testUsersReload)
	t.Run("Withdrawals", testWithdrawalsReload)
}

func TestReloadAll(t *testing.T) {
	t.Run("Accounts", testAccountsReloadAll)
	t.Run("Branches", testBranchesReloadAll)
	t.Run("Brands", testBrandsReloadAll)
	t.Run("Categories", testCategoriesReloadAll)
	t.Run("Customers", testCustomersReloadAll)
	t.Run("Deposits", testDepositsReloadAll)
	t.Run("Payments", testPaymentsReloadAll)
	t.Run("Products", testProductsReloadAll)
	t.Run("ProductCategories", testProductCategoriesReloadAll)
	t.Run("Sales", testSalesReloadAll)
	t.Run("SaleItems", testSaleItemsReloadAll)
	t.Run("Stocks", testStocksReloadAll)
	t.Run("Users", testUsersReloadAll)
	t.Run("Withdrawals", testWithdrawalsReloadAll)
}

func TestSelect(t *testing.T) {
	t.Run("Accounts", testAccountsSelect)
	t.Run("Branches", testBranchesSelect)
	t.Run("Brands", testBrandsSelect)
	t.Run("Categories", testCategoriesSelect)
	t.Run("Customers", testCustomersSelect)
	t.Run("Deposits", testDepositsSelect)
	t.Run("Payments", testPaymentsSelect)
	t.Run("Products", testProductsSelect)
	t.Run("ProductCategories", testProductCategoriesSelect)
	t.Run("Sales", testSalesSelect)
	t.Run("SaleItems", testSaleItemsSelect)
	t.Run("Stocks", testStocksSelect)
	t.Run("Users", testUsersSelect)
	t.Run("Withdrawals", testWithdrawalsSelect)
}

func TestUpdate(t *testing.T) {
	t.Run("Accounts", testAccountsUpdate)
	t.Run("Branches", testBranchesUpdate)
	t.Run("Brands", testBrandsUpdate)
	t.Run("Categories", testCategoriesUpdate)
	t.Run("Customers", testCustomersUpdate)
	t.Run("Deposits", testDepositsUpdate)
	t.Run("Payments", testPaymentsUpdate)
	t.Run("Products", testProductsUpdate)
	t.Run("ProductCategories", testProductCategoriesUpdate)
	t.Run("Sales", testSalesUpdate)
	t.Run("SaleItems", testSaleItemsUpdate)
	t.Run("Stocks", testStocksUpdate)
	t.Run("Users", testUsersUpdate)
	t.Run("Withdrawals", testWithdrawalsUpdate)
}

func TestSliceUpdateAll(t *testing.T) {
	t.Run("Accounts", testAccountsSliceUpdateAll)
	t.Run("Branches", testBranchesSliceUpdateAll)
	t.Run("Brands", testBrandsSliceUpdateAll)
	t.Run("Categories", testCategoriesSliceUpdateAll)
	t.Run("Customers", testCustomersSliceUpdateAll)
	t.Run("Deposits", testDepositsSliceUpdateAll)
	t.Run("Payments", testPaymentsSliceUpdateAll)
	t.Run("Products", testProductsSliceUpdateAll)
	t.Run("ProductCategories", testProductCategoriesSliceUpdateAll)
	t.Run("Sales", testSalesSliceUpdateAll)
	t.Run("SaleItems", testSaleItemsSliceUpdateAll)
	t.Run("Stocks", testStocksSliceUpdateAll)
	t.Run("Users", testUsersSliceUpdateAll)
	t.Run("Withdrawals", testWithdrawalsSliceUpdateAll)
}