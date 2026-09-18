package authservice

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type userContextKey struct{}

var claimsKey = userContextKey{}

type User struct {
	ID           string `json:"id"`
	Email        string `json:"email"`
	PasswordHash string `json:"-"`
	Role         string `json:"role"`
}

type Claims struct {
	Subject string `json:"sub"`
	Role    string `json:"role"`
	Expires int64  `json:"exp"`
}

type AuthService struct {
	secret []byte
	users  map[string]*User // keyed by email
	byID   map[string]*User // keyed by ID
	mu     sync.RWMutex
	router http.Handler
}

func NewAuthService(secret []byte) *AuthService {
	s := &AuthService{
		secret: secret,
		users:  make(map[string]*User),
		byID:   make(map[string]*User),
	}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /auth/register", s.handleRegister)
	mux.HandleFunc("POST /auth/login", s.handleLogin)

	// Protected routes
	authMW := s.authMiddleware
	adminMW := s.requireRole("admin")

	mux.Handle("GET /api/profile", authMW(http.HandlerFunc(s.handleProfile)))
	mux.Handle("GET /api/admin/metrics", authMW(adminMW(http.HandlerFunc(s.handleAdminMetrics))))

	s.router = mux
	return s
}

func (s *AuthService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

// Password hashing & token helpers
func generateRandomID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func (s *AuthService) createToken(claims Claims) (string, error) {
	headerEnc := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"HS256","typ":"JWT"}`))
	payloadBytes, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}
	payloadEnc := base64.RawURLEncoding.EncodeToString(payloadBytes)

	signedData := headerEnc + "." + payloadEnc
	h := hmac.New(sha256.New, s.secret)
	h.Write([]byte(signedData))
	sigEnc := base64.RawURLEncoding.EncodeToString(h.Sum(nil))

	return signedData + "." + sigEnc, nil
}

func (s *AuthService) verifyToken(tokenStr string) (*Claims, error) {
	parts := strings.Split(tokenStr, ".")
	if len(parts) != 3 {
		return nil, errors.New("malformed token")
	}

	signedData := parts[0] + "." + parts[1]
	expectedSig, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return nil, errors.New("invalid signature encoding")
	}

	h := hmac.New(sha256.New, s.secret)
	h.Write([]byte(signedData))
	calculatedSig := h.Sum(nil)

	if subtle.ConstantTimeCompare(expectedSig, calculatedSig) != 1 {
		return nil, errors.New("invalid signature")
	}

	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, errors.New("invalid payload encoding")
	}

	var claims Claims
	if err := json.Unmarshal(payloadBytes, &claims); err != nil {
		return nil, errors.New("invalid claims json")
	}

	if claims.Expires > 0 && time.Now().Unix() > claims.Expires {
		return nil, errors.New("token expired")
	}

	return &claims, nil
}

// Handlers
type registerRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

func (s *AuthService) handleRegister(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid json body"}`, http.StatusBadRequest)
		return
	}

	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	if req.Email == "" || !strings.Contains(req.Email, "@") {
		http.Error(w, `{"error":"valid email is required"}`, http.StatusBadRequest)
		return
	}
	if len(req.Password) < 8 {
		http.Error(w, `{"error":"password must be at least 8 characters"}`, http.StatusBadRequest)
		return
	}
	if req.Role == "" {
		req.Role = "user"
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.users[req.Email]; exists {
		http.Error(w, `{"error":"user already exists"}`, http.StatusConflict)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, `{"error":"failed to hash password"}`, http.StatusInternalServerError)
		return
	}

	newUser := &User{
		ID:           generateRandomID(),
		Email:        req.Email,
		PasswordHash: string(hash),
		Role:         req.Role,
	}

	s.users[req.Email] = newUser
	s.byID[newUser.ID] = newUser

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(newUser)
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token     string `json:"token"`
	ExpiresIn int64  `json:"expires_in"`
}

func (s *AuthService) handleLogin(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid json body"}`, http.StatusBadRequest)
		return
	}

	req.Email = strings.TrimSpace(strings.ToLower(req.Email))

	s.mu.RLock()
	user, exists := s.users[req.Email]
	s.mu.RUnlock()

	// Generic error response to prevent user enumeration
	invalidCredentialsMsg := `{"error":"invalid email or password"}`
	if !exists {
		http.Error(w, invalidCredentialsMsg, http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		http.Error(w, invalidCredentialsMsg, http.StatusUnauthorized)
		return
	}

	expires := time.Now().Add(15 * time.Minute).Unix()
	token, err := s.createToken(Claims{
		Subject: user.ID,
		Role:    user.Role,
		Expires: expires,
	})
	if err != nil {
		http.Error(w, `{"error":"failed to create token"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(loginResponse{
		Token:     token,
		ExpiresIn: 900,
	})
}

// Middleware
func (s *AuthService) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := s.verifyToken(tokenStr)
		if err != nil {
			http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), claimsKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *AuthService) requireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			val := r.Context().Value(claimsKey)
			if val == nil {
				http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
				return
			}
			claims := val.(*Claims)
			if claims.Role != role {
				http.Error(w, `{"error":"forbidden: insufficient role"}`, http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func (s *AuthService) handleProfile(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(claimsKey).(*Claims)

	s.mu.RLock()
	user, exists := s.byID[claims.Subject]
	s.mu.RUnlock()

	if !exists {
		http.Error(w, `{"error":"user not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(user)
}

func (s *AuthService) handleAdminMetrics(w http.ResponseWriter, r *http.Request) {
	s.mu.RLock()
	total := len(s.users)
	s.mu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"total_users":%d}`, total)
}
