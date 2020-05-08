package main

type Server struct {
	Name        string
	ID          int64
	Enabled     bool    `json:"Enabled,omitempty"`
	users       []string // not exported
}
