// Code generated by chaingen. DO NOT EDIT.

package models

func NewUser() *User {
	return &User{}
}

func (u *User) NameSet(Name string) *User {
	u.Name = Name
	return u
}

func (u *User) PasswordSet(Password string) *User {
	u.Password = Password
	return u
}

func (u *User) EmailSet(Email string) *User {
	u.Email = Email
	return u
}

func (u *User) FullNameSet(FullName string) *User {
	u.FullName = FullName
	return u
}

func (u *User) CreatedAtSet(CreatedAt int64) *User {
	u.CreatedAt = CreatedAt
	return u
}
