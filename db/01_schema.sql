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
