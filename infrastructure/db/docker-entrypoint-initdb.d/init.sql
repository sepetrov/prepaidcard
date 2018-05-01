CREATE TABLE card (
    uuid CHAR(128) NOT NULL PRIMARY KEY,
    available_balance BIGINT UNSIGNED NOT NULL,
    blocked_balance BIGINT UNSIGNED NOT NULL
)