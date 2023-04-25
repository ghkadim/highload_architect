CREATE TABLE `users`
(
    `id` bigint AUTO_INCREMENT PRIMARY KEY,
    `uuid` binary(16) DEFAULT (uuid_to_bin(uuid())) NOT NULL,
    `first_name` varchar(255) DEFAULT NULL,
    `second_name` varchar(255) DEFAULT NULL,
    `age` smallint DEFAULT NULL,
    `city` varchar(255) DEFAULT NULL,
    `biography` text DEFAULT NULL,
    `password_hash` binary(100) NOT NULL,
    KEY `users_first_second_name_IDX` (`first_name`, `second_name`) USING BTREE,
    KEY `users_uuid_IDX` (`uuid`) USING BTREE
);

CREATE TABLE `friends`
(
    `id` BIGINT AUTO_INCREMENT PRIMARY KEY,
    `user1_id` BIGINT NOT NULL,
    `user2_id` BIGINT NOT NULL,
    KEY `friends_user1_id_IDX` (`user1_id`, `user2_id`) USING BTREE,
    CONSTRAINT `friends_FK_1` FOREIGN KEY (`user1_id`) REFERENCES users(`id`) ON DELETE CASCADE,
    CONSTRAINT `friends_FK_2` FOREIGN KEY (`user2_id`) REFERENCES users(`id`) ON DELETE CASCADE
);

CREATE TABLE `posts`
(
    `id` BIGINT AUTO_INCREMENT PRIMARY KEY,
    `uuid` binary(16) DEFAULT (uuid_to_bin(uuid())) NOT NULL,
    `text` BLOB NULL,
    `user_id` BIGINT NOT NULL,
    KEY `posts_uuid_IDX` (`uuid`) USING BTREE,
    KEY `posts_user_id_IDX` (`user_id`) USING BTREE,
    CONSTRAINT `posts_FK` FOREIGN KEY (`user_id`) REFERENCES users(`id`) ON DELETE CASCADE
);
