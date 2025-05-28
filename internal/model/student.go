package model

type Student struct {
	Name    string   `json:"name"`
	ID      string   `json:"id"`
	Courses []Course `json:"courses"`
}

type Course struct {
	Name string `json:"name,omitempty"`
	ID   string `json:"id,omitempty"`
}
