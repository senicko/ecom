package users

import (
	"fmt"

	"github.com/jackc/pgx/v4"
)

type UserCreateParams struct {
	Email     string `json:"email"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Password  string `json:"password"`
}

type User struct {
	ID        int    `json:"id"`
	Email     string `json:"email"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Password  []byte `json:"-"`
}

// scan scans sql query result into the user.
func (u *User) scan(r pgx.Row) error {
	if err := r.Scan(&u.ID, &u.Email, &u.Firstname, &u.Lastname, &u.Password); err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
