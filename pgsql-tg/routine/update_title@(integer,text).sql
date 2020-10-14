CREATE FUNCTION update_title(good_id_in integer, title_in text) RETURNS text
	LANGUAGE plpgsql
AS
$$
BEGIN
	IF (SELECT title FROM good WHERE id = good_id_in) != title_in THEN
		UPDATE good
		SET title=title_in
			WHERE id = good_id_in;

	END IF;
END
$$;

ALTER FUNCTION update_title(INTEGER, TEXT) OWNER TO postgres;

