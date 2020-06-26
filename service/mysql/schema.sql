USE enodb;
CREATE TABLE users (
        name char(64),
        password char(64),
        status varchar(255),
        admin boolean,
        PRIMARY KEY (name),
        UNIQUE INDEX (name)
);
CREATE TABLE messages (
        name char(64),
        sender char(64),
        hash char(64),
        message varchar(512),
        INDEX (hash),
        INDEX (name)
);
CREATE TABLE comments (
        id int NOT NULL AUTO_INCREMENT,
        created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
        name char(64),
        product char(64),
        content varchar(512),
        PRIMARY KEY (id),
        INDEX (id)
);
CREATE TABLE tickets (
        name char(64),
        subject char(64),
        hash char(64),
        PRIMARY KEY (hash),
        INDEX (hash)
);
INSERT INTO enodb.users
VALUES("admin", "root", "some status", true);
DELIMITER | CREATE EVENT ttl_delete ON SCHEDULE EVERY 300 SECOND DO BEGIN
DELETE FROM comments
WHERE created_at < NOW() - INTERVAL 1800 SECOND;
END | DELIMITER;