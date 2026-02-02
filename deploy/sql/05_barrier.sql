-- Barrier 表用于 DTM 幂等控制
-- 需要在每个业务数据库中创建

USE account;
CREATE TABLE IF NOT EXISTS barrier (
  id BIGINT AUTO_INCREMENT PRIMARY KEY,
  trans_type VARCHAR(45),
  gid VARCHAR(128),
  branch_id VARCHAR(128),
  op VARCHAR(45),
  barrier_id VARCHAR(45),
  reason VARCHAR(45),
  create_time DATETIME DEFAULT CURRENT_TIMESTAMP,
  update_time DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  KEY (gid),
  KEY (create_time),
  KEY (update_time),
  UNIQUE KEY (gid, branch_id, op, barrier_id)
);

USE transaction;
CREATE TABLE IF NOT EXISTS barrier (
  id BIGINT AUTO_INCREMENT PRIMARY KEY,
  trans_type VARCHAR(45),
  gid VARCHAR(128),
  branch_id VARCHAR(128),
  op VARCHAR(45),
  barrier_id VARCHAR(45),
  reason VARCHAR(45),
  create_time DATETIME DEFAULT CURRENT_TIMESTAMP,
  update_time DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  KEY (gid),
  KEY (create_time),
  KEY (update_time),
  UNIQUE KEY (gid, branch_id, op, barrier_id)
);
