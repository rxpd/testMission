-- auto-generated definition
CREATE TABLE good_user
(
	id      SERIAL  NOT NULL
		CONSTRAINT good_user_pk
			PRIMARY KEY,
	good_id INTEGER NOT NULL
		CONSTRAINT good_user_good_id_fk
			REFERENCES good
			ON UPDATE RESTRICT ON DELETE RESTRICT,
	user_id INTEGER NOT NULL
		CONSTRAINT good_user_users_id_fk
			REFERENCES users
			ON UPDATE RESTRICT ON DELETE RESTRICT,
	price   INTEGER NOT NULL
);

ALTER TABLE good_user
	OWNER TO postgres;

