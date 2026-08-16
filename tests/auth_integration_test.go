package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ararext/Go-JWT-Authentication-API/internal/handler"
	"github.com/ararext/Go-JWT-Authentication-API/internal/middleware"
	"github.com/ararext/Go-JWT-Authentication-API/internal/repository/mock"
	"github.com/ararext/Go-JWT-Authentication-API/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

const testJWTSecret = "test-secret"

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)

	userRepo := mock.NewUserRepository()
	authService := service.NewAuthService(userRepo, testJWTSecret, 15*time.Minute, 7*24*time.Hour)
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userRepo)

	router := gin.New()
	v1 := router.Group("/api/v1")

	auth := v1.Group("/auth")
	auth.POST("/signup", authHandler.Signup)
	auth.POST("/login", authHandler.Login)

	users := v1.Group("/users")
	users.Use(middleware.JWTAuth(testJWTSecret))
	users.GET("/me", userHandler.Me)

	return router
}

func doJSONRequest(router *gin.Engine, method, path string, body interface{}) *httptest.ResponseRecorder {
	var reqBody bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&reqBody).Encode(body); err != nil {
			panic(err) // test helper — encoding failure means the test itself is broken
		}
	}

	req := httptest.NewRequest(method, path, &reqBody)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

func TestSignup_ValidRequest_Returns201(t *testing.T) {
	router := setupTestRouter()

	w := doJSONRequest(router, http.MethodPost, "/api/v1/auth/signup", map[string]string{
		"name":     "Ararext",
		"email":    "int-test1@example.com",
		"password": "securepass123",
	})

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestSignup_DuplicateEmail_Returns409(t *testing.T) {
	router := setupTestRouter()

	payload := map[string]string{
		"name":     "Ararext",
		"email":    "int-test2@example.com",
		"password": "securepass123",
	}

	w1 := doJSONRequest(router, http.MethodPost, "/api/v1/auth/signup", payload)
	assert.Equal(t, http.StatusCreated, w1.Code)

	w2 := doJSONRequest(router, http.MethodPost, "/api/v1/auth/signup", payload)
	assert.Equal(t, http.StatusConflict, w2.Code)
}

func TestSignup_InvalidEmail_Returns400(t *testing.T) {
	router := setupTestRouter()

	w := doJSONRequest(router, http.MethodPost, "/api/v1/auth/signup", map[string]string{
		"name":     "Ararext",
		"email":    "not-an-email",
		"password": "securepass123",
	})

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestSignup_WeakPassword_Returns400(t *testing.T) {
	router := setupTestRouter()

	w := doJSONRequest(router, http.MethodPost, "/api/v1/auth/signup", map[string]string{
		"name":     "Ararext",
		"email":    "int-test3@example.com",
		"password": "short",
	})

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestLogin_ValidCredentials_Returns200(t *testing.T) {
	router := setupTestRouter()

	signupPayload := map[string]string{
		"name":     "Ararext",
		"email":    "int-test4@example.com",
		"password": "securepass123",
	}
	doJSONRequest(router, http.MethodPost, "/api/v1/auth/signup", signupPayload)

	w := doJSONRequest(router, http.MethodPost, "/api/v1/auth/login", map[string]string{
		"email":    "int-test4@example.com",
		"password": "securepass123",
	})

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestLogin_WrongPassword_Returns401(t *testing.T) {
	router := setupTestRouter()

	doJSONRequest(router, http.MethodPost, "/api/v1/auth/signup", map[string]string{
		"name":     "Ararext",
		"email":    "int-test5@example.com",
		"password": "securepass123",
	})

	w := doJSONRequest(router, http.MethodPost, "/api/v1/auth/login", map[string]string{
		"email":    "int-test5@example.com",
		"password": "wrongpassword",
	})

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestLogin_UserNotFound_Returns401(t *testing.T) {
	router := setupTestRouter()

	w := doJSONRequest(router, http.MethodPost, "/api/v1/auth/login", map[string]string{
		"email":    "ghost@example.com",
		"password": "whatever123",
	})

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestProtectedRoute_NoToken_Returns401(t *testing.T) {
	router := setupTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/me", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestProtectedRoute_MalformedToken_Returns401(t *testing.T) {
	router := setupTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/me", nil)
	req.Header.Set("Authorization", "Bearer not-a-real-token")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestProtectedRoute_ValidToken_Returns200(t *testing.T) {
	router := setupTestRouter()

	signupPayload := map[string]string{
		"name":     "Ararext",
		"email":    "int-test6@example.com",
		"password": "securepass123",
	}
	w := doJSONRequest(router, http.MethodPost, "/api/v1/auth/signup", signupPayload)
	assert.Equal(t, http.StatusCreated, w.Code)

	var signupResp struct {
		Data struct {
			AccessToken string `json:"accessToken"`
		} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &signupResp); err != nil {
		t.Fatalf("failed to unmarshal signup response: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/me", nil)
	req.Header.Set("Authorization", "Bearer "+signupResp.Data.AccessToken)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req)

	assert.Equal(t, http.StatusOK, w2.Code)
}
