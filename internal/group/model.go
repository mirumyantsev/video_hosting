package group

type Group struct {
	Id   int    `json:"id"   db:"id"`
	Name string `json:"name" db:"name"`
}
