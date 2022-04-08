package model

type Credential struct {
	Host      string  `json:"host" validate:"required"`
	Port      int     `json:"port,omitempty"`
	User      string  `json:"user" validate:"required"`
	PrivatKey *string `json:"private_key,omitempty"`
	Password  *string `json:"password,omitempty"`
}
