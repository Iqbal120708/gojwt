package validation

import (
    "gojwt/internal/entity"
    "testing"
)

func TestEmail_Success(t *testing.T) {
    err := ValidateEmail("test@example.com")
    
    if err != nil {
        t.Fatalf("unexpexted err: %v", err)
    }
}

func TestEmail_RequireError(t *testing.T) {
    err := ValidateEmail("")
    
    ve, ok := err.(*entity.ValidationError)
    if !ok {
        t.Fatalf("expected ValidationError, got %T", err)
    }
    
    if ve.Message != "email is required" {
        t.Fatalf("unexpexted err: %v", err)
    }
}

func TestEmail_Invalid(t *testing.T) {
    err := ValidateEmail("test")
    
    ve, ok := err.(*entity.ValidationError)
    if !ok {
        t.Fatalf("expected ValidationError, got %T", err)
    }
    
    if ve.Message != "email is invalid" {
        t.Fatalf("unexpexted err: %v", err)
    }
}

func TestPassword_Success(t *testing.T) {
    err := ValidatePassword("secret")
    
    if err != nil {
        t.Fatalf("unexpexted err: %v", err)
    }
}

func TestPassword_RequireError(t *testing.T) {
    err := ValidatePassword("")
    
    ve, ok := err.(*entity.ValidationError)
    if !ok {
        t.Fatalf("expected ValidationError, got %T", err)
    }
    
    if ve.Field != "password" {
        t.Fatalf("unexpexted err: %v", err)
    }
}

func TestRole_Success(t *testing.T) {
    regular := "regular"
    superuser := "superuser"
    err := ValidateRole(&regular)
    err2 := ValidateRole(&superuser)
    
    if err != nil && err2 != nil {
        t.Fatalf("unexpexted err: %v", err)
    }
}

func TestRole_Error(t *testing.T) {
    invalidRole := "invalid role"
    err := ValidateRole(&invalidRole)
    
    ve, ok := err.(*entity.ValidationError)
    if !ok {
        t.Fatalf("expected ValidationError, got %T", err)
    }
    
    if ve.Field != "role" {
        t.Fatalf("unexpexted err: %v", err)
    }
}