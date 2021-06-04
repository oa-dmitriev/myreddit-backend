package user

import "encoding/json"

type User struct {
	Id       string `json:"id,omitempty" pg:"id"`
	Login    string `json:"login,omitempty" pg:"login"`
	Password string `json:"password,omitempty" pg:password,omitempty"`
}

func (u *User) Unpack(m interface{}) error {
	claims, err := json.Marshal(m)
	if err != nil {
		return err
	}
	err = json.Unmarshal(claims, u)
	if err != nil {
		return err
	}
	return nil
}
