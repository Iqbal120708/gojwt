package handler

import (
    "testing"
    "net/http"
    "gojwt/internal/entity"
    "gojwt/internal/usecase"
    "strings"
    "net/http/httptest"
    "encoding/json"
    "errors"
)

func TestLogin_InvalidContentType(t *testing.T) {
	mockUC := &usecase.MockUserUseCase{}
	handler := NewUserHandler(mockUC)

	body := strings.NewReader(`{"email":"test@gmail.com", "password":"test123"}`)
	req := httptest.NewRequest(http.MethodPost, "/login", body)
	req.Header.Set("Content-Type", "text/plain")

	rr := httptest.NewRecorder()

	handler.Login(rr, req)

	if rr.Code != http.StatusUnsupportedMediaType {
		t.Errorf("expected %d got %d", http.StatusUnsupportedMediaType, rr.Code)
	}
}

func TestLogin_InvalidBody(t *testing.T) {
	mockUC := &usecase.MockUserUseCase{}
	handler := NewUserHandler(mockUC)

	body := strings.NewReader(``)
	req := httptest.NewRequest(http.MethodPost, "/login", body)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler.Login(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected %d got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestLogin_ValidationError(t *testing.T) {
    mockUC := &usecase.MockUserUseCase{
        LoginFn: func(login *entity.Login) (*entity.Token, error) {
			return nil, &entity.ValidationError{
				Field: "email",
				Message: "email invalid",
			}
		},
    }
	handler := NewUserHandler(mockUC)
	
	body := strings.NewReader(`{"email":"invalidemail.com", "password":"test123"}`)
	req := httptest.NewRequest(http.MethodPost, "/login", body)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler.Login(rr, req)
	
	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected %d got %d", http.StatusBadRequest, rr.Code)
	}

	var resp entity.ValidationError
	json.NewDecoder(rr.Body).Decode(&resp)

	if resp.Field != "email" {
		t.Errorf("expected field email got %s", resp.Field)
	}
}

func TestLogin_InternalServerError(t *testing.T) {
    mockUC := &usecase.MockUserUseCase{
        LoginFn: func(login *entity.Login) (*entity.Token, error) {
			return nil, errors.New("server error")
		},
    }
	handler := NewUserHandler(mockUC)
	
	body := strings.NewReader(`{"email":"test@gmail.com", "password":"test123"}`)
	req := httptest.NewRequest(http.MethodPost, "/login", body)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler.Login(rr, req)
	
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

func TestLogin_Success(t *testing.T) {
    mockUC := &usecase.MockUserUseCase{
        LoginFn: func(login *entity.Login) (*entity.Token, error) {
			return &entity.Token{
			    Access: "access.token",
			    Refresh: "refresh.token",
			}, nil
		},
    }
	handler := NewUserHandler(mockUC)
	
	body := strings.NewReader(`{"email":"test@gmail.com", "password":"test123"}`)
	req := httptest.NewRequest(http.MethodPost, "/login", body)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler.Login(rr, req)
	
	if rr.Code != http.StatusOK {
		t.Errorf("expected %d got %d", http.StatusOK, rr.Code)
	}

    var resBody entity.Token
    json.NewDecoder(rr.Body).Decode(&resBody)
	
	if resBody.Access != "access.token" {
		t.Errorf("expected access token to be access.token got %s", resBody.Access)
	}
	
	if resBody.Refresh != "" {
		t.Errorf("expected refresh token to be empty, got %q", resBody.Refresh)
	}
	
	setCookie := rr.Header().Get("Set-Cookie")

    if !strings.Contains(setCookie, "refresh_token=") {
        t.Fatalf("expected refresh_token cookie, got: %s", setCookie)
    }
}