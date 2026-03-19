package services

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

// DatabaseService manages SQLite database operations.
type DatabaseService struct {
	db *sql.DB
}

// NewDatabaseService opens SQLite and creates tables.
func NewDatabaseService(path string) (*DatabaseService, error) {
	if path == "" {
		path = "rena.db"
	}
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}
	db.SetMaxOpenConns(1)
	if err := initSchema(db); err != nil {
		db.Close()
		return nil, err
	}
	return &DatabaseService{db: db}, nil
}

func initSchema(db *sql.DB) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY AUTOINCREMENT, email TEXT NOT NULL UNIQUE, password_hash TEXT NOT NULL, created_at DATETIME NOT NULL)`,
		`CREATE TABLE IF NOT EXISTS projects (id INTEGER PRIMARY KEY AUTOINCREMENT, user_id INTEGER NOT NULL, name TEXT NOT NULL, package_name TEXT NOT NULL, platform TEXT NOT NULL, color TEXT, icon_url TEXT, version_code INTEGER DEFAULT 1, version_name TEXT DEFAULT '1.0.0', workspace_xml TEXT, generated_code TEXT, status TEXT, created_at DATETIME NOT NULL, updated_at DATETIME NOT NULL, FOREIGN KEY(user_id) REFERENCES users(id))`,
		`CREATE TABLE IF NOT EXISTS user_keystores (id INTEGER PRIMARY KEY AUTOINCREMENT, user_id INTEGER NOT NULL, keystore_name TEXT NOT NULL, alias TEXT, created_at DATETIME NOT NULL, FOREIGN KEY(user_id) REFERENCES users(id))`,
	}
	for _, q := range queries {
		if _, err := db.Exec(q); err != nil {
			return err
		}
	}
	return nil
}

// User data
type User struct {
	ID           int64
	Email        string
	PasswordHash string
	CreatedAt    time.Time
}

// Project data
type Project struct {
	ID            int64
	UserID        int64
	Name          string
	PackageName   string
	Platform      string
	Color         string
	IconURL       string
	VersionCode   int
	VersionName   string
	WorkspaceXML  string
	GeneratedCode string
	Status        string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// CreateUser inserts a new user
func (s *DatabaseService) CreateUser(email, passwordHash string) (*User, error) {
	res, err := s.db.Exec("INSERT INTO users (email,password_hash,created_at) VALUES (?, ?, datetime('now'))", email, passwordHash)
	if err != nil {
		return nil, err
	}
	id, _ := res.LastInsertId()
	return &User{ID: id, Email: email, PasswordHash: passwordHash, CreatedAt: time.Now()}, nil
}

func (s *DatabaseService) GetUserByEmail(email string) (*User, error) {
	r := s.db.QueryRow("SELECT id,email,password_hash,created_at FROM users WHERE email = ?", email)
	u := &User{}
	var created string
	if err := r.Scan(&u.ID, &u.Email, &u.PasswordHash, &created); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	u.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", created)
	return u, nil
}

func (s *DatabaseService) GetUserByID(id int64) (*User, error) {
	r := s.db.QueryRow("SELECT id,email,password_hash,created_at FROM users WHERE id = ?", id)
	u := &User{}
	var created string
	if err := r.Scan(&u.ID, &u.Email, &u.PasswordHash, &created); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	u.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", created)
	return u, nil
}

func (s *DatabaseService) CreateProject(p *Project) (*Project, error) {
	res, err := s.db.Exec(`INSERT INTO projects (user_id,name,package_name,platform,color,icon_url,version_code,version_name,workspace_xml,generated_code,status,created_at,updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, datetime('now'), datetime('now'))`, p.UserID, p.Name, p.PackageName, p.Platform, p.Color, p.IconURL, p.VersionCode, p.VersionName, p.WorkspaceXML, p.GeneratedCode, p.Status)
	if err != nil {
		return nil, err
	}
	id, _ := res.LastInsertId()
	p.ID = id
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
	return p, nil
}

func (s *DatabaseService) GetProjectsByUserID(userID int64) ([]Project, error) {
	rs, err := s.db.Query("SELECT id,user_id,name,package_name,platform,color,icon_url,version_code,version_name,workspace_xml,generated_code,status,created_at,updated_at FROM projects WHERE user_id = ?", userID)
	if err != nil {
		return nil, err
	}
	defer rs.Close()
	projects := []Project{}
	for rs.Next() {
		p := Project{}
		var ca, ua string
		if err := rs.Scan(&p.ID, &p.UserID, &p.Name, &p.PackageName, &p.Platform, &p.Color, &p.IconURL, &p.VersionCode, &p.VersionName, &p.WorkspaceXML, &p.GeneratedCode, &p.Status, &ca, &ua); err != nil {
			return nil, err
		}
		p.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", ca)
		p.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", ua)
		projects = append(projects, p)
	}
	return projects, nil
}

func (s *DatabaseService) GetProjectByIDAndUserID(projectID, userID int64) (*Project, error) {
	p := &Project{}
	var ca, ua string
	err := s.db.QueryRow("SELECT id,user_id,name,package_name,platform,color,icon_url,version_code,version_name,workspace_xml,generated_code,status,created_at,updated_at FROM projects WHERE id=? AND user_id=?", projectID, userID).Scan(&p.ID, &p.UserID, &p.Name, &p.PackageName, &p.Platform, &p.Color, &p.IconURL, &p.VersionCode, &p.VersionName, &p.WorkspaceXML, &p.GeneratedCode, &p.Status, &ca, &ua)
	if err != nil {
		return nil, err
	}
	p.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", ca)
	p.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", ua)
	return p, nil
}

func (s *DatabaseService) UpdateProject(projectID, userID int64, updates map[string]interface{}) (*Project, error) {
	setParts := []string{}
	args := []interface{}{}
	for k, v := range updates {
		setParts = append(setParts, fmt.Sprintf("%s = ?", k))
		args = append(args, v)
	}
	if len(setParts) == 0 {
		return nil, fmt.Errorf("no updates")
	}
	args = append(args, projectID, userID)
	query := fmt.Sprintf("UPDATE projects SET %s, updated_at = datetime('now') WHERE id=? AND user_id=?", strings.Join(setParts, ", "))
	res, err := s.db.Exec(query, args...)
	if err != nil {
		return nil, err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return nil, fmt.Errorf("not found")
	}
	return s.GetProjectByIDAndUserID(projectID, userID)
}

func (s *DatabaseService) DeleteProject(projectID, userID int64) error {
	res, err := s.db.Exec("DELETE FROM projects WHERE id=? AND user_id=?", projectID, userID)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("not found")
	}
	return nil
}

func (s *DatabaseService) Close() error { return s.db.Close() }
