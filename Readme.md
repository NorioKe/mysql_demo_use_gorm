# 一个简单的使用gorm库进行mysql操作的demo

## 简介
本demo使用Go编写了一个简单的批跑demo，主要实现用户表（users）和订单表（orders）之间的简单交互，具体包含：
1. 创建用户
2. 创建订单（关联用户ID）
3. 设置无效订单（同步用户的消费总额）
代码按照领域驱动设计(DDD)设计。不过目前仅仅只是针对本人对于DDD的理解而设计，因此作为工作前的温习，后续会随着工作需求不断加深这个设计。
此外，实现本demo是巩固go编程的同时，也是在学习使用gorm库，不过目前仅仅只是作为入门练习。

gorm库的中文文档：https://gorm.io/zh_CN/docs/

另附本demo涉及的两个表：
```sql
Create Table: CREATE TABLE `users` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(100) DEFAULT NULL,
  `email` varchar(255) DEFAULT NULL,
  `total_consumption` decimal(12,2) NOT NULL DEFAULT '0.00' COMMENT '用户消费总额',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_users_email` (`email`)
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=latin1

Create Table: CREATE TABLE `orders` (
  `order_id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint(20) unsigned NOT NULL COMMENT '关联users.id',
  `amount` decimal(12,2) NOT NULL COMMENT '订单金额',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `is_valid` tinyint(1) NOT NULL DEFAULT '1' COMMENT '有效性标识(0:无效 1:有效)',
  PRIMARY KEY (`order_id`),
  KEY `idx_user_id` (`user_id`),
  CONSTRAINT `fk_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=latin1
```

## 更新日志

### v1
实现了基本框架以及基本功能。
后续需要完善的点：
1. 需要提供一个可查询的接口。
2. 目前的 `./infrastructure/transaction_manager.go` 仅仅只是“形式化”的事务，目前的实现**丢失了上下文，并不能保证原子性操作**，换句话说目前的实现无法做到“要么全成功，要么全部失败回滚”。后续考虑引入上下文 context 对其进行优化。

### v1.0.1
修改了go.mod文件，更正模块名以及取消replace命令。
