package postgres

import (
	"context"
	"errors"
	"github.com/evanebb/gobble/auth"
	"github.com/evanebb/gobble/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ApiUserRepository struct {
	db *pgxpool.Pool
}

func NewApiUserRepository(db *pgxpool.Pool) (ApiUserRepository, error) {
	return ApiUserRepository{db: db}, nil
}

type postgresApiUser struct {
	Id       uint
	UUID     uuid.UUID
	Name     string
	Password []byte
}

func (r ApiUserRepository) GetApiUsers() ([]auth.ApiUser, error) {
	var users []auth.ApiUser

	stmt := "SELECT id, uuid, name, password FROM api_user"
	rows, err := r.db.Query(context.Background(), stmt)
	if err != nil {
		return users, err
	}

	for rows.Next() {
		var u auth.ApiUser
		var pu postgresApiUser

		err = rows.Scan(&pu.Id, &pu.UUID, &pu.Name, &pu.Password)
		if err != nil {
			return users, err
		}

		u = auth.NewApiUser(pu.UUID, pu.Name, pu.Password)
		users = append(users, u)
	}

	return users, nil
}

func (r ApiUserRepository) GetApiUserById(id uuid.UUID) (auth.ApiUser, error) {
	var a auth.ApiUser
	var pa postgresApiUser

	stmt := "SELECT id, uuid, name, password FROM api_user WHERE uuid = $1"
	err := r.db.QueryRow(context.Background(), stmt, id).Scan(&pa.Id, &pa.UUID, &pa.Name, &pa.Password)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return a, repository.ErrNotFound
		}
		return a, err
	}

	return auth.NewApiUser(pa.UUID, pa.Name, pa.Password), nil
}

func (r ApiUserRepository) GetApiUserByName(name string) (auth.ApiUser, error) {
	var a auth.ApiUser
	var pa postgresApiUser

	stmt := "SELECT id, uuid, name, password FROM api_user WHERE name = $1"
	err := r.db.QueryRow(context.Background(), stmt, name).Scan(&pa.Id, &pa.UUID, &pa.Name, &pa.Password)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return a, repository.ErrNotFound
		}
		return a, err
	}

	return auth.NewApiUser(pa.UUID, pa.Name, pa.Password), nil
}

func (r ApiUserRepository) SetApiUser(a auth.ApiUser) error {
	stmt := "INSERT INTO api_user (uuid, name, password) VALUES ($1, $2, $3) ON CONFLICT (uuid) DO UPDATE SET name = $2, password = $3"
	_, err := r.db.Exec(context.Background(), stmt, a.Id, a.Name, a.Password)
	return err
}

func (r ApiUserRepository) DeleteApiUserById(id uuid.UUID) error {
	stmt := "DELETE FROM api_user WHERE uuid = $1"
	_, err := r.db.Exec(context.Background(), stmt, id)
	return err
}
