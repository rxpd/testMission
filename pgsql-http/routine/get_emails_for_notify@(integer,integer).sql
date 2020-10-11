CREATE FUNCTION get_emails_for_notify(good_id_in integer, price_in integer)
	RETURNS TABLE(email character varying)
	LANGUAGE plpgsql
AS
$$
BEGIN
	RETURN QUERY SELECT u.email
		             FROM good_user
			                  INNER JOIN "user" u ON u.id = good_user.user_id
		             WHERE good_user.good_id = good_id_in
			           AND good_user.price != price_in;
END
$$;

ALTER FUNCTION get_emails_for_notify(INTEGER, INTEGER) OWNER TO postgres;

