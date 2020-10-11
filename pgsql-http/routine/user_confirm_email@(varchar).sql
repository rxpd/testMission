CREATE FUNCTION user_confirm_email(token_in character varying) RETURNS character varying
	LANGUAGE plpgsql
AS
$$
DECLARE
	email_vr VARCHAR;
BEGIN
	email_vr = (SELECT email
		            FROM "user"
			                 INNER JOIN email_confirmation_token ect ON "user".id = ect.user_id
		            WHERE ect.token = token_in);
	IF email_vr ISNULL THEN
		RETURN 'token does not exists';
	END IF;
	UPDATE "user"
	SET verified_email = TRUE
		WHERE email = email_vr;
	DELETE FROM email_confirmation_token WHERE token = token_in;
	RETURN '';
END
$$;

ALTER FUNCTION user_confirm_email(VARCHAR) OWNER TO postgres;

