
DROP DATABASE IF EXISTS stream;
CREATE DATABASE stream;
USE stream;

CREATE TABLE `words` (
  `word_id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '单词ID 自增主键',
  `user_id` bigint DEFAULT NULL COMMENT '用户ID',
  `word_name` varchar(256) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '单词名',
  `description` longtext COLLATE utf8mb4_general_ci COMMENT '描述',
  `explains` longtext COLLATE utf8mb4_general_ci COMMENT '解释',
  `pronounce_us` longtext COLLATE utf8mb4_general_ci COMMENT '美式发音',
  `pronounce_uk` longtext COLLATE utf8mb4_general_ci COMMENT '英式发音',
  `youdao_url` longtext COLLATE utf8mb4_general_ci COMMENT '有道词典的url',
  `tag_id` bigint DEFAULT NULL COMMENT 'tag ID',
  `phonetic_us` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci COMMENT '美式音标',
  `phonetic_uk` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci COMMENT '英式音标',
  PRIMARY KEY (`word_id`),
  UNIQUE KEY `uniq_words_user_id_word_name` (`user_id`,`word_name`),
  KEY `idx_words_word_id` (`word_id`),
  KEY `idx_words_user_id_word_id_desc` (`user_id`,`word_id` DESC)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='词语标签';

CREATE TABLE `word_tags` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '词典ID',
  `tag_name` varchar(50) COLLATE utf8mb4_general_ci NOT NULL COMMENT '标签名称',
  `question_types` int NOT NULL COMMENT '题型组合（位运算）',
  `max_score` int NOT NULL COMMENT '该标签的最大分数',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_tag_name` (`tag_name`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='词语标签';

CREATE TABLE `words_risite_record` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '自增ID 主键',
  `word_id` int NOT NULL COMMENT '单词ID',
  `level` int NOT NULL COMMENT '难度等级',
  `next_review_time` bigint NOT NULL COMMENT '下一次review的时间',
  `downgrade_step` int NOT NULL COMMENT '降级step',
  `total_correct` int NOT NULL COMMENT '全部正确',
  `total_wrong` int NOT NULL COMMENT '全部错误',
  `score` int NOT NULL COMMENT '得分',
  `user_id` bigint NOT NULL COMMENT '用户id',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_words_risite_user_id_next_review_time` (`user_id`,`next_review_time`),
  KEY `idx_words_risite_user_id_word_id` (`user_id`,`word_id`),
  KEY `idx_words_risite_user_id_level` (`user_id`,`level`)
) ENGINE=InnoDB AUTO_INCREMENT=22323 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='risite记录';

CREATE TABLE `review_progress` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `user_id` bigint NOT NULL COMMENT '用户ID',
  `pending_review_count` int NOT NULL DEFAULT '0' COMMENT '待复习单词数量',
  `completed_review_count` int NOT NULL DEFAULT '0' COMMENT '已完成复习单词数量',
  `last_update_time` bigint NOT NULL COMMENT '最后更新时间戳',
  `all_completed_count` int NOT NULL DEFAULT '0' COMMENT '总的完成数',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_id` (`user_id`)
) ENGINE=InnoDB AUTO_INCREMENT=6196 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='用户复习进度表';

CREATE TABLE `answer_list` (
  `answer_id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '答案ID',
  `word_id` int NOT NULL COMMENT '单词ID',
  `user_id` bigint DEFAULT NULL COMMENT '用户ID',
  `word_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '单词名',
  `description` varchar(1024) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '描述',
  `order_id` int NOT NULL COMMENT '订单ID',
  PRIMARY KEY (`answer_id`),
  KEY `idx_answer_list_user_id` (`user_id`),
  KEY `idx_answer_list_user_id_order_id` (`user_id`,`order_id`),
  KEY `idx_answer_list_user_id_word_id` (`user_id`,`word_id`),
  KEY `idx_answer_list_user_id_answer_id` (`user_id`,`answer_id`)
) ENGINE=InnoDB AUTO_INCREMENT=22362 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='答案列表';

INSERT INTO `word_tags` (`id`, `tag_name`, `question_types`, `max_score`) VALUES
(1, '熟练掌握', 15, 15),
(2, '能听懂', 7, 7),
(3, '能看懂', 3, 3);
