create database sysreqs;
use sysreqs;

create table package(
  id integer not null auto_increment primary key,
  name varchar(25) not null,
  version varchar(11) not null,
  architecture varchar(7) not null,
  description text,
  platforms JSON not null,
  dependencies JSON
);

create unique index idx_package_name_version on package (name, version);
