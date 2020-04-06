package main

type Server struct {
	Name        string
	ID          int
	Enabled     bool
	users       []string // not exported
}

