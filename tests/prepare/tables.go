package prepare

const TableUser = "user"

type User struct {
	Id   int    `json:"id,omitempty" orm:"id"`
	Name string `json:"name,omitempty" orm:"name"`
	Age  int    `json:"age,omitempty" orm:"age"`
}

func (u *User) TableName() string {
	return TableUser
}
