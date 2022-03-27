create database test;
use test;

CREATE TABLE users
(
	id INTEGER AUTO_INCREMENT,
	username varchar(256),
	password varchar(256),
	status varchar(256) NULL DEFAULT 'offline',
	primary key(id)
);
