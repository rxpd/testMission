-- auto-generated definition
CREATE TABLE good_user
(
	good_id    INTEGER                                NOT NULL
		CONSTRAINT good_user_good_id_fk
			REFERENCES good
			ON UPDATE RESTRICT ON DELETE RESTRICT,
	user_id    INTEGER                                NOT NULL
		CONSTRAINT good_user_user_id_fk
			REFERENCES "user"
			ON UPDATE RESTRICT ON DELETE RESTRICT,
	last_check TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
	price      INTEGER
);

ALTER TABLE good_user
	OWNER TO postgres;

