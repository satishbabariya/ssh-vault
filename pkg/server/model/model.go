package model

import "time"

type Identity struct {
	ID        int64     `bun:"id,pk,autoincrement"`
	CreatedAt time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
	GithubID  string    `bun:"github_id,unique" validate:"required"`
	Name      *string   `bun:"name"`
}

type Remote struct {
	ID          int64        `bun:"id,pk,autoincrement"`
	Name        string       `bun:"name,unique" validate:"required"`
	Host        string       `bun:"host,unique" validate:"required"`
	Port        int          `bun:"port" validate:"required"`
	CreatedAt   time.Time    `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt   time.Time    `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
	Permissions []Permission `bun:"m2m:permissions,join:Remote=Identity"`
}

type Credential struct {
	ID        int64     `bun:"id,pk,autoincrement"`
	CreatedAt time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
	RemoteID  int64     `bun:"remote_id,notnull"`
}

type Permission struct {
	ID          int64       `bun:"id,pk,autoincrement"`
	Remote      *Remote     `bun:"rel:belongs-to,join:remote_id=id"`
	RemoteID    int64       `bun:"remote_id"`
	Identity    *Identity   `bun:"rel:belongs-to,join:identity_id=id"`
	IdentityID  int64       `bun:"identity_id"`
	Permissions Permissions `bun:"embed:roles_"`
}

type Permissions struct {
	Read  bool
	Write bool
}
