CREATE FUNCTION select_users_for_notify(good_id_in INTEGER, price_in INTEGER)
	RETURNS TABLE
	        (
		        chat_id INTEGER
	        )
	LANGUAGE plpgsql
AS
$$
BEGIN
	RETURN QUERY SELECT u.chat_id
		             FROM good_user
			                  INNER JOIN tg_user u ON u.chat_id = good_user.user_id
		             WHERE good_user.good_id = good_id_in
			           AND good_user.price != price_in;
END
$$;

ALTER FUNCTION select_users_for_notify(INTEGER, INTEGER) OWNER TO postgres;

