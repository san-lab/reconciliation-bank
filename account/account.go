package account

import (
	"crypto/sha256"
	"errors"
	"fmt"

	. "github.com/san-lab/immudb-tests/datastructs"
)

type Account struct {
	Suspended bool
	Bic       string
	Iban      string
	Balance   float32
	Holder    string
	Currency  string
	IsCA      bool
	IsMirror  bool
	CABank    string
}

// Opens a new bank account with the specified parameters
func SetAccount(bic, iban, holder, currency, cABank string, balance float32, isCA, isMirror bool) *Account {
	return &Account{
		Suspended: false,
		Bic:       bic,
		Iban:      iban,
		Balance:   balance,
		Holder:    holder,
		Currency:  currency,
		IsCA:      isCA,
		IsMirror:  isMirror,
		CABank:    cABank,
	}
}

// Suspends an account
func (account *Account) Suspend() {
	account.Suspended = true
}

func (account *Account) Unsuspend() {
	account.Suspended = false
}

// Deposits an amount into an account
func (account *Account) Deposit(amount float32) error {
	if !account.Suspended {
		account.Balance += amount
	} else {
		return errors.New("account suspended: operation not performed")
	}
	return nil
}

// Withdraws an amount from an account
func (account *Account) Withdraw(amount float32) error {
	if !account.Suspended {
		// Will allow negative balances for a moment
		// if amount <= account.Balance {
		account.Balance -= amount
		// } else {
		//	  return errors.New("balance cannot be negative")
		// }
	} else {
		return errors.New("account suspended: operation not performed")
	}
	return nil
}

// Returns current balance of an account
func (account *Account) GetBalance() float32 {
	return account.Balance
}

func (account *Account) SetBalance(newBalance float32) error {
	if newBalance < 0 {
		return errors.New("balance cannot be negative")
	}
	account.Balance = newBalance
	return nil
}

func (account *Account) GetIsCA() bool {
	return account.IsCA
}

func (account *Account) GetIsMirror() bool {
	return account.IsMirror
}

func (account *Account) GetCABank() (string, error) {
	if !(account.IsCA || account.IsMirror) {
		return "", errors.New("account is not CA or Mirror")
	}
	return account.CABank, nil
}

func (account *Account) GetDigest() (string, error) {
	// TODO: add more fields to the digest (making sure both banks have will have the same values!)
	fields := fmt.Sprintf("%f", account.Balance)
	fmt.Println("debug fields:", fields)

	sum := sha256.Sum256([]byte(fields))
	return fmt.Sprintf("%x", sum), nil
}

// Print account details
func (account *Account) PrintAccount(spacing bool) {
	holder := account.Holder
	if account.IsCA {
		holder = account.CABank + " - CA"
	} else if account.IsMirror {
		holder = THIS_BANK.Name + "@" + account.CABank + " - Mirror"
	}
	if spacing {
		fmt.Println(" -----------------")
		fmt.Printf("| IBAN: %s\n| Holder: %s\n| Balance: %.2f\n| Currency: %s\n| BIC: %s\n| Suspended: %t\n",
			account.Iban, holder, account.Balance, account.Currency, account.Bic, account.Suspended)
		fmt.Println(" -----------------")
	} else {
		fmt.Printf("| IBAN: %s | Holder: %s | Balance: %.2f | Currency: %s | BIC: %s | Suspended: %t\n",
			account.Iban, holder, account.Balance, account.Currency, account.Bic, account.Suspended)
	}
}

func PrintAllAccounts(accounts []*Account) {
	fmt.Println(" -----------------")
	for _, account := range accounts {
		account.PrintAccount(false)
	}
	fmt.Println(" -----------------")
}
