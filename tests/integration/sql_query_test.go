//go:build integration
// +build integration

package integration

const (
	tableCreate = `CREATE TABLE IF NOT EXISTS public.users (
name          varchar(30) NOT NULL CONSTRAINT name_right CHECK ( name ~ '^[A-Za-z0-9_\.]+$' ) PRIMARY KEY,
password      varchar(30) NOT NULL,
email         varchar(50) NOT NULL UNIQUE CONSTRAINT email_right CHECK(email ~ '^.*@[A-Za-z0-9\-_\.]*$'),
full_name     varchar(255) NOT NULL,
created_at    integer
);`

	insertUsers = `INSERT INTO public.users (name, password, email, full_name, created_at)
VALUES ('Piter','123','piter@email.com','Piter Parker',1659447420),
('Sara','321','sara@email.com','Sara Conor',1659447430),
('Tony','456','tony@email.com','Tony Stark',1659447440),
('Berta','654','berta@email.com','Big Berta',1659447450);`

	deleteUsers = `DELETE FROM public.users;`

	selectUser = `SELECT name, password, email, full_name, created_at FROM users WHERE name=$1`
)
