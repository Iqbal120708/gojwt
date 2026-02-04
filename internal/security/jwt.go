package security

import (
    "github.com/golang-jwt/jwt/v5"
    "time"
    "gojwt/pkg/config"
    "gojwt/internal/entity"
    "errors"
    "net/http"
)

func GenerateTokens(userID int64, email, role string) (*entity.Token, error) {
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

func ParseAndValidateRefreshToken(refreshToken string) (*entity.RefreshClaims, error) {
    token, err := jwt.ParseWithClaims(
        refreshToken,
        &entity.RefreshClaims{},
        func(token *jwt.Token) (interface{}, error) {

            // Cek signing method
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, errors.New("unexpected signing method")
            }

            // Kembalikan secret key
            return []byte(config.Get().RefreshSecret), nil
        },
    )
    
    // Cek error parsing
    if err != nil {
        return nil, err
    }

    // Cek status valid token
    if !token.Valid {
        return nil, &entity.AppError{
            Code: "token_invalid",
            Message: "token is invalid or expired",
        }
    }

    // Ambil claims
    claims, ok := token.Claims.(*entity.RefreshClaims)
    if !ok {
        return nil, errors.New("invalid token claims")
    }
    return claims, err
}

func GenerateRefreshCookie(w http.ResponseWriter, token *entity.Token) {
    http.SetCookie(w, &http.Cookie{
        Name:     "refresh_token",
        Value:    token.Refresh,
        Path:     "/refresh",
        HttpOnly: true,
        Secure:   false,
        SameSite: http.SameSiteLaxMode,
    })
    
    token.Refresh = ""
}
