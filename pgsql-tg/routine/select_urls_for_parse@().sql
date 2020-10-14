CREATE FUNCTION select_urls_for_parse()
	RETURNS TABLE(url text, price integer, chat_id integer, good_id integer)
	LANGUAGE plpgsql
AS
$$
BEGIN
	RETURN QUERY SELECT good.url   AS url
	                  , good.price AS price
	                  , tu.chat_id    chat_id
	                  , gu.good_id AS good_id
		             FROM good
			                  INNER JOIN good_user gu ON good.id = gu.good_id
			                  INNER JOIN tg_user   tu ON gu.user_id = tu.chat_id;
END
$$;

ALTER FUNCTION select_urls_for_parse() OWNER TO postgres;

