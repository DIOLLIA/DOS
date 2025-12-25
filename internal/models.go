package internal

type User struct {
	Name string `json:"username"`
}

type Entry struct {
	Noun      string `json:"noun"`
	Adjective string `json:"adjective"`
}
