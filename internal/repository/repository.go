package repository

import (
	"context"
	"errors"

	"github.com/P04KA/API/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo struct {
	conn *pgxpool.Pool
}

func New(conn *pgxpool.Pool) *UserRepo {
	return &UserRepo{conn: conn}
}

func (r *UserRepo) CreateUser(ctx context.Context, user models.User) (string, error) {
	user.ID = uuid.New().String()
	_, err := r.conn.Exec(ctx, "INSERT INTO users (id, name, age, country) VALUES ($1, $2, $3, $4)", user.ID, user.Name, user.Age, user.Country)
	if err != nil {
		return "", err
	}

	return user.ID, nil
}

func (r *UserRepo) GetUser(ctx context.Context, id string) (*models.User, error) {
	var user models.User
	err := r.conn.QueryRow(ctx, "SELECT id, name, age, country FROM users WHERE id = $1", id).Scan(&user.ID, &user.Name, &user.Age, &user.Country)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepo) UpdateUser(ctx context.Context, user models.User) error {
	rows, err := r.conn.Exec(ctx, "UPDATE users SET name = $2, age = $3, country = $4 WHERE id = $1", user.ID, user.Name, user.Age, user.Country)
	if err != nil {
		return err
	}
	if rows.RowsAffected() == 0 {
		return errors.New("user not found")
	}
	return nil
}
func (r *UserRepo) DeleteUser(ctx context.Context, id string) error {
	rows, err := r.conn.Exec(ctx, "DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return err
	}

	if rows.RowsAffected() == 0 {
		return errors.New("user not found")
	}
	return nil
}
