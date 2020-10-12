CREATE FUNCTION manual_check_validation(user_id_in integer, good_id_in integer, cooldown_in_minutes integer) RETURNS TABLE(exists bool, url VARCHAR, old_price VARCHAR)
	LANGUAGE plpgsql
AS
$$
BEGIN
	IF NOT exists(SELECT FROM "user" WHERE "user".id = user_id_in AND verified_email = TRUE) THEN
		RETURN 'Пожалуйста подтвердите вашу электронную почту';
	END IF;
	IF NOT exists(
			SELECT
				FROM good_user
				WHERE good_user.user_id = user_id_in
				  AND good_user.good_id = good_id_in
				  AND now() > good_user.last_check + (cooldown_in_minutes * INTERVAL '1 minute')
		) THEN
		RETURN format('Лимит запросов превышен, последняя цена - %s',
		              (SELECT price FROM good WHERE id = good_id_in)::VARCHAR);
	END IF;
	RETURN (SELECT url FROM good WHERE id = good_id_in);
END
$$;

ALTER FUNCTION manual_check_validation(INTEGER, INTEGER, INTEGER) OWNER TO postgres;

