package app

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/SYSTEMTerror/GoEDU/cmd/app/middleware"
	"github.com/SYSTEMTerror/GoEDU/pkg/users"
	"github.com/gorilla/mux"
)

//Server is structure for server with mux from gorilla/mux
type Server struct {
	mux          *mux.Router
	usersSvc     *users.Service
}

//NewServer creates new server with mux from gorilla/mux
func NewServer(mux *mux.Router, usersSvc *users.Service) *Server {
	return &Server{mux: mux, usersSvc: usersSvc}
}

// ServeHTTP
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

//Init initializes server
func (s *Server) Init() {
	s.mux.Use(middleware.LoggersFuncs)

	usersAuthenticateMd := middleware.Authenticate(s.usersSvc.IDByToken)

	mainSubrouter := s.mux.PathPrefix("/api/v1").Subrouter()
	mainSubrouter.Use(usersAuthenticateMd)

	mainSubrouter.HandleFunc("/login", s.handleLoginUser).Methods("POST")
	mainSubrouter.HandleFunc("/register", s.handleRegisterUser).Methods("POST")
	mainSubrouter.HandleFunc("/admin", s.handleMakeAdmin).Methods("POST")
	mainSubrouter.HandleFunc("/subscribe", s.handleSubscribe).Methods("POST")
	mainSubrouter.HandleFunc("/user/{id}", s.handleGetUserByID).Methods("GET")

	courseSubrouter := mainSubrouter.PathPrefix("/course").Subrouter()
	courseSubrouter.HandleFunc("/create", s.handleCreateCourse).Methods("POST")
	courseSubrouter.HandleFunc("/user/all", s.handleGetAllUsers).Methods("GET")
	courseSubrouter.HandleFunc("/user/{id}", s.handleUserCourses).Methods("GET")
	courseSubrouter.HandleFunc("/subscribers/{id}", s.handleCourseSubscribes).Methods("GET")
}

//function jsoner marshal interfaces to json and write to response writer
func jsoner(w http.ResponseWriter, v interface{}, code int) error {
	data, err := json.Marshal(v)
	if err != nil {
		log.Println("jsoner json.Marshal error:", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err = w.Write(data)
	if err != nil {
		log.Println("jsoner w.Write error:", err)
		return err
	}
	return nil
}
