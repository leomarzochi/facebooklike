CREATE DATABASE IF NOT EXISTS facebooklike;

USE facebooklike;

DROP TABLE IF EXISTS usuarios;

-- Active: 1684881718745@@127.0.0.1@3306@facebooklike
CREATE TABLE users(
    id INT AUTO_INCREMENT primary key,
    nome VARCHAR(50) NOT NULL,
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(50) NOT NULL UNIQUE,
    password VARCHAR(50) NOT NULL,
    createdAt TIMESTAMP default current_timestamp()
) ENGINE=INNODB;

CREATE TABLE followers(
    user_id INT NOT NULL,
    FOREIGN KEY(user_id)
    REFERENCES users(id)
    ON DELETE CASCADE,

    follower_id INT NOT NULL,
    FOREIGN KEY(follower_id)
    REFERENCES users(id)
    ON DELETE CASCADE,

    PRIMARY KEY(user_id, follower_id)
) ENGINE=INNODB;