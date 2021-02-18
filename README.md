
Пример запуска
go run cmd/api/main.go -user="имя пользователя" -p="пароль" --db="productdb"
, где
"имя пользователя" - логин на сервере PostgreSQL
"пароль" - пароль этого пользователя
"productdb" - имя базы данных

Клиент: http://localhost:8080/

На сервере PostgreSQL будет база данных productdb, в которой есть таблица Products, описываемая следующим скриптом: 

CREATE TABLE product(
  vendor_id bigint, 
  offer_id bigint,
  name varchar(255),
  price real,
  quantity integer,
  PRIMARY KEY(vendor_id, offer_id)
);