CREATE FUNCTION select_subscribes_list(chat_id_in INTEGER)
	RETURNS TABLE
	        (
		        title   TEXT,
		        good_id INTEGER
	        )
	LANGUAGE plpgsql
AS
$$
BEGIN
	RETURN QUERY SELECT good.title AS title
	                  , good.id    AS good_id
		             FROM good
			                  INNER JOIN good_user gu ON good.id = gu.good_id
		             WHERE gu.user_id = chat_id_in;
END
$$;

ALTER FUNCTION select_subscribes_list(INTEGER) OWNER TO postgres;

