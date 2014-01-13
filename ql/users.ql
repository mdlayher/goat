BEGIN TRANSACTION;

CREATE TABLE users (
	username      string,
	passkey       string,
	torrent_limit int
);

COMMIT;
