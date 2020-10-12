CREATE FUNCTION get_subscribes_list(chat_id_in integer)
	RETURNS TABLE(title character varying, good_id integer)
	LANGUAGE plpgsql
AS
$$
BEGIN
	RETURN QUERY SELECT good.title AS title
	                  , good.id    AS good_id
		             FROM good
			                  INNER JOIN good_user gu ON good.id = gu.good_id
			                  INNER JOIN users     u ON u.id = gu.user_id
		             WHERE u.chat_id = chat_id_in;
END
$$;

ALTER FUNCTION get_subscribes_list(INTEGER) OWNER TO postgres;

