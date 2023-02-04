DROP TABLE IF EXISTS distro;
CREATE TABLE distro
(
    id               int PRIMARY KEY,
    uuid             uuid UNIQUE,
    name             varchar(64),
    description      varchar(128),
    kernel           varchar(128),
    initrd           varchar(128),
    kernelParameters varchar(128)[]
);

DROP TABLE IF EXISTS profile;
CREATE TABLE profile
(
    id               int PRIMARY KEY,
    uuid             uuid UNIQUE,
    name             varchar(64),
    description      varchar(128),
    distro           int REFERENCES distro (uuid),
    kernelParameters varchar(128)[]
);

DROP TABLE IF EXISTS system;
CREATE TABLE system
(
    id               int PRIMARY KEY,
    uuid             uuid UNIQUE,
    name             varchar(64),
    description      varchar(128),
    profile          int REFERENCES profile (uuid),
    mac              macaddr,
    kernelParameters varchar(128)[]
);
