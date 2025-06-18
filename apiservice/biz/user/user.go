// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: MIT

package user

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

// User contains user information
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
}

// LoginRequest contains login request parameters
type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// RegisterRequest contains registration request parameters
type RegisterRequest struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// Claims represents JWT claims
type Claims struct {
	UserID int `json:"user_id"`
	jwt.StandardClaims
}

var (
	// jwtSecret load from config
	jwtSecret  string
	users      = make(map[int]User)
	nextUserID = 1
	mu         sync.RWMutex
)

func InitJWTSecret(jwt string) {
	jwtSecret = jwt
}

// RegisterHandler handles user registration requests.
// It validates the registration request, checks if the username already exists,
// hashes the password, creates a new user, saves the user to memory,
// generates a JWT token, and returns the user information and token.
//
// Parameters:
//   - ctx: The context for the request.
//   - c: A pointer to the Hertz `app.RequestContext`, used to handle HTTP requests and responses.
func RegisterHandler(ctx context.Context, c *app.RequestContext) {
	var registerRequest RegisterRequest
	err := c.BindAndValidate(&registerRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, "Invalid request body")
		return
	}

	// check if username exist
	mu.RLock()
	for _, user := range users {
		if user.Username == registerRequest.Username {
			mu.RUnlock()
			c.JSON(http.StatusConflict, "Username already exists")
			return
		}
	}
	mu.RUnlock()

	// encrypt password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Error hashing password")
		return
	}

	// create user
	mu.Lock()
	userID := nextUserID
	nextUserID++
	mu.Unlock()

	user := User{
		ID:       userID,
		Username: registerRequest.Username,
		Email:    registerRequest.Email,
		Password: string(hashedPassword),
	}

	// save user to memory
	mu.Lock()
	users[user.ID] = user
	mu.Unlock()

	// generate JWT token
	token, err := generateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Error generating token")
		return
	}

	// return user info and token
	response := map[string]interface{}{
		"user":  user,
		"token": token,
	}

	c.SetContentType("application/json; charset=utf-8")
	c.JSON(http.StatusOK, response)
}

// LoginHandler handles user login requests.
// It validates the login request, checks if the user exists,
// verifies the password, generates a JWT token, and returns the user information and token.
//
// Parameters:
//   - ctx: The context for the request.
//   - c: A pointer to the Hertz `app.RequestContext`, used to handle HTTP requests and responses.
func LoginHandler(ctx context.Context, c *app.RequestContext) {
	var loginRequest LoginRequest
	err := c.BindAndValidate(&loginRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, "Invalid request body")
		return
	}
	// get user info
	var user User
	var found bool

	mu.RLock()
	for _, u := range users {
		if u.Username == loginRequest.Username {
			user = u
			found = true
			break
		}
	}
	mu.RUnlock()

	if !found {
		c.JSON(http.StatusUnauthorized, "Invalid username or password")
		return
	}

	// verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, "Invalid username or password")
		return
	}

	// generate JWT token
	token, err := generateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Error generating token")
		return
	}

	// return user info and token
	response := map[string]interface{}{
		"user":  user,
		"token": token,
	}

	c.SetContentType("application/json; charset=utf-8")
	c.JSON(http.StatusOK, response)
}

// GetUserHandler retrieves user information based on the JWT token in the request header.
// It extracts the token, decodes it, verifies its validity,
// fetches the user information from memory, and returns the user information.
//
// Parameters:
//   - ctx: The context for the request.
//   - c: A pointer to the Hertz `app.RequestContext`, used to handle HTTP requests and responses.
func GetUserHandler(ctx context.Context, c *app.RequestContext) {
	// get token from request header
	authHeader := string(c.GetHeader("Authorization"))
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, "Authorization header is missing")
		return
	}

	// decode token
	tokenString := authHeader[7:] // remove "Bearer " prefix
	claims := &Claims{}

	tkn, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	if err != nil || !tkn.Valid {
		c.JSON(http.StatusUnauthorized, "Invalid token")
		return
	}

	// get user info from memory by userID
	mu.RLock()
	user, found := users[claims.UserID]
	mu.RUnlock()

	if !found {
		c.JSON(http.StatusNotFound, "User not found")
		return
	}

	// return user info
	c.SetContentType("application/json; charset=utf-8")
	c.JSON(http.StatusOK, user)
}

// generateToken generates a JWT token for the given user ID.
// The token expires in 24 hours.
//
// Parameters:
//   - userID: The ID of the user for whom the token is generated.
//
// Returns:
//   - string: The generated JWT token.
//   - error: An error object if an unexpected error occurs during token generation.
func generateToken(userID int) (string, error) {
	// set token expiration time to 24 hours
	expirationTime := time.Now().Add(24 * time.Hour)

	// create claims
	claims := &Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// create a new token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// sign the token with the secret key
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

var ErrInvalidToken = errors.New("invalid token")

// GetUserIDFromContext extracts the user ID from the JWT token in the request header.
// It checks if the authorization header is present, verifies the token's format and validity,
// and returns the user ID if successful.
//
// Parameters:
//   - c: A pointer to the Hertz `app.RequestContext`, used to handle HTTP requests and responses.
//
// Returns:
//   - int64: The user ID extracted from the token.
//   - error: An error object if an unexpected error occurs during extraction.
func GetUserIDFromContext(c *app.RequestContext) (int64, error) {
	authHeader := string(c.GetHeader("Authorization"))
	if authHeader == "" {
		return 0, errors.New("authorization header is missing")
	}

	// check Bearer prefix
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return 0, ErrInvalidToken
	}

	// get token from header
	tokenString := authHeader[7:]
	claims := &Claims{}

	// decoding the token
	tkn, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	if err != nil || !tkn.Valid {
		return 0, ErrInvalidToken
	}

	return int64(claims.UserID), nil
}
