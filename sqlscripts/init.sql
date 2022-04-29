CREATE TABLE IF NOT EXISTS `user`
(
    `id`        INT         NOT NULL AUTO_INCREMENT,
    `username`  VARCHAR(45) NOT NULL COMMENT 'username',
    `password`  VARCHAR(60) NOT NULL COMMENT 'password',
    `name`      VARCHAR(45) NOT NULL COMMENT 'real name',
    `age`       INT         NOT NULL DEFAULT '0' COMMENT 'age',
    `phone`     VARCHAR(45) NOT NULL COMMENT 'phone number',
    `create_at` DATETIME    NULL DEFAULT current_timestamp,
    `update_at` DATETIME    NULL DEFAULT current_timestamp on update current_timestamp,
    `delete_at` DATETIME    NULL,
    PRIMARY KEY (`id`)
);
