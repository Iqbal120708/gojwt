package usecase

import "gojwt/internal/entity"

type MockUserUseCase struct {
	CreateFn func(user *entity.User) (*entity.User, error)
	LoginFn func(login *entity.Login) (*entity.Token, error)
	UserFn func(email string) (*entity.User, error)
	RefreshTokenFn func(refreshToken string) (*entity.Token, error) 
}

func (m *MockUserUseCase) Create(user *entity.User) (*entity.User, error) {
	return m.CreateFn(user)
}

func (m *MockUserUseCase) Login(login *entity.Login) (*entity.Token, error) {
	return m.LoginFn(login)
}

func (m *MockUserUseCase) User(email string) (*entity.User, error) {
	return m.UserFn(email)
}

func (m *MockUserUseCase) RefreshToken(refreshToken string) (*entity.Token, error) {
	return m.RefreshTokenFn(refreshToken)
}
