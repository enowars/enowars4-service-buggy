USE enodb;

CREATE TABLE users (
name char(64),
password char(64),
status varchar(255),
admin boolean
);

CREATE TABLE messages (
name char(64),
sender char(64),
message varchar(255)
);

CREATE TABLE comments (
name char(64),
product char(64),
content varchar(255)
);

INSERT INTO enodb.users VALUES("admin", "root", "some status", true);
