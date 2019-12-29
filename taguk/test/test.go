// Test exports various example structs and methods that define a sample API
// They are used by Taguk testing packages to test Taguk
package test

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

func (b Branch) GetAll() []Branch {
	return []Branch{
		Branch{
			Name:     "test",
			Location: "Chevy Chase",
		},
		Branch{
			Name:     "test2",
			Location: "California",
		},
	}
}

func (b Branch) Get(id int64) Branch {
	return Branch{
		Name:     "test",
		Location: "Chevy Chase",
	}
}

func (b Branch) NoArgs() Branch {
	return Branch{
		Name:     "test",
		Location: "Chevy Chase",
	}
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
