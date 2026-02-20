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

func TestCreate_InvalidContentType(t *testing.T) {
	mockUC := &usecase.MockUserUseCase{}
	handler := NewUserHandler(mockUC)

	body := strings.NewReader(`{"name":"andi"}`)
	req := httptest.NewRequest(http.MethodPost, "/signup", body)
	req.Header.Set("Content-Type", "text/plain")

	rr := httptest.NewRecorder()

	handler.Create(rr, req)

	if rr.Code != http.StatusUnsupportedMediaType {
		t.Errorf("expected %d got %d", http.StatusUnsupportedMediaType, rr.Code)
	}
}

func TestCreate_InvalidBody(t *testing.T) {
	mockUC := &usecase.MockUserUseCase{}
	handler := NewUserHandler(mockUC)

	body := strings.NewReader(``)
	req := httptest.NewRequest(http.MethodPost, "/signup", body)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler.Create(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected %d got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestCreate_ValidationError(t *testing.T) {
    mockUC := &usecase.MockUserUseCase{
        CreateFn: func(user *entity.User) (*entity.User, error) {
			return nil, &entity.ValidationError{
				Field: "email",
				Message: "email invalid",
			}
		},
    }
	handler := NewUserHandler(mockUC)
	
	body := strings.NewReader(`{"email":"invalidemail.com"}`)
	req := httptest.NewRequest(http.MethodPost, "/signup", body)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler.Create(rr, req)
	
	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected %d got %d", http.StatusBadRequest, rr.Code)
	}

	var resp entity.ValidationError
	json.NewDecoder(rr.Body).Decode(&resp)

	if resp.Field != "email" {
		t.Errorf("expected field email got %s", resp.Field)
	}
}

func TestCreate_InternalServerError(t *testing.T) {
    mockUC := &usecase.MockUserUseCase{
        CreateFn: func(user *entity.User) (*entity.User, error) {
			return nil, errors.New("server error")
		},
    }
	handler := NewUserHandler(mockUC)
	
	body := strings.NewReader(`{"email":"test@email.com"}`)
	req := httptest.NewRequest(http.MethodPost, "/signup", body)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler.Create(rr, req)
	
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

func TestCreate_Success(t *testing.T) {
    mockUC := &usecase.MockUserUseCase{
        CreateFn: func(user *entity.User) (*entity.User, error) {
			return &entity.User{
			    Name: "test",
			    Email: "test@email.com",
			    Password: "test123",
			}, nil
		},
    }
	handler := NewUserHandler(mockUC)
	
	body := strings.NewReader(`{"email":"test@email.com"}`)
	req := httptest.NewRequest(http.MethodPost, "/signup", body)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler.Create(rr, req)
	
	if rr.Code != http.StatusCreated {
		t.Errorf("expected %d got %d", http.StatusCreated, rr.Code)
	}

    type ResBody struct {
        Message string `json:"message"`
    }	
	var resBody ResBody
	json.NewDecoder(rr.Body).Decode(&resBody)
	
	if resBody.Message != "user created successfully" {
		t.Errorf("expected message user created successfully got %s", resBody.Message)
	}
}