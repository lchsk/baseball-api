package main

type Error struct {
	Message string `json:"message"`
}

type ResponseErrors struct {
	Errors []Error `json:"errors"`
}

type Person struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	// URL?
}
