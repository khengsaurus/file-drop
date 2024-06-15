DROP DATABASE IF EXISTS db1;
CREATE DATABASE db1;
USE db1;
CREATE TABLE IF NOT EXISTS Urls (
  CreatedAt INT(11),
  Id VARCHAR(36) NOT NULL,
  Link TEXT
);
-- SET SQL_SAFE_UPDATES = 0;
-- ALTER USER 'root' @'localhost' IDENTIFIED WITH mysql_native_password BY 'password';
-- flush privileges;