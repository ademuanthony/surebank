package models

import "fmt"

func (m BankAccount) String() string {
	return fmt.Sprintf("%s (%s) - %s", m.AccountName, m.AccountNumber, m.Bank)
}

