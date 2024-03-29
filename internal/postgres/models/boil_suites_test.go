// Code generated by SQLBoiler 4.6.0 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
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
	t.Run("BankAccounts", testBankAccounts)
	t.Run("BankDeposits", testBankDeposits)
	t.Run("Branches", testBranches)
	t.Run("Brands", testBrands)
	t.Run("Categories", testCategories)
	t.Run("Customers", testCustomers)
	t.Run("DailySummaries", testDailySummaries)
	t.Run("DSCommissions", testDSCommissions)
	t.Run("Expenditures", testExpenditures)
	t.Run("Inventories", testInventories)
	t.Run("Payments", testPayments)
	t.Run("Products", testProducts)
	t.Run("ProductCategories", testProductCategories)
	t.Run("Profits", testProfits)
	t.Run("RepsExpenses", testRepsExpenses)
	t.Run("Sales", testSales)
	t.Run("SaleItems", testSaleItems)
	t.Run("Transactions", testTransactions)
	t.Run("Users", testUsers)
}

func TestDelete(t *testing.T) {
	t.Run("Accounts", testAccountsDelete)
	t.Run("BankAccounts", testBankAccountsDelete)
	t.Run("BankDeposits", testBankDepositsDelete)
	t.Run("Branches", testBranchesDelete)
	t.Run("Brands", testBrandsDelete)
	t.Run("Categories", testCategoriesDelete)
	t.Run("Customers", testCustomersDelete)
	t.Run("DailySummaries", testDailySummariesDelete)
	t.Run("DSCommissions", testDSCommissionsDelete)
	t.Run("Expenditures", testExpendituresDelete)
	t.Run("Inventories", testInventoriesDelete)
	t.Run("Payments", testPaymentsDelete)
	t.Run("Products", testProductsDelete)
	t.Run("ProductCategories", testProductCategoriesDelete)
	t.Run("Profits", testProfitsDelete)
	t.Run("RepsExpenses", testRepsExpensesDelete)
	t.Run("Sales", testSalesDelete)
	t.Run("SaleItems", testSaleItemsDelete)
	t.Run("Transactions", testTransactionsDelete)
	t.Run("Users", testUsersDelete)
}

func TestQueryDeleteAll(t *testing.T) {
	t.Run("Accounts", testAccountsQueryDeleteAll)
	t.Run("BankAccounts", testBankAccountsQueryDeleteAll)
	t.Run("BankDeposits", testBankDepositsQueryDeleteAll)
	t.Run("Branches", testBranchesQueryDeleteAll)
	t.Run("Brands", testBrandsQueryDeleteAll)
	t.Run("Categories", testCategoriesQueryDeleteAll)
	t.Run("Customers", testCustomersQueryDeleteAll)
	t.Run("DailySummaries", testDailySummariesQueryDeleteAll)
	t.Run("DSCommissions", testDSCommissionsQueryDeleteAll)
	t.Run("Expenditures", testExpendituresQueryDeleteAll)
	t.Run("Inventories", testInventoriesQueryDeleteAll)
	t.Run("Payments", testPaymentsQueryDeleteAll)
	t.Run("Products", testProductsQueryDeleteAll)
	t.Run("ProductCategories", testProductCategoriesQueryDeleteAll)
	t.Run("Profits", testProfitsQueryDeleteAll)
	t.Run("RepsExpenses", testRepsExpensesQueryDeleteAll)
	t.Run("Sales", testSalesQueryDeleteAll)
	t.Run("SaleItems", testSaleItemsQueryDeleteAll)
	t.Run("Transactions", testTransactionsQueryDeleteAll)
	t.Run("Users", testUsersQueryDeleteAll)
}

func TestSliceDeleteAll(t *testing.T) {
	t.Run("Accounts", testAccountsSliceDeleteAll)
	t.Run("BankAccounts", testBankAccountsSliceDeleteAll)
	t.Run("BankDeposits", testBankDepositsSliceDeleteAll)
	t.Run("Branches", testBranchesSliceDeleteAll)
	t.Run("Brands", testBrandsSliceDeleteAll)
	t.Run("Categories", testCategoriesSliceDeleteAll)
	t.Run("Customers", testCustomersSliceDeleteAll)
	t.Run("DailySummaries", testDailySummariesSliceDeleteAll)
	t.Run("DSCommissions", testDSCommissionsSliceDeleteAll)
	t.Run("Expenditures", testExpendituresSliceDeleteAll)
	t.Run("Inventories", testInventoriesSliceDeleteAll)
	t.Run("Payments", testPaymentsSliceDeleteAll)
	t.Run("Products", testProductsSliceDeleteAll)
	t.Run("ProductCategories", testProductCategoriesSliceDeleteAll)
	t.Run("Profits", testProfitsSliceDeleteAll)
	t.Run("RepsExpenses", testRepsExpensesSliceDeleteAll)
	t.Run("Sales", testSalesSliceDeleteAll)
	t.Run("SaleItems", testSaleItemsSliceDeleteAll)
	t.Run("Transactions", testTransactionsSliceDeleteAll)
	t.Run("Users", testUsersSliceDeleteAll)
}

func TestExists(t *testing.T) {
	t.Run("Accounts", testAccountsExists)
	t.Run("BankAccounts", testBankAccountsExists)
	t.Run("BankDeposits", testBankDepositsExists)
	t.Run("Branches", testBranchesExists)
	t.Run("Brands", testBrandsExists)
	t.Run("Categories", testCategoriesExists)
	t.Run("Customers", testCustomersExists)
	t.Run("DailySummaries", testDailySummariesExists)
	t.Run("DSCommissions", testDSCommissionsExists)
	t.Run("Expenditures", testExpendituresExists)
	t.Run("Inventories", testInventoriesExists)
	t.Run("Payments", testPaymentsExists)
	t.Run("Products", testProductsExists)
	t.Run("ProductCategories", testProductCategoriesExists)
	t.Run("Profits", testProfitsExists)
	t.Run("RepsExpenses", testRepsExpensesExists)
	t.Run("Sales", testSalesExists)
	t.Run("SaleItems", testSaleItemsExists)
	t.Run("Transactions", testTransactionsExists)
	t.Run("Users", testUsersExists)
}

func TestFind(t *testing.T) {
	t.Run("Accounts", testAccountsFind)
	t.Run("BankAccounts", testBankAccountsFind)
	t.Run("BankDeposits", testBankDepositsFind)
	t.Run("Branches", testBranchesFind)
	t.Run("Brands", testBrandsFind)
	t.Run("Categories", testCategoriesFind)
	t.Run("Customers", testCustomersFind)
	t.Run("DailySummaries", testDailySummariesFind)
	t.Run("DSCommissions", testDSCommissionsFind)
	t.Run("Expenditures", testExpendituresFind)
	t.Run("Inventories", testInventoriesFind)
	t.Run("Payments", testPaymentsFind)
	t.Run("Products", testProductsFind)
	t.Run("ProductCategories", testProductCategoriesFind)
	t.Run("Profits", testProfitsFind)
	t.Run("RepsExpenses", testRepsExpensesFind)
	t.Run("Sales", testSalesFind)
	t.Run("SaleItems", testSaleItemsFind)
	t.Run("Transactions", testTransactionsFind)
	t.Run("Users", testUsersFind)
}

func TestBind(t *testing.T) {
	t.Run("Accounts", testAccountsBind)
	t.Run("BankAccounts", testBankAccountsBind)
	t.Run("BankDeposits", testBankDepositsBind)
	t.Run("Branches", testBranchesBind)
	t.Run("Brands", testBrandsBind)
	t.Run("Categories", testCategoriesBind)
	t.Run("Customers", testCustomersBind)
	t.Run("DailySummaries", testDailySummariesBind)
	t.Run("DSCommissions", testDSCommissionsBind)
	t.Run("Expenditures", testExpendituresBind)
	t.Run("Inventories", testInventoriesBind)
	t.Run("Payments", testPaymentsBind)
	t.Run("Products", testProductsBind)
	t.Run("ProductCategories", testProductCategoriesBind)
	t.Run("Profits", testProfitsBind)
	t.Run("RepsExpenses", testRepsExpensesBind)
	t.Run("Sales", testSalesBind)
	t.Run("SaleItems", testSaleItemsBind)
	t.Run("Transactions", testTransactionsBind)
	t.Run("Users", testUsersBind)
}

func TestOne(t *testing.T) {
	t.Run("Accounts", testAccountsOne)
	t.Run("BankAccounts", testBankAccountsOne)
	t.Run("BankDeposits", testBankDepositsOne)
	t.Run("Branches", testBranchesOne)
	t.Run("Brands", testBrandsOne)
	t.Run("Categories", testCategoriesOne)
	t.Run("Customers", testCustomersOne)
	t.Run("DailySummaries", testDailySummariesOne)
	t.Run("DSCommissions", testDSCommissionsOne)
	t.Run("Expenditures", testExpendituresOne)
	t.Run("Inventories", testInventoriesOne)
	t.Run("Payments", testPaymentsOne)
	t.Run("Products", testProductsOne)
	t.Run("ProductCategories", testProductCategoriesOne)
	t.Run("Profits", testProfitsOne)
	t.Run("RepsExpenses", testRepsExpensesOne)
	t.Run("Sales", testSalesOne)
	t.Run("SaleItems", testSaleItemsOne)
	t.Run("Transactions", testTransactionsOne)
	t.Run("Users", testUsersOne)
}

func TestAll(t *testing.T) {
	t.Run("Accounts", testAccountsAll)
	t.Run("BankAccounts", testBankAccountsAll)
	t.Run("BankDeposits", testBankDepositsAll)
	t.Run("Branches", testBranchesAll)
	t.Run("Brands", testBrandsAll)
	t.Run("Categories", testCategoriesAll)
	t.Run("Customers", testCustomersAll)
	t.Run("DailySummaries", testDailySummariesAll)
	t.Run("DSCommissions", testDSCommissionsAll)
	t.Run("Expenditures", testExpendituresAll)
	t.Run("Inventories", testInventoriesAll)
	t.Run("Payments", testPaymentsAll)
	t.Run("Products", testProductsAll)
	t.Run("ProductCategories", testProductCategoriesAll)
	t.Run("Profits", testProfitsAll)
	t.Run("RepsExpenses", testRepsExpensesAll)
	t.Run("Sales", testSalesAll)
	t.Run("SaleItems", testSaleItemsAll)
	t.Run("Transactions", testTransactionsAll)
	t.Run("Users", testUsersAll)
}

func TestCount(t *testing.T) {
	t.Run("Accounts", testAccountsCount)
	t.Run("BankAccounts", testBankAccountsCount)
	t.Run("BankDeposits", testBankDepositsCount)
	t.Run("Branches", testBranchesCount)
	t.Run("Brands", testBrandsCount)
	t.Run("Categories", testCategoriesCount)
	t.Run("Customers", testCustomersCount)
	t.Run("DailySummaries", testDailySummariesCount)
	t.Run("DSCommissions", testDSCommissionsCount)
	t.Run("Expenditures", testExpendituresCount)
	t.Run("Inventories", testInventoriesCount)
	t.Run("Payments", testPaymentsCount)
	t.Run("Products", testProductsCount)
	t.Run("ProductCategories", testProductCategoriesCount)
	t.Run("Profits", testProfitsCount)
	t.Run("RepsExpenses", testRepsExpensesCount)
	t.Run("Sales", testSalesCount)
	t.Run("SaleItems", testSaleItemsCount)
	t.Run("Transactions", testTransactionsCount)
	t.Run("Users", testUsersCount)
}

func TestInsert(t *testing.T) {
	t.Run("Accounts", testAccountsInsert)
	t.Run("Accounts", testAccountsInsertWhitelist)
	t.Run("BankAccounts", testBankAccountsInsert)
	t.Run("BankAccounts", testBankAccountsInsertWhitelist)
	t.Run("BankDeposits", testBankDepositsInsert)
	t.Run("BankDeposits", testBankDepositsInsertWhitelist)
	t.Run("Branches", testBranchesInsert)
	t.Run("Branches", testBranchesInsertWhitelist)
	t.Run("Brands", testBrandsInsert)
	t.Run("Brands", testBrandsInsertWhitelist)
	t.Run("Categories", testCategoriesInsert)
	t.Run("Categories", testCategoriesInsertWhitelist)
	t.Run("Customers", testCustomersInsert)
	t.Run("Customers", testCustomersInsertWhitelist)
	t.Run("DailySummaries", testDailySummariesInsert)
	t.Run("DailySummaries", testDailySummariesInsertWhitelist)
	t.Run("DSCommissions", testDSCommissionsInsert)
	t.Run("DSCommissions", testDSCommissionsInsertWhitelist)
	t.Run("Expenditures", testExpendituresInsert)
	t.Run("Expenditures", testExpendituresInsertWhitelist)
	t.Run("Inventories", testInventoriesInsert)
	t.Run("Inventories", testInventoriesInsertWhitelist)
	t.Run("Payments", testPaymentsInsert)
	t.Run("Payments", testPaymentsInsertWhitelist)
	t.Run("Products", testProductsInsert)
	t.Run("Products", testProductsInsertWhitelist)
	t.Run("ProductCategories", testProductCategoriesInsert)
	t.Run("ProductCategories", testProductCategoriesInsertWhitelist)
	t.Run("Profits", testProfitsInsert)
	t.Run("Profits", testProfitsInsertWhitelist)
	t.Run("RepsExpenses", testRepsExpensesInsert)
	t.Run("RepsExpenses", testRepsExpensesInsertWhitelist)
	t.Run("Sales", testSalesInsert)
	t.Run("Sales", testSalesInsertWhitelist)
	t.Run("SaleItems", testSaleItemsInsert)
	t.Run("SaleItems", testSaleItemsInsertWhitelist)
	t.Run("Transactions", testTransactionsInsert)
	t.Run("Transactions", testTransactionsInsertWhitelist)
	t.Run("Users", testUsersInsert)
	t.Run("Users", testUsersInsertWhitelist)
}

// TestToOne tests cannot be run in parallel
// or deadlocks can occur.
func TestToOne(t *testing.T) {
	t.Run("AccountToBranchUsingBranch", testAccountToOneBranchUsingBranch)
	t.Run("AccountToCustomerUsingCustomer", testAccountToOneCustomerUsingCustomer)
	t.Run("AccountToUserUsingSalesRep", testAccountToOneUserUsingSalesRep)
	t.Run("BankDepositToBankAccountUsingBankAccount", testBankDepositToOneBankAccountUsingBankAccount)
	t.Run("CustomerToBranchUsingBranch", testCustomerToOneBranchUsingBranch)
	t.Run("CustomerToUserUsingSalesRep", testCustomerToOneUserUsingSalesRep)
	t.Run("DSCommissionToAccountUsingAccount", testDSCommissionToOneAccountUsingAccount)
	t.Run("DSCommissionToCustomerUsingCustomer", testDSCommissionToOneCustomerUsingCustomer)
	t.Run("InventoryToBranchUsingBranch", testInventoryToOneBranchUsingBranch)
	t.Run("InventoryToProductUsingProduct", testInventoryToOneProductUsingProduct)
	t.Run("InventoryToUserUsingSalesRep", testInventoryToOneUserUsingSalesRep)
	t.Run("PaymentToSaleUsingSale", testPaymentToOneSaleUsingSale)
	t.Run("PaymentToUserUsingSalesRep", testPaymentToOneUserUsingSalesRep)
	t.Run("ProductToUserUsingArchivedBy", testProductToOneUserUsingArchivedBy)
	t.Run("ProductToBrandUsingBrand", testProductToOneBrandUsingBrand)
	t.Run("ProductToCategoryUsingCategory", testProductToOneCategoryUsingCategory)
	t.Run("ProductToUserUsingCreatedBy", testProductToOneUserUsingCreatedBy)
	t.Run("ProductToUserUsingUpdatedBy", testProductToOneUserUsingUpdatedBy)
	t.Run("ProductCategoryToCategoryUsingCategory", testProductCategoryToOneCategoryUsingCategory)
	t.Run("ProductCategoryToProductUsingProduct", testProductCategoryToOneProductUsingProduct)
	t.Run("RepsExpenseToUserUsingSalesRep", testRepsExpenseToOneUserUsingSalesRep)
	t.Run("SaleToUserUsingArchivedBy", testSaleToOneUserUsingArchivedBy)
	t.Run("SaleToBranchUsingBranch", testSaleToOneBranchUsingBranch)
	t.Run("SaleToUserUsingCreatedBy", testSaleToOneUserUsingCreatedBy)
	t.Run("SaleToUserUsingUpdatedBy", testSaleToOneUserUsingUpdatedBy)
	t.Run("SaleItemToProductUsingProduct", testSaleItemToOneProductUsingProduct)
	t.Run("SaleItemToSaleUsingSale", testSaleItemToOneSaleUsingSale)
	t.Run("TransactionToAccountUsingAccount", testTransactionToOneAccountUsingAccount)
	t.Run("TransactionToUserUsingSalesRep", testTransactionToOneUserUsingSalesRep)
	t.Run("UserToBranchUsingBranch", testUserToOneBranchUsingBranch)
}

// TestOneToOne tests cannot be run in parallel
// or deadlocks can occur.
func TestOneToOne(t *testing.T) {}

// TestToMany tests cannot be run in parallel
// or deadlocks can occur.
func TestToMany(t *testing.T) {
	t.Run("AccountToDSCommissions", testAccountToManyDSCommissions)
	t.Run("AccountToTransactions", testAccountToManyTransactions)
	t.Run("BankAccountToBankDeposits", testBankAccountToManyBankDeposits)
	t.Run("BranchToAccounts", testBranchToManyAccounts)
	t.Run("BranchToCustomers", testBranchToManyCustomers)
	t.Run("BranchToInventories", testBranchToManyInventories)
	t.Run("BranchToSales", testBranchToManySales)
	t.Run("BranchToUsers", testBranchToManyUsers)
	t.Run("BrandToProducts", testBrandToManyProducts)
	t.Run("CategoryToProducts", testCategoryToManyProducts)
	t.Run("CategoryToProductCategories", testCategoryToManyProductCategories)
	t.Run("CustomerToAccounts", testCustomerToManyAccounts)
	t.Run("CustomerToDSCommissions", testCustomerToManyDSCommissions)
	t.Run("ProductToInventories", testProductToManyInventories)
	t.Run("ProductToProductCategories", testProductToManyProductCategories)
	t.Run("ProductToSaleItems", testProductToManySaleItems)
	t.Run("SaleToPayments", testSaleToManyPayments)
	t.Run("SaleToSaleItems", testSaleToManySaleItems)
	t.Run("UserToSalesRepAccounts", testUserToManySalesRepAccounts)
	t.Run("UserToSalesRepCustomers", testUserToManySalesRepCustomers)
	t.Run("UserToSalesRepInventories", testUserToManySalesRepInventories)
	t.Run("UserToSalesRepPayments", testUserToManySalesRepPayments)
	t.Run("UserToArchivedByProducts", testUserToManyArchivedByProducts)
	t.Run("UserToCreatedByProducts", testUserToManyCreatedByProducts)
	t.Run("UserToUpdatedByProducts", testUserToManyUpdatedByProducts)
	t.Run("UserToSalesRepRepsExpenses", testUserToManySalesRepRepsExpenses)
	t.Run("UserToArchivedBySales", testUserToManyArchivedBySales)
	t.Run("UserToCreatedBySales", testUserToManyCreatedBySales)
	t.Run("UserToUpdatedBySales", testUserToManyUpdatedBySales)
	t.Run("UserToSalesRepTransactions", testUserToManySalesRepTransactions)
}

// TestToOneSet tests cannot be run in parallel
// or deadlocks can occur.
func TestToOneSet(t *testing.T) {
	t.Run("AccountToBranchUsingAccounts", testAccountToOneSetOpBranchUsingBranch)
	t.Run("AccountToCustomerUsingAccounts", testAccountToOneSetOpCustomerUsingCustomer)
	t.Run("AccountToUserUsingSalesRepAccounts", testAccountToOneSetOpUserUsingSalesRep)
	t.Run("BankDepositToBankAccountUsingBankDeposits", testBankDepositToOneSetOpBankAccountUsingBankAccount)
	t.Run("CustomerToBranchUsingCustomers", testCustomerToOneSetOpBranchUsingBranch)
	t.Run("CustomerToUserUsingSalesRepCustomers", testCustomerToOneSetOpUserUsingSalesRep)
	t.Run("DSCommissionToAccountUsingDSCommissions", testDSCommissionToOneSetOpAccountUsingAccount)
	t.Run("DSCommissionToCustomerUsingDSCommissions", testDSCommissionToOneSetOpCustomerUsingCustomer)
	t.Run("InventoryToBranchUsingInventories", testInventoryToOneSetOpBranchUsingBranch)
	t.Run("InventoryToProductUsingInventories", testInventoryToOneSetOpProductUsingProduct)
	t.Run("InventoryToUserUsingSalesRepInventories", testInventoryToOneSetOpUserUsingSalesRep)
	t.Run("PaymentToSaleUsingPayments", testPaymentToOneSetOpSaleUsingSale)
	t.Run("PaymentToUserUsingSalesRepPayments", testPaymentToOneSetOpUserUsingSalesRep)
	t.Run("ProductToUserUsingArchivedByProducts", testProductToOneSetOpUserUsingArchivedBy)
	t.Run("ProductToBrandUsingProducts", testProductToOneSetOpBrandUsingBrand)
	t.Run("ProductToCategoryUsingProducts", testProductToOneSetOpCategoryUsingCategory)
	t.Run("ProductToUserUsingCreatedByProducts", testProductToOneSetOpUserUsingCreatedBy)
	t.Run("ProductToUserUsingUpdatedByProducts", testProductToOneSetOpUserUsingUpdatedBy)
	t.Run("ProductCategoryToCategoryUsingProductCategories", testProductCategoryToOneSetOpCategoryUsingCategory)
	t.Run("ProductCategoryToProductUsingProductCategories", testProductCategoryToOneSetOpProductUsingProduct)
	t.Run("RepsExpenseToUserUsingSalesRepRepsExpenses", testRepsExpenseToOneSetOpUserUsingSalesRep)
	t.Run("SaleToUserUsingArchivedBySales", testSaleToOneSetOpUserUsingArchivedBy)
	t.Run("SaleToBranchUsingSales", testSaleToOneSetOpBranchUsingBranch)
	t.Run("SaleToUserUsingCreatedBySales", testSaleToOneSetOpUserUsingCreatedBy)
	t.Run("SaleToUserUsingUpdatedBySales", testSaleToOneSetOpUserUsingUpdatedBy)
	t.Run("SaleItemToProductUsingSaleItems", testSaleItemToOneSetOpProductUsingProduct)
	t.Run("SaleItemToSaleUsingSaleItems", testSaleItemToOneSetOpSaleUsingSale)
	t.Run("TransactionToAccountUsingTransactions", testTransactionToOneSetOpAccountUsingAccount)
	t.Run("TransactionToUserUsingSalesRepTransactions", testTransactionToOneSetOpUserUsingSalesRep)
	t.Run("UserToBranchUsingUsers", testUserToOneSetOpBranchUsingBranch)
}

// TestToOneRemove tests cannot be run in parallel
// or deadlocks can occur.
func TestToOneRemove(t *testing.T) {
	t.Run("ProductToUserUsingArchivedByProducts", testProductToOneRemoveOpUserUsingArchivedBy)
	t.Run("ProductToBrandUsingProducts", testProductToOneRemoveOpBrandUsingBrand)
	t.Run("SaleToUserUsingArchivedBySales", testSaleToOneRemoveOpUserUsingArchivedBy)
	t.Run("SaleToUserUsingUpdatedBySales", testSaleToOneRemoveOpUserUsingUpdatedBy)
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
	t.Run("AccountToDSCommissions", testAccountToManyAddOpDSCommissions)
	t.Run("AccountToTransactions", testAccountToManyAddOpTransactions)
	t.Run("BankAccountToBankDeposits", testBankAccountToManyAddOpBankDeposits)
	t.Run("BranchToAccounts", testBranchToManyAddOpAccounts)
	t.Run("BranchToCustomers", testBranchToManyAddOpCustomers)
	t.Run("BranchToInventories", testBranchToManyAddOpInventories)
	t.Run("BranchToSales", testBranchToManyAddOpSales)
	t.Run("BranchToUsers", testBranchToManyAddOpUsers)
	t.Run("BrandToProducts", testBrandToManyAddOpProducts)
	t.Run("CategoryToProducts", testCategoryToManyAddOpProducts)
	t.Run("CategoryToProductCategories", testCategoryToManyAddOpProductCategories)
	t.Run("CustomerToAccounts", testCustomerToManyAddOpAccounts)
	t.Run("CustomerToDSCommissions", testCustomerToManyAddOpDSCommissions)
	t.Run("ProductToInventories", testProductToManyAddOpInventories)
	t.Run("ProductToProductCategories", testProductToManyAddOpProductCategories)
	t.Run("ProductToSaleItems", testProductToManyAddOpSaleItems)
	t.Run("SaleToPayments", testSaleToManyAddOpPayments)
	t.Run("SaleToSaleItems", testSaleToManyAddOpSaleItems)
	t.Run("UserToSalesRepAccounts", testUserToManyAddOpSalesRepAccounts)
	t.Run("UserToSalesRepCustomers", testUserToManyAddOpSalesRepCustomers)
	t.Run("UserToSalesRepInventories", testUserToManyAddOpSalesRepInventories)
	t.Run("UserToSalesRepPayments", testUserToManyAddOpSalesRepPayments)
	t.Run("UserToArchivedByProducts", testUserToManyAddOpArchivedByProducts)
	t.Run("UserToCreatedByProducts", testUserToManyAddOpCreatedByProducts)
	t.Run("UserToUpdatedByProducts", testUserToManyAddOpUpdatedByProducts)
	t.Run("UserToSalesRepRepsExpenses", testUserToManyAddOpSalesRepRepsExpenses)
	t.Run("UserToArchivedBySales", testUserToManyAddOpArchivedBySales)
	t.Run("UserToCreatedBySales", testUserToManyAddOpCreatedBySales)
	t.Run("UserToUpdatedBySales", testUserToManyAddOpUpdatedBySales)
	t.Run("UserToSalesRepTransactions", testUserToManyAddOpSalesRepTransactions)
}

// TestToManySet tests cannot be run in parallel
// or deadlocks can occur.
func TestToManySet(t *testing.T) {
	t.Run("BrandToProducts", testBrandToManySetOpProducts)
	t.Run("UserToArchivedByProducts", testUserToManySetOpArchivedByProducts)
	t.Run("UserToArchivedBySales", testUserToManySetOpArchivedBySales)
	t.Run("UserToUpdatedBySales", testUserToManySetOpUpdatedBySales)
}

// TestToManyRemove tests cannot be run in parallel
// or deadlocks can occur.
func TestToManyRemove(t *testing.T) {
	t.Run("BrandToProducts", testBrandToManyRemoveOpProducts)
	t.Run("UserToArchivedByProducts", testUserToManyRemoveOpArchivedByProducts)
	t.Run("UserToArchivedBySales", testUserToManyRemoveOpArchivedBySales)
	t.Run("UserToUpdatedBySales", testUserToManyRemoveOpUpdatedBySales)
}

func TestReload(t *testing.T) {
	t.Run("Accounts", testAccountsReload)
	t.Run("BankAccounts", testBankAccountsReload)
	t.Run("BankDeposits", testBankDepositsReload)
	t.Run("Branches", testBranchesReload)
	t.Run("Brands", testBrandsReload)
	t.Run("Categories", testCategoriesReload)
	t.Run("Customers", testCustomersReload)
	t.Run("DailySummaries", testDailySummariesReload)
	t.Run("DSCommissions", testDSCommissionsReload)
	t.Run("Expenditures", testExpendituresReload)
	t.Run("Inventories", testInventoriesReload)
	t.Run("Payments", testPaymentsReload)
	t.Run("Products", testProductsReload)
	t.Run("ProductCategories", testProductCategoriesReload)
	t.Run("Profits", testProfitsReload)
	t.Run("RepsExpenses", testRepsExpensesReload)
	t.Run("Sales", testSalesReload)
	t.Run("SaleItems", testSaleItemsReload)
	t.Run("Transactions", testTransactionsReload)
	t.Run("Users", testUsersReload)
}

func TestReloadAll(t *testing.T) {
	t.Run("Accounts", testAccountsReloadAll)
	t.Run("BankAccounts", testBankAccountsReloadAll)
	t.Run("BankDeposits", testBankDepositsReloadAll)
	t.Run("Branches", testBranchesReloadAll)
	t.Run("Brands", testBrandsReloadAll)
	t.Run("Categories", testCategoriesReloadAll)
	t.Run("Customers", testCustomersReloadAll)
	t.Run("DailySummaries", testDailySummariesReloadAll)
	t.Run("DSCommissions", testDSCommissionsReloadAll)
	t.Run("Expenditures", testExpendituresReloadAll)
	t.Run("Inventories", testInventoriesReloadAll)
	t.Run("Payments", testPaymentsReloadAll)
	t.Run("Products", testProductsReloadAll)
	t.Run("ProductCategories", testProductCategoriesReloadAll)
	t.Run("Profits", testProfitsReloadAll)
	t.Run("RepsExpenses", testRepsExpensesReloadAll)
	t.Run("Sales", testSalesReloadAll)
	t.Run("SaleItems", testSaleItemsReloadAll)
	t.Run("Transactions", testTransactionsReloadAll)
	t.Run("Users", testUsersReloadAll)
}

func TestSelect(t *testing.T) {
	t.Run("Accounts", testAccountsSelect)
	t.Run("BankAccounts", testBankAccountsSelect)
	t.Run("BankDeposits", testBankDepositsSelect)
	t.Run("Branches", testBranchesSelect)
	t.Run("Brands", testBrandsSelect)
	t.Run("Categories", testCategoriesSelect)
	t.Run("Customers", testCustomersSelect)
	t.Run("DailySummaries", testDailySummariesSelect)
	t.Run("DSCommissions", testDSCommissionsSelect)
	t.Run("Expenditures", testExpendituresSelect)
	t.Run("Inventories", testInventoriesSelect)
	t.Run("Payments", testPaymentsSelect)
	t.Run("Products", testProductsSelect)
	t.Run("ProductCategories", testProductCategoriesSelect)
	t.Run("Profits", testProfitsSelect)
	t.Run("RepsExpenses", testRepsExpensesSelect)
	t.Run("Sales", testSalesSelect)
	t.Run("SaleItems", testSaleItemsSelect)
	t.Run("Transactions", testTransactionsSelect)
	t.Run("Users", testUsersSelect)
}

func TestUpdate(t *testing.T) {
	t.Run("Accounts", testAccountsUpdate)
	t.Run("BankAccounts", testBankAccountsUpdate)
	t.Run("BankDeposits", testBankDepositsUpdate)
	t.Run("Branches", testBranchesUpdate)
	t.Run("Brands", testBrandsUpdate)
	t.Run("Categories", testCategoriesUpdate)
	t.Run("Customers", testCustomersUpdate)
	t.Run("DailySummaries", testDailySummariesUpdate)
	t.Run("DSCommissions", testDSCommissionsUpdate)
	t.Run("Expenditures", testExpendituresUpdate)
	t.Run("Inventories", testInventoriesUpdate)
	t.Run("Payments", testPaymentsUpdate)
	t.Run("Products", testProductsUpdate)
	t.Run("ProductCategories", testProductCategoriesUpdate)
	t.Run("Profits", testProfitsUpdate)
	t.Run("RepsExpenses", testRepsExpensesUpdate)
	t.Run("Sales", testSalesUpdate)
	t.Run("SaleItems", testSaleItemsUpdate)
	t.Run("Transactions", testTransactionsUpdate)
	t.Run("Users", testUsersUpdate)
}

func TestSliceUpdateAll(t *testing.T) {
	t.Run("Accounts", testAccountsSliceUpdateAll)
	t.Run("BankAccounts", testBankAccountsSliceUpdateAll)
	t.Run("BankDeposits", testBankDepositsSliceUpdateAll)
	t.Run("Branches", testBranchesSliceUpdateAll)
	t.Run("Brands", testBrandsSliceUpdateAll)
	t.Run("Categories", testCategoriesSliceUpdateAll)
	t.Run("Customers", testCustomersSliceUpdateAll)
	t.Run("DailySummaries", testDailySummariesSliceUpdateAll)
	t.Run("DSCommissions", testDSCommissionsSliceUpdateAll)
	t.Run("Expenditures", testExpendituresSliceUpdateAll)
	t.Run("Inventories", testInventoriesSliceUpdateAll)
	t.Run("Payments", testPaymentsSliceUpdateAll)
	t.Run("Products", testProductsSliceUpdateAll)
	t.Run("ProductCategories", testProductCategoriesSliceUpdateAll)
	t.Run("Profits", testProfitsSliceUpdateAll)
	t.Run("RepsExpenses", testRepsExpensesSliceUpdateAll)
	t.Run("Sales", testSalesSliceUpdateAll)
	t.Run("SaleItems", testSaleItemsSliceUpdateAll)
	t.Run("Transactions", testTransactionsSliceUpdateAll)
	t.Run("Users", testUsersSliceUpdateAll)
}
