package storage

var users = []struct {
	name     string
	password string
}{
	{
		name:     "Piter",
		password: "123",
	},
	{
		name:     "Maria",
		password: "321",
	},
	{
		name:     "Sunny",
		password: "12345678",
	},
}

func migrateUsers() error {
	for _, u := range users {
		user, err := NewUser(u.name, u.password)
		if err != nil {
			return err
		}
		if err = Add(user); err != nil {
			return err
		}
	}
	return nil
}
