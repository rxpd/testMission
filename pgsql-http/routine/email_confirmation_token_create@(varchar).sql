CREATE FUNCTION email_confirmation_token_create(email_in character varying) RETURNS character varying
	LANGUAGE plpgsql
AS
$$
DECLARE
	generated_uuid_vr UUID;
	user_id_vr        INT;
BEGIN
	generated_uuid_vr = uuid.uuid_generate_v4();
	user_id_vr = (SELECT id FROM "user" WHERE email = email_in);
	INSERT
		INTO email_confirmation_token
			(user_id, token)
		VALUES
			(user_id_vr, generated_uuid_vr);
	RETURN generated_uuid_vr::VARCHAR;
END
$$;

ALTER FUNCTION email_confirmation_token_create(VARCHAR) OWNER TO postgres;

