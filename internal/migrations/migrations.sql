DROP TABLE IF EXISTS cars, infoAPI;

CREATE TABLE IF NOT EXISTS cars (
  id SERIAL PRIMARY KEY, 
  regNum  VARCHAR (150) NOT NULL UNIQUE,
  mark VARCHAR (150) NOT NULL,
  model VARCHAR (150) NOT NULL,
  year INT,
  name VARCHAR (150) NOT NULL,
  surname VARCHAR (150) NOT NULL,
  patronymic VARCHAR (150));

  --owner_id INT NOT NULL REFERENCES owners ON DELETE CASCADE);

/*
Денормализация выбрана для удобства заполнения базы при запросе во внешний ресурс.
Для нормализации потребуется уникальный столбец у владельца - ФИО может совпадать.

CREATE TABLE IF NOT EXISTS owners (
  id SERIAL PRIMARY KEY,
  name VARCHAR (150) NOT NULL,
  surname VARCHAR (150) NOT NULL,
  patronymic VARCHAR (150));
  */

-- For testing api
CREATE TABLE infoAPI (
  id SERIAL PRIMARY KEY, 
  regNum  VARCHAR (150) NOT NULL,
  mark VARCHAR (150) NOT NULL,
  model VARCHAR (150) NOT NULL,
  year INT,
  name VARCHAR (150) NOT NULL,
  surname VARCHAR (150) NOT NULL,
  patronymic VARCHAR (150)
);


INSERT INTO infoAPI VALUES (0, 'X123XX150', 'Lada', 'Kalina', 2016, 'Viktor', 'Mikhaylov');
INSERT INTO infoAPI VALUES (1, 'P223XX150', 'Lada', 'Vesta', 2015, 'Andrew', 'Denn');
INSERT INTO infoAPI VALUES (2, 'K003HT127', 'Opel', 'Astra', 2005, 'Willy', 'Dyson');
INSERT INTO infoAPI VALUES (3, 'O853HT177', 'Opel', 'Corsa', 2000, 'John', 'Doe');
INSERT INTO infoAPI VALUES (4, 'H903MT137', 'Opel', 'Mokka', 2012, 'Daniel', 'Kraig');
INSERT INTO infoAPI VALUES (5, 'C563AT056', 'BMW', 'X3', 2009, 'Sanchez', 'Pereiro'); 

INSERT INTO cars (regnum, mark, model, year, name, surname) VALUES ('C513AT056', 'BMW', 'X5', 2009, 'Andy', 'Pereiro'); 
INSERT INTO cars (regnum, mark, model, year, name, surname) VALUES ('X513AT056', 'BMW', 'X5', 2019, 'Dann', 'Donn'); 
INSERT INTO cars (regnum, mark, model, year, name, surname) VALUES ('X513AR056', 'Opel', 'Test', 2000, 'Test', 'Test'); 
INSERT INTO cars (regnum, mark, model, year, name, surname) VALUES ('T213AR056', 'Volga', '3110', 2000, 'Nicolas', 'Tesla'); 


