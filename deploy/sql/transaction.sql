CREATE DATABASE transaction;
USE transaction;

DROP TABLE IF EXISTS transaction_record;
CREATE TABLE transaction_record
(
    transaction_id   VARCHAR(20) PRIMARY KEY,
    account_from     VARCHAR(20)    NOT NULL,
    account_to       VARCHAR(20),
    transaction_type TINYINT        NOT NULL, #0:transfer, 1:deposit, 2:withdraw
    amount           DECIMAL(15, 2) NOT NULL,
    created_at       TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    status           TINYINT   DEFAULT 1      #0:failed, 1:successful
);