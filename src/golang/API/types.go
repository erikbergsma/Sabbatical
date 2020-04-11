package main

type Server struct {
	Name        string
	ID          int64
	Enabled     bool
	users       []string // not exported
}

