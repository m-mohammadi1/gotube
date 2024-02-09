package userrepo

import (
	"context"
	"database/sql"
	"fmt"
	"gotube/internal/utils/passwordutil"
	"gotube/pkg/model"
)

type PostgresRepository struct {
	db *sql.DB
}

func New(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (repo *PostgresRepository) All(ctx context.Context) ([]*model.User, error) {

	return nil, nil
}

func (repo *PostgresRepository) Find(ctx context.Context, id int) (*model.User, error) {
	query := "SELECT * from users WHERE id = $1"

	var user model.User
	err := repo.db.QueryRow(query, id).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.CreatedAt)

	if err != nil {
		return nil, err
	}

	return &user, err
}

func (repo *PostgresRepository) Delete(ctx context.Context, id int) error {

	return nil
}

func (repo *PostgresRepository) Create(ctx context.Context, user model.User) error {
	hashed, err := passwordutil.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashed

	query := "INSERT INTO users (name, email, password) VALUES ($1, $2, $3)"

	stmt, err := repo.db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.Name, user.Email, user.Password)
	if err != nil {
		return err
	}

	return nil
}

func (repo *PostgresRepository) EmailExists(email string) (bool, error) {
	var exists bool
	err := repo.db.QueryRow("SELECT EXISTS (SELECT 1 FROM users WHERE email = $1)", email).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (repo *PostgresRepository) FindByField(ctx context.Context, field string, value string) (*model.User, error) {
	var user model.User
	err := repo.db.QueryRow(fmt.Sprintf("SELECT * FROM users WHERE %s = $1 Limit 1", field), value).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.CreatedAt)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
