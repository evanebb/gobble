drop table if exists distro;
create table distro
(
    id               int primary key,
    name             varchar(64),
    description      varchar(128),
    kernel           varchar(128),
    initrd           varchar(128),
    kernelParameters varchar(128)[]
)

drop table if exists profile;
create table profile
(
    id               int primary key,
    name             varchar(64),
    description      varchar(128),
    distro           int references distro (id),
    kernelParameters varchar(128)[]
)

drop table if exists system;
create table system
(
    id               int primary key,
    name             varchar(64),
    description      varchar(128),
    profile          int references profile (id),
    mac              macaddr,
    kernelParameters varchar(128)[]
)
