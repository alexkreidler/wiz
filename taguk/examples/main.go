package main

import (
	"github.com/alexkreidler/wiz/taguk"
)

// Bank represents a bank
type Bank struct {
	Name     string
	Users    []User
	Accounts []*Account
	Branches []Branch
}

// Branch represents a local, physical branch of a bank
type Branch struct {
	Name     string
	Location string
}

// User represents a single bank customer
type User struct {
	Name     string
	Age      int
	Accounts []*Account
}

// Account represents a bank account and can be owned by multiple users
type Account struct {
	Name    string
	Balance float32
	Owners  []*User
}

// TransactionStream is a client-to-server stream of transactions
type TransactionStream struct {
	ToServer <-chan Transaction
	// ToClient chan Transaction
}

// OnStream handles a new stream
func (t TransactionStream) OnStream() {
	// Do stuff with the stream
	// t.ToServer
}

// Transaction represents a change in balance on an account
type Transaction struct {
	PreviousBalance float32
	NewBalance      float32
	Date            string
	Name            string
}

func main() {
	s := taguk.NewServer()
	s.AddResources(Bank{}, Branch{}, User{}, Account{})

	// Serves an HTTP API on port, options allow serving a schema
	s.ServeHTTP(taguk.JSONSchema)

	// Serves gRPC API
	s.ServeGRPC()
}
