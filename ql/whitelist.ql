BEGIN TRANSACTION;

CREATE TABLE whitelist (
	client   string,
	approved bool
);

COMMIT;
