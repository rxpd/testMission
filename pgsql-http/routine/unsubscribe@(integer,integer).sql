CREATE FUNCTION unsubscribe(user_id_in integer, good_id_in integer) RETURNS void
	LANGUAGE plpgsql
AS
$$
BEGIN
	DELETE FROM good_user WHERE good_id = good_id_in AND user_id = user_id_in;
	SELECT delete_non_subscribed_goods(good_id_in);
END
$$;

ALTER FUNCTION unsubscribe(INTEGER, INTEGER) OWNER TO postgres;

