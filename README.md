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
