-- auto-generated definition
CREATE TABLE tg_user
(
	chat_id  INTEGER NOT NULL
		CONSTRAINT tg_user_pk
			PRIMARY KEY,
	username TEXT    NOT NULL
);

ALTER TABLE tg_user
	OWNER TO postgres;

