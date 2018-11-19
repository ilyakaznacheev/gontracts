package test

import (
	"errors"

	"bitbucket.org/ilyakaznacheev/gontracts/model"
)

// Set of DB mock objects for testing

var ErrTest = errors.New("test error")

type TestCompany struct {
	CL []*model.Company
}

func (t TestCompany) GetList() ([]*model.Company, error) {
	return t.CL, nil
}
func (t TestCompany) GetItem(id int) (*model.Company, error) {
	for _, c := range t.CL {
		if c.ID == id {
			return c, nil
		}
	}
	return nil, ErrTest
}

func (t TestCompany) CreateItem(comp *model.Company) (int, error) {
	t.CL = append(t.CL, comp)
	return len(t.CL), nil
}

func (t TestCompany) UpdateItem(comp *model.Company) error {
	for _, c := range t.CL {
		if c.ID == comp.ID {
			c = comp
			return nil
		}
	}
	return ErrTest
}

func (t TestCompany) DeleteItem(id int) error {
	for idx, c := range t.CL {
		if c.ID == id {
			t.CL = append(t.CL[:idx], t.CL[idx+1:]...)
			return nil
		}
	}
	return ErrTest
}

func (t TestCompany) CheckExist(id int) bool {
	for _, c := range t.CL {
		if c.ID == id {
			return true
		}
	}
	return false
}

type TestCompanyErr struct {
}

func (t TestCompanyErr) GetList() ([]*model.Company, error)          { return nil, ErrTest }
func (t TestCompanyErr) GetItem(id int) (*model.Company, error)      { return nil, ErrTest }
func (t TestCompanyErr) CreateItem(comp *model.Company) (int, error) { return 0, ErrTest }
func (t TestCompanyErr) UpdateItem(comp *model.Company) error        { return ErrTest }
func (t TestCompanyErr) DeleteItem(id int) error                     { return ErrTest }
func (t TestCompanyErr) CheckExist(id int) bool                      { return false }

type TestContract struct {
	CL []*model.Contract
}

func (t TestContract) GetList() ([]*model.Contract, error) {
	return t.CL, nil
}

func (t TestContract) GetItem(id int) (*model.Contract, error) {
	for _, c := range t.CL {
		if c.ID == id {
			return c, nil
		}
	}
	return nil, ErrTest
}

func (t TestContract) CreateItem(contr *model.Contract) (int, error) {
	t.CL = append(t.CL, contr)
	return len(t.CL), nil
}

func (t TestContract) UpdateItem(contr *model.Contract) error {
	for _, c := range t.CL {
		if c.ID == contr.ID {
			c = contr
			return nil
		}
	}
	return ErrTest
}

func (t TestContract) DeleteItem(id int) error {
	for idx, c := range t.CL {
		if c.ID == id {
			t.CL = append(t.CL[:idx], t.CL[idx+1:]...)
			return nil
		}
	}
	return ErrTest
}

func (t TestContract) CheckExist(id int) bool {
	for _, c := range t.CL {
		if c.ID == id {
			return true
		}
	}
	return false
}

type TestContractErr struct {
}

func (t TestContractErr) GetList() ([]*model.Contract, error)           { return nil, ErrTest }
func (t TestContractErr) GetItem(id int) (*model.Contract, error)       { return nil, ErrTest }
func (t TestContractErr) CreateItem(contr *model.Contract) (int, error) { return 0, ErrTest }
func (t TestContractErr) UpdateItem(contr *model.Contract) error        { return ErrTest }
func (t TestContractErr) DeleteItem(id int) error                       { return ErrTest }
func (t TestContractErr) CheckExist(id int) bool                        { return false }

type TestPurchase struct {
	CL []*model.Purchase
}

func (t TestPurchase) AddItem(pur *model.Purchase) (int, error) {
	t.CL = append(t.CL, pur)
	return len(t.CL), nil
}

func (t TestPurchase) GetContractHistory(id int) ([]*model.Purchase, error) {
	var hist []*model.Purchase
	for _, c := range t.CL {
		if c.ContractID == id {
			hist = append(hist, c)
		}
	}
	if len(hist) == 0 {
		return nil, ErrTest
	}
	return hist, nil
}

func (t TestPurchase) GetContractSum(id int) int {
	var sum int
	for _, c := range t.CL {
		if c.ContractID == id {
			sum += c.CreditSpent
		}
	}
	return sum
}

type TestPurchaseErr struct {
}

func (t TestPurchaseErr) AddItem(pur *model.Purchase) (int, error)             { return 0, ErrTest }
func (t TestPurchaseErr) GetContractHistory(id int) ([]*model.Purchase, error) { return nil, ErrTest }
func (t TestPurchaseErr) GetContractSum(id int) int                            { return 0 }
