package db

import (
	"database/sql"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"groupie-tracker/internal/models"
	"log"
)

type Storage interface {
	CreateUser(*models.User) (int, error)
	UpdateUser(*models.User) (*models.User, error)
	DeleteUser(int) (*models.User, error)
	UserByID(int) (*models.User, error)
	UserByEmail(string, string) (*models.User, error)
	Users() ([]*models.User, error)
}

type PostgreSQL struct {
	DBSql *sql.DB
}

func (db *PostgreSQL) CreateUser(user *models.User) (int, error) {
	var ID int
	query := `INSERT INTO users (name, email, password, is_admin, created_at) VALUES ($1,$2,$3,$4,$5) RETURNING id`
	if err := db.DBSql.QueryRow(query, user.Name, user.Email, user.HashPassword, user.IsAdmin, user.CreatedAt).Scan(&ID); err != nil {
		return -1, err
	}

	return ID, nil
}
func (db *PostgreSQL) UserByEmail(email string, password string) (*models.User, error) {
	query := `SELECT id, name , email, password, is_admin, created_at FROM users WHERE email = $1`

	row := db.DBSql.QueryRow(query, email)
	var user models.User
	if err := row.Scan(&user.ID, &user.Name, &user.Email, &user.HashPassword, &user.IsAdmin, &user.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			log.Println("NOT FOUND IN DB: ", err)
			return nil, fmt.Errorf("NOT FOUND IN DB: ", err)
		}
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.HashPassword), []byte(password)); err != nil {
		return nil, err // UPDATE PASSWORD ???
	}
	return &user, nil
}

func (db *PostgreSQL) UpdateUser(update *models.User) (*models.User, error) {
	updateQuery := `UPDATE users SET name = $2,email=$3, password= $4, is_admin = $5 WHERE id = $1`
	if _, err := db.DBSql.Exec(updateQuery, update.ID, update.Name, update.Email, update.HashPassword, update.IsAdmin); err != nil {
		return nil, err
	}
	log.Println("User updated successfully")
	return update, nil
}
func (db *PostgreSQL) DeleteUser(id int) (*models.User, error) {
	deleteQuery := `WITH deleted AS (DELETE FROM users WHERE id = $1 RETURNING id , name, email, password, is_admin, created_at) SELECT * FROM deleted`
	var user models.User
	if err := db.DBSql.QueryRow(deleteQuery, id).Scan(&user.ID, &user.Name, &user.Email, &user.HashPassword, &user.IsAdmin, &user.CreatedAt); err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}
	log.Println("User DELETED SUCCESSFULLY")
	return &user, nil
}
func (db *PostgreSQL) UserByID(id int) (*models.User, error) {
	query := `SELECT id, name , email, password , is_admin, created_at FROM users WHERE id = $1 `
	row := db.DBSql.QueryRow(query, id)
	var user models.User
	if err := row.Scan(&user.ID, &user.Name, &user.Email, &user.HashPassword, &user.IsAdmin, &user.CreatedAt); err != nil {
		return nil, err
	}
	return &user, nil
}
func (db *PostgreSQL) Users() ([]*models.User, error) {
	var users []*models.User
	query := `SELECT id, name, email, password, is_admin, created_at FROM users`
	rows, err := db.DBSql.Query(query) // Используйте db.DB для выполнения запроса
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.HashPassword, &u.IsAdmin, &u.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, &u)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

func (db *PostgreSQL) CreatePost(post *models.Post) error {
	query := `INSERT INTO posts (user_id, title, description, created_at, updated_at)
    VALUES ($1, $2, $3, $4, $5)`
	if _, err := db.DBSql.Exec(query, post.UserID, post.Title, post.Description, post.CreatedAt, post.UpdatedAt); err != nil {
		return err
	}
	log.Println("Post created successfully")
	return nil
}

func (db *PostgreSQL) Post(id int) (*models.Post, error) {
	query := `SELECT id, user_id,  title, description, updated_at FROM posts WHERE id=$1`
	var post models.Post
	if err := db.DBSql.QueryRow(query, id).Scan(&post.ID, &post.UserID, &post.Title, &post.Description, &post.UpdatedAt); err != nil {
		return nil, err
	}
	return &post, nil
}

func (db *PostgreSQL) UpdatePost(post *models.Post) error {
	query := `UPDATE posts SET title=$2, description=$3, updated_at=$4 WHERE id = $1`
	if _, err := db.DBSql.Exec(query, post.ID, post.Title, post.Description, post.UpdatedAt); err != nil {
		return err
	}
	log.Println("Post updated successfully.")
	return nil
}

func (db *PostgreSQL) DeletePost(id int) (*models.Post, error) {
	query := `WITH deleted AS  (DELETE FROM posts WHERE id=$1 RETURNING id, user_id, title, description, created_at, updated_at) SELECT * FROM deleted`
	var post models.Post
	if err := db.DBSql.QueryRow(query, id).Scan(&post.ID, &post.UserID, &post.Title, &post.Description, &post.CreatedAt, &post.UpdatedAt); err != nil {
		return nil, err
	}
	return &post, nil

}

func (db *PostgreSQL) Posts(id int) ([]*models.Post, error) {
	var posts []*models.Post
	query := `SELECT id, user_id, title, description, created_at, updated_at FROM posts WHERE user_id=$1`
	rows, err := db.DBSql.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var p models.Post
		if err := rows.Scan(&p.ID, &p.UserID, &p.Title, &p.Description, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		posts = append(posts, &p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return posts, nil
}
