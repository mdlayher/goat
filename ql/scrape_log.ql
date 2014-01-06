BEGIN TRANSACTION;

CREATE TABLE scrape_log (
	info_hash string,
	passkey   string,
	ip        string,
	ts        time,
);

COMMIT;
