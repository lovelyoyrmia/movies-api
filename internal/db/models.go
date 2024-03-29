// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package db

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Movie struct {
	ID          int32            `json:"id"`
	Title       string           `json:"title"`
	Description string           `json:"description"`
	Image       pgtype.Text      `json:"image"`
	Rating      pgtype.Float8    `json:"rating"`
	CreatedAt   pgtype.Timestamp `json:"created_at"`
	UpdatedAt   pgtype.Timestamp `json:"updated_at"`
}
