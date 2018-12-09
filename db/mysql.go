package db

import (
	"database/sql"
	"sync"

	_ "github.com/go-sql-driver/mysql" //use MySQL driver

	"github.com/ilyakaznacheev/gontracts/model"
)

func readIndex(db *sql.DB) (int, error) {
	rows, err := db.Query(
		`SELECT LAST_INSERT_ID()`,
	)
	if err != nil {
		return 0, err
	}

	rows.Next()
	var idx int
	err = rows.Scan(&idx)

	return idx, err
}

// CompanyDAC is a company table data access class
type CompanyDAC struct {
	db *sql.DB
	mx *sync.Mutex
}

// NewCompanyDAC creates new company DAC
func NewCompanyDAC(db *sql.DB) *CompanyDAC {
	return &CompanyDAC{
		db: db,
		mx: &sync.Mutex{},
	}
}

// GetList returns list of all companies
func (dac *CompanyDAC) GetList() ([]*model.Company, error) {
	rows, err := dac.db.Query(
		`SELECT id, name, regcode
			FROM company`,
	)
	if err != nil {
		return nil, err
	}
	compList := make([]*model.Company, 0)

	for rows.Next() {
		compItem := &model.Company{}
		err = rows.Scan(
			&compItem.ID,
			&compItem.Name,
			&compItem.RegCode,
		)
		if err != nil {
			return nil, err
		}
		compList = append(compList, compItem)
	}
	rows.Close()

	return compList, nil
}

// GetItem returns company by id
func (dac *CompanyDAC) GetItem(id int) (*model.Company, error) {
	rows, err := dac.db.Query(
		`SELECT id, name, regcode
			FROM company
			WHERE
				id = ?`,
		id,
	)
	if err != nil {
		return nil, err
	}
	compItem := &model.Company{}
	rows.Next()
	err = rows.Scan(
		&compItem.ID,
		&compItem.Name,
		&compItem.RegCode,
	)

	return compItem, err
}

// CreateItem creates new company
func (dac *CompanyDAC) CreateItem(company *model.Company) (int, error) {
	dac.mx.Lock()
	defer dac.mx.Unlock()
	_, err := dac.db.Exec(
		`INSERT 
			INTO company (name, regcode) 
			VALUES (?, ?)`,
		company.Name,
		company.RegCode,
	)
	if err != nil {
		return 0, err
	}

	return readIndex(dac.db)
}

// UpdateItem updates company
func (dac *CompanyDAC) UpdateItem(company *model.Company) error {
	dac.mx.Lock()
	_, err := dac.db.Exec(
		`UPDATE company
			SET
				name=?,
				regcode=?
			WHERE
				id=?`,
		company.Name,
		company.RegCode,
		company.ID,
	)
	dac.mx.Unlock()
	return err
}

// DeleteItem removes company
func (dac *CompanyDAC) DeleteItem(id int) error {
	dac.mx.Lock()
	_, err := dac.db.Exec(
		`DELETE FROM company
			WHERE
				id=?`,
		id,
	)
	dac.mx.Unlock()
	return err
}

// CheckExist checks are company with id exists
func (dac *CompanyDAC) CheckExist(id int) bool {
	rows, err := dac.db.Query(
		`SELECT EXISTS(
			SELECT 1 
				FROM company 
				WHERE id=?
		)`,
		id,
	)
	if err != nil {
		return false
	}
	var exist bool
	rows.Next()
	err = rows.Scan(&exist)
	if err != nil {
		return false
	}
	return exist
}

// ContractDAC is a company table data access class
type ContractDAC struct {
	db *sql.DB
	mx *sync.Mutex
}

// NewContractDAC creates new company DAC
func NewContractDAC(db *sql.DB) *ContractDAC {
	return &ContractDAC{
		db: db,
		mx: &sync.Mutex{},
	}
}

// GetList returns list of all contracts
func (dac *ContractDAC) GetList() ([]*model.Contract, error) {
	rows, err := dac.db.Query(
		`SELECT id, clientid, sellerid, validfrom, validto, creditamount
			FROM contract`,
	)
	if err != nil {
		return nil, err
	}
	contrList := make([]*model.Contract, 0)

	for rows.Next() {
		contrItem := &model.Contract{}
		err = rows.Scan(
			&contrItem.ID,
			&contrItem.ClientID,
			&contrItem.SellerID,
			&contrItem.ValidFrom,
			&contrItem.ValidTo,
			&contrItem.CreditAmount,
		)
		if err != nil {
			return nil, err
		}
		contrList = append(contrList, contrItem)
	}
	rows.Close()

	return contrList, nil
}

// GetItem returns contract by id
func (dac *ContractDAC) GetItem(id int) (*model.Contract, error) {
	rows, err := dac.db.Query(
		`SELECT id, clientid, sellerid, validfrom, validto, creditamount
			FROM contract
			WHERE
				id = ?`,
		id,
	)
	if err != nil {
		return nil, err
	}
	contrItem := &model.Contract{}
	rows.Next()
	err = rows.Scan(
		&contrItem.ID,
		&contrItem.ClientID,
		&contrItem.SellerID,
		&contrItem.ValidFrom,
		&contrItem.ValidTo,
		&contrItem.CreditAmount,
	)

	return contrItem, err
}

// CreateItem creates new contract
func (dac *ContractDAC) CreateItem(contract *model.Contract) (int, error) {
	dac.mx.Lock()
	defer dac.mx.Unlock()
	_, err := dac.db.Exec(
		`INSERT 
			INTO contract (clientid, sellerid, validfrom, validto, creditamount) 
			VALUES (?, ?, ?, ?, ?)`,
		contract.ClientID,
		contract.SellerID,
		contract.ValidFrom,
		contract.ValidTo,
		contract.CreditAmount,
	)
	if err != nil {
		return 0, err
	}

	return readIndex(dac.db)
}

// UpdateItem updates contract
func (dac *ContractDAC) UpdateItem(contract *model.Contract) error {
	dac.mx.Lock()
	_, err := dac.db.Exec(
		`UPDATE contract
			SET
				clientid=?, 
				sellerid=?, 
				validfrom=?, 
				validto=?, 
				creditamount=?
			WHERE
				id=?`,
		contract.ClientID,
		contract.SellerID,
		contract.ValidFrom,
		contract.ValidTo,
		contract.CreditAmount,
		contract.ID,
	)
	dac.mx.Unlock()
	return err
}

// DeleteItem removes contract
func (dac *ContractDAC) DeleteItem(id int) error {
	dac.mx.Lock()
	_, err := dac.db.Exec(
		`DELETE FROM contract
			WHERE
				id=?`,
		id,
	)
	dac.mx.Unlock()
	return err
}

// CheckExist checks are company with id exists
func (dac *ContractDAC) CheckExist(id int) bool {
	rows, err := dac.db.Query(
		`SELECT EXISTS(
			SELECT 1 
				FROM contract 
				WHERE id=?
		)`,
		id,
	)
	if err != nil {
		return false
	}
	var exist bool
	rows.Next()
	err = rows.Scan(&exist)
	if err != nil {
		return false
	}
	return exist
}

// PurchaseDAC is a purchase table data access class
type PurchaseDAC struct {
	db *sql.DB
	mx *sync.Mutex
}

// NewPurchaseDAC creates new company DAC
func NewPurchaseDAC(db *sql.DB) *PurchaseDAC {
	return &PurchaseDAC{
		db: db,
		mx: &sync.Mutex{},
	}
}

// AddItem creates new purchase document
func (dac *PurchaseDAC) AddItem(purchase *model.Purchase) (int, error) {
	dac.mx.Lock()
	defer dac.mx.Unlock()
	_, err := dac.db.Exec(
		`INSERT 
			INTO purchase (contractid, purchasedatetime, creditspent) 
			VALUES (?, ?, ?)`,
		purchase.ContractID,
		purchase.PurchaseDateTime,
		purchase.CreditSpent,
	)
	if err != nil {
		return 0, err
	}

	return readIndex(dac.db)
}

// GetContractHistory returns purchase history of contract
func (dac *PurchaseDAC) GetContractHistory(id int) ([]*model.Purchase, error) {
	rows, err := dac.db.Query(
		`SELECT id, contractid, purchasedatetime, creditspent
			FROM purchase
			WHERE 
				contractid=?
			ORDER BY
				purchasedatetime`,
		id,
	)
	if err != nil {
		return nil, err
	}

	purList := make([]*model.Purchase, 0)

	for rows.Next() {
		purItem := &model.Purchase{}
		err = rows.Scan(
			&purItem.ID,
			&purItem.ContractID,
			&purItem.PurchaseDateTime,
			&purItem.CreditSpent,
		)
		if err != nil {
			return nil, err
		}
		purList = append(purList, purItem)
	}
	rows.Close()
	return purList, nil
}

// GetContractSum returns purchase sum of contract
func (dac *PurchaseDAC) GetContractSum(id int) int {
	rows, err := dac.db.Query(
		`SELECT sum(creditspent) as credit
			FROM purchase
			WHERE 
				contractid=?`,
		id,
	)
	if err != nil {
		return 0
	}

	var sum int
	rows.Next()
	_ = rows.Scan(&sum)

	return sum
}

// Connection to DB

// Connect returns connection to MySQL
func Connect() (*sql.DB, error) {
	// hardcoded, change to config in real app
	db, err := sql.Open("mysql", "default:1234@/gontracts?parseTime=true")
	if err != nil {
		defer db.Close()
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
