CREATE FUNCTION new_message_handler(chat_id_in integer, user_name_in character varying) RETURNS character varying
	LANGUAGE plpgsql
AS
$$
BEGIN
	IF exists(SELECT FROM users WHERE chat_id = chat_id_in) -- если пользователь в БД
	THEN
		IF (SELECT user_name FROM users WHERE chat_id = chat_id_in) != user_name_in -- если пользователь поменял имя
		THEN
			UPDATE users
			SET user_name = user_name_in
				WHERE chat_id = chat_id_in;
		END IF;
		RETURN 'user exists';
	END IF;
	INSERT
		INTO users
			(chat_id, user_name)
		VALUES
			(chat_id_in, user_name_in);
	RETURN 'new user';
END
$$;

ALTER FUNCTION new_message_handler(INTEGER, VARCHAR) OWNER TO postgres;

