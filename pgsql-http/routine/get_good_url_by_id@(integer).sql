CREATE FUNCTION get_good_url_by_id(id_in integer) RETURNS character varying
	LANGUAGE plpgsql
AS
$$
BEGIN
	RETURN (SELECT good.url FROM good WHERE id = id_in);
END
$$;

ALTER FUNCTION get_good_url_by_id(INTEGER) OWNER TO postgres;

