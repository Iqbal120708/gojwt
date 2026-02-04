package usecase

import (
    "gojwt/internal/entity"
    "gojwt/internal/repository"
    "gojwt/internal/security"
    "gojwt/internal/validation"
)



type UserUseCase interface {
    Create(user *entity.User) (*entity.User, error)
    Login(login *entity.Login) (*entity.Token, error)
    User(email string) (*entity.User, error)
    RefreshToken(refreshToken string) (*entity.Token, error)
}

type userUseCase struct {
    userRepo repository.UserRepository
}

func NewUserUseCase(userRepo repository.UserRepository) *userUseCase {
    return &userUseCase{userRepo: userRepo}
}

func (uc *userUseCase) Create(user *entity.User) (*entity.User, error) {
    errInput := validation.ValidateEmail(user.Email)
    if errInput != nil {
        return nil, errInput
    }
    
    errInput = validation.ValidatePassword(user.Password)
    if errInput != nil {
        return nil, errInput
    }
    
    errInput = validation.ValidateRole(user.Role)
    if errInput != nil {
        return nil, errInput
    }
    
    hashPwd, err := security.HashPassword(user.Password) 
    if err != nil {
        return nil, err
    }
    user.Password = hashPwd

    if user.Role == nil {
        defaultRole := "regular"
        user.Role = &defaultRole
    }
    
    return uc.userRepo.Create(user)
}

func (uc *userUseCase) Login(login *entity.Login) (*entity.Token, error) {
    errInput := validation.ValidateEmail(login.Email)
    if errInput != nil {
        return nil, errInput
    }
    
    errInput = validation.ValidatePassword(login.Password)
    if errInput != nil {
        return nil, errInput
    }
    
    user, err := uc.userRepo.GetByEmail(login.Email)
    if err != nil {
		return nil, err
	}
	
	checkPwd := security.CheckPasswordHash(login.Password, user.Password)
	if !checkPwd {
	    return nil, &entity.ValidationError{
	        Field: "password",
	        Message: "password tidak valid",
	    }
	}
	
	token, err := security.GenerateTokens(*user.ID, user.Email, *user.Role)
	if err != nil {
	    return nil, err
	}
	
	return token, nil
}

func (uc *userUseCase) User(email string) (*entity.User, error) {
    errInput := validation.ValidateEmail(email)
    if errInput != nil {
        return nil, errInput
    }
    
    user, err := uc.userRepo.GetByEmail(email)
    if err != nil {
		return nil, err
	}
	
	return user, nil
}

func (uc *userUseCase) RefreshToken(refreshToken string) (*entity.Token, error) {
    isBlacklisted, err := uc.userRepo.IsRefreshTokenBlacklisted(refreshToken)
    if err != nil {
        return nil, err
    }
    
    if isBlacklisted {
        return nil, &entity.AppError{
            Code: "token_invalid",
            Message: "Token is blacklisted.",
        }
    }
    
    claims, err := security.ParseAndValidateRefreshToken(refreshToken)
    if err != nil {
        return nil, err
    }
    
    user, err := uc.userRepo.GetByID(claims.UserID)
    if err != nil {
        return nil, err
    }
    
    token, err := security.GenerateTokens(*user.ID, user.Email, *user.Role)
	if err != nil {
	    return nil, err
	}
	
	err = uc.userRepo.AddBlacklistToken(*user.ID, refreshToken)
	if err != nil {
	    return nil, err
	}
	
	return token, nil
}

func ToUserResponse(u *entity.User) entity.UserResponse {
	return entity.UserResponse{
		ID:    *u.ID,
		Name:  u.Name,
		Email: u.Email,
		Role:  *u.Role,
	}
}

