package model

import "github.com/francoispqt/gojay"

// User is the model used to communicate on the wire
type User struct {
	ID      int   `json:"user_id,omitempty"`
	Friends []int `json:"friends,omitempty"`
	Online  bool  `json:"online,omitempty"`
}

// UserStream implements gojay.UnmarshalerStream
type UserStream chan *User

// UnmarshalStream implements gojay.UnmarshalerStream
func (s UserStream) UnmarshalStream(dec *gojay.StreamDecoder) error {
	u := &User{}
	if err := dec.Object(u); err != nil {
		return err
	}
	s <- u
	return nil
}
