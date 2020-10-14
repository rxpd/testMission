CREATE FUNCTION new_message_handler(chat_id_in integer, username_in character varying) RETURNS character varying
	LANGUAGE plpgsql
AS
$$
BEGIN
	IF exists(SELECT FROM tg_user WHERE chat_id = chat_id_in) -- если пользователь в БД
	THEN
		IF (SELECT username FROM tg_user WHERE chat_id = chat_id_in) != username_in -- если пользователь поменял имя
		THEN
			UPDATE tg_user
			SET username = username_in
				WHERE chat_id = chat_id_in;
		END IF;
		RETURN 'user exists';
	END IF;
	INSERT
		INTO tg_user
			(chat_id, username)
		VALUES
			(chat_id_in, username_in);
	RETURN 'new user';
END
$$;

ALTER FUNCTION new_message_handler(INTEGER, VARCHAR) OWNER TO postgres;

