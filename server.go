package gontracts

import (
	"crypto/rand"
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

	// generate random secret key for sesstion
	key := make([]byte, 64)
	_, err := rand.Read(key)
	if err != nil {
		return err
	}

	dbConn, err := db.Connect()
	if err != nil {
		return err
	}

	h := NewHandler(
		db.NewCompanyDAC(dbConn),
		db.NewContractDAC(dbConn),
		db.NewPurchaseDAC(dbConn),
	)

	a := NewAuthHandler(key)

	r := mux.NewRouter()

	// setup uri handlers
	r.Handle("/company/{id:[0-9]+}", a.HandlerFunc(h.GetCompany)).Methods("GET")
	r.Handle("/company", a.HandlerFunc(h.CreateCompany)).Methods("POST")
	r.Handle("/company", a.HandlerFunc(h.UpdateCompany)).Methods("PUT")
	r.Handle("/company/{id:[0-9]+}", a.HandlerFunc(h.DeleteCompany)).Methods("DELETE")
	r.Handle("/company", a.HandlerFunc(h.GetCompanyList)).Methods("GET")
	r.Handle("/contract/{id:[0-9]+}", a.HandlerFunc(h.GetContract)).Methods("GET")
	r.Handle("/contract", a.HandlerFunc(h.CreateContract)).Methods("POST")
	r.Handle("/contract", a.HandlerFunc(h.UpdateContract)).Methods("PUT")
	r.Handle("/contract/{id:[0-9]+}", a.HandlerFunc(h.DeleteContract)).Methods("DELETE")
	r.Handle("/contract/{id:[0-9]+}/purchase", a.HandlerFunc(h.GetPurchaseHistory)).Methods("GET")
	r.Handle("/contract", a.HandlerFunc(h.GetContractList)).Methods("GET")
	r.Handle("/purchase", a.HandlerFunc(h.Purchase)).Methods("POST")

	r.HandleFunc("/get-token", a.GenerateToken).Methods("GET")

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
