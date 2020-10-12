CREATE FUNCTION update_last_manual_check(chat_id_in integer, good_id_in integer, price_in integer) RETURNS void
	LANGUAGE plpgsql
AS
$$
BEGIN
	UPDATE good_user
	SET last_check= now()
	  , price=price_in
		WHERE user_id = (SELECT id FROM users WHERE chat_id = chat_id_in)
		  AND good_id = good_id_in;

	-- 	UPDATE good
-- 	SET price = price_in
-- 		WHERE id = good_id_in;
END
$$;

ALTER FUNCTION update_last_manual_check(INTEGER, INTEGER, INTEGER) OWNER TO postgres;

