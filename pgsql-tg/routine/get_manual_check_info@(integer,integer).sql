CREATE FUNCTION get_manual_check_info(chat_id_in integer, good_id_in integer)
	RETURNS TABLE(title_u text, url_u text, price_u integer)
	LANGUAGE plpgsql
AS
$$
BEGIN
	RETURN QUERY SELECT title                            AS title_u
	                  , url                              AS url_u
	                  , (SELECT price
		                     FROM good_user
		                     WHERE good_id = good_id_in
			                   AND user_id = chat_id_in) AS price_u
		             FROM good
		             WHERE id = good_id_in;
END
$$;

ALTER FUNCTION get_manual_check_info(INTEGER, INTEGER) OWNER TO postgres;

