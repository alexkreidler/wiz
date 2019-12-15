package taguk

// import "testing"

import (
	// "github.com/alexkreidler/wiz/taguk"
	"testing"
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

func Test(t *testing.T) {

	s := NewServer()
	s.AddResources(Bank{}, Branch{}, User{}, Account{})
}
