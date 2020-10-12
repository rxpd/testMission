CREATE FUNCTION delete_non_subscribed_goods() RETURNS void
	LANGUAGE plpgsql
AS
$$
BEGIN
	DELETE FROM good WHERE id NOT IN (SELECT good_id FROM good_user);
END
$$;

ALTER FUNCTION delete_non_subscribed_goods() OWNER TO postgres;

