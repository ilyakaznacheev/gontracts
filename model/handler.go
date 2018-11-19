package model

// ModelHandler is a persistent data interaction object
type ModelHandler struct {
	company  CompanyModel
	contract ContractModel
	purchase PurchaseModel
}

// NewModelHandler creates model handler
func NewModelHandler(
	company CompanyModel,
	contract ContractModel,
	purchase PurchaseModel,
) *ModelHandler {
	return &ModelHandler{
		company:  company,
		contract: contract,
		purchase: purchase,
	}
}

var mh *ModelHandler

// GetModelHandler returns model singleton
func GetModelHandler(
	company CompanyModel,
	contract ContractModel,
	purchase PurchaseModel,
) *ModelHandler {
	if mh == nil {
		mh = NewModelHandler(company, contract, purchase)
	}
	return mh
}

// GetCompanyList returns list of all companies
func (m *ModelHandler) GetCompanyList() ([]*Company, error) {
	return m.company.GetList()
}

// GetCompany returns company by id
func (m *ModelHandler) GetCompany(id int) (*Company, error) {
	return m.company.GetItem(id)
}

// CreateCompany creates new company
func (m *ModelHandler) CreateCompany(c *Company) (int, error) {
	return m.company.CreateItem(c)
}

// UpdateCompany updates company
func (m *ModelHandler) UpdateCompany(c *Company) error {
	return m.company.UpdateItem(c)
}

// DeleteCompany removes company
func (m *ModelHandler) DeleteCompany(id int) error {
	return m.company.DeleteItem(id)
}

// CheckCompanyExist checks are company with id  exists
func (m *ModelHandler) CheckCompanyExist(id int) bool {
	return m.company.CheckExist(id)
}

// GetContractList returns list of all contracts
func (m *ModelHandler) GetContractList() ([]*Contract, error) {
	return m.contract.GetList()
}

// GetContract returns contract by id
func (m *ModelHandler) GetContract(id int) (*Contract, error) {
	return m.contract.GetItem(id)
}

// CreateContract creates new contract
func (m *ModelHandler) CreateContract(c *Contract) (int, error) {
	return m.contract.CreateItem(c)
}

// UpdateContract updates contract
func (m *ModelHandler) UpdateContract(c *Contract) error {
	return m.contract.UpdateItem(c)
}

// DeleteContract removes contract
func (m *ModelHandler) DeleteContract(id int) error {
	return m.contract.DeleteItem(id)
}

// CheckContractsExist checks are company with id  exists
func (m *ModelHandler) CheckContractsExist(id int) bool {
	return m.contract.CheckExist(id)
}

// CreatePurchase creates new purchase document
func (m *ModelHandler) CreatePurchase(p *Purchase) (int, error) {
	return m.purchase.AddItem(p)
}

// GetContractPurchaseSum returns purchase sum of contract
func (m *ModelHandler) GetContractPurchaseSum(id int) int {
	return m.purchase.GetContractSum(id)
}

// GetContractPurchaseHistory returns purchase history of contract
func (m *ModelHandler) GetContractPurchaseHistory(id int) ([]*Purchase, error) {
	return m.purchase.GetContractHistory(id)
}
