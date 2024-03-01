package storage

import (
	"database/sql"
	"errors"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"log"
)

type User struct {
	Login string
	Tasks []Task
}

type Storage struct {
	db *sql.DB
}

type Task struct {
	ID          int
	Title       string
	Description string
}

func New(storagePath string) *Storage {
	db, err := sql.Open("pgx", storagePath)
	if err != nil {
		log.Fatalln("Невозможно открыть базу данных:", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalln("Невозможно подключиться к базе данных:", err)
	}

	if err = goose.SetDialect("postgres"); err != nil {
		log.Fatalln("Невозможно задать драйвер базы данных:", err)
	}

	if err = goose.Up(db, "migrations"); err != nil {
		log.Fatalln("Невозможно выполнить миграции:", err)
	}

	return &Storage{db: db}
}

func (s *Storage) IsUserInDB(login string) (bool, error) {
	var count int
	query := "SELECT COUNT(*) FROM users WHERE login = $1"
	row := s.db.QueryRow(query, login)
	if err := row.Scan(&count); err != nil {
		return false, err
	}
	return count > 0, nil
}

func (s *Storage) AddUser(login, password string) error {

	isExists, err := s.IsUserInDB(login)
	if err != nil {
		return err
	}

	if isExists {
		return errors.New("user already exists")
	}

	query := "INSERT INTO users (login, password) VALUES ($1, $2)"

	_, err = s.db.Exec(query, login, password)

	return err
}

func (s *Storage) GetPassword(login string) (string, error) {
	var password string

	query := "SELECT password FROM users WHERE login = $1"

	row := s.db.QueryRow(query, login)

	if err := row.Scan(&password); err != nil {
		return "", err
	}

	return password, nil
}

func (s *Storage) GetUserID(login string) (int, error) {
	var userID int

	query := "SELECT id FROM users WHERE login = $1"

	row := s.db.QueryRow(query, login)

	if err := row.Scan(&userID); err != nil {
		return 0, err
	}

	return userID, nil
}

func (s *Storage) GetAllTasks(login string) (*User, error) {
	userID, err := s.GetUserID(login)
	if err != nil {
		return nil, err
	}
	tasks := make([]Task, 0)
	query := "SELECT id, title, description FROM tasks WHERE user_id = $1"
	rows, err := s.db.Query(query, userID)
	if err != nil {
		return nil, err
	}

	task := Task{}

	for rows.Next() {
		if err := rows.Scan(&task.ID, &task.Title, &task.Description); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	user := User{
		Login: login,
		Tasks: tasks,
	}

	return &user, nil
}

func (s *Storage) AddTask(login string, title, description string) error {
	userID, err := s.GetUserID(login)
	if err != nil {
		return err
	}
	query := "INSERT INTO tasks (user_id, title, description) VALUES ($1, $2, $3)"
	if _, err := s.db.Exec(query, userID, title, description); err != nil {
		return err
	}

	return nil
}

func (s *Storage) DeleteTask(taskID int) error {
	query := "DELETE FROM tasks WHERE id = $1"
	if _, err := s.db.Exec(query, taskID); err != nil {
		return err
	}
	return nil
}
