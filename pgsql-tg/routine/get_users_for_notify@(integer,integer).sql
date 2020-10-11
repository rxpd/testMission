CREATE FUNCTION get_users_for_notify(good_id_in integer, price_in integer)
	RETURNS TABLE(chat_id integer)
	LANGUAGE plpgsql
AS
$$
BEGIN
	RETURN QUERY SELECT u.chat_id
		             FROM good_user
			                  INNER JOIN users u ON u.id = good_user.user_id
		             WHERE good_user.good_id = good_id_in
			           AND good_user.price != price_in;
END
$$;

ALTER FUNCTION get_users_for_notify(INTEGER, INTEGER) OWNER TO postgres;

