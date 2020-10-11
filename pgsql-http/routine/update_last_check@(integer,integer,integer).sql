CREATE FUNCTION update_last_check(user_id_in integer, good_id_in integer, price_in integer) RETURNS void
	LANGUAGE plpgsql
AS
$$
BEGIN
	UPDATE good_user
	SET last_check= now()
	  , price=price_in
		WHERE user_id = user_id_in
		  AND good_id = good_id_in;

	UPDATE good
	SET price = price_in
		WHERE id = good_id_in;
END
$$;

ALTER FUNCTION update_last_check(INTEGER, INTEGER, INTEGER) OWNER TO postgres;

