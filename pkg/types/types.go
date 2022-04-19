package types

import "strings"

type Credential struct {
	Host      string   `json:"host" validate:"required"`
	Port      int      `json:"port,omitempty"`
	User      string   `json:"user" validate:"required"`
	PrivatKey []byte   `json:"private_key,omitempty"`
	Password  *string  `json:"password,omitempty"`
	Tags      []string `json:"tags,omitempty"`
}

func (c *Credential) TagsString() string {
	if len(c.Tags) == 0 {
		return ""
	}

	return strings.Join(c.Tags, ",")
}

func (c *Credential) HasTag(tag string) bool {
	for _, t := range c.Tags {
		if t == tag {
			return true
		}
	}

	return false
}
