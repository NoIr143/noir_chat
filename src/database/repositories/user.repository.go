package repositories

import (
	"database/sql"
	"errors"

	"github.com/noir143/noir_chat/src/database/entities"
)

type UserRepository struct {
	db *sql.DB
}

func UserRepositoryConstructor(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user entities.User) (entities.User, error) {
	query := `
		INSERT INTO users (first_name, last_name, email, hashed_password) 
		VALUES ($1, $2, $3, $4) 
		RETURNING id, first_name, last_name, email, hashed_password, created_at, updated_at`

	var createdUser entities.User
	err := r.db.QueryRow(query, user.FirstName, user.LastName, user.Email, user.HashedPassword).Scan(
		&createdUser.ID,
		&createdUser.FirstName,
		&createdUser.LastName,
		&createdUser.Email,
		&createdUser.HashedPassword,
		&createdUser.CreatedAt,
		&createdUser.UpdatedAt,
	)

	if err != nil {
		return entities.User{}, err
	}

	return createdUser, nil
}

func (r *UserRepository) GetByID(id int) (entities.User, error) {
	query := `SELECT id, first_name, last_name, email, hashed_password, created_at, updated_at FROM users WHERE id = $1`

	var user entities.User
	err := r.db.QueryRow(query, id).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.HashedPassword, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return entities.User{}, errors.New("user not found")
		}
		return entities.User{}, err
	}

	return user, nil
}

func (r *UserRepository) GetByEmail(email string) (entities.User, error) {
	query := `SELECT id, first_name, last_name, email, hashed_password, created_at, updated_at FROM users WHERE email = $1`

	var user entities.User
	err := r.db.QueryRow(query, email).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.HashedPassword, &user.CreatedAt, &user.UpdatedAt)

	return user, err
}

func (r *UserRepository) GetAll() ([]entities.User, error) {
	query := `SELECT id, first_name, last_name, email, hashed_password, created_at, updated_at FROM users ORDER BY id`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []entities.User
	for rows.Next() {
		var user entities.User
		err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.HashedPassword, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserRepository) Update(id int, user entities.User) (entities.User, error) {
	query := `
		UPDATE users 
		SET first_name = $1, last_name = $2, email = $3, hashed_password = $4
		WHERE id = $5
		RETURNING id, first_name, last_name, email, hashed_password, created_at, updated_at`

	var updatedUser entities.User
	err := r.db.QueryRow(query, user.FirstName, user.LastName, user.Email, user.HashedPassword, id).Scan(
		&updatedUser.ID,
		&updatedUser.FirstName,
		&updatedUser.LastName,
		&updatedUser.Email,
		&updatedUser.HashedPassword,
		&updatedUser.CreatedAt,
		&updatedUser.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return entities.User{}, errors.New("user not found")
		}
		return entities.User{}, err
	}

	return updatedUser, nil
}

func (r *UserRepository) Delete(id int) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}

// CountUsers returns the total number of users
func (r *UserRepository) CountUsers() (int, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// UserExists checks if a user exists by email
func (r *UserRepository) UserExists(email string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`
	err := r.db.QueryRow(query, email).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
