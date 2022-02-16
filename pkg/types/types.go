package types

import (
	"log"
	"os"
	"time"
)

var (
	// Logger for INFO messages
	InfoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	// Logger for ERROR messages
	ErrorLog = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Llongfile)
)

// Type RegInfo is structure for registration info
type RegInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SubscribeInfo struct {
	UserID   int64 `json:"user_id"`
	CourseID int64 `json:"course_id"`
}

// Type User is structure with user data
type User struct {
	ID       int64     `json:"id"`
	Username string    `json:"username"`
	Password string    `json:"password"`
	IsAdmin  bool      `json:"is_admin"`
	Active   bool      `json:"active"`
	Created  time.Time `json:"created"`
}

// Type TokenInfo is structure of token info
type TokenInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Type Token is structure for token
type Token struct {
	Token   string    `json:"token"`
	UserID  int64     `json:"user_id"`
	Expires time.Time `json:"expires"`
	Created time.Time `json:"created"`
}

// MakeAdminInfo contains information for s.custumersSvc.MakeAdmin method
type MakeAdminInfo struct {
	ID          int64 `json:"id"`
	AdminStatus bool  `json:"adminStatus"`
}

// Course is structure for course
type Course struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Status      string    `json:"status"`
	Description string    `json:"description"`
	Created     time.Time `json:"created"`
}
