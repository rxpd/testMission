CREATE FUNCTION get_subscribes(email_in character varying) RETURNS json
	LANGUAGE plpgsql
AS
$$
DECLARE
	user_id_vr INT;
BEGIN
	user_id_vr = (SELECT id FROM "user" WHERE email = email_in);
	RETURN (SELECT json_build_object(
			               'userID', user_id_vr,
			               'subscribes', json_agg(sub.good_id))
		        FROM (
			             (SELECT good_id AS good_id
				              FROM good_user
				              WHERE user_id = user_id_vr)
		             ) sub);

END
$$;

ALTER FUNCTION get_subscribes(VARCHAR) OWNER TO postgres;

