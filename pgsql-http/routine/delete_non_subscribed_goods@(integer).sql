CREATE FUNCTION delete_non_subscribed_goods(good_id_in integer) RETURNS void
	LANGUAGE plpgsql
AS
$$
BEGIN
	DELETE FROM good WHERE id NOT IN (SELECT good_id FROM good_user WHERE good_id = good_id_in);
END
$$;

ALTER FUNCTION delete_non_subscribed_goods(INTEGER) OWNER TO postgres;

