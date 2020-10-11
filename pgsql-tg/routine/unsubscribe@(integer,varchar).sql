CREATE FUNCTION unsubscribe(chat_id_in integer, url_in character varying) RETURNS void
	LANGUAGE plpgsql
AS
$$
BEGIN
	DELETE FROM good_user WHERE good_id = (select id from good WHERE url = url_in) AND user_id = (SELECT id FROM users WHERE chat_id = chat_id_in);
END
$$;

ALTER FUNCTION unsubscribe(INTEGER, VARCHAR) OWNER TO postgres;

