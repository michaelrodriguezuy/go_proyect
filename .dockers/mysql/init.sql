SET @MYSQLDUMP_TEMP_LOG_BIN = @@SESSION.sql_log_bin;
SET @@SESSION.sql_log_bin = 0;

SET @@GLOBAL.GTID_PURGED= /*!80000 '+'*/ '';

CREATE DATABASE IF NOT EXISTS `go_users`;

CREATE TABLE IF NOT EXISTS `go_users`.`users` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `first_name` VARCHAR(45) NOT NULL,
  `last_name` VARCHAR(45) NOT NULL,
  `age` INT NOT NULL,  
  PRIMARY KEY (`id`))
