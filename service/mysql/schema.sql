USE enodb;

CREATE TABLE users (
name char(64),
password char(64),
status char(250),
admin boolean
);

INSERT INTO enodb.users VALUES("admin", "root", "some status", true);
