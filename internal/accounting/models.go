package accounting

import (
	"github.com/volatiletech/null"
)

// CreateBankAccount holds req for creating a new bank account
type CreateBankAccount struct {
	Bank   string `json:"bank" validate:"required"`
	Name   string `json:"name" validate:"required"`
	Number string `json:"number" validate:"required"`
}

// CreateBankDeposit holds req for creating a new bank deposit
type CreateBankDeposit struct {
	BankID string  `json:"bank_id" validate:"required"`
	Amount float64 `json:"amount" validate:"required"`
}

// CreateExpenditure holds req for creating a new bank expenditure
type CreateExpenditure struct {
	Amount float64 `json:"amount" validate:"required"`
	Info   string  `json:"info"`
}

type RepsSummary struct {
	SalesRepID  string           `json:"sales_rep_id"`
	SalesRep    string           `json:"sales_rep"`
	Income      null.Float64     `json:"income"`
	Expenditure null.Float64     `json:"expenditure"`
}
