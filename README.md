Проект отправки логов(журналирование)
======================

loging настройка и запуск
-----------------------------

1. Необходимо перейти в каталог проекта:
`cd ./путь_до_проекта/loging_module`

2. Потом нужно собрать Docker с помощью команды:
`sudo docker build -t loging_module .`

3. Запустить проект:
`sudo docker run -p -d 50051:50051 loging_module`  

api_getwey настройка и запуск
-----------------------------

1. Необходимо перейти в каталог проекта:
`cd ./путь_до_проекта/api_gateway`

2. Потом нужно собрать Docker с помощью команды:
`sudo docker build -t api_gateway .`

3. Запустить проект внутри контейнера:
`sudo docker run -p 8080:8080 api_gateway`  

