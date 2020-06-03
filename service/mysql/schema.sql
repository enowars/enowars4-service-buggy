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
message varchar(255),
INDEX (hash)
);

CREATE TABLE comments (
name char(64),
product char(64),
content varchar(255),
INDEX (product)
);

CREATE TABLE tickets (
name char(64),
subject char(64),
hash char(64),
PRIMARY KEY (hash),
UNIQUE INDEX (hash)
);

INSERT INTO enodb.users VALUES("admin", "root", "some status", true);
