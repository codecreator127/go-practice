package user

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(
	db *pgxpool.Pool,
) *PostgresRepository {
	return &PostgresRepository{
		db: db,
	}
}

func (r *PostgresRepository) Create(ctx context.Context, user *User) error {
	_, err := r.db.Exec(
		ctx,
		"INSERT INTO users (id, name, email) VALUES ($1, $2, $3)",
		user.ID,
		user.Name,
		user.Email,
	)

	return err
}

func (r *PostgresRepository) GetByID(ctx context.Context, id int64) (*User, error) {

	row := r.db.QueryRow(
		ctx,
		"SELECT id, name, email FROM users WHERE id = $1",
		id,
	)
	user := &User{}

	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
	)

	if err != nil {
		return nil, ErrUserNotFound
	}

	return user, nil
}

func (r *PostgresRepository) Update(ctx context.Context, user *User) error {

	_, err := r.db.Exec(
		ctx,
		"UPDATE users SET name = $1, email = $2 WHERE id = $3",
		user.Name,
		user.Email,
		user.ID,
	)

	return err
}
