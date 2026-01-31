package middleware

import (
    "net/http"
    "gojwt/internal/entity"
    "github.com/golang-jwt/jwt/v5"
    "strings"
    "errors"
    "encoding/json"
    "gojwt/pkg/config"
    "context"
)

func JSONError(w http.ResponseWriter, status int, message string) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)

    json.NewEncoder(w).Encode(map[string]string{
        "message": message,
    })
}

func ParseAndValidateToken(tokenString string) (*entity.AccessClaims, error) {
    token, err := jwt.ParseWithClaims(
        tokenString,
        &entity.AccessClaims{},
        func(token *jwt.Token) (interface{}, error) {

            // Cek signing method
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, errors.New("unexpected signing method")
            }

            // Kembalikan secret key
            return []byte(config.Get().AccessSecret), nil
        },
    )
    
    // Cek error parsing
    if err != nil {
        return nil, err
    }

    // Cek status valid token
    if !token.Valid {
        return nil, errors.New("token is invalid or expired")
    }

    // Ambil claims
    claims, ok := token.Claims.(*entity.AccessClaims)
    if !ok {
        return nil, errors.New("invalid token claims")
    }
    return claims, err
}

func ValidateBearerToken(r *http.Request) (*entity.AccessClaims, error) {
    // 1. Ambil header Authorization
    authHeader := r.Header.Get("Authorization")
    if authHeader == "" {
        return nil, errors.New("authorization header missing")
    }

    // 2. Pisahkan "Bearer" dan token
    parts := strings.Split(authHeader, " ")
    if len(parts) != 2 || parts[0] != "Bearer" {
        return nil, errors.New("invalid authorization format")
    }


    // 3. Parse & verifikasi JWT
    tokenString := parts[1]
    claims, err := ParseAndValidateToken(tokenString)
    
    if err != nil {
        return nil, err
    }

    return claims, nil
}

func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

        claims, err := ValidateBearerToken(r)
        if err != nil {
            JSONError(w, http.StatusUnauthorized, err.Error())
            return
        }
        
        ctx := context.WithValue(r.Context(), "user_claims", claims)

        next.ServeHTTP(w, r.WithContext(ctx))
    })
}