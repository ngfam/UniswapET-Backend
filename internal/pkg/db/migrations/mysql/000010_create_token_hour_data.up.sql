create table if not exists TokenHourData(
    token varchar(20) not null, 
    hourId int not null,
    tokenPrice float(30, 10) not null,
    foreign key (token) references Token(ID) on delete cascade,
    unique(token, hourId)
);