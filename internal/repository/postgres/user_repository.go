package postgres

import (
	"github.com/breamon/sinav-bilgi-sistemi/internal/domain"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *domain.User) error {
	query := `
		INSERT INTO users (full_name, email, password_hash, role)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`

	return r.db.QueryRowx(
		query,
		user.FullName,
		user.Email,
		user.PasswordHash,
		user.Role,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}

func (r *UserRepository) GetByEmail(email string) (*domain.User, error) {
	var user domain.User

	query := `
		SELECT id, full_name, email, password_hash, role, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	if err := r.db.Get(&user, query, email); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetByID(id int64) (*domain.User, error) {
	var user domain.User

	query := `
		SELECT id, full_name, email, password_hash, role, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	if err := r.db.Get(&user, query, id); err != nil {
		return nil, err
	}

	return &user, nil
}
