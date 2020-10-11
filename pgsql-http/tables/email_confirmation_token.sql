-- auto-generated definition
CREATE TABLE email_confirmation_token
(
	id      SERIAL  NOT NULL
		CONSTRAINT email_confirmation_token_pk
			PRIMARY KEY,
	user_id INTEGER NOT NULL
		CONSTRAINT email_confirmation_token_user_id_fk
			REFERENCES "user"
			ON UPDATE RESTRICT ON DELETE RESTRICT,
	token   VARCHAR NOT NULL
);

ALTER TABLE email_confirmation_token
	OWNER TO postgres;

CREATE UNIQUE INDEX email_confirmation_token_user_id_uindex
	ON email_confirmation_token (user_id);

CREATE UNIQUE INDEX email_confirmation_token_user_id_uindex_2
	ON email_confirmation_token (user_id);

CREATE UNIQUE INDEX email_confirmation_token_token_uindex
	ON email_confirmation_token (token);

