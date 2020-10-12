CREATE FUNCTION get_manual_check_info(good_id_in integer, chat_id_in integer)
	RETURNS TABLE(title_u character varying, url_u character varying, price_u integer)
	LANGUAGE plpgsql
AS
$$
	-- DECLARE
-- 	title_out VARCHAR;
-- 	url_out   VARCHAR;
-- 	price_out INT;
BEGIN
	-- 	title_out = (SELECT title FROM good WHERE id = good_id_in);
-- 	url_out = (SELECT url FROM good WHERE id = good_id_in);
-- 	price_out = (SELECT price FROM good WHERE id = good_id_in);
	-- 	DELETE
-- 		FROM good_user
-- 		WHERE good_id = good_id_in
-- 		  AND user_id = (SELECT id FROM users WHERE chat_id = chat_id_in);
	RETURN QUERY SELECT title                                                                   AS title_u
	                  , url                                                                     AS url_u
	                  , (SELECT price
		                     FROM good_user
		                     WHERE good_id = good_id_in
			                   AND user_id = (SELECT id FROM users WHERE chat_id = chat_id_in)) AS price_u
		             FROM good
		             WHERE id = good_id_in;
END
$$;

ALTER FUNCTION get_manual_check_info(INTEGER, INTEGER) OWNER TO postgres;

