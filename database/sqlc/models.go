// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package sqlc

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Users struct {
	ID        uuid.UUID    `db:"id" json:"id"`
	Name      string       `db:"name" json:"name"`
	Email     string       `db:"email" json:"email"`
	CreatedAt time.Time    `db:"created_at" json:"created_at"`
	UpdatedAt time.Time    `db:"updated_at" json:"updated_at"`
	DeletedAt sql.NullTime `db:"deleted_at" json:"deleted_at"`
}
