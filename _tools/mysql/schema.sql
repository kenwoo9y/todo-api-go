CREATE TABLE `users`
(
    `id`       BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `user_name`     varchar(30) NOT NULL,
    `email` VARCHAR(80) NOT NULL,
    `last_name` VARCHAR(40),
    `first_name` VARCHAR(40),
    `password` VARCHAR(80) NOT NULL,
    `created_at`  DATETIME(6) NOT NULL,
    `updated_at` DATETIME(6) NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uix_name` (`user_name`) USING BTREE
) Engine=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `tasks`
(
    `id`       BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `title`    VARCHAR(30) NOT NULL,
    `description` TEXT,
    `status`   VARCHAR(10)  NOT NULL,
    `created_at`  DATETIME(6) NOT NULL,
    `updated_at` DATETIME(6) NOT NULL,
    PRIMARY KEY (`id`)
) Engine=InnoDB DEFAULT CHARSET=utf8mb4;