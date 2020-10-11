CREATE FUNCTION manual_checks_select()
	RETURNS TABLE(price integer, good_id character varying, user_id character varying)
	LANGUAGE plpgsql
AS
$$
BEGIN
	RETURN QUERY SELECT price     AS price
	                  , good_id   AS good_id
	                  , user_id AS user_id
		             FROM good_user;

END
$$;

ALTER FUNCTION manual_checks_select() OWNER TO postgres;

