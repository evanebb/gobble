DROP TABLE IF EXISTS distro;
CREATE TABLE distro
(
    id               serial PRIMARY KEY,
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
    id               serial PRIMARY KEY,
    uuid             uuid UNIQUE,
    name             varchar(64),
    description      varchar(128),
    distro           uuid REFERENCES distro (uuid),
    kernelParameters varchar(128)[]
);

DROP TABLE IF EXISTS system;
CREATE TABLE system
(
    id               serial PRIMARY KEY,
    uuid             uuid UNIQUE,
    name             varchar(64),
    description      varchar(128),
    profile          uuid REFERENCES profile (uuid),
    mac              macaddr,
    kernelParameters varchar(128)[]
);
