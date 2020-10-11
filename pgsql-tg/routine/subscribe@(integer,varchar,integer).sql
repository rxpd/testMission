CREATE FUNCTION subscribe(chat_id_in integer, good_url_in character varying, price_in integer) RETURNS character varying
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
	user_id_vr = (SELECT id FROM users WHERE chat_id = chat_id_in); -- заполняю переменные
	IF exists(SELECT FROM good_user WHERE good_user.good_id = good_id_vr AND good_user.user_id = user_id_vr) THEN
		RETURN 'this subscribe already in database';
	END IF;
	IF good_id_vr ISNULL THEN
		INSERT
			INTO good
				(url, price)
			VALUES
				(good_url_in, price_in)
			RETURNING id INTO good_id_vr;
	END IF;

	INSERT
		INTO good_user
			(good_id, user_id, price)
		VALUES
			(good_id_vr, user_id_vr, price_in);
	RETURN response_message;
END
$$;

ALTER FUNCTION subscribe(INTEGER, VARCHAR, INTEGER) OWNER TO postgres;

