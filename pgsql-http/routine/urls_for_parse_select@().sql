CREATE FUNCTION urls_for_parse_select()
	RETURNS TABLE
	        (
		        url     CHARACTER VARYING,
		        price   INTEGER,
		        email   CHARACTER VARYING,
		        user_id INTEGER,
		        good_id INTEGER
	        )
	LANGUAGE plpgsql
AS
$$
BEGIN
	RETURN QUERY SELECT good.url   AS url
	                  , good.price AS price
	                  , u.email    AS email
	                  , u.id       AS user_id
	                  , gu.good_id AS good_id
		             FROM good
			                  INNER JOIN good_user gu ON good.id = gu.good_id
			                  INNER JOIN "user"    u ON u.id = gu.user_id;
END
$$;

ALTER FUNCTION urls_for_parse_select() OWNER TO postgres;

