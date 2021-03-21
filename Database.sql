DROP DATABASE ProjetoUnivesp2021;

CREATE DATABASE IF NOT EXISTS ProjetoUnivesp2021;

-- ProjetoUnivesp2021.HOSPITAL definition

CREATE TABLE `HOSPITAL`
(
    `HID`              int(10) unsigned    NOT NULL AUTO_INCREMENT,
    `HNOME`            varchar(255)        NOT NULL,
    `HUF`              char(2)             NOT NULL,
    `HCIDADE`          varchar(50)         NOT NULL,
    `HCEP`             varchar(8)          NOT NULL,
    `HENDERECO`        varchar(150)        NOT NULL,
    `HCOMPLEMENTO`     varchar(150)                 DEFAULT NULL,
    `HTELEFONE`        bigint(14) unsigned NOT NULL,
    `HISPRONTOSOCORRO` tinyint(1)          NOT NULL DEFAULT 0,
    `HATIVADO`         enum ('T','F','D')           DEFAULT 'T',
    PRIMARY KEY (`HID`)
);

-- ProjetoUnivesp2021.ESPECIALIDADES definition

CREATE TABLE `ESPECIALIDADES`
(
    `EID`   int(10) unsigned NOT NULL AUTO_INCREMENT,
    `HNOME` varchar(30) DEFAULT NULL,
    PRIMARY KEY (`EID`),
    UNIQUE KEY `HNOME` (`HNOME`)
);

-- ProjetoUnivesp2021.MEDICOS definition

CREATE TABLE `MEDICOS`
(
    `MID`      bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `MNOME`    varchar(255)        NOT NULL,
    `EID`      int(10) unsigned    NOT NULL,
    `HID`      int(10) unsigned    NOT NULL,
    `MATIVADO` enum ('T','F','D')  NOT NULL DEFAULT 'T',
    PRIMARY KEY (`MID`)
);

-- ProjetoUnivesp2021.USUARIOS definition

CREATE TABLE `USUARIOS`
(
    `UID`          bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `UNOME`        varchar(255)        NOT NULL,
    `UEMAIL`       varchar(100)        NOT NULL,
    `UPASSWORD`    varchar(80)         NOT NULL,
    `UTOKEN`       char(16)            NOT NULL,
    `UCPF`         char(11)            NOT NULL,
    `UUF`          char(2)             NOT NULL,
    `UCIDADE`      varchar(50)         NOT NULL,
    `UCEP`         varchar(8)          NOT NULL,
    `UENDERECO`    varchar(150)        NOT NULL,
    `UCOMPLEMENTO` varchar(150)                 DEFAULT NULL,
    `UATIVADO`     enum ('T','F','D')  NOT NULL DEFAULT 'T',
    PRIMARY KEY (`UID`),
    UNIQUE KEY `UEMAIL` (`UEMAIL`),
    UNIQUE KEY `UTOKEN` (`UTOKEN`),
    UNIQUE KEY `UCPF` (`UCPF`)
);

-- ProjetoUnivesp2021.AGENDAMENTOS definition

CREATE TABLE `AGENDAMENTOS`
(
    `AID`      bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `UID`      bigint(20) unsigned NOT NULL,
    `MID`      bigint(20) unsigned NOT NULL,
    `ADATA`    datetime            NOT NULL,
    `AATIVADO` enum ('T','D')      NOT NULL DEFAULT 'T',
    PRIMARY KEY (`AID`),
    UNIQUE KEY `ADATA` (`ADATA`)
);