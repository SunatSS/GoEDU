package app

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/SYSTEMTerror/GoEDU/cmd/app/middleware"
	"github.com/SYSTEMTerror/GoEDU/pkg/types"
	"github.com/gorilla/mux"
)

//handleRegisterUser
func (s *Server) handleRegisterUser(w http.ResponseWriter, r *http.Request) {
	loggers, err := middleware.GetLoggers(r.Context())
	if err != nil {
		log.Println("LOGGERS DON'T WORK!!!")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	loggers.InfoLogger.Println("handleRegisterUser started")

	var item *types.RegInfo
	err = json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		loggers.ErrorLogger.Println("handleRegisterUser json.NewDecoder error:", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	user, statusCode, err := s.usersSvc.RegisterUser(r.Context(), item)
	if err != nil {
		loggers.ErrorLogger.Println("handleRegisterUser s.usersSvc.RegisterUser error:", err)
		http.Error(w, http.StatusText(statusCode), statusCode)
		return
	}

	err = jsoner(w, user, statusCode)
	if err != nil {
		loggers.ErrorLogger.Println("handleRegisterUser jsoner error:", err)
		return
	}
	loggers.InfoLogger.Println("handleRegisterUser finished with any error!")
}

//handleLoginUser
func (s *Server) handleLoginUser(w http.ResponseWriter, r *http.Request) {
	loggers, err := middleware.GetLoggers(r.Context())
	if err != nil {
		log.Println("LOGGERS DON'T WORK!!!")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	loggers.InfoLogger.Println("handleLoginUser started")

	var item *types.TokenInfo
	err = json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		loggers.ErrorLogger.Println("handleLoginUser json.NewDecoder error:", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	token, statusCode, err := s.usersSvc.Token(r.Context(), item)
	if err != nil {
		loggers.ErrorLogger.Println("handleLoginUser s.usersSvc.Token error:", err)
		http.Error(w, http.StatusText(statusCode), statusCode)
		return
	}

	err = jsoner(w, token, statusCode)
	if err != nil {
		loggers.ErrorLogger.Println("handleLoginUser jsoner error:", err)
		return
	}
	loggers.InfoLogger.Println("handleLoginUser finished with any error!")
}

//handleSubscribe subscribes a user to a course
func (s *Server) handleSubscribe(w http.ResponseWriter, r *http.Request) {
	loggers, err := middleware.GetLoggers(r.Context())
	if err != nil {
		log.Println("LOGGERS DON'T WORK!!!")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	userId, err := middleware.Authentication(r.Context())
	if err != nil {
		loggers.ErrorLogger.Println("handleSubscribe middleware.Authentication error:", err)
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	var item *types.SubscribeInfo
	err = json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		loggers.ErrorLogger.Println("handleSubscribe json.NewDecoder error:", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	item.UserID = userId

	statusCode, err := s.usersSvc.Subscribe(r.Context(), item)
	if err != nil {
		loggers.ErrorLogger.Println("handleSubscribe s.usersSvc.Subscribe error:", err)
		http.Error(w, http.StatusText(statusCode), statusCode)
		return
	}

	err = jsoner(w, nil, statusCode)
	if err != nil {
		loggers.ErrorLogger.Println("handleSubscribe jsoner error:", err)
		return
	}
	loggers.InfoLogger.Println("handleSubscribe finished with any error!")
}

//handleMakeAdmin makes a user with id an admin
func (s *Server) handleMakeAdmin(w http.ResponseWriter, r *http.Request) {
	loggers, err := middleware.GetLoggers(r.Context())
	if err != nil {
		log.Println("LOGGERS DON'T WORK!!!")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	loggers.InfoLogger.Println("handleMakeAdmin started")

	adminId, err := middleware.Authentication(r.Context())
	if err != nil {
		loggers.ErrorLogger.Println("handleMakeAdmin middleware.Authentication error:", err)
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	isAdmin, statusCode, err := s.usersSvc.IsAdmin(r.Context(), adminId)
	if err != nil {
		loggers.ErrorLogger.Println("handleMakeAdmin s.usersSvc.IsAdmin error:", err)
		http.Error(w, http.StatusText(statusCode), statusCode)
		return
	}
	
	if !isAdmin {
		loggers.ErrorLogger.Println("handleMakeAdmin s.usersSvc.IsAdmin isAdmin:", isAdmin)
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	var makeAdminInfo *types.MakeAdminInfo
	err = json.NewDecoder(r.Body).Decode(&makeAdminInfo)
	if err != nil {
		loggers.ErrorLogger.Println("handleMakeAdmin json.NewDecoder error:", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	statusCode, err = s.usersSvc.MakeAdmin(r.Context(), makeAdminInfo)
	if err != nil {
		loggers.ErrorLogger.Println("handleMakeAdmin s.usersSvc.MakeAdmin error:", err)
		http.Error(w, http.StatusText(statusCode), statusCode)
		return
	}

	err = jsoner(w, makeAdminInfo, statusCode)
	if err != nil {
		loggers.ErrorLogger.Println("handleMakeAdmin jsoner error:", err)
		return
	}
	loggers.InfoLogger.Println("handleMakeAdmin finished with any error!")
}

//handleCreateCourse creates a course
func (s *Server) handleCreateCourse(w http.ResponseWriter, r *http.Request) {
	loggers, err := middleware.GetLoggers(r.Context())
	if err != nil {
		log.Println("LOGGERS DON'T WORK!!!")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	loggers.InfoLogger.Println("handleCreateCourse started")

	adminId, err := middleware.Authentication(r.Context())
	if err != nil {
		loggers.ErrorLogger.Println("handleCreateCourse middleware.Authentication error:", err)
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	isAdmin, statusCode, err := s.usersSvc.IsAdmin(r.Context(), adminId)
	if err != nil {
		loggers.ErrorLogger.Println("handleCreateCourse s.usersSvc.IsAdmin error:", err)
		http.Error(w, http.StatusText(statusCode), statusCode)
		return
	}

	if !isAdmin {
		loggers.ErrorLogger.Println("handleCreateCourse s.usersSvc.IsAdmin isAdmin:", isAdmin)
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	var course *types.Course
	err = json.NewDecoder(r.Body).Decode(&course)
	if err != nil {
		loggers.ErrorLogger.Println("handleCreateCourse json.NewDecoder error:", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	statusCode, err = s.usersSvc.CreateCourse(r.Context(), course)
	if err != nil {
		loggers.ErrorLogger.Println("handleCreateCourse s.usersSvc.CreateCourse error:", err)
		http.Error(w, http.StatusText(statusCode), statusCode)
		return
	}

	err = jsoner(w, course, statusCode)
	if err != nil {
		loggers.ErrorLogger.Println("handleCreateCourse jsoner error:", err)
		return
	}
	loggers.InfoLogger.Println("handleCreateCourse finished with any error!")
}

//handleGetUserByID
func (s *Server) handleGetUserByID(w http.ResponseWriter, r *http.Request) {
	loggers, err := middleware.GetLoggers(r.Context())
	if err != nil {
		log.Println("LOGGERS DON'T WORK!!!")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	loggers.InfoLogger.Println("handleGetUserByID started")

	adminId, err := middleware.Authentication(r.Context())
	if err != nil {
		loggers.ErrorLogger.Println("handleMakeAdmin middleware.Authentication error:", err)
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	isAdmin, statusCode, err := s.usersSvc.IsAdmin(r.Context(), adminId)
	if err != nil {
		loggers.ErrorLogger.Println("handleMakeAdmin s.usersSvc.IsAdmin error:", err)
		http.Error(w, http.StatusText(statusCode), statusCode)
		return
	}

	if !isAdmin {
		loggers.ErrorLogger.Println("handleMakeAdmin s.usersSvc.IsAdmin isAdmin:", isAdmin)
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	idParam, ok := mux.Vars(r)["id"]
	if !ok {
		loggers.ErrorLogger.Println("handleGetUserByID mux.Vars(r) ID not found")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		loggers.ErrorLogger.Println("handleGetUserByID strconv.ParseInt error:", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	user, statusCode, err := s.usersSvc.GetUserByID(r.Context(), id)
	if err != nil {
		loggers.ErrorLogger.Println("handleGetUserByID s.usersSvc.GetUserByID error:", err)
		http.Error(w, http.StatusText(statusCode), statusCode)
		return
	}
	err = jsoner(w, user, statusCode)
	if err != nil {
		loggers.ErrorLogger.Println("handleGetUserByID jsoner error:", err)
		return
	}
	loggers.InfoLogger.Println("handleGetUserByID finished with any error!")
}

//handleGetAllUsers
func (s *Server) handleGetAllUsers(w http.ResponseWriter, r *http.Request) {
	loggers, err := middleware.GetLoggers(r.Context())
	if err != nil {
		log.Println("LOGGERS DON'T WORK!!!")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	loggers.InfoLogger.Println("handleGetAllUsers started")

	adminId, err := middleware.Authentication(r.Context())
	if err != nil {
		loggers.ErrorLogger.Println("handleMakeAdmin middleware.Authentication error:", err)
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	isAdmin, statusCode, err := s.usersSvc.IsAdmin(r.Context(), adminId)
	if err != nil {
		loggers.ErrorLogger.Println("handleMakeAdmin s.usersSvc.IsAdmin error:", err)
		http.Error(w, http.StatusText(statusCode), statusCode)
		return
	}
	if !isAdmin {
		loggers.ErrorLogger.Println("handleMakeAdmin s.usersSvc.IsAdmin isAdmin:", isAdmin)
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	usersArr, statusCode, err := s.usersSvc.GetAllUsers(r.Context())
	if err != nil {
		loggers.ErrorLogger.Println("handleGetAllUsers s.usersSvc.GetAllUsers error:", err)
		http.Error(w, http.StatusText(statusCode), statusCode)
		return
	}

	err = jsoner(w, usersArr, statusCode)
	if err != nil {
		loggers.ErrorLogger.Println("handleGetAllUsers jsoner error:", err)
		return
	}
	loggers.InfoLogger.Println("handleGetAllUsers finished with any error!")
}

//handleCourseSubscribes returns all users that subscribed to the course
func (s *Server) handleCourseSubscribes(w http.ResponseWriter, r *http.Request) {
	loggers, err := middleware.GetLoggers(r.Context())
	if err != nil {
		log.Println("LOGGERS DON'T WORK!!!")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	adminId, err := middleware.Authentication(r.Context())
	if err != nil {
		loggers.ErrorLogger.Println("handleCourseSubscribes middleware.Authentication error:", err)
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	isAdmin, statusCode, err := s.usersSvc.IsAdmin(r.Context(), adminId)
	if err != nil {
		loggers.ErrorLogger.Println("handleCourseSubscribes s.usersSvc.IsAdmin error:", err)
		http.Error(w, http.StatusText(statusCode), statusCode)
		return
	}

	if !isAdmin {
		loggers.ErrorLogger.Println("handleCourseSubscribes s.usersSvc.IsAdmin isAdmin:", isAdmin)
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	courseIDParam, ok := mux.Vars(r)["id"]
	if !ok {
		loggers.ErrorLogger.Println("handleCourseSubscribes mux.Vars(r) ID not found")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	courseID, err := strconv.ParseInt(courseIDParam, 10, 64)
	if err != nil {
		loggers.ErrorLogger.Println("handleCourseSubscribes strconv.ParseInt error:", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	usersArr, statusCode, err := s.usersSvc.CourseSubscribes(r.Context(), courseID)
	if err != nil {
		loggers.ErrorLogger.Println("handleCourseSubscribes s.usersSvc.CourseSubscribes error:", err)
		http.Error(w, http.StatusText(statusCode), statusCode)
		return
	}

	err = jsoner(w, usersArr, statusCode)
	if err != nil {
		loggers.ErrorLogger.Println("handleCourseSubscribes jsoner error:", err)
		return
	}
	loggers.InfoLogger.Println("handleCourseSubscribes finished with any error!")
}

//handleUserCourses returns all courses of a user
func (s *Server) handleUserCourses(w http.ResponseWriter, r *http.Request) {
	loggers, err := middleware.GetLoggers(r.Context())
	if err != nil {
		log.Println("LOGGERS DON'T WORK!!!")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	userIDParam, ok := mux.Vars(r)["id"]
	if !ok {
		loggers.ErrorLogger.Println("handleUserCourses mux.Vars(r) ID not found")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	userID, err := strconv.ParseInt(userIDParam, 10, 64)
	if err != nil {
		loggers.ErrorLogger.Println("handleUserCourses strconv.ParseInt error:", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	coursesArr, statusCode, err := s.usersSvc.UserCourses(r.Context(), userID)
	if err != nil {
		loggers.ErrorLogger.Println("handleUserCourses s.usersSvc.UserCourses error:", err)
		http.Error(w, http.StatusText(statusCode), statusCode)
		return
	}

	err = jsoner(w, coursesArr, statusCode)
	if err != nil {
		loggers.ErrorLogger.Println("handleUserCourses jsoner error:", err)
		return
	}

	loggers.InfoLogger.Println("handleUserCourses finished with any error!")
}
