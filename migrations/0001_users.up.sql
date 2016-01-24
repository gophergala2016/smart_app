CREATE TABLE users (
	id         serial NOT NULL,
	email      character varying(254) NOT NULL,
	mobile     character varying(16) NOT NULL,
	password   text NOT NULL,
	created_at timestamp with time zone NOT NULL,
	updated_at timestamp with time zone NULL,
	deleted_at timestamp with time zone NULL,
	PRIMARY KEY (id),
	UNIQUE (email),
	UNIQUE (mobile)
);
