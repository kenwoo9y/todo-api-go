CREATE TABLE `users` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `username` VARCHAR(30) UNIQUE,
    `email` VARCHAR(80) UNIQUE,
    `first_name` VARCHAR(40),
    `last_name` VARCHAR(40),
    `created_at` DATETIME(6) NOT NULL,
    `updated_at` DATETIME(6) NOT NULL,
    PRIMARY KEY (`id`)
);

CREATE TABLE `tasks` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `title` VARCHAR(30),
    `description` VARCHAR(255),
    `due_date` DATE,
    `status` VARCHAR(10),
    `owner_id` BIGINT,
    `created_at` DATETIME(6) NOT NULL,
    `updated_at` DATETIME(6) NOT NULL,
    PRIMARY KEY (`id`),
    FOREIGN KEY (`owner_id`) REFERENCES `users`(`id`)
);