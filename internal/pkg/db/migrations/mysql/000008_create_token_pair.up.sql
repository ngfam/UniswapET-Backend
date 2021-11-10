create table if not exists Pair(
    ID int not null unique auto_increment,
    token0 varchar (20) not null,
    token1 varchar (20) not null,
    balance0 float(30, 10) not null,
    balance1 float(30, 10) not null,
    marketCap float(30, 10) not null,
    totalVolumeRecorded float(30, 10) not null,
    foreign key (token0) references Token(ID) on delete cascade,
    foreign key (token1) references Token(ID) on delete cascade,
    unique (token0, token1),
    primary key (ID)
);
