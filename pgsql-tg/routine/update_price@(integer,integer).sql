CREATE FUNCTION update_price(good_id_in integer, price_in integer) RETURNS void
	LANGUAGE plpgsql
AS
$$
BEGIN
	UPDATE good_user
	SET price=price_in
	  , last_check =now()
		WHERE good_id = good_id_in;

	UPDATE good
	SET price = price_in
		WHERE id = good_id_in;
END
$$;

ALTER FUNCTION update_price(INTEGER, INTEGER) OWNER TO postgres;

