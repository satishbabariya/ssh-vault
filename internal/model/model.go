package model

import "time"

// type Credential struct {
// 	Host      string  `json:"host" validate:"required"`
// 	Port      int     `json:"port,omitempty"`
// 	User      string  `json:"user" validate:"required"`
// 	PrivatKey *string `json:"private_key,omitempty"`
// 	Password  *string `json:"password,omitempty"`
// }

// type Credentials []Credential

// type Remote struct {
// 	Host string `json:"host" validate:"required"`
// 	Port int    `json:"port" validate:"required"`
// }

// type Remotes []Remote

type Remote struct {
	ID        int64     `bun:"id,pk,autoincrement"`
	Host      string    `bun:"host" validate:"required"`
	Port      int       `bun:"port" validate:"required"`
	CreatedAt time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
}
