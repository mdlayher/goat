BEGIN TRANSACTION;

CREATE TABLE users (
	username      string,
	password      string,
	passkey       string,
	torrent_limit int
);

COMMIT;
