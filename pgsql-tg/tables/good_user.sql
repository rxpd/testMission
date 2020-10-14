-- auto-generated definition
CREATE TABLE good_user
(
	id         SERIAL                                 NOT NULL
		CONSTRAINT good_user_pk
			PRIMARY KEY,
	user_id    INTEGER                                NOT NULL
		CONSTRAINT good_user_user_id_fkey
			REFERENCES tg_user
			ON UPDATE CASCADE ON DELETE CASCADE,
	good_id    INTEGER                                NOT NULL
		CONSTRAINT good_user_good_id_fkey
			REFERENCES good
			ON UPDATE CASCADE ON DELETE CASCADE,
	price      INTEGER                                NOT NULL,
	last_check TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL
);

ALTER TABLE good_user
	OWNER TO postgres;

