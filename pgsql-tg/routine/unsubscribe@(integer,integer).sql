CREATE FUNCTION unsubscribe(chat_id_in integer, good_id_in integer)
	RETURNS TABLE(title_u character varying, url_u character varying, price_u integer)
	LANGUAGE plpgsql
AS
$$
DECLARE
	title_out VARCHAR;
	url_out   VARCHAR;
	price_out INT;
BEGIN
	title_out = (SELECT title FROM good WHERE id = good_id_in);
	url_out = (SELECT url FROM good WHERE id = good_id_in);
	price_out = (SELECT price FROM good WHERE id = good_id_in);
	DELETE
		FROM good_user
		WHERE good_id = good_id_in
		  AND user_id = chat_id_in;
	RETURN QUERY (SELECT title_out, url_out, price_out);
END
$$;

ALTER FUNCTION unsubscribe(INTEGER, INTEGER) OWNER TO postgres;

