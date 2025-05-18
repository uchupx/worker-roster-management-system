package model

import (
	"database/sql"
	"time"
)

type Shift struct {
	ID        int64        `json:"id" db:"id"`
	UserID    int64        `json:"user_id" db:"user_id"`
	StartTime string       `json:"start_time" db:"start_time"` // TEXT in SQL
	EndTime   string       `json:"end_time" db:"end_time"`     // TEXT in SQL
	ShiftDate time.Time    `json:"shift_date" db:"shift_date"` // DATE in SQL
	Status    int8         `json:"status" db:"status"`         // TEXT in SQL
	CreatedAt time.Time    `json:"created_at" db:"created_at"` // DATETIME in SQL
	UpdatedAt sql.NullTime `json:"updated_at" db:"updated_at"` // Nullable DATETIME
	DeletedAt sql.NullTime `json:"deleted_at" db:"deleted_at"` // Nullable DATETIME
}

type ShiftUser struct {
	ID        int64        `json:"id" db:"id"`
	UserID    int64        `json:"user_id" db:"user_id"`
	Name      string       `json:"name" db:"name"`
	StartTime string       `json:"start_time" db:"start_time"` // TEXT in SQL
	EndTime   string       `json:"end_time" db:"end_time"`     // TEXT in SQL
	ShiftDate time.Time    `json:"shift_date" db:"shift_date"` // DATE in SQL
	Status    int8         `json:"status" db:"status"`         // TEXT in SQL
	CreatedAt time.Time    `json:"created_at" db:"created_at"` // DATETIME in SQL
	UpdatedAt sql.NullTime `json:"updated_at" db:"updated_at"` // Nullable DATETIME
	DeletedAt sql.NullTime `json:"deleted_at" db:"deleted_at"` // Nullable DATETIME
}
