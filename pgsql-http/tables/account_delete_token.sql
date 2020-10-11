-- auto-generated definition
CREATE TABLE account_delete_token
(
	id      SERIAL  NOT NULL
		CONSTRAINT account_delete_token_pk
			PRIMARY KEY,
	user_id INTEGER
		CONSTRAINT account_delete_token_user_id_fk
			REFERENCES "user"
			ON UPDATE RESTRICT ON DELETE RESTRICT,
	token   VARCHAR NOT NULL
);

ALTER TABLE account_delete_token
	OWNER TO postgres;

CREATE UNIQUE INDEX account_delete_token_token_uindex
	ON account_delete_token (token);

