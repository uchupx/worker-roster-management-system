package model

import (
	"database/sql"
	"time"
)

type User struct {
	ID        int64        `json:"id" db:"id"`                 // Maps to 'id INTEGER PRIMARY KEY AUTOINCREMENT'
	Name      string       `json:"name" db:"name"`             // Maps to 'name TEXT NOT NULL'
	Email     string       `json:"email" db:"email"`           // Maps to 'email TEXT NOT NULL UNIQUE'
	Password  string       `json:"password" db:"password"`     // Maps to 'password TEXT NOT NULL'
	CreatedAt time.Time    `json:"created_at" db:"created_at"` // Maps to 'created_at DATETIME DEFAULT CURRENT_TIMESTAMP'
	UpdatedAt sql.NullTime `json:"updated_at" db:"updated_at"` // Maps to 'updated_at DATETIME', nullable
	DeletedAt sql.NullTime `json:"deleted_at" db:"deleted_at"` // Maps to 'deleted_at DATETIME', nullable
}

type UserRoles struct {
	ID        int64        `json:"id" db:"id"`                 // Maps to 'id INTEGER PRIMARY KEY AUTOINCREMENT'
	Name      string       `json:"name" db:"name"`             // Maps to 'name TEXT NOT NULL'
	Email     string       `json:"email" db:"email"`           // Maps to 'email TEXT NOT NULL UNIQUE'
	Password  string       `json:"password" db:"password"`     // Maps to 'password TEXT NOT NULL'
	Roles     []Role       `json:"roles" db:"roles"`           // Maps to 'roles TEXT NOT NULL'
	CreatedAt time.Time    `json:"created_at" db:"created_at"` // Maps to 'created_at DATETIME DEFAULT CURRENT_TIMESTAMP'
	UpdatedAt sql.NullTime `json:"updated_at" db:"updated_at"` // Maps to 'updated_at DATETIME', nullable
	DeletedAt sql.NullTime `json:"deleted_at" db:"deleted_at"` // Maps to 'deleted_at DATETIME', nullable
}
