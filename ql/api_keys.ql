BEGIN TRANSACTION;

CREATE TABLE api_keys (
	user_id int,
	key string,
);

COMMIT;
