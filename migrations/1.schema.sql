DROP TABLE IF EXISTS `user`;
CREATE TABLE `user` (
    `id` INTEGER PRIMARY KEY,
    `name` VARCHAR(32) NOT NULL,
    `password` VARCHAR(64) NOT NULL,
    `avatar` VARCHAR(256),

    `created_at` DATETIME DEFAULT current_timestamp,
    `updated_at` DATETIME DEFAULT current_timestamp
);

DROP TABLE IF EXISTS `user_friend`;
CREATE TABLE `user_friend` (
    `user_a_id` INTEGER NOT NULL,
    `user_b_id` INTEGER NOT NULL,
    
    `created_at` DATETIME DEFAULT current_timestamp
);

DROP TABLE IF EXISTS `token`;
CREATE TABLE `token` (
    `user_id` INTEGER NOT NULL,
    `code` VARCHAR(64) NOT NULL,

    `created_at` DATETIME DEFAULT current_timestamp,
    `updated_at` DATETIME DEFAULT current_timestamp
);

DROP TABLE IF EXISTS `t`;
CREATE TABLE `t` (
    `id` INTEGER PRIMARY KEY,
    `user_id` INTEGER NOT NULL,
    `text` VARCHAR(5000),
    `image` VARCHAR(256),

    `created_at` DATETIME DEFAULT current_timestamp
);

DROP TABLE IF EXISTS `t_like`;
CREATE TABLE `t_like` (
    `t_id` INTEGER NOT NULL,
    `user_id` INTEGER NOT NULL,

    `created_at` DATETIME DEFAULT current_timestamp
);

DROP TABLE IF EXISTS `t_comment`;
CREATE TABLE `t_comment` (
    `id` INTEGER PRIMARY KEY,
    `t_id` INTEGER NOT NULL,
    `user_id` INTEGER NOT NULL,
    `content` VARCHAR(5000) NOT NULL,

    `created_at` DATETIME DEFAULT current_timestamp
);

/*
 * NOTE:
 *
 *      1. no foreign key
 *      2. hashed field should VARCHAR(64)
 *      3. datetime fields should be `:name:_at`
 *      4. resource fields should be VARCHAR(256)
 **/
