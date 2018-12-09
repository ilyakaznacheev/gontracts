package gontracts

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
	"github.com/ilyakaznacheev/gontracts/db"
)

// Server is an main application server
type Server struct {
	stop func()
}

// Start runs the server
func (s *Server) Start() error {

	dbConn, err := db.Connect()
	if err != nil {
		return err
	}

	h := NewHandler(
		db.NewCompanyDAC(dbConn),
		db.NewContractDAC(dbConn),
		db.NewPurchaseDAC(dbConn),
	)

	r := mux.NewRouter()
	// setup uri handlers
	r.HandleFunc("/company/{id:[0-9]+}", h.GetCompany).Methods("GET")
	r.HandleFunc("/company", h.CreateCompany).Methods("POST")
	r.HandleFunc("/company", h.UpdateCompany).Methods("PUT")
	r.HandleFunc("/company/{id:[0-9]+}", h.DeleteCompany).Methods("DELETE")
	r.HandleFunc("/companies", h.GetCompanyList).Methods("GET")
	r.HandleFunc("/contract/{id:[0-9]+}", h.GetContract).Methods("GET")
	r.HandleFunc("/contract", h.CreateContract).Methods("POST")
	r.HandleFunc("/contract", h.UpdateContract).Methods("PUT")
	r.HandleFunc("/contract/{id:[0-9]+}", h.DeleteContract).Methods("DELETE")
	r.HandleFunc("/contract/{id:[0-9]+}/purchase", h.GetPurchaseHistory).Methods("GET")
	r.HandleFunc("/contracts", h.GetContractList).Methods("GET")
	r.HandleFunc("/purchase", h.Purchase).Methods("POST")

	// handle keyboard interrupt
	s.stop = func() {
		dbConn.Close()
	}

	s.handleInterrupt()

	// start server
	log.Println("starting server")
	err = http.ListenAndServe(":8000", r)
	log.Fatal(err)
	return err

}

// Stop shuts the server down
func (s *Server) Stop() {
	log.Println("shutdown server")
	s.stop()
}

func (s *Server) handleInterrupt() {
	var signalChan = make(chan os.Signal)
	signal.Notify(signalChan, syscall.SIGTERM)
	signal.Notify(signalChan, syscall.SIGINT)

	go func() {
		<-signalChan
		s.Stop()
		os.Exit(0)
	}()
}
