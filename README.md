# peer2peer


#### Create DB and connect to it using PSQL
```
createdb peer2peer
psql peer2peer
```



Note: For my MAC starting postgresql
```
postgres -D db_lbapp
```

DB Queries
```
create table if not exists Visitor(id SERIAL PRIMARY KEY, firstname varchar(50) not null, lastname varchar(50), age int not null, gender varchar(1) not null, email varchar(100) not null, phonenumber varchar(15), university varchar(100), creationtime timestamp default current_timestamp);
```

Insert Visitor Sample
```
CURL -X POST -d '{"firstname":"Sau", "lastname":"Josh", "age":"22", "gender":"M", "email":"sau@vj.com", "phonenumber":"9192222222", "university":"NCSU"}' http://localhost:8000/v1/visitor
```

Get a visitor using REST call
```
CURL -X GET http://localhost:8000/v1/visitor/1
```
