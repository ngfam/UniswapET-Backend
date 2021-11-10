create table if not exists Token(
    ID varchar(20) not null unique,
    name varchar(10) not null unique,
    totalSupply float(30, 10) not null,
    iconURL varchar(1000),
    price float(30, 10),
    primary key (ID)
);