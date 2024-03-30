package user

import (
	"context"
	"errors"

	"github.com/gabriel-tama/banking-app/common/db"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Repository interface {
	Create(ctx context.Context, user *User) error
	GetSalt() int
	FindByCredential(ctx context.Context, user *User) error
	PingDB(ctx context.Context) error
}

type dbRepository struct {
	db          *db.DB
	BCRYPT_SALT int
}

func NewRepository(db *db.DB, BCRYPT_SALT int) Repository {
	return &dbRepository{db: db, BCRYPT_SALT: BCRYPT_SALT}
}

func (d *dbRepository) GetSalt() int {
	return d.BCRYPT_SALT
}

func (d *dbRepository) Create(ctx context.Context, user *User) error {
	stmt := `
        INSERT INTO users (
            email, name, password) 
    `

	stmt = stmt + `VALUES ($1, $2, $3) RETURNING id`
	row := d.db.Pool.QueryRow(ctx, stmt, user.Email, user.Name, user.Password)
	err := row.Scan(&user.ID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505":
				return ErrEmailAlreadyExists
			default:
				return err
			}
		}
		return err
	}
	return nil
}

func (d *dbRepository) FindByCredential(ctx context.Context, user *User) error {
	stmt := `SELECT id, email, name, password FROM users WHERE email=$1 `

	row := d.db.Pool.QueryRow(ctx, stmt, user.Email)
	err := row.Scan(&user.ID, &user.Email, &user.Name, &user.Password)
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrUserNotFound
	}
	if err != nil {
		return err
	}

	return nil

}

func (d *dbRepository) PingDB(ctx context.Context) error {
	return d.db.Pool.Ping(ctx)
}
