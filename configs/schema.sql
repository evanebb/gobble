DROP TABLE IF EXISTS profile;
CREATE TABLE profile
(
    id               serial PRIMARY KEY,
    uuid             uuid UNIQUE,
    name             varchar(64) UNIQUE,
    description      varchar(128),
    kernel           varchar(128),
    initrd           varchar(128),
    kernelParameters varchar(128)[]
);

DROP TABLE IF EXISTS system;
CREATE TABLE system
(
    id               serial PRIMARY KEY,
    uuid             uuid UNIQUE,
    name             varchar(64) UNIQUE,
    description      varchar(128),
    profile          uuid REFERENCES profile (uuid) ON DELETE CASCADE,
    mac              macaddr UNIQUE,
    kernelParameters varchar(128)[]
);

DROP TABLE IF EXISTS api_user;
CREATE TABLE api_user
(
    id       serial PRIMARY KEY,
    uuid     uuid UNIQUE,
    name     varchar(64) UNIQUE,
    password varchar
);

INSERT INTO api_user (uuid, name, password) VALUES ('62fb65af-2d12-4758-93d6-7b58eadde3f1', 'admin', '$2a$10$dYnBNGXrDH/1Rf75zqkENelFhrmPEQrUTARkgYOFhKyGJn/nvi90e');
