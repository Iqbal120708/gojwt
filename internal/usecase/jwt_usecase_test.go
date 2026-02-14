package usecase

import (
    "gojwt/internal/entity"
    "gojwt/internal/security"
    "gojwt/internal/repository"
	"gojwt/internal/middleware"
	"database/sql"
	"gojwt/pkg/config"
    "testing"
    "errors"
    "github.com/DATA-DOG/go-sqlmock"
)

func TestCreate_Success(t *testing.T) {
    db, mock, _ := sqlmock.New()
    defer db.Close()

    repo := repository.NewUserRepo(db)
    uc := NewUserUseCase(repo)

    rows := sqlmock.NewRows(
        []string{"id", "name", "email", "role"},
    ).AddRow(
        1, "usertest", "test@example.com", "regular",
    )
    
    mock.ExpectExec(`insert into user`).
        WithArgs("usertest", "test@example.com", "regular", sqlmock.AnyArg()).
        WillReturnResult(sqlmock.NewResult(1, 1))
    
    mock.ExpectQuery(`select id, name, email, role from user where id = \?`).
        WithArgs(1).
        WillReturnRows(rows)
        
    user, err := uc.Create(&entity.User{
        Name:     "usertest",
        Email:    "test@example.com",
        Password: "secret",
    })
    
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    
    if *user.ID != 1 {
        t.Fatalf("unexpected id: %v", user.ID)
    }
    
    // test default role
    if *user.Role != "regular" {
        t.Fatalf("unexpected role, got regular")
    }
    
    isHash := security.CheckPasswordHash("secret", user.Password)
    if isHash {
        t.Fatal("Not hash user password")
    }

    if err := mock.ExpectationsWereMet(); err != nil {
        t.Fatal(err)
    }
}

func TestCreate_Error(t *testing.T) {
    db, mock, _ := sqlmock.New()
    defer db.Close()
    
    repo := repository.NewUserRepo(db)
    uc := NewUserUseCase(repo)

    _, err := uc.Create(&entity.User{
        Name:     "usertest",
        Email:    "test@",
        Password: "secret",
    })
    
    if err == nil {
        t.Fatal("expected error, got nil")
    }
    
    ve, ok := err.(*entity.ValidationError)
    if !ok {
        t.Fatalf("expected ValidationError, got %T", err)
    }
    
    if ve.Field != "email" {
        t.Fatalf("expected field email, got %s", ve.Field)
    }
    
    if err = mock.ExpectationsWereMet(); err != nil {
        t.Fatal(err)
    }
}

func TestLogin_Success(t *testing.T) {
    config.Load()
    
    db, mock, _ := sqlmock.New()
    defer db.Close()

    repo := repository.NewUserRepo(db)
    uc := NewUserUseCase(repo)

    pwd := "secret"
    pwdHash, _ := security.HashPassword(pwd)
    rows := sqlmock.NewRows(
        []string{"id", "name", "email", "role", "password"},
    ).AddRow(
        1, "usertest", "test@example.com", "regular", pwdHash,
    )
    
    mock.ExpectQuery(`select id, name, email, role, password from user where email = \?`).
        WithArgs("test@example.com").
        WillReturnRows(rows)
        
    token, err := uc.Login(&entity.Login{
        Email: "test@example.com",
        Password: pwd,
    })
    
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    
    // validate Access Token
    claims, err := middleware.ParseAndValidateToken(token.Access)
    if err != nil {
        t.Fatalf("unexpected error access token: %v", err)
    }
    if claims.Email != "test@example.com" {
        t.Fatalf("unexpected claims email: %v", claims.Email)
    }
    
    if err := mock.ExpectationsWereMet(); err != nil {
        t.Fatal(err)
    }
}

func TestLogin_ErrorPassword(t *testing.T) {
    config.Load()
    
    db, mock, _ := sqlmock.New()
    defer db.Close()

    repo := repository.NewUserRepo(db)
    uc := NewUserUseCase(repo)

    pwdReal := "secret"
    pwdHash, _ := security.HashPassword(pwdReal)
    rows := sqlmock.NewRows(
        []string{"id", "name", "email", "role", "password"},
    ).AddRow(
        1, "usertest", "test@example.com", "regular", pwdHash,
    )
    
    mock.ExpectQuery(`select id, name, email, role, password from user where email = \?`).
        WithArgs("test@example.com").
        WillReturnRows(rows)
        
    pwdFake := "fake secret"
    token, err := uc.Login(&entity.Login{
        Email: "test@example.com",
        Password: pwdFake,
    })
    
    if err == nil {
        t.Fatalf("unexpected error, got nil")
    }
    
    if token != nil {
        t.Fatal("unexpected token, got nil")
    }
    
    ve, ok := err.(*entity.ValidationError)
    if !ok {
        t.Fatalf("expected ValidationError, got %T", err)
    }
    
    if ve.Field != "password" {
        t.Fatalf("expected field password, got %s", ve.Field)
    } 
    
    if err := mock.ExpectationsWereMet(); err != nil {
        t.Fatal(err)
    }
}

func TestGetUser_Success(t *testing.T) {
    db, mock, _ := sqlmock.New()
    defer db.Close()

    repo := repository.NewUserRepo(db)
    uc := NewUserUseCase(repo)
    
    pwdReal := "secret"
    pwdHash, _ := security.HashPassword(pwdReal)
    rows := sqlmock.NewRows(
        []string{"id", "name", "email", "role", "password"},
    ).AddRow(
        1, "usertest", "test@example.com", "regular", pwdHash,
    )
    
    mock.ExpectQuery(`select id, name, email, role, password from user where email = \?`).
        WithArgs("test@example.com").
        WillReturnRows(rows)
        
    user, err := uc.User("test@example.com")
    
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    
    if user.Email != "test@example.com" {
        t.Fatalf("unexpected email: %v", user.Email)
    }
    
    if err := mock.ExpectationsWereMet(); err != nil {
        t.Fatal(err)
    }
}

func TestGetUser_Error(t *testing.T) {
    db, mock, _ := sqlmock.New()
    defer db.Close()

    repo := repository.NewUserRepo(db)
    uc := NewUserUseCase(repo)
        
    user, err := uc.User("test@.com")
    
    if err == nil {
        t.Fatalf("unexpected error, got nil")
    }
    
    if user != nil {
        t.Fatalf("unexpected user: %v", user)
    }
    
    if err := mock.ExpectationsWereMet(); err != nil {
        t.Fatal(err)
    }
}

func TestRefreshToken_Success(t *testing.T) {
    db, mock, _ := sqlmock.New()
    defer db.Close()

    refreshToken, _ := security.GenerateRefreshToken(1)
    
    repo := repository.NewUserRepo(db)
    uc := NewUserUseCase(repo)
    
    mock.ExpectQuery(`SELECT 1 FROM blacklist_token WHERE refresh_token = \? limit 1`).
        WithArgs(*refreshToken).
        WillReturnError(sql.ErrNoRows)
        
    // repo get by id
    rows := sqlmock.NewRows(
        []string{"id", "name", "email", "role"},
    ).AddRow(
        1, "usertest", "test@example.com", "regular",
    )

    mock.ExpectQuery(`select id, name, email, role from user where id = \?`).
        WithArgs(1).
        WillReturnRows(rows)
        
    mock.ExpectExec(`insert into blacklist_token`).
        WithArgs(1, *refreshToken).
        WillReturnResult(sqlmock.NewResult(1, 1))
        
    
    token, err := uc.RefreshToken(*refreshToken)
    
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    
    // validate Access Token
    claims, err := middleware.ParseAndValidateToken(token.Access)
    if err != nil {
        t.Fatalf("unexpected error access token: %v", err)
    }
    if claims.Email != "test@example.com" {
        t.Fatalf("unexpected claims email: %v", claims.Email)
    }
    
    // validate RefreshToken
    claimsRefresh, err := security.ParseAndValidateRefreshToken(token.Refresh)
    if err != nil {
        t.Fatalf("unexpected error access token: %v", err)
    }
    if claimsRefresh.UserID != 1 {
        t.Fatalf("unexpected claims ID: %v", claims.UserID)
    }
    
    if err := mock.ExpectationsWereMet(); err != nil {
        t.Fatal(err)
    }
}

func TestRefreshToken_Blacklist(t *testing.T) {
    db, mock, _ := sqlmock.New()
    defer db.Close()

    refreshToken, _ := security.GenerateRefreshToken(1)
    
    repo := repository.NewUserRepo(db)
    uc := NewUserUseCase(repo)
    
    rows := sqlmock.NewRows([]string{"1"}).AddRow(1)
    
    mock.ExpectQuery(`SELECT 1 FROM blacklist_token WHERE refresh_token = \? limit 1`).
        WithArgs(*refreshToken).
        WillReturnRows(rows)
        
    token, err := uc.RefreshToken(*refreshToken)
    
    if token != nil {
        t.Fatalf("unexpected token: %v", token)
    }
    
    if err == nil {
        t.Fatal("unexpected error, got nil")
    }
    
    appErr, ok := err.(*entity.AppError)
    if !ok {
        t.Fatalf("expected AppError, got %T", appErr)
    }
    
    if appErr.Code != "token_invalid" {
        t.Fatalf("expected AppError Code, got %T", appErr.Code)
    }
    
    if err := mock.ExpectationsWereMet(); err != nil {
        t.Fatal(err)
    }
}

func TestRefreshToken_Error(t *testing.T) {
    db, mock, _ := sqlmock.New()
    defer db.Close()

    refreshToken, _ := security.GenerateRefreshToken(1)
    
    repo := repository.NewUserRepo(db)
    uc := NewUserUseCase(repo)
    
    mock.ExpectQuery(`SELECT 1 FROM blacklist_token WHERE refresh_token = \? limit 1`).
        WithArgs(*refreshToken).
        WillReturnError(errors.New("database connection lost"))
        
    token, err := uc.RefreshToken(*refreshToken)
    
    if token != nil {
        t.Fatalf("unexpected token: %v", token)
    }
    
    if err == nil {
        t.Fatal("unexpected error, got nil")
    }
    
    // pastikan error bukan AppError
    err, ok := err.(*entity.AppError)
    if ok {
        t.Fatalf("expected err, got AppError")
    }
    
    if err := mock.ExpectationsWereMet(); err != nil {
        t.Fatal(err)
    }
}
