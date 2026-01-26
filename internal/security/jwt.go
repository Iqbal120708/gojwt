package security

import (
    "github.com/golang-jwt/jwt/v5"
    "time"
    "gojwt/pkg/config"
    "gojwt/internal/entity"
)

func GenerateTokens(userID int, email, role string) (*entity.Token, error) {
    cfg := config.Get()
    
    accessExp := time.Now().Add(15 * time.Minute)
    accessClaims := &entity.AccessClaims{
        UserID: userID,
        Email:  email,
        Role:   role,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(accessExp),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            Issuer:    "myapp",
            Subject:   "access_token",
        },
    }
    
    accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).
        SignedString([]byte(cfg.AccessSecret))
    if err != nil {
        return nil, err
    }
    
    refreshExp := time.Now().Add(7 * 24 * time.Hour)
    refreshClaims := &entity.RefreshClaims{
        UserID: userID,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(refreshExp),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            Issuer:    "myapp",
            Subject:   "refresh_token",
        },
    }
    
    refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).
        SignedString([]byte(cfg.RefreshSecret))
    if err != nil {
        return nil, err
    }
    
    return &entity.Token{
        Access: accessToken,
        Refresh: refreshToken,
    }, nil
}

