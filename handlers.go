package gontracts

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
	"github.com/ilyakaznacheev/gontracts/model"
)

var (
	// ErrContractNotFound contract doesn't exist in DB
	ErrContractNotFound = errors.New("contract doesn't exist")
	// ErrSellerNotExist seller company doesn't exist in DB
	ErrSellerNotExist = errors.New("seller company doesn't exist")
	// ErrClientNotExist client company doesn't exist in DB
	ErrClientNotExist = errors.New("client company doesn't exist")
	// ErrDateNotValid purchase date is outside the contract date range
	ErrDateNotValid = errors.New("purchase date is outside the contract date range")
	// ErrNotEnoughMoney not enough money for the purchase
	ErrNotEnoughMoney = errors.New("not enough money for the purchase")
)

// ResponseID represents id in response
type ResponseID struct {
	ID int
}

// Handler is a request handler
type Handler struct {
	mh         *model.ModelHandler
	purchaseMX *sync.Mutex
}

// NewHandler returns new request handler
func NewHandler(company model.CompanyModel, contract model.ContractModel, purchase model.PurchaseModel) *Handler {
	return &Handler{
		mh:         model.GetModelHandler(company, contract, purchase),
		purchaseMX: &sync.Mutex{},
	}
}

// GetCompany returns company info
func (h *Handler) GetCompany(w http.ResponseWriter, r *http.Request) {
	// get id from request params
	rvars := mux.Vars(r)
	id, err := strconv.Atoi(rvars["id"])
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// read data from DB
	c, err := h.mh.GetCompany(id)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// fill response json
	resp, err := json.Marshal(*c)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// setup response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

// GetCompanyList returns list of companies
func (h *Handler) GetCompanyList(w http.ResponseWriter, r *http.Request) {
	// read data from DB
	c, err := h.mh.GetCompanyList()
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	compList := make([]model.Company, 0, len(c))
	for _, comp := range c {
		compList = append(compList, *comp)
	}

	// fill response json
	resp, err := json.Marshal(compList)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// setup response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

// CreateCompany creates new company
func (h *Handler) CreateCompany(w http.ResponseWriter, r *http.Request) {
	var company model.Company

	// read request body
	dc := json.NewDecoder(r.Body)
	err := dc.Decode(&company)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// create new company in DB
	idx, err := h.mh.CreateCompany(&company)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// fill response json
	resp, err := json.Marshal(&ResponseID{idx})
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// setup response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(resp)
}

// UpdateCompany updates company data
func (h *Handler) UpdateCompany(w http.ResponseWriter, r *http.Request) {
	var company model.Company
	var okStatus int

	// read request body
	dc := json.NewDecoder(r.Body)
	err := dc.Decode(&company)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if company.ID == 0 {
		// if id is empty, create new company
		idx, err := h.mh.CreateCompany(&company)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		company.ID = idx
		okStatus = http.StatusCreated
	} else {
		// if id is set, updete existing company
		err := h.mh.UpdateCompany(&company)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		okStatus = http.StatusOK
	}

	// fill response json
	resp, err := json.Marshal(&company)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// setup response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(okStatus)
	w.Write(resp)
}

// DeleteCompany removes company
func (h *Handler) DeleteCompany(w http.ResponseWriter, r *http.Request) {
	// get id from request params
	rvars := mux.Vars(r)
	id, err := strconv.Atoi(rvars["id"])
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// delete company from DB
	err = h.mh.DeleteCompany(id)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// setup response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

// GetContract returns contract data
func (h *Handler) GetContract(w http.ResponseWriter, r *http.Request) {
	// get id from request params
	rvars := mux.Vars(r)
	id, err := strconv.Atoi(rvars["id"])
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// read data from DB
	c, err := h.mh.GetContract(id)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// fill response json
	resp, err := json.Marshal(*c)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// setup response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

// GetContractList returns contract list
func (h *Handler) GetContractList(w http.ResponseWriter, r *http.Request) {
	// read data from DB
	c, err := h.mh.GetContractList()
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	contrList := make([]model.Contract, 0, len(c))
	for _, comp := range c {
		contrList = append(contrList, *comp)
	}

	// fill response json
	resp, err := json.Marshal(contrList)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// setup response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

// CreateContract creates new contract
func (h *Handler) CreateContract(w http.ResponseWriter, r *http.Request) {
	var contract model.Contract

	// read request body
	dc := json.NewDecoder(r.Body)
	err := dc.Decode(&contract)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// validity checks

	// chech seller company exists in DB
	if !h.mh.CheckCompanyExist(contract.SellerID) {
		log.Println(ErrSellerNotExist)
		http.Error(w, ErrSellerNotExist.Error(), http.StatusBadRequest)
		return
	}

	// chech client company exists in DB
	if !h.mh.CheckCompanyExist(contract.ClientID) {
		log.Println(ErrClientNotExist)
		http.Error(w, ErrClientNotExist.Error(), http.StatusBadRequest)
		return
	}

	// create new contract in DB
	idx, err := h.mh.CreateContract(&contract)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// fill response json
	resp, err := json.Marshal(&ResponseID{idx})
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// setup response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(resp)
}

// UpdateContract updates contract
func (h *Handler) UpdateContract(w http.ResponseWriter, r *http.Request) {
	var contract model.Contract
	var okStatus int

	// read request body
	dc := json.NewDecoder(r.Body)
	err := dc.Decode(&contract)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// validity checks

	// chech seller company exists in DB
	if !h.mh.CheckCompanyExist(contract.SellerID) {
		log.Println(ErrSellerNotExist)
		http.Error(w, ErrSellerNotExist.Error(), http.StatusBadRequest)
		return
	}

	// chech client company exists in DB
	if !h.mh.CheckCompanyExist(contract.ClientID) {
		log.Println(ErrClientNotExist)
		http.Error(w, ErrClientNotExist.Error(), http.StatusBadRequest)
		return
	}

	if contract.ID == 0 {
		// if id is empty create new contract
		idx, err := h.mh.CreateContract(&contract)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		contract.ID = idx
		okStatus = http.StatusCreated
	} else {
		// if id is set update existing contract
		err := h.mh.UpdateContract(&contract)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		okStatus = http.StatusOK
	}

	// fill response json
	resp, err := json.Marshal(&contract)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// setup response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(okStatus)
	w.Write(resp)
}

// DeleteContract removes contract
func (h *Handler) DeleteContract(w http.ResponseWriter, r *http.Request) {
	// get id from request params
	rvars := mux.Vars(r)
	id, err := strconv.Atoi(rvars["id"])
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// delete contract from DB
	err = h.mh.DeleteContract(id)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// setup response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

// Purchase creates new purchase document
func (h *Handler) Purchase(w http.ResponseWriter, r *http.Request) {
	var purchase model.Purchase

	// read request body
	dc := json.NewDecoder(r.Body)
	err := dc.Decode(&purchase)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// read contract data from DB
	contract, err := h.mh.GetContract(purchase.ContractID)
	if err != nil {
		log.Println(ErrContractNotFound)
		http.Error(w, ErrContractNotFound.Error(), http.StatusBadRequest)
		return
	}

	// check if purchase document in valud date range of contract
	if purchase.PurchaseDateTime.Before(contract.ValidFrom) || purchase.PurchaseDateTime.After(contract.ValidTo) {
		log.Println(ErrDateNotValid)
		http.Error(w, ErrDateNotValid.Error(), http.StatusBadRequest)
		return
	}

	h.purchaseMX.Lock()
	defer h.purchaseMX.Unlock()
	// read sum of existing purchase documents
	sum := h.mh.GetContractPurchaseSum(purchase.ContractID)

	// calculate remain credits and check
	// if there is enough money to process new payment
	remain := contract.CreditAmount - sum
	if remain < purchase.CreditSpent {
		log.Println(ErrNotEnoughMoney)
		http.Error(w, ErrNotEnoughMoney.Error(), http.StatusInternalServerError)
		return
	}

	// create new payment document in DB
	idx, err := h.mh.CreatePurchase(&purchase)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// fill response json
	resp, err := json.Marshal(&ResponseID{idx})
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// setup response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(resp)
}

// GetPurchaseHistory returns purcase history of contract
func (h *Handler) GetPurchaseHistory(w http.ResponseWriter, r *http.Request) {
	// get id from request params
	rvars := mux.Vars(r)
	id, err := strconv.Atoi(rvars["id"])
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// read purchase history of contract
	p, _ := h.mh.GetContractPurchaseHistory(id)

	if len(p) == 0 {
		if !h.mh.CheckContractsExist(id) {
			log.Println(ErrContractNotFound)
			http.Error(w, ErrContractNotFound.Error(), http.StatusNotFound)
			return
		}
	}

	purList := make([]model.Purchase, 0, len(p))
	for _, comp := range p {
		purList = append(purList, *comp)
	}

	// fill response json
	resp, err := json.Marshal(purList)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// setup response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
