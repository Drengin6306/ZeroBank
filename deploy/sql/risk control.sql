CREATE DATABASE risk_control;
USE risk_control;

DROP TABLE IF EXISTS risk_record;
CREATE TABLE risk_record
(
    id             INT AUTO_INCREMENT PRIMARY KEY,
    account_id     VARCHAR(20)    NOT NULL,
    transaction_id VARCHAR(20)    NOT NULL,
    risk_type      tinyint        NOT NULL, # 1:账户冻结 2:日转账限额 3:单笔转账限额 4:日提现限额 5:单笔提现限额
    risk_value     DECIMAL(15, 2) NOT NULL,
    created_at     TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

# 限额金额
DROP TABLE IF EXISTS limit_amount;
CREATE TABLE limit_amount
(
    id           INT AUTO_INCREMENT PRIMARY KEY,
    account_type tinyint        NOT NULL, # 0:个人账户 1:企业账户
    limit_type   tinyint        NOT NULL, # 1:日转账限额 2:单笔转账限额 3:日提现限额 4:单笔提现限额
    amount       DECIMAL(15, 2) NOT NULL,
    updated_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

INSERT INTO limit_amount (account_type, limit_type, amount)
VALUES
    -- 个人账户 0
    (0, 1, 500000.00),  -- 个人 日转账限额
    (0, 2, 100000.00),  -- 个人 单笔转账限额
    (0, 3, 100000.00),  -- 个人 日提现限额
    (0, 4, 20000.00),   -- 个人 单笔提现限额

    -- 企业账户 1
    (1, 1, 5000000.00), -- 企业 日转账限额
    (1, 2, 500000.00),  -- 企业 单笔转账限额
    (1, 3, 500000.00),  -- 企业 日提现限额
    (1, 4, 100000.00); -- 企业 单笔提现限额
