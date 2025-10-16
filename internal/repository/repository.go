package repository

import (
	"context"
	"time"

	"github.com/P04KA/API/internal/apperr"
	"github.com/P04KA/API/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type UserRepo struct {
	conn *pgxpool.Pool
}

func New(conn *pgxpool.Pool) *UserRepo {
	return &UserRepo{conn: conn}
}

func (r *UserRepo) CreateUser(ctx context.Context, user models.User) (*models.User, error) {
	user.ID = uuid.New().String()
	_, err := r.conn.Exec(ctx, "INSERT INTO users (id, name, age, country) VALUES ($1, $2, $3, $4)", user.ID, user.Name, user.Age, user.Country)
	if err != nil {
		return nil, errors.Wrap(err, "insert user")
	}

	return &user, nil
}

func (r *UserRepo) GetUser(ctx context.Context, id string) (*models.User, error) {
	var user models.User
	err := r.conn.QueryRow(ctx, "SELECT id, name, age, country FROM users WHERE id = $1 AND deleted_at IS NULL", id).Scan(&user.ID, &user.Name, &user.Age, &user.Country)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, apperr.ErrNotFound
	}
	if err != nil {
		return nil, errors.Wrap(err, "get user")
	}

	return &user, nil
}

func (r *UserRepo) UpdateUser(ctx context.Context, user models.User) error {
	rows, err := r.conn.Exec(ctx, "UPDATE users SET name = $2, age = $3, country = $4, updated_at = NOW() WHERE id = $1", user.ID, user.Name, user.Age, user.Country)
	if err != nil {
		return errors.Wrap(err, "update user")
	}
	if rows.RowsAffected() == 0 {
		return apperr.ErrNotFound
	}
	return nil
}
func (r *UserRepo) DeleteUser(ctx context.Context, id string) error {
	rows, err := r.conn.Exec(ctx, "UPDATE users SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL", id)
	if err != nil {
		return errors.Wrap(err, "delete user")
	}

	if rows.RowsAffected() == 0 {
		return apperr.ErrNotFound
	}
	return nil
}
func (r *UserRepo) CountUserCreated(ctx context.Context, time time.Time) (int64, error) {
	var count int64
	err := r.conn.QueryRow(ctx, `SELECT COUNT(*) FROM users WHERE created_at > $1 AND deleted_at IS NULL`, time).Scan(&count)
	return count, err
}

func (r *UserRepo) CountUserUpdated(ctx context.Context, time time.Time) (int64, error) {
	var count int64
	err := r.conn.QueryRow(ctx, `SELECT COUNT(*) FROM users WHERE updated_at > $1 AND created_at != updated_at AND deleted_at IS NULL`, time).Scan(&count)
	return count, err
}

func (r *UserRepo) CountUserDeleted(ctx context.Context, time time.Time) (int64, error) {
	var count int64
	err := r.conn.QueryRow(ctx, `SELECT COUNT(*) FROM users WHERE deleted_at > $1`, time).Scan(&count)
	return count, err
	//что будет возвращать queryRow, если нет записей. мб pgx вернется
}

//запилить новую репу statistics, там запускаю сервер, с соблюдение архитектуры user stats server - usecase, его надо запихнуть в internal/usecase. - Это сделали
// в api gateway, клиент(убрать его в Handle и добавить зависимость по типу internal/usecase/stat.go)
// stats.prorto запихнуть в каждый репозиторий
//запилить новый репозиторий all-in-one, там пропаисать docker compose, который будет поднимать весь сервис/проект(там же конфиги для графаны)
