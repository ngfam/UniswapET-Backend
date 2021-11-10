create table if not exists User(
    ID int not null unique auto_increment,
    username varchar(127) not null unique,
    hashedPassword varchar(127) not null,
    primary key (ID)
);