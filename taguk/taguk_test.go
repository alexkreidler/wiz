package taguk_test

// import "testing"

import (
	. "github.com/alexkreidler/wiz/taguk"
	. "github.com/alexkreidler/wiz/taguk/test"
	"testing"
)

func TestAddResources(t *testing.T) {
	s := NewServer()
	s.AddResources(Bank{}, Branch{}, User{}, Account{})
}
