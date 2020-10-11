-- auto-generated definition
CREATE TABLE "user"
(
	id             SERIAL                NOT NULL
		CONSTRAINT user_pk
			PRIMARY KEY,
	email          VARCHAR(320)          NOT NULL,
	verified_email BOOLEAN DEFAULT FALSE NOT NULL
);

ALTER TABLE "user"
	OWNER TO postgres;

CREATE UNIQUE INDEX user_email_uindex
	ON "user" (email);

