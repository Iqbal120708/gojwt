package repository

import (
    "gojwt/internal/entity"
    "testing"
    "database/sql"
    "github.com/go-sql-driver/mysql"
    "github.com/DATA-DOG/go-sqlmock"
)

func TestGetByEmail_Success(t *testing.T) {
    db, mock, _ := sqlmock.New()
    defer db.Close()

    repo := NewUserRepo(db)

    rows := sqlmock.NewRows(
        []string{"id", "name", "email", "role", "password"},
    ).AddRow(
        1, "usertest", "test@example.com", "regular", "secret",
    )

    mock.ExpectQuery(`select id, name, email, role, password from user where email = \?`).
        WithArgs("test@example.com").
        WillReturnRows(rows)

    user, err := repo.GetByEmail("test@example.com")

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

func TestGetByEmail_NotFound(t *testing.T) {
    db, mock, _ := sqlmock.New()
    defer db.Close()

    repo := NewUserRepo(db)

    mock.ExpectQuery(`select id, name, email, role, password from user where email = \?`).
        WithArgs("notfound@test.com").
        WillReturnError(sql.ErrNoRows)

    user, err := repo.GetByEmail("notfound@test.com")

    if err == nil {
        t.Fatal("expected error, got nil")
    }

    if user != nil {
        t.Fatal("expected user nil")
    }
    
    if err := mock.ExpectationsWereMet(); err != nil {
        t.Fatal(err)
    }
}

func TestGetByID_Success(t *testing.T) {
    db, mock, _ := sqlmock.New()
    defer db.Close()

    repo := NewUserRepo(db)

    rows := sqlmock.NewRows(
        []string{"id", "name", "email", "role"},
    ).AddRow(
        1, "usertest", "test@example.com", "regular",
    )

    mock.ExpectQuery(`select id, name, email, role from user where id = \?`).
        WithArgs(1).
        WillReturnRows(rows)
        
    user, err := repo.GetByID(1)
    
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }

    if *user.ID != 1 {
        t.Fatalf("unexpected id: %v", user.ID)
    }

    if err := mock.ExpectationsWereMet(); err != nil {
        t.Fatal(err)
    }
}

func TestGetByID_NotFound(t *testing.T) {
    db, mock, _ := sqlmock.New()
    defer db.Close()

    repo := NewUserRepo(db)

    mock.ExpectQuery(`select id, name, email, role from user where id = \?`).
        WithArgs(99).
        WillReturnError(sql.ErrNoRows)

    user, err := repo.GetByID(99)

    if err == nil {
        t.Fatal("expected error, got nil")
    }

    if user != nil {
        t.Fatal("expected user nil")
    }
    
    if err := mock.ExpectationsWereMet(); err != nil {
        t.Fatal(err)
    }
}

func TestCreateUser_Success(t *testing.T) {
    db, mock, _ := sqlmock.New()
    defer db.Close()

    repo := NewUserRepo(db)

    rows := sqlmock.NewRows(
        []string{"id", "name", "email", "role"},
    ).AddRow(
        1, "usertest", "test@example.com", "regular",
    )
    
    mock.ExpectExec(`insert into user`).
        WithArgs("usertest", "test@example.com", "regular", "secret").
        WillReturnResult(sqlmock.NewResult(1, 1))
    
    mock.ExpectQuery(`select id, name, email, role from user where id = \?`).
        WithArgs(1).
        WillReturnRows(rows)
        

    role := "regular"
    user, err := repo.Create(&entity.User{
        Name:     "usertest",
        Email:    "test@example.com",
        Role:     &role,
        Password: "secret",
    })

    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    
    if *user.ID != 1 {
        t.Fatalf("unexpected id: %v", user.Email)
    }

    if err := mock.ExpectationsWereMet(); err != nil {
        t.Fatal(err)
    }
}

func TestCreateUser_DuplicateKey(t *testing.T) {
    db, mock, _ := sqlmock.New()
    defer db.Close()

    repo := NewUserRepo(db)
    
    sqlmock.NewRows(
        []string{"id", "name", "email", "role", "password"},
    ).AddRow(
        1, "usertest", "test@example.com", "regular", "secret",
    )
    
    mock.ExpectExec(`insert into user`).
        WithArgs("usertest2", "test@example.com", "regular", "secret2").
        WillReturnError(&mysql.MySQLError{
            Number:  1062,
            Message: "Duplicate entry 'test@example.com' for key 'email'",
        })

    role := "regular"
    user, err := repo.Create(&entity.User{
        Name:     "usertest2",
        Email:    "test@example.com",
        Role:     &role,
        Password: "secret2",
    })

    if err == nil {
        t.Fatalf("unexpected error, got nil: %v", err)
    }

    if user != nil {
        t.Fatal("expected user nil")
    }
    
    ve, ok := err.(*entity.ValidationError)
    if !ok {
        t.Fatalf("expected ValidationError, got %T", err)
    }
    
    if ve.Field != "email" {
        t.Fatalf("expected field email, got %s", ve.Field)
    }
    
    if err := mock.ExpectationsWereMet(); err != nil {
        t.Fatal(err)
    }
}