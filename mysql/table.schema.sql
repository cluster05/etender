DROP DATABASE tenderDB;

CREATE DATABASE tenderDB;

use tenderDB;

CREATE TABLE auth(
    userId INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) NOT NULL,
    email VARCHAR(50) NOT NULL, 
    pwd VARCHAR(150) NOT NULL
);

CREATE TABLE division(
    divisionId INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE
);

CREATE TABLE ssg(
    ssgId INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    station VARCHAR(50) NOT NULL,
    sector VARCHAR(50) NOT NULL,
    pgroup VARCHAR(50) NOT NULL,
    uniquestream VARCHAR(150) NOT NULL UNIQUE,
    divisionId INT,
    constraint divisionId
        FOREIGN KEY (divisionId)
        REFERENCES division (divisionId)
        ON DELETE CASCADE
);



CREATE TABLE fre(
    freId INT  NOT NULL AUTO_INCREMENT PRIMARY KEY,
    flatno VARCHAR(50) NOT NULL,
    reserveprice VARCHAR(50) NOT NULL,
    emd VARCHAR(50) NOT NULL,
    uniquefre VARCHAR(50) NOT NULL UNIQUE,
    ssgId INT,
    constraint ssgId
        FOREIGN KEY (ssgId)
        REFERENCES ssg (ssgId)
        ON DELETE CASCADE
);

use tenderDB;

