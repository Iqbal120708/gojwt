package handler

import (
    "testing"
    "net/http"
    "gojwt/internal/entity"
    "gojwt/internal/usecase"
    "net/http/httptest"
    "encoding/json"
    "github.com/golang-jwt/jwt/v5"
    "time"
    "errors"
    "context"
)

func TestProfile_Success(t *testing.T) {
    id := int64(1)
    role := "regular"
    
    user := entity.User{
        ID:       &id,
        Name:     "test",
        Email:    "test@gmail.com",
        Role:     &role,
        Password: "secret123",
    }
    
    mockUC := &usecase.MockUserUseCase{
        UserFn: func(email string) (*entity.User, error) {
            return &user, nil
        },
    }
    
	handler := NewUserHandler(mockUC)

	req := httptest.NewRequest(http.MethodGet, "/api/profile", nil)
	
	accessExp := time.Now().Add(15 * time.Minute)
	claims := &entity.AccessClaims{
	    UserID: 1,
        Email:  "test@gmail.com",
        Role:   "regular",
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(accessExp),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            Issuer:    "myapp",
            Subject:   "access_token",
        },
	}
	
	ctx := context.WithValue(req.Context(), "user_claims", claims)
	req = req.WithContext(ctx)
	
	rr := httptest.NewRecorder()
	
	handler.Profile(rr, req)
	
	if rr.Code != http.StatusOK {
	    t.Errorf("expected %d got %d", http.StatusOK, rr.Code)
	}
	
	var resBody entity.User
	json.NewDecoder(rr.Body).Decode(&resBody)
	
	if *resBody.ID != 1 {
	    t.Errorf("expected id to be 1, got %d", resBody.ID)
	}
}

func TestProfile_Unauthorized(t *testing.T) {
    mockUC := &usecase.MockUserUseCase{}
	handler := NewUserHandler(mockUC)

	req := httptest.NewRequest(http.MethodGet, "/api/profile", nil)
	
	rr := httptest.NewRecorder()
	
	handler.Profile(rr, req)
	
	if rr.Code != http.StatusUnauthorized {
	    t.Errorf("expected %d got %d", http.StatusUnauthorized, rr.Code)
	}
}

func TestProfile_ValidationError(t *testing.T) {
    mockUC := &usecase.MockUserUseCase{
        UserFn: func(email string) (*entity.User, error) {
			return nil, &entity.ValidationError{
				Field: "email",
				Message: "email invalid",
			}
		},
    }
	handler := NewUserHandler(mockUC)
	
	req := httptest.NewRequest(http.MethodGet, "/api/profile", nil)
    
    accessExp := time.Now().Add(15 * time.Minute)
	claims := &entity.AccessClaims{
	    UserID: 1,
        Email:  "test@gmail.com",
        Role:   "regular",
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(accessExp),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            Issuer:    "myapp",
            Subject:   "access_token",
        },
	}
	
	ctx := context.WithValue(req.Context(), "user_claims", claims)
	req = req.WithContext(ctx)
	
	rr := httptest.NewRecorder()

	handler.Profile(rr, req)
	
	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected %d got %d", http.StatusBadRequest, rr.Code)
	}

	var resp entity.ValidationError
	json.NewDecoder(rr.Body).Decode(&resp)

	if resp.Field != "email" {
		t.Errorf("expected field email got %s", resp.Field)
	}
}

func TestProfile_ServerError(t *testing.T) {
    mockUC := &usecase.MockUserUseCase{
        UserFn: func(email string) (*entity.User, error) {
			return nil, errors.New("server error")
		},
    }
	handler := NewUserHandler(mockUC)
	
	req := httptest.NewRequest(http.MethodGet, "/api/profile", nil)
    
    accessExp := time.Now().Add(15 * time.Minute)
	claims := &entity.AccessClaims{
	    UserID: 1,
        Email:  "test@gmail.com",
        Role:   "regular",
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(accessExp),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            Issuer:    "myapp",
            Subject:   "access_token",
        },
	}
	
	ctx := context.WithValue(req.Context(), "user_claims", claims)
	req = req.WithContext(ctx)
	
	rr := httptest.NewRecorder()

	handler.Profile(rr, req)
	
	if rr.Code != http.StatusInternalServerError {
		t.Errorf("expected %d got %d", http.StatusInternalServerError, rr.Code)
	}

    type ResBody struct {
        Message string `json:"message"`
    }	
	var resBody ResBody
	json.NewDecoder(rr.Body).Decode(&resBody)
	
	if resBody.Message != "internal server error" {
		t.Errorf("expected error message internal server error got %s", resBody.Message)
	}
}