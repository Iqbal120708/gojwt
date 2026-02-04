package handler

import (
    "encoding/json"
	"net/http"
	"gojwt/internal/usecase"
	"gojwt/internal/entity"
	"gojwt/internal/security"
	"gojwt/internal/middleware"
	"strings"
	"errors"
)

type userHandler struct {
    userUseCase usecase.UserUseCase
}

func NewUserHandler(userUseCase usecase.UserUseCase) *userHandler {
    return &userHandler{userUseCase: userUseCase}
}

func (uh *userHandler) Create(w http.ResponseWriter, r *http.Request) {
    contentType := r.Header.Get("Content-Type")

	if !strings.HasPrefix(contentType, "application/json") {
		http.Error(w, "Content-Type harus application/json", http.StatusUnsupportedMediaType)
		return
	}
	
	var user entity.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "body tidak valid", http.StatusBadRequest)
		return
	}

	_, errUser := uh.userUseCase.Create(&user)
    if errUser != nil {
        var invalidErr *entity.ValidationError
        if errors.As(errUser, &invalidErr) {
            w.Header().Set("Content-Type", "application/json")
            w.WriteHeader(http.StatusBadRequest)
            json.NewEncoder(w).Encode(*invalidErr)
            return
        } else {
            w.Header().Set("Content-Type", "application/json")
            w.WriteHeader(http.StatusBadRequest)
            json.NewEncoder(w).Encode(map[string]string{
                "message": "internal server error",
            })
            return
        }
    }
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusBadRequest)
    json.NewEncoder(w).Encode(map[string]string{
        "message": "user created successfully",
    })
}

func (uh *userHandler) Login(w http.ResponseWriter, r *http.Request) {
    contentType := r.Header.Get("Content-Type")

	if !strings.HasPrefix(contentType, "application/json") {
		http.Error(w, "Content-Type harus application/json", http.StatusUnsupportedMediaType)
		return
	}
	
	var login entity.Login
	err := json.NewDecoder(r.Body).Decode(&login)
	if err != nil {
	   // fmt.Println(err)
		http.Error(w, "body tidak valid", http.StatusBadRequest)
		return
	}

    token, errUser := uh.userUseCase.Login(&login)
    if errUser != nil {
        var invalidErr *entity.ValidationError
        if errors.As(errUser, &invalidErr) {
            w.Header().Set("Content-Type", "application/json")
            w.WriteHeader(http.StatusBadRequest)
            json.NewEncoder(w).Encode(*invalidErr)
            return
        } else {
            // fmt.Println(errUser)
            w.Header().Set("Content-Type", "application/json")
            w.WriteHeader(http.StatusBadRequest)
            json.NewEncoder(w).Encode(map[string]string{
                "message": "internal server error",
            })
            return
        }
    }
    
    security.GenerateRefreshCookie(w, token)
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(*token)
}

func (uh *userHandler) Profile(w http.ResponseWriter, r *http.Request) {
    claims, ok := r.Context().Value("user_claims").(*entity.AccessClaims)
    if !ok {
        middleware.JSONError(w, http.StatusUnauthorized, "user tidak terautentikasi")
        return
    }
    
    
    user, err := uh.userUseCase.User(claims.Email)
    if err != nil {
        var invalidErr *entity.ValidationError
        w.Header().Set("Content-Type", "application/json")
        if errors.As(err, &invalidErr) {
            w.WriteHeader(http.StatusBadRequest)
            json.NewEncoder(w).Encode(*invalidErr)
            return
        } else {
            w.WriteHeader(http.StatusInternalServerError)
            json.NewEncoder(w).Encode(map[string]string{
                "message": "internal server error",
            })
            return
        }
    }
    
    userRes := usecase.ToUserResponse(user)
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(userRes)
}

func (uh *userHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
    cookie, err := r.Cookie("refresh_token")
    if err != nil {
        w.Header().Set("Content-Type", "application/json")
        if err == http.ErrNoCookie {
            w.WriteHeader(http.StatusUnauthorized)
            json.NewEncoder(w).Encode(map[string]string{
                "message": "no cookies found",
            })
            return
        }
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]string{
            "message": "failed to read cookie",
        })
        return
    }

    refreshToken := cookie.Value
    
    token, err := uh.userUseCase.RefreshToken(refreshToken)
    if err != nil {
        w.Header().Set("Content-Type", "application/json")
        
        var appErr *entity.AppError
        if errors.As(err, &appErr) {
            w.WriteHeader(http.StatusBadRequest)
            json.NewEncoder(w).Encode(*appErr)
            return
        }
        
        var invalidErr *entity.ValidationError
        if errors.As(err, &invalidErr) {
            w.WriteHeader(http.StatusBadRequest)
            json.NewEncoder(w).Encode(*invalidErr)
            return
        } else {
            w.WriteHeader(http.StatusInternalServerError)
            json.NewEncoder(w).Encode(map[string]string{
                "message": "internal server error",
            })
            return
        }
    }
    
    security.GenerateRefreshCookie(w, token)
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(*token)
}