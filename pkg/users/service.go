package users

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/SYSTEMTerror/GoEDU/pkg/types"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

var (
	//ErrNotFound is returned when a user is not found
	ErrNotFound = errors.New("user not found")
	//ErrInvalidPassword is returned when password is incorrect
	ErrInvalidPassword = errors.New("invalid password")
	//ErrInternal is returned when an internal error occurs
	ErrInternal = errors.New("internal error")
	//ErrExpired is returned when token is expired
	ErrExpired = errors.New("expired")
	//ErrEmptyPassword is returned when password is empty
	ErrEmptyPassword = errors.New("empty password")
)

//Service is a users service
type Service struct {
	pool *pgxpool.Pool
}

//NewService creates new users service
func NewService(pool *pgxpool.Pool) *Service {
	return &Service{pool: pool}
}

// RegisterUser registers user
func (s *Service) RegisterUser(ctx context.Context, item *types.RegInfo) (*types.User, int, error) {
	user := &types.User{}

	hash, err := bcrypt.GenerateFromPassword([]byte(item.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Save bcrypt.GenerateFromPassword Error:", err)
		return nil, http.StatusInternalServerError, ErrInternal
	}

	item.Password = string(hash)
	err = s.pool.QueryRow(ctx, `
			INSERT INTO users (username, password) VALUES ($1, $2)
			ON CONFLICT (username) DO NOTHING
			RETURNING id, username, password, active, created
		`, item.Username, item.Password).Scan(
		&user.ID, &user.Username, &user.Password,
		&user.Active, &user.Created)
	if err != nil {
		log.Println("Register s.pool.QueryRow error:", err)
		return nil, http.StatusInternalServerError, ErrInternal
	}

	return user, http.StatusOK, nil
}

// Token generates token for user
func (s *Service) Token(ctx context.Context, item *types.TokenInfo) (*types.Token, int, error) {
	var hash string
	token := &types.Token{}
	err := s.pool.QueryRow(ctx, `SELECT id, password FROM users WHERE username = $1`, item.Username).Scan(&token.UserID, &hash)
	if err == pgx.ErrNoRows {
		log.Println("Token s.pool.QueryRow error:", err)
		return nil, http.StatusNotFound, ErrNotFound
	}
	if err != nil {
		log.Println(err)
		return nil, http.StatusInternalServerError, ErrInternal
	}

	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(item.Password))
	if err != nil {
		log.Println("Token bcrypt.CompareHashAndPassword error:", err)
		return nil, http.StatusUnauthorized, ErrInvalidPassword
	}

	buffer := make([]byte, 256)
	n, err := rand.Read(buffer)
	if n != len(buffer) || err != nil {
		log.Println("Token rand.Read len : %w (must be 256), error: %w", n, err)
		return nil, http.StatusInternalServerError, ErrInternal
	}

	token.Token = hex.EncodeToString(buffer)
	_, err = s.pool.Exec(ctx, `INSERT INTO users_tokens (user_id, token) VALUES ($1, $2)`, token.UserID, token.Token)
	if err != nil {
		log.Println("Token s.pool.Exec error:", err)
		return nil, http.StatusInternalServerError, ErrInternal
	}

	return token, http.StatusOK, nil
}

// IDByToken returns user id by token
func (s *Service) IDByToken(ctx context.Context, token string) (int64, error) {
	var id int64
	var expires time.Time

	err := s.pool.QueryRow(ctx, `SELECT user_id, expires FROM users_tokens WHERE token = $1`, token).Scan(&id, &expires)
	if err == pgx.ErrNoRows {
		log.Println("IDByToken s.pool.QueryRow No rows:", err)
		return 0, nil
	}
	if err != nil || expires.Before(time.Now()) {
		log.Println("IDByToken s.pool.QueryRow error:", err)
		return 0, ErrInternal
	}

	return id, nil
}

// IsAdmin checks if user is admin
func (s *Service) IsAdmin(ctx context.Context, id int64) (bool, int, error) {
	var isAdmin bool
	err := s.pool.QueryRow(ctx, `SELECT is_admin FROM users WHERE id = $1`, id).Scan(&isAdmin)
	if err == pgx.ErrNoRows {
		log.Println("IsAdmin s.pool.QueryRow No rows:", err)
		return false, http.StatusNotFound, ErrNotFound
	}
	if err != nil {
		log.Println("IsAdmin s.pool.QueryRow error:", err)
		return false, http.StatusInternalServerError, ErrInternal
	}

	return isAdmin, http.StatusOK, nil
}

// MakeAdmin makes user admin
func (s *Service) MakeAdmin(ctx context.Context, makeAdminInfo *types.MakeAdminInfo) (int, error) {
	_, err := s.pool.Exec(ctx, `UPDATE users SET is_admin = $2 WHERE id = $1`, makeAdminInfo.ID, makeAdminInfo.AdminStatus)
	if err != nil {
		log.Println("MakeAdmin s.pool.Exec error:", err)
		return http.StatusInternalServerError, ErrInternal
	}

	return http.StatusOK, nil
}

// CreateCourse creates course
func (s *Service) CreateCourse(ctx context.Context, course *types.Course) (int, error) {
	if course.ID == 0 {
		_, err := s.pool.Exec(ctx, `
			INSERT INTO courses (name, description, status) VALUES ($1, $2, $3)
		`, course.Name, course.Description, course.Status)
		if err != nil {
			log.Println("CreateCourse s.pool.Exec error:", err)
			return http.StatusInternalServerError, ErrInternal
		}
	}
	if course.ID != 0 {
		_, err := s.pool.Exec(ctx, `
			UPDATE courses SET name = $1, description = $2, status = $3 WHERE id = $4
		`, course.Name, course.Description, course.Status, course.ID)
		if err != nil {
			log.Println("CreateCourse s.pool.Exec error:", err)
			return http.StatusInternalServerError, ErrInternal
		}
	}

	return http.StatusOK, nil
}

// Subscribe subscribes user to course
func (s *Service) Subscribe(ctx context.Context, subscribeInfo *types.SubscribeInfo) (int, error) {
	_, err := s.pool.Exec(ctx, `
		INSERT INTO users_courses (user_id, course_id) VALUES ($1, $2)
	`, subscribeInfo.UserID, subscribeInfo.CourseID)
	if err != nil {
		log.Println("Subscribe s.pool.Exec error:", err)
		return http.StatusInternalServerError, ErrInternal
	}

	return http.StatusOK, nil
}

// GetUserById returns user by id
func (s *Service) GetUserByID(ctx context.Context, id int64) (*types.User, int, error) {
	user := &types.User{}
	err := s.pool.QueryRow(ctx, `SELECT id, username, password, is_admin, active, created FROM users WHERE id = $1`, id).Scan(
		&user.ID, &user.Username, &user.Password, &user.IsAdmin, &user.Active, &user.Created)
	if err == pgx.ErrNoRows {
		log.Println("GetUserByID s.pool.QueryRow No rows:", err)
		return nil, http.StatusNotFound, ErrNotFound
	}
	if err != nil {
		log.Println("GetUserByID s.pool.QueryRow error:", err)
		return nil, http.StatusInternalServerError, ErrInternal
	}

	return user, http.StatusOK, nil
}

//GetAllUsers returns all users
func (s *Service) GetAllUsers(ctx context.Context) ([]*types.User, int, error) {
	var users []*types.User
	rows, err := s.pool.Query(ctx, `SELECT id, username, password, is_admin, active, created FROM users`)
	if err != nil {
		log.Println("GetAllUsers s.pool.Query error:", err)
		return nil, http.StatusInternalServerError, ErrInternal
	}
	defer rows.Close()

	for rows.Next() {
		user := &types.User{}
		err := rows.Scan(
			&user.ID, &user.Username, &user.Password, &user.IsAdmin, &user.Active, &user.Created)
		if err != nil {
			log.Println("GetAllUsers rows.Scan error:", err)
			return nil, http.StatusInternalServerError, ErrInternal
		}
		users = append(users, user)
	}

	return users, http.StatusOK, nil
}

// CourseSubscribes returns course subscribes
func (s *Service) CourseSubscribes(ctx context.Context, courseID int64) ([]*types.User, int, error) {
	var users []*types.User
	rows, err := s.pool.Query(ctx, `
		SELECT users.id, users.username, users.password, users.is_admin, users.active, users.created
		FROM users
		JOIN users_courses ON users_courses.user_id = users.id
		WHERE users_courses.course_id = $1
	`, courseID)
	if err != nil {
		log.Println("CourseSubscribes s.pool.Query error:", err)
		return nil, http.StatusInternalServerError, ErrInternal
	}

	defer rows.Close()

	for rows.Next() {
		user := &types.User{}
		err := rows.Scan(&user.ID, &user.Username, &user.Password, &user.IsAdmin, &user.Active, &user.Created)
		if err != nil {
			log.Println("CourseSubscribes rows.Scan error:", err)
			return nil, http.StatusInternalServerError, ErrInternal
		}
		users = append(users, user)
	}

	return users, http.StatusOK, nil
}

//UserCourses returns users courses
func (s *Service) UserCourses(ctx context.Context, userID int64) ([]*types.Course, int, error) {
	var courses []*types.Course
	rows, err := s.pool.Query(ctx, `
		SELECT courses.id, courses.name, courses.description, courses.status
		FROM courses
		JOIN users_courses ON users_courses.course_id = courses.id
		WHERE users_courses.user_id = $1
	`, userID)
	if err != nil {
		log.Println("UsersCourses s.pool.Query error:", err)
		return nil, http.StatusInternalServerError, ErrInternal
	}

	defer rows.Close()

	for rows.Next() {
		course := &types.Course{}
		err := rows.Scan(&course.ID, &course.Name, &course.Description, &course.Status)
		if err != nil {
			log.Println("UsersCourses rows.Scan error:", err)
			return nil, http.StatusInternalServerError, ErrInternal
		}
		courses = append(courses, course)
	}

	return courses, http.StatusOK, nil
}