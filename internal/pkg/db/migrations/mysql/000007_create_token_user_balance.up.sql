create table if not exists UserBalance(
    userID int not null,
    tokenID varchar(20) not null,
    balance float(30, 10) not null,
    foreign key (userID) references User(ID) on delete cascade,
    foreign key (tokenID) references Token(ID) on delete cascade,
    Unique (userID, tokenID)
);
