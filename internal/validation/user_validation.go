package validation

import (
    "net/mail"
    "gojwt/internal/entity"
)

func ValidateEmail(email string) error {
    if email == "" {
        return &entity.ValidationError{
            Field:   "email",
            Message: "email is required",
        }
    }

    if _, err := mail.ParseAddress(email); err != nil {
        return &entity.ValidationError{
            Field:   "email",
            Message: "email is invalid",
        }
    }

    return nil
}

func ValidatePassword(password string) error {
    if password == "" {
        return &entity.ValidationError{
            Field:   "password",
            Message: "password is required",
        }
    }
    
    return nil
}

func ValidateRole(role *string) error {
    if role == nil {
        return nil
    }
    
    if *role != "regular" && *role != "superuser" {
        return &entity.ValidationError{
            Field:   "role",
            Message: "role is not regular or superuser",
        }
    }
    return nil
}
