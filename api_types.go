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

type Player struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	PositionNumber int    `json:"position_number"`
	PositionSymbol string `json:"position_symbol"`
	PositionName   string `json:"position_name"`
	// URL
}
