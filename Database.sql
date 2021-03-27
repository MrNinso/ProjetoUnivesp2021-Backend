DROP DATABASE ProjetoUnivesp2021;

CREATE DATABASE IF NOT EXISTS ProjetoUnivesp2021;

USE ProjetoUnivesp2021;

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
    `ENOME` varchar(30) DEFAULT NULL,
    PRIMARY KEY (`EID`),
    UNIQUE KEY `ENOME` (`ENOME`)
);

-- ProjetoUnivesp2021.MEDICOS definition

CREATE TABLE `MEDICOS`
(
    `MID`      bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `MNOME`    varchar(255)        NOT NULL,
    `EID`      int(10) unsigned    NOT NULL,
    `HID`      int(10) unsigned    NOT NULL,
    `MATIVADO` enum ('T','F','D')  NOT NULL DEFAULT 'T',
    PRIMARY KEY (`MID`),
    FOREIGN KEY (`HID`) REFERENCES `HOSPITAL` (`HID`),
    FOREIGN KEY (`EID`) REFERENCES `ESPECIALIDADES` (`EID`)
);

-- ProjetoUnivesp2021.USUARIOS definition

CREATE TABLE `USUARIOS`
(
    `UID`          bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `UNOME`        varchar(255)        NOT NULL,
    `UEMAIL`       varchar(100)        NOT NULL,
    `UPASSWORD`    varchar(80)         NOT NULL,
    `UTOKEN`       char(36),
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
    FOREIGN KEY (`UID`) REFERENCES `USUARIOS` (`UID`),
    FOREIGN KEY (`MID`) REFERENCES `MEDICOS` (`MID`)
);

-- Get Login Password from email

CREATE PROCEDURE GetLoginUsuario(email VARCHAR(100))
BEGIN
    SELECT UPASSWORD FROM USUARIOS WHERE UATIVADO = 'T' AND UEMAIL = email LIMIT 1;
END;

-- Add a User

CREATE PROCEDURE ProjetoUnivesp2021.AddUsuario(nome VARCHAR(255), email VARCHAR(100),
                                               password VARCHAR(80), cpf varchar(11),
                                               uf varchar(2), cidade varchar(50), cep varchar(8),
                                               endereco varchar(150), complemento varchar(150))

BEGIN
    INSERT INTO USUARIOS (UNOME, UEMAIL, UPASSWORD, UCPF, UUF, UCIDADE, UCEP, UENDERECO, UCOMPLEMENTO)
    VALUES (nome, email, password, cpf, uf, cidade, cep, endereco, complemento);
END;

-- Valida o token de um email

CREATE PROCEDURE ProjetoUnivesp2021.ValidarToken(email varchar(100), token char(36))
BEGIN
    SELECT IF(UTOKEN = token, 1, 0) FROM USUARIOS WHERE UEMAIL = email AND UATIVADO = 'T' LIMIT 1;
END;

-- Atualiza um token de um email

CREATE PROCEDURE ProjetoUnivesp2021.RegistrarToken(email varchar(100), token char(36))
BEGIN
    UPDATE USUARIOS SET UTOKEN = token WHERE UEMAIL = email AND UATIVADO = 'T';
END;

CREATE PROCEDURE ProjetoUnivesp2021.LogOff(email varchar(100), token varchar(36))
BEGIN
    UPDATE USUARIOS SET UTOKEN = NULL WHERE UEMAIL = email AND UTOKEN = token AND UATIVADO = 'T';
END;