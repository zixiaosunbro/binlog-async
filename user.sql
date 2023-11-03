CREATE TABLE `user_info` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `name` varchar(64) NOT NULL,
    `age` int unsigned NOT NULL,
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;