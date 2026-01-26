package repository

import (
    "gojwt/internal/entity"
    "database/sql"
    "github.com/go-sql-driver/mysql"
    "errors"
)

type userRepo struct {
    db *sql.DB
}

type UserRepository interface {
	GetByEmail(email string) (*entity.User, error)
	GetByID(id int64) (*entity.User, error)
	Create(data *entity.User) (*entity.User, error)
}

func NewUserRepo(db *sql.DB) *userRepo {
    return &userRepo{db:db}
}

func (u *userRepo) GetByEmail(email string) (*entity.User,  error) {
    query := "select id, name, email, role, password from user where email = ?"
    
    var user entity.User
    err := u.db.QueryRow(query, email).Scan(&user.ID, &user.Name, &user.Email, &user.Role, &user.Password)
	if err != nil {
		return nil, err
	}

    return &user, nil
}

func (u *userRepo) GetByID(id int64) (*entity.User, error) {
    var user entity.User
    err := u.db.QueryRow("select id, name, email, role from user where id = ?", id).Scan(&user.ID, &user.Name, &user.Email, &user.Role)
    if (err != nil) {
        return nil, err
    };
    return &user, nil
}

func (u *userRepo) Create(data *entity.User) (*entity.User, error) {
    query := "insert into user (name, email, role, password) values (?, ?, ?, ?)"
    
    user, err := u.db.Exec(query, data.Name, data.Email, data.Role, data.Password)
    
    if (err != nil) {
        if isDuplicateKey(err) {
            return nil, &entity.ValidationError{
                Field: "email",
                Message: "email already exists",
            }
        }
        return nil, err
    }
    
    id, err2 := user.LastInsertId()
    if err2 != nil {
    	return nil, err
    }
    
    return u.GetByID(id)
}

func isDuplicateKey(err error) bool {
    var mysqlErr *mysql.MySQLError
    if errors.As(err, &mysqlErr) {
        return mysqlErr.Number == 1062
    }
    return false
}