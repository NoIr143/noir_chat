package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

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
		INSERT INTO users (name, email) 
		VALUES ($1, $2) 
		RETURNING id, name, email, created_at, updated_at`

	var createdUser entities.User
	err := r.db.QueryRow(query, user.Name, user.Email).Scan(
		&createdUser.ID,
		&createdUser.Name,
		&createdUser.Email,
		&createdUser.CreatedAt,
		&createdUser.UpdatedAt,
	)

	if err != nil {
		return entities.User{}, err
	}

	return createdUser, nil
}

func (r *UserRepository) GetByID(id int) (entities.User, error) {
	query := `SELECT id, name, email, created_at, updated_at FROM users WHERE id = $1`

	var user entities.User
	err := r.db.QueryRow(query, id).Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return entities.User{}, errors.New("user not found")
		}
		return entities.User{}, err
	}

	return user, nil
}

func (r *UserRepository) GetByEmail(email string) (entities.User, error) {
	query := `SELECT id, name, email, created_at, updated_at FROM users WHERE email = $1`

	var user entities.User
	err := r.db.QueryRow(query, email).Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return entities.User{}, errors.New("user not found")
		}
		return entities.User{}, err
	}

	return user, nil
}

func (r *UserRepository) GetAll() ([]entities.User, error) {
	query := `SELECT id, name, email, created_at, updated_at FROM users ORDER BY id`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []entities.User
	for rows.Next() {
		var user entities.User
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)
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
		SET name = $1, email = $2 
		WHERE id = $3 
		RETURNING id, name, email, created_at, updated_at`

	var updatedUser entities.User
	err := r.db.QueryRow(query, user.Name, user.Email, id).Scan(
		&updatedUser.ID,
		&updatedUser.Name,
		&updatedUser.Email,
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

// GetPaginated returns paginated users
func (r *UserRepository) GetPaginated(page, pageSize int) ([]entities.User, int, error) {
	offset := (page - 1) * pageSize

	query := `SELECT id, name, email, created_at, updated_at FROM users ORDER BY id LIMIT $1 OFFSET $2`

	rows, err := r.db.Query(query, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []entities.User
	for rows.Next() {
		var user entities.User
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	// Get total count
	var total int
	err = r.db.QueryRow("SELECT COUNT(*) FROM users").Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
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

// FindByConditions finds users based on multiple conditions
func (r *UserRepository) FindByConditions(conditions map[string]interface{}, orderBy string, limit int) ([]entities.User, error) {
	query := `SELECT id, name, email, created_at, updated_at FROM users`
	var values []interface{}

	if len(conditions) > 0 {
		whereClauses := make([]string, 0, len(conditions))
		paramIndex := 1
		for column, value := range conditions {
			whereClauses = append(whereClauses, fmt.Sprintf("%s = $%d", column, paramIndex))
			values = append(values, value)
			paramIndex++
		}
		query += fmt.Sprintf(" WHERE %s", strings.Join(whereClauses, " AND "))
	}

	if orderBy != "" {
		query += fmt.Sprintf(" ORDER BY %s", orderBy)
	}

	if limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", limit)
	}

	rows, err := r.db.Query(query, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []entities.User
	for rows.Next() {
		var user entities.User
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)
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
