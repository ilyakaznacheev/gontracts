package model

import (
	"time"
)

// Company represent company DB table structure
type Company struct {
	ID      int     `json:"ID"`
	Name    string  `json:"name"`
	RegCode *string `json:"regcode"`
}

// Contract represent contract DB table structure
type Contract struct {
	ID           int       `json:"ID"`
	SellerID     int       `json:"sellerID"`
	ClientID     int       `json:"clientID"`
	ValidFrom    time.Time `json:"validFrom"`
	ValidTo      time.Time `json:"validTo"`
	CreditAmount int       `json:"amount"`
}

// Purchase represent purchase DB table structure
type Purchase struct {
	ID               int       `json:"ID"`
	ContractID       int       `json:"contractID"`
	PurchaseDateTime time.Time `json:"datetime"`
	CreditSpent      int       `json:"amount"`
}

// CompanyModel represents company interaction scheme
type CompanyModel interface {
	GetList() ([]*Company, error)
	GetItem(int) (*Company, error)
	CreateItem(*Company) (int, error)
	UpdateItem(*Company) error
	DeleteItem(int) error
	CheckExist(int) bool
}

// ContractModel represents contract interaction scheme
type ContractModel interface {
	GetList() ([]*Contract, error)
	GetItem(int) (*Contract, error)
	CreateItem(*Contract) (int, error)
	UpdateItem(*Contract) error
	DeleteItem(int) error
	CheckExist(int) bool
}

// PurchaseModel represents purchase interaction scheme
type PurchaseModel interface {
	AddItem(*Purchase) (int, error)
	GetContractHistory(int) ([]*Purchase, error)
	GetContractSum(int) int
}
