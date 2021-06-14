DROP DATABASE IF EXISTS ProjetoUnivesp2021;

CREATE DATABASE IF NOT EXISTS ProjetoUnivesp2021;

USE ProjetoUnivesp2021;

-- ProjetoUnivesp2021.CONVENIO definition

CREATE TABLE `CONVENIO`
(
    `CID`       int(10) unsigned    NOT NULL AUTO_INCREMENT,
    `CNOME`     varchar(255)        NOT NULL,
    `CATIVADO`  enum('T', 'F')      NOT NULL DEFAULT 'T',
    PRIMARY KEY (`CID`),
    UNIQUE (`CNOME`)
);

-- ProjetoUnivesp2021.CONVENIO_PLANOS definition

CREATE TABLE `CONVENIO_PLANOS`
(
    `CPID`    bigint(20) unsigned    NOT NULL AUTO_INCREMENT,
    `CPNOME`  varchar(255)           NOT NULL,
    `CID`     int(10) unsigned       NOT NULL,
    `CPATIVADO`  enum('T', 'F')      NOT NULL DEFAULT 'T',
    PRIMARY KEY (`CPID`),
    UNIQUE (`CPNOME`, `CID`),
    FOREIGN KEY (`CID`) REFERENCES `CONVENIO` (`CID`)
);

-- ProjetoUnivesp2021.HOSPITAL definition

CREATE TABLE `HOSPITAIS`
(
    `HID`              int(10) unsigned    NOT NULL AUTO_INCREMENT,
    `HNOME`            varchar(255)        NOT NULL,
    `HUF`              char(2)             NOT NULL,
    `HCIDADE`          varchar(50)         NOT NULL,
    `HCEP`             char(8)             NOT NULL,
    `HENDERECO`        varchar(150)        NOT NULL,
    `HCOMPLEMENTO`     varchar(150)                 DEFAULT NULL,
    `HTELEFONE`        bigint(14) unsigned NOT NULL,
    `HISPRONTOSOCORRO` tinyint(1)          NOT NULL DEFAULT 0,
    `HATIVADO`         enum ('T','F','D')           DEFAULT 'T',
    PRIMARY KEY (`HID`)
);

-- ProjetoUnivesp2021.HOSPITAIS_PLANOS_CONVENIOS definition

CREATE TABLE `HOSPITAIS_PLANOS_CONVENIOS`
(
    `CPID`      bigint(20) unsigned     NOT NULL,
    `HID`       int(10) unsigned        NOT NULL,
    UNIQUE (`CPID`, `HID`),
    FOREIGN KEY (`CPID`) REFERENCES `CONVENIO_PLANOS` (`CPID`),
    FOREIGN KEY (`HID`) REFERENCES `HOSPITAIS` (`HID`)
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
    FOREIGN KEY (`HID`) REFERENCES `HOSPITAIS` (`HID`),
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
    `UNASCIMENTO`  date                NOT NULL,
    `USEXO`        enum ('M', 'F')     NOT NULL,
    `UTELEFONE`    bigint(14) unsigned NOT NULL,
    `CPID`         bigint(20) unsigned NOT NULL,
    `UATIVADO`     enum ('T','F','D')  NOT NULL DEFAULT 'T',
    PRIMARY KEY (`UID`),
    UNIQUE KEY `UEMAIL` (`UEMAIL`),
    UNIQUE KEY `UTOKEN` (`UTOKEN`),
    UNIQUE KEY `UCPF` (`UCPF`),
    FOREIGN KEY (`CPID`) REFERENCES `CONVENIO_PLANOS` (`CPID`)
);

-- ProjetoUnivesp2021.HOSPITAIS_FAVORITOS definition

CREATE TABLE `HOSPITAIS_FAVORITOS` (
                                       `HID`          int(10) unsigned    NOT NULL,
                                       `UID`          bigint(20) unsigned NOT NULL,
                                       FOREIGN KEY (`HID`) REFERENCES `HOSPITAIS` (`HID`),
                                       FOREIGN KEY (`UID`) REFERENCES `USUARIOS` (`UID`),
                                       UNIQUE (`HID`, `UID`)
);


-- ProjetoUnivesp2021.DEPENDENTES definition

CREATE TABLE `DEPENDENTES`
(
    `DID`          bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `UID`          bigint(20) unsigned NOT NULL,
    `DNOME`        varchar(255)        NOT NULL,
    `DNASCIMENTO`   date                NOT NULL,
    `DSEXO`        enum ('M', 'F')     NOT NULL,
    `DATIVADO`     enum ('T','F')  NOT NULL DEFAULT 'T',
    PRIMARY KEY (`DID`),
    FOREIGN KEY (`UID`) REFERENCES `USUARIOS` (`UID`)
);

-- ProjetoUnivesp2021.AGENDAMENTOS definition

CREATE TABLE `AGENDAMENTOS`
(
    `AID`      bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `UID`      bigint(20) unsigned NOT NULL,
    `DID`      bigint(20) unsigned,
    `MID`      bigint(20) unsigned NOT NULL,
    `ADATA`    datetime            NOT NULL,
    `AATIVADO` enum ('T','D')      NOT NULL DEFAULT 'T',
    PRIMARY KEY (`AID`),
    FOREIGN KEY (`UID`) REFERENCES `USUARIOS` (`UID`),
    FOREIGN KEY (`MID`) REFERENCES `MEDICOS` (`MID`),
    FOREIGN KEY (`DID`) REFERENCES `DEPENDENTES` (`DID`),
    UNIQUE(`MID`, `ADATA`)
);

-- Get UID from TOKEN
CREATE FUNCTION getUIDFromToken(token char(36)) RETURNS bigint(20)
BEGIN
RETURN (SELECT UID FROM USUARIOS WHERE UTOKEN = token AND UATIVADO = 'T' LIMIT 1);
END;

-- Get Login Password from email

CREATE PROCEDURE GetLoginUsuario(email VARCHAR(100))
BEGIN
SELECT UPASSWORD FROM USUARIOS WHERE UATIVADO = 'T' AND UEMAIL = email LIMIT 1;
END;

-- Add a User

CREATE PROCEDURE ProjetoUnivesp2021.RegistrarUsuario(nome VARCHAR(255), email VARCHAR(100),
                                                     password VARCHAR(80), cpf varchar(11),
                                                     nascimento date, sexo enum('M', 'F'),
                                                     telefone bigint(14), pcpid bigint(20))

BEGIN
INSERT INTO USUARIOS (UNOME, UEMAIL, UPASSWORD, UCPF, UNASCIMENTO, USEXO, UTELEFONE, CPID)
VALUES (nome, email, password, cpf, nascimento, sexo, telefone, pcpid);
END;

-- Adiciona um convenio

CREATE PROCEDURE ProjetoUnivesp2021.RegistrarConvenio(nome varchar(255))
BEGIN
INSERT INTO `CONVENIO`(`CNOME`)
VALUES (nome);
END;

-- Adicionar um plano convenio

CREATE PROCEDURE ProjetoUnivesp2021.RegistarPlanoEmConvenio(pcid int(10), nome varchar(255))
BEGIN
INSERT INTO `CONVENIO_PLANOS`(`CID`, `CPNOME`)
VALUES (pcid, nome);
END;

-- Listar Convenios

CREATE PROCEDURE ProjetoUnivesp2021.ListarConvenios()
BEGIN
SELECT CID, CNOME FROM CONVENIO WHERE CATIVADO = 'T';
END;

-- Listar Planos de um Convenio

CREATE PROCEDURE ProjetoUnivesp2021.ListarPlanos(pcid int(10))
BEGIN
SELECT CPID, CPNOME, CID FROM CONVENIO_PLANOS WHERE CID = pcid;
END;

-- Adiciona hospital nos favoritos

CREATE PROCEDURE ProjetoUnivesp2021.FavoritarHospital(token char(36), phid int(10))
BEGIN
INSERT INTO HOSPITAIS_FAVORITOS (UID, HID)
VALUES (getUIDFromToken(token), phid);
END;

-- Lista hospitais nos favoritos

CREATE PROCEDURE ProjetoUnivesp2021.ListarHospitaisFavoritos(token char(36))
BEGIN
SELECT HID, HNOME, HUF, HCIDADE, HCEP, HENDERECO, HCOMPLEMENTO, HTELEFONE, HISPRONTOSOCORRO
FROM HOSPITAIS WHERE HATIVADO = 'T' AND HID IN (
    SELECT HID FROM HOSPITAIS_FAVORITOS WHERE UID = getUIDFromToken(token)
);
END;

-- Listar Agendamentos usuario

CREATE PROCEDURE ProjetoUnivesp2021.ListarAgendamentosUsuario(token char(36))
BEGIN
SELECT AID, ADATA, MID, DID
FROM AGENDAMENTOS WHERE UID = getUIDFromToken(token) AND AATIVADO = 'T'
ORDER BY ADATA DESC;
END;

-- Registar dependente

CREATE PROCEDURE ProjetoUnivesp2021.RegistrarDependete(token char(36), nome varchar(255),
                                                       nascimento date, sexo enum('M', 'F'))
BEGIN
INSERT INTO DEPENDENTES (UID, DNOME, DNASCIMENTO, DSEXO)
VALUES (getUIDFromToken(token), nome, nascimento, sexo);
END;

-- Listar dependetes

CREATE PROCEDURE ProjetoUnivesp2021.ListarDependentes(token char(36))
BEGIN
SELECT DID, DNOME, DNASCIMENTO, DSEXO
FROM DEPENDENTES WHERE UID = getUIDFromToken(token) AND DATIVADO = 'T';
END;

-- Desativar dependete

CREATE PROCEDURE ProjetoUnivesp2021.DesativarDepente(token char(36), pdid bigint(20))
BEGIN
UPDATE DEPENDENTES SET DATIVADO = 'F' WHERE UID = getUIDFromToken(token) AND DID = pdid;
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

-- Seta como NULL o token do usuario

CREATE PROCEDURE ProjetoUnivesp2021.LogOff(email varchar(100), token varchar(36))
BEGIN
UPDATE USUARIOS SET UTOKEN = NULL WHERE UEMAIL = email AND UTOKEN = token AND UATIVADO = 'T';
END;

-- Adiciona uma especialidade no banco

CREATE PROCEDURE ProjetoUnivesp2021.RegistrarEspecialidade(nome varchar(30))
BEGIN
INSERT INTO ESPECIALIDADES (ENOME) VALUES (nome);
END;

-- Lista especialidades

CREATE PROCEDURE ProjetoUnivesp2021.ListarEspecialidades()
BEGIN
SELECT EID, ENOME FROM ESPECIALIDADES;
END;

-- Cadastra Um hospital

CREATE PROCEDURE ProjetoUnivesp2021.RegistrarHospital(
    nome varchar(255), uf char(2), cidade varchar(50),
    cep char(8), endereco varchar(150), complemento varchar(150),
    telefone bigint(14), isProntoSocorro tinyint(1))
BEGIN
INSERT INTO HOSPITAIS (HNOME, HUF, HCIDADE, HCEP, HENDERECO, HCOMPLEMENTO, HTELEFONE, HISPRONTOSOCORRO)
VALUES (nome, uf, cidade, cep, endereco, complemento, telefone, isProntoSocorro);
END;

-- Cadastra um convenio em um hospital

CREATE PROCEDURE ProjetoUnivesp2021.RegistrarConvenioHospital(phid int(10), pcpid bigint(20))
BEGIN
INSERT INTO HOSPITAIS_PLANOS_CONVENIOS(CPID, HID)
VALUES (pcpid, phid);
END;

-- remove um convenio em um hospital

CREATE PROCEDURE ProjetoUnivesp2021.RemoverConvenioHospital(pcpid bigint(20), phid int(10))
BEGIN
DELETE FROM HOSPITAIS_PLANOS_CONVENIOS WHERE CPID = pcpid AND HID = phid LIMIT 1;
END;

-- Lista Hospitais

CREATE PROCEDURE ProjetoUnivesp2021.ListarHospitais()
BEGIN
SELECT HID, HNOME, HUF, HCIDADE, HCEP, HENDERECO, HCOMPLEMENTO, HTELEFONE, HISPRONTOSOCORRO
FROM HOSPITAIS WHERE HATIVADO = 'T';
END;

-- Listar Hospitais Por Planos de convenio

CREATE PROCEDURE ProjetoUnivesp2021.ListarHospitaisPorPlanoConvenio(pcpid bigint(20))
BEGIN
SELECT HID FROM HOSPITAIS_PLANOS_CONVENIOS WHERE CPID = pcpid;
END;

-- Lista Especialiadedes de um hospital

CREATE PROCEDURE ProjetoUnivesp2021.ListarEspecialidadesHospital(phid int(10))
BEGIN
SELECT EID FROM MEDICOS WHERE HID = phid GROUP BY EID;
END;

-- Registra medico em um hospital

CREATE PROCEDURE ProjetoUnivesp2021.RegistrarMedico(phid int(10), peid int(10), nome varchar(255))
BEGIN
INSERT INTO MEDICOS (`MNOME`, `EID`, `HID`) VALUES (nome, peid, phid);
END;

-- Lista Medicos por especialidade

CREATE PROCEDURE ProjetoUnivesp2021.ListarMedicosPorEspecialidade(peid int(10))
BEGIN
SELECT MID, HID, MNOME FROM MEDICOS WHERE EID = peid AND MATIVADO = 'T';
END;

-- Lista Agendamentos de um medico

CREATE PROCEDURE ProjetoUnivesp2021.ListarAgendamentosMedico(pmid bigint(20))
BEGIN
SELECT ADATA FROM AGENDAMENTOS WHERE MID = pmid AND AATIVADO = 'T' AND ADATA >= DATE(NOW());
END;

-- Registra um novo agendamento

CREATE PROCEDURE ProjetoUnivesp2021.RegistrarAgendamento(token char(36), pdid bigint(20), pmid bigint(20), padata datetime)
BEGIN
INSERT INTO AGENDAMENTOS (UID, DID, MID, ADATA)
VALUES (getUIDFromToken(token), pdid, pmid, padata);
END;
