
CREATE TABLE `words` (
  `word_id` INTEGER PRIMARY KEY,
  `user_id` bigint DEFAULT NULL,
  `word_name` varchar(256) DEFAULT NULL,
  `description` longtext ,
  `explains` longtext,
  `pronounce_us` longtext ,
  `pronounce_uk` longtext ,
  `youdao_url` longtext ,
  `tag_id` bigint DEFAULT NULL,
  `phonetic_us` longtext,
  `phonetic_uk` longtext
);

CREATE TABLE `word_tags` (
  `id` INTEGER PRIMARY KEY ,
  `tag_name` varchar(50) NOT NULL,
  `question_types` int NOT NULL,
  `max_score` int NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
); 

CREATE TRIGGER update_word_tags_updated_at
AFTER UPDATE ON word_tags
BEGIN
  UPDATE word_tags
  SET updated_at = CURRENT_TIMESTAMP
  WHERE id = NEW.id;
END;

CREATE TABLE `words_risite_record` (
  `id` INTEGER PRIMARY KEY,
  `word_id` int NOT NULL,
  `level` int NOT NULL,
  `next_review_time` bigint NOT NULL,
  `downgrade_step` int NOT NULL,
  `total_correct` int NOT NULL ,
  `total_wrong` int NOT NULL ,
  `score` int NOT NULL,
  `user_id` bigint NOT NULL
);

CREATE TABLE `review_progress` (
  `id` INTEGER PRIMARY KEY,
  `user_id` bigint NOT NULL,
  `pending_review_count` int NOT NULL DEFAULT '0',
  `completed_review_count` int NOT NULL DEFAULT '0',
  `last_update_time` bigint NOT NULL,
  `all_completed_count` int NOT NULL DEFAULT '0'
);

CREATE TABLE `answer_list` (
  `answer_id` INTEGER PRIMARY KEY ,
  `word_id` int NOT NULL ,
  `user_id` bigint DEFAULT NULL ,
  `word_name` varchar(255) DEFAULT NULL ,
  `description` varchar(1024) DEFAULT NULL,
  `order_id` int NOT NULL
);

INSERT INTO `word_tags` (`id`, `tag_name`, `question_types`, `max_score`) VALUES
(1, '熟练掌握', 15, 15),
(2, '能听懂', 7, 7),
(3, '能看懂', 3, 3);
