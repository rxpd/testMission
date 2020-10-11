CREATE FUNCTION urls_for_parse_select()
	RETURNS TABLE(url character varying, price integer, chat_id integer, good_id integer)
	LANGUAGE plpgsql
AS
$$
BEGIN
	RETURN QUERY SELECT good.url   AS url
	                  , good.price AS price
	                  , u.chat_id     chat_id
	                  , gu.good_id AS good_id
		             FROM good
			                  INNER JOIN good_user gu ON good.id = gu.good_id
			                  INNER JOIN users     u ON u.id = gu.user_id;
END
$$;

ALTER FUNCTION urls_for_parse_select() OWNER TO postgres;

