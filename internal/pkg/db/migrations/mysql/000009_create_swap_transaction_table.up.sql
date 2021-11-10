create table if not exists SwapTransaction(
    ID int not null unique auto_increment,
    userID int not null,
    inToken varchar(20) not null,
    outToken varchar(20) not null,
    inAmount float(30, 10) not null,
    outAmount float(30, 10) not null,
    volumeContributed float(30, 10) not null,
    foreign key (inToken) references Token(ID) on delete cascade,
    foreign key (outToken) references Token(ID) on delete cascade,
    foreign key (userID) references User(ID) on delete cascade,
    primary key (ID)
);