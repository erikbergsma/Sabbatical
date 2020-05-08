package main

type Server struct {
	Name        string
	ID          int64
	Enabled     bool
	users       []string // not exported
}

type ServerString struct {
	Name        string
	ID          int64
	Enabled     string
	users       []string // not exported
}
