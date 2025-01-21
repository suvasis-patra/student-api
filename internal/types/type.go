package types

type Student struct {
	Id    int64 `json:"id"`
	Name  string `validate:"required" json:"name"`
	Email string `validate:"required" json:"email"`
	Age   int    `validate:"required" json:"age"`
}
