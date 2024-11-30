package main

// Account represents a user or customer in the system.
type Account string

// NewAccount creates a new Account.
func NewAccount(name string) Account {
	return Account(name)
}
