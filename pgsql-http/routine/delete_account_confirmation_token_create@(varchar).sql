CREATE FUNCTION delete_account_confirmation_token_create(user_id_in character varying) RETURNS character varying
	LANGUAGE plpgsql
AS
$$
DECLARE
	generated_uuid_vr UUID;
BEGIN
	generated_uuid_vr = uuid.uuid_generate_v4();

	INSERT
		INTO account_delete_token
			(user_id, token)
		VALUES
			(user_id_in, generated_uuid_vr);
	RETURN generated_uuid_vr::VARCHAR;
END
$$;

ALTER FUNCTION delete_account_confirmation_token_create(VARCHAR) OWNER TO postgres;

