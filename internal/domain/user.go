package domain

type User struct {
	ID        uint64 `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name,omitempty"`
	Age       uint8  `json:"age"`
}
