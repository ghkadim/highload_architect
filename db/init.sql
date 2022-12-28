CREATE TABLE `users`
(
    `id` bigint NOT NULL AUTO_INCREMENT,
    `first_name` varchar(255) DEFAULT NULL,
    `last_name` varchar(255) DEFAULT NULL,
    `sex` enum ('male', 'female') DEFAULT NULL,
    `age` smallint DEFAULT NULL,
    `city` varchar(255) DEFAULT NULL,
    `biography` text DEFAULT NULL,
    PRIMARY KEY (`id`),
    KEY `users_first_last_name_IDX` (`first_name`,`last_name`) USING BTREE
);
