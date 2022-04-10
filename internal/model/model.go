package model

type Credential struct {
	Host      string  `json:"host" validate:"required"`
	Port      int     `json:"port,omitempty"`
	User      string  `json:"user" validate:"required"`
	PrivatKey *string `json:"private_key,omitempty"`
	Password  *string `json:"password,omitempty"`
}

type Credentials []Credential

type Remote struct {
	Host string `json:"host" validate:"required"`
	Port int    `json:"port" validate:"required"`
}

type Remotes []Remote
