-- auto-generated definition
CREATE TABLE users
(
	id        SERIAL  NOT NULL
		CONSTRAINT users_pk
			PRIMARY KEY,
	chat_id   INTEGER NOT NULL,
	user_name VARCHAR NOT NULL
);

ALTER TABLE users
	OWNER TO postgres;

