package entity

import (
    "github.com/golang-jwt/jwt/v5"
)

type AccessClaims struct {
    UserID int64    `json:"user_id"`
    Email  string `json:"email"`
    Role   string `json:"role"`
    jwt.RegisteredClaims
}

type RefreshClaims struct {
    UserID int64 `json:"user_id"`
    jwt.RegisteredClaims
}

type User struct {
    ID       *int64  `json:"id"`
    Name     string  `json:"name"`
    Email    string  `json:"email"`
    Role     *string `json:"role"`
    Password string  `json:"password"`
}

type UserResponse struct {
    ID int64 `json:"id"`
    Name string `json:"name"`
    Email string `json:"email"`
    Role string `json:"role"`
}

type ValidationError struct {
    Field string `json:"field"`
    Message string `json:"message"`
}

func (e *ValidationError) Error() string {
    return e.Message
}

type AppError struct {
    Code    string `json:"code"`
    Message string `json:"message"`
}

func (e *AppError) Error() string {
    return e.Message
}

type Token struct {
    Access string `json:"access_token"`
    Refresh string `json:"refresh_token"`
}

type Login struct {
    Email string `json:"email"`
    Password string `json:"password"`
}

type BlacklistToken struct {
    ID int
    userID int64
    refreshToken string
}
