BEGIN TRANSACTION;

CREATE TABLE api_keys (
	user_id int64,
	pubkey string,
	secret string,
	expire int64,
);

COMMIT;
