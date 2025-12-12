CREATE DATABASE account;
USE account;
DROP TABLE IF EXISTS customer_individual;
CREATE TABLE customer_individual
(
    id_card    VARCHAR(20) PRIMARY KEY,
    name       VARCHAR(50) NOT NULL,
    email      VARCHAR(50) NOT NULL UNIQUE,
    phone      VARCHAR(15) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

DROP TABLE IF EXISTS customer_enterprise;
CREATE TABLE customer_enterprise
(
    # 社会统一信用代码
    credit_code   VARCHAR(20) PRIMARY KEY,
    company_name  VARCHAR(100) NOT NULL,
    # 法定代表人
    legal_id_card VARCHAR(20)  NOT NULL,
    legal_name    VARCHAR(50)  NOT NULL,
    email         VARCHAR(50)  NOT NULL UNIQUE,
    phone         VARCHAR(15)  NOT NULL UNIQUE,
    created_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

DROP TABLE IF EXISTS account;
CREATE TABLE account
(
    account_id   VARCHAR(20) PRIMARY KEY,
    account_type tinyint      NOT NULL default 0, #0:individual, 1:enterprise
    customer_id  VARCHAR(20)  NOT NULL,
    password     VARCHAR(255) NOT NULL,
    balance      DECIMAL(15, 2)        DEFAULT 0.00,
    status       TINYINT               DEFAULT 1, #1:active, 2:frozen
    created_at   TIMESTAMP             DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP             DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);