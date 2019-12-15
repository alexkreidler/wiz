package httpjson

import (
	"fmt"
	"github.com/alexkreidler/wiz/taguk"
	"github.com/alexkreidler/wiz/taguk/test"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCRUD(t *testing.T) {
	api := taguk.NewServer()
	api.AddResources(test.Branch{})
	s := NewServer()
	s.Configure(api.Resources)

	r := s.Router

	testCases := []struct {
		desc	string
		method string
		path string
	}{
		{
			desc: "Create Branch",
			method: "POST",
			path: "/branch",
		},
		{
			desc: "Get All Branches",
			method: "GET",
			path: "/branch",
		},
		{
			desc: "Read Single Branch",
			method: "GET",
			path: "/branch/1",
		},
		{
			desc: "Update Single Branch",
			method: "POST",
			path: "/branch/1",
		},
		{
			desc: "Delete Single Branch",
			method: "DELETE",
			path: "/branch/1",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {		
			req, err := http.NewRequest(tC.method, tC.path, nil)
			test.Ok(t, err)

			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)
		
			// Check the status code is what we expect.
			if status := rr.Code; status != http.StatusOK {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, http.StatusOK)
			}
			fmt.Println(rr.Body.String())
		})
	}
}
func TestServe(t *testing.T) {
	api := taguk.NewServer()
	api.AddResources(test.Branch{})
	s := NewServer()
	s.Configure(api.Resources)

	r := s.Router

	req, err := http.NewRequest("GET", "/branch/12", nil)
	test.Ok(t, err)

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()

	// Our router satisfies http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	r.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	fmt.Println(rr.Body.String())
}
