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

func TestRefresh_NoCookie(t *testing.T) {
    mockUC := &usecase.MockUserUseCase{}
	handler := NewUserHandler(mockUC)

	req := httptest.NewRequest(http.MethodPost, "/refresh", nil)
	
	rr := httptest.NewRecorder()
	
	handler.RefreshToken(rr, req)
	
	if rr.Code != http.StatusUnauthorized {
	    t.Errorf("expected %d got %d", http.StatusUnauthorized, rr.Code)
	}
}

func TestRefresh_AppError(t *testing.T) {
    mockUC := &usecase.MockUserUseCase{
        RefreshTokenFn: func(refreshToken string) (*entity.Token, error) {
            return nil, &entity.AppError{
                Code: "token_invalid",
                Message: "Token is blacklisted.",
            }
        },
    }
	handler := NewUserHandler(mockUC)

	req := httptest.NewRequest(http.MethodPost, "/refresh", nil)
	
	cookie := &http.Cookie{
        Name:     "refresh_token",
        Value:    "refresh.token",
        Path:     "/refresh",
        HttpOnly: true,
        Secure:   false,
        SameSite: http.SameSiteLaxMode,
    }
    req.AddCookie(cookie)
    
	rr := httptest.NewRecorder()
	handler.RefreshToken(rr, req)
	
	if rr.Code != http.StatusBadRequest {
	    t.Errorf("expected %d got %d", http.StatusBadRequest, rr.Code)
	}
	
	var resBody entity.AppError
	json.NewDecoder(rr.Body).Decode(&resBody)
	
	if resBody.Code != "token_invalid" {
	    t.Errorf("expected AppError token_invalid, got %s", resBody.Code)
	}
}

func TestRefresh_ValidationError(t *testing.T) {
    mockUC := &usecase.MockUserUseCase{
        RefreshTokenFn: func(refreshToken string) (*entity.Token, error) {
            return nil, &entity.ValidationError{
                Field: "token_invalid",
                Message: "Token is blacklisted.",
            }
        },
    }
	handler := NewUserHandler(mockUC)

	req := httptest.NewRequest(http.MethodPost, "/refresh", nil)
	
	cookie := &http.Cookie{
        Name:     "refresh_token",
        Value:    "refresh.token",
        Path:     "/refresh",
        HttpOnly: true,
        Secure:   false,
        SameSite: http.SameSiteLaxMode,
    }
    req.AddCookie(cookie)
    
	rr := httptest.NewRecorder()
	handler.RefreshToken(rr, req)
	
	if rr.Code != http.StatusBadRequest {
	    t.Errorf("expected %d got %d", http.StatusBadRequest, rr.Code)
	}
	
	var resBody entity.ValidationError
	json.NewDecoder(rr.Body).Decode(&resBody)
	
	if resBody.Field != "token_invalid" {
	    t.Errorf("expected ValidationError token_invalid, got %s", resBody.Field)
	}
}

func TestRefresh_ServerError(t *testing.T) {
    mockUC := &usecase.MockUserUseCase{
        RefreshTokenFn: func(refreshToken string) (*entity.Token, error) {
            return nil, errors.New("server error")
        },
    }
	handler := NewUserHandler(mockUC)

	req := httptest.NewRequest(http.MethodPost, "/refresh", nil)
	
	cookie := &http.Cookie{
        Name:     "refresh_token",
        Value:    "refresh.token",
        Path:     "/refresh",
        HttpOnly: true,
        Secure:   false,
        SameSite: http.SameSiteLaxMode,
    }
    req.AddCookie(cookie)
    
	rr := httptest.NewRecorder()
	handler.RefreshToken(rr, req)
	
	if rr.Code != http.StatusInternalServerError {
	    t.Errorf("expected %d got %d", http.StatusInternalServerError, rr.Code)
	}
}

func TestRefresh_Success(t *testing.T) {
    mockUC := &usecase.MockUserUseCase{
        RefreshTokenFn: func(refreshToken string) (*entity.Token, error) {
            return &entity.Token{
			    Access: "access.token",
			    Refresh: "refresh.token",
			}, nil
        },
    }
	handler := NewUserHandler(mockUC)

	req := httptest.NewRequest(http.MethodPost, "/refresh", nil)
	
	cookie := &http.Cookie{
        Name:     "refresh_token",
        Value:    "refresh.token",
        Path:     "/refresh",
        HttpOnly: true,
        Secure:   false,
        SameSite: http.SameSiteLaxMode,
    }
    req.AddCookie(cookie)
    
	rr := httptest.NewRecorder()
	handler.RefreshToken(rr, req)
	
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

