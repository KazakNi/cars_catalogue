# Тестовое задание для Effective Mobile 
## Cars catalogue

ܸПриложение представляет собой сервис для управления каталогом автомобилей.

Для разворачивания сервиса локально необходимо клонировать репозиторий:
```
git clone https://github.com/KazakNi/cars_catalogue.git
```
Перейти в папку

```
cd cars_catalogue
```
Заполнить переменные окружения в **example.env** в корне папки для Докера и в папке cmd для локального запуска через консоль

Далее: 
Для запуска проекта в Докере:

```
docker-compose up
```

Для запуска локально:

```
go run cmd/main.go -devmode=true
```

Спецификация OpenAPI доступна по localhost:8080/redoc

<div align="center">
  <img src="https://github.com/KazakNi/cars_catalogue/blob/main/redoc.jpg" align="center"> </img>
  </div>
