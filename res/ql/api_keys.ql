BEGIN TRANSACTION;

CREATE TABLE api_keys (
	user_id int64,
	key string,
	expire int64,
);

COMMIT;
