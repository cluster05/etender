CREATE DATABASE tenderDB;


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
    FOREIGN KEY (divisionId) REFERENCES division(divisionId)
);


CREATE TABLE fre(
    freId INT  NOT NULL AUTO_INCREMENT PRIMARY KEY,
    flatno VARCHAR(50) NOT NULL,
    reserveprice VARCHAR(50) NOT NULL,
    emd VARCHAR(50) NOT NULL,
    userId INT,
    ssgId INT,
    FOREIGN KEY (userId) REFERENCES ssg(userId)
    FOREIGN KEY (userId) REFERENCES ssg(userId)
);


