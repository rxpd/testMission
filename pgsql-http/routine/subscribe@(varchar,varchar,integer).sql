CREATE FUNCTION subscribe(email_in character varying, good_url_in character varying, price_in integer) RETURNS character varying
	LANGUAGE plpgsql
AS
$$
DECLARE
	good_id_vr       INTEGER;
	user_id_vr       INTEGER;
	response_message VARCHAR;
BEGIN
	response_message = 'ok';
	good_id_vr = (SELECT id FROM good WHERE url = good_url_in); -- заполняю переменные
	user_id_vr = (SELECT id FROM "user" WHERE email = email_in); -- заполняю переменные
	IF exists(SELECT FROM good_user WHERE good_user.good_id = good_id_vr AND good_user.user_id = user_id_vr) THEN
		RETURN 'this subscribe already in database';
	END IF;
-- 	RAISE NOTICE '% %', good_id_vr, user_id_vr;
	CASE WHEN good_id_vr ISNULL AND user_id_vr ISNULL THEN -- если ни товар ни пользователь не зарегистрирован
		INSERT
			INTO good
				(url, price)
			VALUES
				(good_url_in, price_in)
			RETURNING id INTO good_id_vr;
		INSERT
			INTO "user"
				(email, verified_email)
			VALUES
				(email_in, DEFAULT)
			RETURNING id INTO user_id_vr;
		response_message = 'new user registered';
		WHEN good_id_vr ISNULL AND user_id_vr NOTNULL THEN -- если товар не зарегистрирован, а пользователь зарегистрирован
			INSERT
				INTO good
					(url, price)
				VALUES
					(good_url_in, price_in)
				RETURNING id INTO good_id_vr;
		WHEN good_id_vr NOTNULL AND user_id_vr ISNULL THEN -- если товар зарегистрирован, а пользователь нет
			INSERT
				INTO "user"
					(email, verified_email)
				VALUES
					(email_in, DEFAULT)
				RETURNING id INTO user_id_vr;
			response_message = 'new user registered';
		ELSE
		END CASE;


	INSERT
		INTO good_user
			(good_id, user_id, price)
		VALUES
			(good_id_vr, user_id_vr, price_in);
	RETURN response_message;
END
$$;

ALTER FUNCTION subscribe(VARCHAR, VARCHAR, INTEGER) OWNER TO postgres;

