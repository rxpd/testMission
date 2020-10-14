-- auto-generated definition
CREATE TABLE good
(
	id    SERIAL  NOT NULL
		CONSTRAINT good_pk
			PRIMARY KEY,
	url   TEXT    NOT NULL,
	price INTEGER NOT NULL,
	title TEXT
);

ALTER TABLE good
	OWNER TO postgres;

CREATE UNIQUE INDEX good_url_uindex
	ON good (url);

