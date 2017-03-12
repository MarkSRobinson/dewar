package caledonia

import (
	"io"
	)

type Account struct {
	AccountId string
	Name string
	Mask string
	Type string
	Subtype string
}

type Transaction struct {
	TransactionType string
	Posted		string
	Amount		string
	Fitid		string
	Name 		string
	Memo 		string
}

type OfxWriter interface {
	AddTransaction(accountId string, transaction Transaction)
	AddAccount(account Account)

	Write(writer io.Writer)
}
