-- Phase 2 建表脚本

CREATE TABLE IF NOT EXISTS `pay_callbacks` (
  `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  `order_id` BIGINT UNSIGNED NOT NULL COMMENT '订单ID',
  `transaction_id` VARCHAR(64) DEFAULT '' COMMENT '微信订单号',
  `amount` INT NOT NULL COMMENT '支付金额(分)',
  `status` VARCHAR(20) DEFAULT '' COMMENT '状态',
  `raw_xml` TEXT COMMENT '原始回调XML',
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
  INDEX `idx_order_id` (`order_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='支付回调记录表';

CREATE TABLE IF NOT EXISTS `member_levels` (
  `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  `name` VARCHAR(50) NOT NULL COMMENT '等级名称',
  `min_points` INT NOT NULL DEFAULT 0 COMMENT '最低积分门槛',
  `discount` DECIMAL(3,2) NOT NULL DEFAULT 1.00 COMMENT '折扣率',
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='会员等级表';

INSERT INTO `member_levels` (`name`, `min_points`, `discount`) VALUES
  ('普通会员', 0, 1.00),
  ('银卡会员', 500, 0.95),
  ('金卡会员', 2000, 0.90),
  ('钻石会员', 5000, 0.85);

CREATE TABLE IF NOT EXISTS `members` (
  `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  `user_id` BIGINT UNSIGNED NOT NULL UNIQUE COMMENT '用户ID',
  `level_id` BIGINT UNSIGNED NOT NULL DEFAULT 1 COMMENT '等级ID',
  `points` INT NOT NULL DEFAULT 0 COMMENT '当前积分',
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  INDEX `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='会员表';

CREATE TABLE IF NOT EXISTS `coupons` (
  `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  `name` VARCHAR(100) NOT NULL COMMENT '优惠券名称',
  `type` VARCHAR(20) NOT NULL COMMENT '类型：cash满减, discount折扣',
  `threshold` DECIMAL(10,2) NOT NULL DEFAULT 0 COMMENT '使用门槛金额',
  `discount` DECIMAL(10,2) NOT NULL COMMENT '优惠金额或折扣率',
  `total_count` INT NOT NULL DEFAULT 0 COMMENT '总数量',
  `left_count` INT NOT NULL DEFAULT 0 COMMENT '剩余数量',
  `expire_at` DATETIME NOT NULL COMMENT '过期时间',
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='优惠券表';

CREATE TABLE IF NOT EXISTS `user_coupons` (
  `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
  `coupon_id` BIGINT UNSIGNED NOT NULL COMMENT '优惠券ID',
  `used` TINYINT DEFAULT 0 COMMENT '是否已使用',
  `used_at` DATETIME DEFAULT NULL COMMENT '使用时间',
  `order_id` BIGINT UNSIGNED DEFAULT NULL COMMENT '使用订单ID',
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
  INDEX `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户优惠券表';

-- 初始化一些优惠券
INSERT INTO `coupons` (`name`, `type`, `threshold`, `discount`, `total_count`, `left_count`, `expire_at`) VALUES
  ('新人满减券', 'cash', 100.00, 10.00, 1000, 1000, DATE_ADD(CURDATE(), INTERVAL 90 DAY)),
  ('满200减20', 'cash', 200.00, 20.00, 500, 500, DATE_ADD(CURDATE(), INTERVAL 90 DAY)),
  ('95折优惠券', 'discount', 0.00, 0.95, 300, 300, DATE_ADD(CURDATE(), INTERVAL 90 DAY));
