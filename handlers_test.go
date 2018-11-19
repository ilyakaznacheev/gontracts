package gontracts

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"bitbucket.org/ilyakaznacheev/gontracts/model"
	"bitbucket.org/ilyakaznacheev/gontracts/test"
	"github.com/gorilla/mux"
)

type testModelSet struct {
	company  model.CompanyModel
	contract model.ContractModel
	purchase model.PurchaseModel
}

func testStrPtr(s string) *string {
	return &s
}

func testHandle(path string, w http.ResponseWriter, req *http.Request, h func(http.ResponseWriter, *http.Request)) {
	r := mux.NewRouter()
	r.HandleFunc(path, h).Methods(req.Method)
	r.ServeHTTP(w, req)
}

func testNewHandler(company model.CompanyModel, contract model.ContractModel, purchase model.PurchaseModel) *Handler {
	return &Handler{
		mh:         model.NewModelHandler(company, contract, purchase),
		purchaseMX: &sync.Mutex{},
	}
}

func testCheckResponse(loc string, t *testing.T, w *httptest.ResponseRecorder, respStatus int, respBody string) {
	if w.Code != respStatus {
		t.Errorf("[%s]:\twrong StatusCode: got %d, expected %d",
			loc, w.Code, respStatus)
	}

	body, _ := ioutil.ReadAll(w.Result().Body)
	bodyStr := string(body)
	if bodyStr != respBody {
		t.Errorf("[%s]:\twrong Response: got %s, expected %s",
			loc, bodyStr, respBody)
	}
}
func TestGetCompany(t *testing.T) {
	company := test.TestCompany{
		CL: []*model.Company{
			{1, "test1", testStrPtr("123")},
			{2, "test2", nil},
		},
	}

	cases := []struct {
		Num      string
		ID       int
		Response string
		Status   int
		Models   testModelSet
	}{
		{
			Num:      "1",
			ID:       1,
			Response: `{"ID":1,"name":"test1","regcode":"123"}`,
			Status:   http.StatusOK,
			Models: testModelSet{
				company: company,
			},
		},
		{
			Num:      "2",
			ID:       2,
			Response: `{"ID":2,"name":"test2","regcode":null}`,
			Status:   http.StatusOK,
			Models: testModelSet{
				company: company,
			},
		},
		{
			Num:      "3",
			ID:       3,
			Response: test.ErrTest.Error() + "\n",
			Status:   http.StatusNotFound,
			Models: testModelSet{
				company: company,
			},
		},
	}

	for _, c := range cases {
		h := testNewHandler(c.Models.company, c.Models.contract, c.Models.purchase)

		url := fmt.Sprintf("/company/%d", c.ID)
		req := httptest.NewRequest("GET", url, nil)
		w := httptest.NewRecorder()

		testHandle("/company/{id:[0-9]+}", w, req, h.GetCompany)
		testCheckResponse("GetCompany:"+c.Num, t, w, c.Status, c.Response)
	}
}

func TestGetCompanyList(t *testing.T) {
	company := test.TestCompany{
		CL: []*model.Company{
			{1, "test1", testStrPtr("123")},
			{2, "test2", nil},
		},
	}
	companyError := test.TestCompanyErr{}

	cases := []struct {
		Num      string
		Response string
		Status   int
		Models   testModelSet
	}{
		{
			Num:      "1",
			Response: `[{"ID":1,"name":"test1","regcode":"123"},{"ID":2,"name":"test2","regcode":null}]`,
			Status:   http.StatusOK,
			Models: testModelSet{
				company: company,
			},
		},
		{
			Num:      "2",
			Response: test.ErrTest.Error() + "\n",
			Status:   http.StatusNotFound,
			Models: testModelSet{
				company: companyError,
			},
		},
	}

	for _, c := range cases {
		h := testNewHandler(c.Models.company, c.Models.contract, c.Models.purchase)

		url := "/companies"
		req := httptest.NewRequest("GET", url, nil)
		w := httptest.NewRecorder()

		testHandle("/companies", w, req, h.GetCompanyList)
		testCheckResponse("GetCompanyList:"+c.Num, t, w, c.Status, c.Response)
	}
}

func TestCreateCompany(t *testing.T) {
	cases := []struct {
		Num      string
		Request  string
		Response string
		Status   int
		Models   testModelSet
	}{
		// normal creation with null reg code
		{
			Num:      "1",
			Request:  `{"name":"test2","regcode":null}`,
			Response: `{"ID":1}`,
			Status:   http.StatusCreated,
			Models: testModelSet{
				company: test.TestCompany{
					CL: []*model.Company{},
				},
			},
		},
		// normal creation with filled reg code
		{
			Num:      "2",
			Request:  `{"name":"test2","regcode":"123"}`,
			Response: `{"ID":1}`,
			Status:   http.StatusCreated,
			Models: testModelSet{
				company: test.TestCompany{
					CL: []*model.Company{},
				},
			},
		},
		// normal creation without reg code
		{
			Num:      "3",
			Request:  `{"name":"test2"}`,
			Response: `{"ID":1}`,
			Status:   http.StatusCreated,
			Models: testModelSet{
				company: test.TestCompany{
					CL: []*model.Company{},
				},
			},
		},
		// normal creation when some items exist
		{
			Num:      "4",
			Request:  `{"name":"test2","regcode":null}`,
			Response: `{"ID":2}`,
			Status:   http.StatusCreated,
			Models: testModelSet{
				company: test.TestCompany{
					CL: []*model.Company{
						{ID: 1, Name: "test"},
					},
				},
			},
		},
		// error handling
		{
			Num:      "5",
			Request:  `{"name":"test2","regcode":null}`,
			Response: test.ErrTest.Error() + "\n",
			Status:   http.StatusInternalServerError,
			Models: testModelSet{
				company: test.TestCompanyErr{},
			},
		},
	}

	for _, c := range cases {
		h := testNewHandler(c.Models.company, c.Models.contract, c.Models.purchase)

		url := "/company"
		req := httptest.NewRequest("POST", url, bytes.NewBuffer([]byte(c.Request)))
		w := httptest.NewRecorder()

		testHandle("/company", w, req, h.CreateCompany)
		testCheckResponse("CreateCompany:"+c.Num, t, w, c.Status, c.Response)
	}
}

func TestUpdateCompany(t *testing.T) {
	cases := []struct {
		Num      string
		Request  string
		Response string
		Status   int
		Models   testModelSet
	}{
		// normal creation
		{
			Num:      "1",
			Request:  `{"name":"test2","regcode":null}`,
			Response: `{"ID":1,"name":"test2","regcode":null}`,
			Status:   http.StatusCreated,
			Models: testModelSet{
				company: test.TestCompany{
					CL: []*model.Company{},
				},
			},
		},
		// normal update
		{
			Num:      "2",
			Request:  `{"ID":1,"name":"test2","regcode":null}`,
			Response: `{"ID":1,"name":"test2","regcode":null}`,
			Status:   http.StatusOK,
			Models: testModelSet{
				company: test.TestCompany{
					CL: []*model.Company{
						{ID: 1, Name: "test"},
					},
				},
			},
		},
		// wrong id
		{
			Num:      "3",
			Request:  `{"ID":1,"name":"test2","regcode":null}`,
			Response: test.ErrTest.Error() + "\n",
			Status:   http.StatusInternalServerError,
			Models: testModelSet{
				company: test.TestCompany{
					CL: []*model.Company{},
				},
			},
		},
		// wrong id
		{
			Num:      "4",
			Request:  `{"ID":2,"name":"test2","regcode":null}`,
			Response: test.ErrTest.Error() + "\n",
			Status:   http.StatusInternalServerError,
			Models: testModelSet{
				company: test.TestCompany{
					CL: []*model.Company{
						{ID: 1, Name: "test"},
					},
				},
			},
		},
		// error handling
		{
			Num:      "5",
			Request:  `{"name":"test2","regcode":null}`,
			Response: test.ErrTest.Error() + "\n",
			Status:   http.StatusInternalServerError,
			Models: testModelSet{
				company: test.TestCompanyErr{},
			},
		},
	}

	for _, c := range cases {
		h := testNewHandler(c.Models.company, c.Models.contract, c.Models.purchase)

		url := "/company"
		req := httptest.NewRequest("PUT", url, bytes.NewBuffer([]byte(c.Request)))
		w := httptest.NewRecorder()

		testHandle("/company", w, req, h.UpdateCompany)
		testCheckResponse("UpdateCompany:"+c.Num, t, w, c.Status, c.Response)
	}
}

func TestDeleteCompany(t *testing.T) {
	cases := []struct {
		Num      string
		ID       int
		Response string
		Status   int
		Models   testModelSet
	}{
		// ok
		{
			Num:      "1",
			ID:       1,
			Response: http.StatusText(http.StatusOK),
			Status:   http.StatusOK,
			Models: testModelSet{
				company: test.TestCompany{
					CL: []*model.Company{
						{1, "test", testStrPtr("123")},
					},
				},
			},
		},
		// wrong id
		{
			Num:      "2",
			ID:       2,
			Response: test.ErrTest.Error() + "\n",
			Status:   http.StatusNotFound,
			Models: testModelSet{
				company: test.TestCompany{
					CL: []*model.Company{
						{1, "test", testStrPtr("123")},
					},
				},
			},
		},
		// not exist
		{
			Num:      "3",
			ID:       2,
			Response: test.ErrTest.Error() + "\n",
			Status:   http.StatusNotFound,
			Models: testModelSet{
				company: test.TestCompany{
					CL: []*model.Company{},
				},
			},
		},
		// error handling
		{
			Num:      "4",
			ID:       3,
			Response: test.ErrTest.Error() + "\n",
			Status:   http.StatusNotFound,
			Models: testModelSet{
				company: test.TestCompanyErr{},
			},
		},
	}

	for _, c := range cases {
		h := testNewHandler(c.Models.company, c.Models.contract, c.Models.purchase)

		url := fmt.Sprintf("/company/%d", c.ID)
		req := httptest.NewRequest("DELETE", url, nil)
		w := httptest.NewRecorder()

		testHandle("/company/{id:[0-9]+}", w, req, h.DeleteCompany)
		testCheckResponse("DeleteCompany:"+c.Num, t, w, c.Status, c.Response)
	}
}

func TestGetContract(t *testing.T) {
	time1 := time.Date(2000, 01, 01, 00, 00, 00, 0, time.UTC)
	contract := test.TestContract{
		CL: []*model.Contract{
			{1, 10, 11, time1, time1.AddDate(1, 0, 0), 10},
			{2, 20, 21, time1, time1.AddDate(1, 0, 0), 100},
		},
	}

	cases := []struct {
		Num      string
		ID       int
		Response string
		Status   int
		Models   testModelSet
	}{
		{
			Num:      "1",
			ID:       1,
			Response: `{"ID":1,"sellerID":10,"clientID":11,"validFrom":"2000-01-01T00:00:00Z","validTo":"2001-01-01T00:00:00Z","amount":10}`,
			Status:   http.StatusOK,
			Models: testModelSet{
				contract: contract,
			},
		},
		{
			Num:      "2",
			ID:       2,
			Response: `{"ID":2,"sellerID":20,"clientID":21,"validFrom":"2000-01-01T00:00:00Z","validTo":"2001-01-01T00:00:00Z","amount":100}`,
			Status:   http.StatusOK,
			Models: testModelSet{
				contract: contract,
			},
		},
		{
			Num:      "3",
			ID:       3,
			Response: test.ErrTest.Error() + "\n",
			Status:   http.StatusNotFound,
			Models: testModelSet{
				contract: contract,
			},
		},
		{
			Num:      "4",
			ID:       1,
			Response: test.ErrTest.Error() + "\n",
			Status:   http.StatusNotFound,
			Models: testModelSet{
				contract: test.TestContractErr{},
			},
		},
	}

	for _, c := range cases {
		h := testNewHandler(c.Models.company, c.Models.contract, c.Models.purchase)

		url := fmt.Sprintf("/contract/%d", c.ID)
		req := httptest.NewRequest("GET", url, nil)
		w := httptest.NewRecorder()

		testHandle("/contract/{id:[0-9]+}", w, req, h.GetContract)
		testCheckResponse("GetContract:"+c.Num, t, w, c.Status, c.Response)
	}
}

func TestGetContractList(t *testing.T) {
	time1 := time.Date(2000, 01, 01, 00, 00, 00, 0, time.UTC)
	contract := test.TestContract{
		CL: []*model.Contract{
			{1, 10, 11, time1, time1.AddDate(1, 0, 0), 10},
			{2, 20, 21, time1, time1.AddDate(1, 0, 0), 100},
		},
	}

	cases := []struct {
		Num      string
		Response string
		Status   int
		Models   testModelSet
	}{
		{
			Num:      "1",
			Response: `[{"ID":1,"sellerID":10,"clientID":11,"validFrom":"2000-01-01T00:00:00Z","validTo":"2001-01-01T00:00:00Z","amount":10},{"ID":2,"sellerID":20,"clientID":21,"validFrom":"2000-01-01T00:00:00Z","validTo":"2001-01-01T00:00:00Z","amount":100}]`,
			Status:   http.StatusOK,
			Models: testModelSet{
				contract: contract,
			},
		},
		{
			Num:      "2",
			Response: test.ErrTest.Error() + "\n",
			Status:   http.StatusNotFound,
			Models: testModelSet{
				contract: test.TestContractErr{},
			},
		},
	}

	for _, c := range cases {
		h := testNewHandler(c.Models.company, c.Models.contract, c.Models.purchase)

		url := "/contracts"
		req := httptest.NewRequest("GET", url, nil)
		w := httptest.NewRecorder()

		testHandle("/contracts", w, req, h.GetContractList)
		testCheckResponse("GetContractList:"+c.Num, t, w, c.Status, c.Response)
	}
}

func TestCreateContract(t *testing.T) {
	cases := []struct {
		Num      string
		Request  string
		Response string
		Status   int
		Models   testModelSet
	}{
		// normal creation
		{
			Num:      "1",
			Request:  `{"ID":1,"sellerID":10,"clientID":11,"validFrom":"2000-01-01T00:00:00Z","validTo":"2001-01-01T00:00:00Z","amount":10}`,
			Response: `{"ID":1}`,
			Status:   http.StatusCreated,
			Models: testModelSet{
				company: test.TestCompany{
					CL: []*model.Company{
						{ID: 10, Name: "test1"},
						{ID: 11, Name: "test1"},
					},
				},
				contract: test.TestContract{
					CL: []*model.Contract{},
				},
			},
		},
		// seller company doesn't exist
		{
			Num:      "2",
			Request:  `{"ID":1,"sellerID":10,"clientID":11,"validFrom":"2000-01-01T00:00:00Z","validTo":"2001-01-01T00:00:00Z","amount":10}`,
			Response: ErrSellerNotExist.Error() + "\n",
			Status:   http.StatusBadRequest,
			Models: testModelSet{
				company: test.TestCompany{
					CL: []*model.Company{
						{ID: 11, Name: "test1"},
					},
				},
				contract: test.TestContract{
					CL: []*model.Contract{},
				},
			},
		},
		// client company doesn't exist
		{
			Num:      "3",
			Request:  `{"ID":1,"sellerID":10,"clientID":11,"validFrom":"2000-01-01T00:00:00Z","validTo":"2001-01-01T00:00:00Z","amount":10}`,
			Response: ErrClientNotExist.Error() + "\n",
			Status:   http.StatusBadRequest,
			Models: testModelSet{
				company: test.TestCompany{
					CL: []*model.Company{
						{ID: 10, Name: "test1"},
					},
				},
				contract: test.TestContract{
					CL: []*model.Contract{},
				},
			},
		},
		// error handling
		{
			Num:      "4",
			Request:  `{"ID":1,"sellerID":10,"clientID":11,"validFrom":"2000-01-01T00:00:00Z","validTo":"2001-01-01T00:00:00Z","amount":10}`,
			Response: test.ErrTest.Error() + "\n",
			Status:   http.StatusInternalServerError,
			Models: testModelSet{
				company: test.TestCompany{
					CL: []*model.Company{
						{ID: 10, Name: "test1"},
						{ID: 11, Name: "test1"},
					},
				},
				contract: test.TestContractErr{},
			},
		},
	}

	for _, c := range cases {
		h := testNewHandler(c.Models.company, c.Models.contract, c.Models.purchase)

		url := "/contract"
		req := httptest.NewRequest("POST", url, bytes.NewBuffer([]byte(c.Request)))
		w := httptest.NewRecorder()

		testHandle("/contract", w, req, h.CreateContract)
		testCheckResponse("CreateContract:"+c.Num, t, w, c.Status, c.Response)
	}
}

func TestUpdateContract(t *testing.T) {
	time1 := time.Date(2000, 01, 01, 00, 00, 00, 0, time.UTC)

	cases := []struct {
		Num      string
		Request  string
		Response string
		Status   int
		Models   testModelSet
	}{
		// creation
		{
			Num:      "1",
			Request:  `{"sellerID":10,"clientID":11,"validFrom":"2000-01-01T00:00:00Z","validTo":"2001-01-01T00:00:00Z","amount":150}`,
			Response: `{"ID":1,"sellerID":10,"clientID":11,"validFrom":"2000-01-01T00:00:00Z","validTo":"2001-01-01T00:00:00Z","amount":150}`,
			Status:   http.StatusCreated,
			Models: testModelSet{
				company: test.TestCompany{
					CL: []*model.Company{
						{ID: 10, Name: "test1"},
						{ID: 11, Name: "test1"},
					},
				},
				contract: test.TestContract{
					CL: []*model.Contract{},
				},
			},
		},
		// update existing one
		{
			Num:      "2",
			Request:  `{"ID": 1, "sellerID":10,"clientID":11,"validFrom":"2000-01-01T00:00:00Z","validTo":"2001-01-01T00:00:00Z","amount":150}`,
			Response: `{"ID":1,"sellerID":10,"clientID":11,"validFrom":"2000-01-01T00:00:00Z","validTo":"2001-01-01T00:00:00Z","amount":150}`,
			Status:   http.StatusOK,
			Models: testModelSet{
				company: test.TestCompany{
					CL: []*model.Company{
						{ID: 10, Name: "test1"},
						{ID: 11, Name: "test1"},
					},
				},
				contract: test.TestContract{
					CL: []*model.Contract{
						{1, 10, 11, time1, time1.AddDate(1, 0, 0), 10},
					},
				},
			},
		},
		// seller company doesn't exist
		{
			Num:      "3",
			Request:  `{"ID":1,"sellerID":10,"clientID":11,"validFrom":"2000-01-01T00:00:00Z","validTo":"2001-01-01T00:00:00Z","amount":10}`,
			Response: ErrSellerNotExist.Error() + "\n",
			Status:   http.StatusBadRequest,
			Models: testModelSet{
				company: test.TestCompany{
					CL: []*model.Company{
						{ID: 11, Name: "test1"},
					},
				},
				contract: test.TestContract{
					CL: []*model.Contract{},
				},
			},
		},
		// client company doesn't exist
		{
			Num:      "4",
			Request:  `{"ID":1,"sellerID":10,"clientID":11,"validFrom":"2000-01-01T00:00:00Z","validTo":"2001-01-01T00:00:00Z","amount":10}`,
			Response: ErrClientNotExist.Error() + "\n",
			Status:   http.StatusBadRequest,
			Models: testModelSet{
				company: test.TestCompany{
					CL: []*model.Company{
						{ID: 10, Name: "test1"},
					},
				},
				contract: test.TestContract{
					CL: []*model.Contract{},
				},
			},
		},
		// error handling
		{
			Num:      "5",
			Request:  `{"ID":1,"sellerID":10,"clientID":11,"validFrom":"2000-01-01T00:00:00Z","validTo":"2001-01-01T00:00:00Z","amount":10}`,
			Response: test.ErrTest.Error() + "\n",
			Status:   http.StatusInternalServerError,
			Models: testModelSet{
				company: test.TestCompany{
					CL: []*model.Company{
						{ID: 10, Name: "test1"},
						{ID: 11, Name: "test1"},
					},
				},
				contract: test.TestContractErr{},
			},
		},
	}

	for _, c := range cases {
		h := testNewHandler(c.Models.company, c.Models.contract, c.Models.purchase)

		url := "/contract"
		req := httptest.NewRequest("PUT", url, bytes.NewBuffer([]byte(c.Request)))
		w := httptest.NewRecorder()

		testHandle("/contract", w, req, h.UpdateContract)
		testCheckResponse("UpdateContract:"+c.Num, t, w, c.Status, c.Response)
	}
}

func TestDeleteContract(t *testing.T) {
	time1 := time.Date(2000, 01, 01, 00, 00, 00, 0, time.UTC)

	cases := []struct {
		Num      string
		ID       int
		Response string
		Status   int
		Models   testModelSet
	}{
		// delete existing item
		{
			Num:      "1",
			ID:       1,
			Response: http.StatusText(http.StatusOK),
			Status:   http.StatusOK,
			Models: testModelSet{
				contract: test.TestContract{
					CL: []*model.Contract{
						{1, 10, 11, time1, time1.AddDate(1, 0, 0), 10},
					},
				},
			},
		},
		// update non-existing one
		{
			Num:      "2",
			ID:       1,
			Response: test.ErrTest.Error() + "\n",
			Status:   http.StatusNotFound,
			Models: testModelSet{
				contract: test.TestContract{
					CL: []*model.Contract{},
				},
			},
		},
		// error handling
		{
			Num:      "3",
			ID:       1,
			Response: test.ErrTest.Error() + "\n",
			Status:   http.StatusNotFound,
			Models: testModelSet{
				contract: test.TestContractErr{},
			},
		},
	}

	for _, c := range cases {
		h := testNewHandler(c.Models.company, c.Models.contract, c.Models.purchase)

		url := fmt.Sprintf("/contract/%d", c.ID)
		req := httptest.NewRequest("DELETE", url, nil)
		w := httptest.NewRecorder()

		testHandle("/contract/{id:[0-9]+}", w, req, h.DeleteContract)
		testCheckResponse("DeleteContract:"+c.Num, t, w, c.Status, c.Response)
	}
}

func TestPurchase(t *testing.T) {
	time1 := time.Date(2000, 02, 01, 00, 00, 00, 0, time.UTC)
	time2 := time.Date(2000, 04, 01, 00, 00, 00, 0, time.UTC)
	time3 := time.Date(2000, 03, 05, 00, 00, 00, 0, time.UTC)

	cases := []struct {
		Num      string
		Request  string
		Response string
		Status   int
		Models   testModelSet
	}{
		// normal creation
		{
			Num:      "1",
			Request:  `{"contractID":1,"datetime":"2000-03-01T00:00:00Z","amount":10}`,
			Response: `{"ID":1}`,
			Status:   http.StatusCreated,
			Models: testModelSet{
				contract: test.TestContract{
					CL: []*model.Contract{
						{1, 10, 11, time1, time2, 10},
					},
				},
				purchase: test.TestPurchase{
					CL: []*model.Purchase{},
				},
			},
		},
		// contract doesn't exist
		{
			Num:      "2",
			Request:  `{"contractID":1,"datetime":"2000-03-01T00:00:00Z","amount":10}`,
			Response: ErrContractNotFound.Error() + "\n",
			Status:   http.StatusBadRequest,
			Models: testModelSet{
				contract: test.TestContract{
					CL: []*model.Contract{},
				},
				purchase: test.TestPurchase{
					CL: []*model.Purchase{},
				},
			},
		},
		// date is below validity range
		{
			Num:      "3",
			Request:  `{"contractID":1,"datetime":"2000-01-01T00:00:00Z","amount":10}`,
			Response: ErrDateNotValid.Error() + "\n",
			Status:   http.StatusBadRequest,
			Models: testModelSet{
				contract: test.TestContract{
					CL: []*model.Contract{
						{1, 10, 11, time1, time2, 10},
					},
				},
				purchase: test.TestPurchase{
					CL: []*model.Purchase{},
				},
			},
		},
		// date is over validity range
		{
			Num:      "4",
			Request:  `{"contractID":1,"datetime":"2000-05-01T00:00:00Z","amount":10}`,
			Response: ErrDateNotValid.Error() + "\n",
			Status:   http.StatusBadRequest,
			Models: testModelSet{
				contract: test.TestContract{
					CL: []*model.Contract{
						{1, 10, 11, time1, time2, 10},
					},
				},
				purchase: test.TestPurchase{
					CL: []*model.Purchase{},
				},
			},
		},
		// error not enough money
		{
			Num:      "5",
			Request:  `{"contractID":1,"datetime":"2000-03-01T00:00:00Z","amount":5}`,
			Response: ErrNotEnoughMoney.Error() + "\n",
			Status:   http.StatusInternalServerError,
			Models: testModelSet{
				contract: test.TestContract{
					CL: []*model.Contract{
						{1, 10, 11, time1, time2, 10},
					},
				},
				purchase: test.TestPurchase{
					CL: []*model.Purchase{
						{1, 1, time3, 7},
					},
				},
			},
		},
		// error not enough money
		{
			Num:      "6",
			Request:  `{"contractID":1,"datetime":"2000-03-01T00:00:00Z","amount":5}`,
			Response: ErrNotEnoughMoney.Error() + "\n",
			Status:   http.StatusInternalServerError,
			Models: testModelSet{
				contract: test.TestContract{
					CL: []*model.Contract{
						{1, 10, 11, time1, time2, 10},
					},
				},
				purchase: test.TestPurchase{
					CL: []*model.Purchase{
						{1, 1, time3, 3},
						{1, 1, time3, 3},
					},
				},
			},
		},
		// error not enough money
		{
			Num:      "7",
			Request:  `{"contractID":1,"datetime":"2000-03-01T00:00:00Z","amount":5}`,
			Response: ErrNotEnoughMoney.Error() + "\n",
			Status:   http.StatusInternalServerError,
			Models: testModelSet{
				contract: test.TestContract{
					CL: []*model.Contract{
						{1, 10, 11, time1, time2, 1},
					},
				},
				purchase: test.TestPurchase{
					CL: []*model.Purchase{},
				},
			},
		},
		// error handling
		{
			Num:      "8",
			Request:  `{"contractID":1,"datetime":"2000-03-01T00:00:00Z","amount":5}`,
			Response: test.ErrTest.Error() + "\n",
			Status:   http.StatusInternalServerError,
			Models: testModelSet{
				contract: test.TestContract{
					CL: []*model.Contract{
						{1, 10, 11, time1, time2, 10},
					},
				},
				purchase: test.TestPurchaseErr{},
			},
		},
	}

	for _, c := range cases {
		h := testNewHandler(c.Models.company, c.Models.contract, c.Models.purchase)

		url := "/purchase"
		req := httptest.NewRequest("POST", url, bytes.NewBuffer([]byte(c.Request)))
		w := httptest.NewRecorder()

		testHandle("/purchase", w, req, h.Purchase)
		testCheckResponse("Purchase:"+c.Num, t, w, c.Status, c.Response)
	}
}

func TestPurchaseHistory(t *testing.T) {
	time1 := time.Date(2000, 01, 01, 00, 00, 00, 0, time.UTC)

	cases := []struct {
		Num      string
		ID       int
		Response string
		Status   int
		Models   testModelSet
	}{
		// normal request
		{
			Num:      "1",
			ID:       1,
			Response: `[{"ID":1,"contractID":1,"datetime":"2000-01-01T00:00:00Z","amount":3},{"ID":1,"contractID":1,"datetime":"2000-01-01T00:00:00Z","amount":3}]`,
			Status:   http.StatusOK,
			Models: testModelSet{
				contract: test.TestContract{
					CL: []*model.Contract{
						{1, 10, 11, time1, time1.AddDate(1, 0, 0), 10},
					},
				},
				purchase: test.TestPurchase{
					CL: []*model.Purchase{
						{1, 1, time1, 3},
						{1, 1, time1, 3},
					},
				},
			},
		},
		// request with empty history
		{
			Num:      "2",
			ID:       1,
			Response: `[]`,
			Status:   http.StatusOK,
			Models: testModelSet{
				contract: test.TestContract{
					CL: []*model.Contract{
						{1, 10, 11, time1, time1.AddDate(1, 0, 0), 10},
					},
				},
				purchase: test.TestPurchase{
					CL: []*model.Purchase{},
				},
			},
		},
		// handle error
		{
			Num:      "3",
			ID:       1,
			Response: ErrContractNotFound.Error() + "\n",
			Status:   http.StatusNotFound,
			Models: testModelSet{
				contract: test.TestContractErr{},
				purchase: test.TestPurchaseErr{},
			},
		},
	}

	for _, c := range cases {
		h := testNewHandler(c.Models.company, c.Models.contract, c.Models.purchase)

		url := fmt.Sprintf("/contract/%d/purchase", c.ID)
		req := httptest.NewRequest("GET", url, nil)
		w := httptest.NewRecorder()

		testHandle("/contract/{id:[0-9]+}/purchase", w, req, h.GetPurchaseHistory)
		testCheckResponse("GetPurchaseHistory:"+c.Num, t, w, c.Status, c.Response)
	}
}

func TestNewHandler(t *testing.T) {
	h := NewHandler(test.TestCompanyErr{}, test.TestContractErr{}, test.TestPurchaseErr{})

	if h == nil {
		t.Error("[NewHandler]:\tempty handler")
	}
}
