CREATE TABLE snippet (
	id serial primary key,
	title varchar(150) not null,
	content text not null,
	created timestamp not null,
	expired timestamp not null
);