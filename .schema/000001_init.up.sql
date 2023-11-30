CREATE TABLE Message
(
    MessageId         serial not null unique,
    OrderUid          varchar(255),
    TrackNumber       varchar(255),
    Entry             varchar(255),
    Locale            varchar(255),
    InternalSignature varchar(255),
    CustomerId        varchar(255),
    DeliveryService   varchar(255),
    Shardkey          varchar(255),
    SmId              int,
    DateCreated       time,
    OofShard          varchar(255)
);

CREATE TABLE delivery
(
    DeliveryId serial not null unique,
    Name       varchar(255),
    Phone      varchar(255),
    Zip        varchar(255),
    City       varchar(255),
    Address    varchar(255),
    Region     varchar(255),
    Email      varchar(255)
);

CREATE TABLE payment
(
    PaymentId    serial not null unique,
    Transaction  varchar(255),
    RequestId    varchar(255),
    Currency     varchar(255),
    Provider     varchar(255),
    Amount       int,
    PaymentDt    int,
    Bank         varchar(255),
    DeliveryCost int,
    GoodsTotal   int,
    CustomFee    int
);

CREATE TABLE item
(
    ItemId      serial not null unique,
    ChrtId      int,
    TrackNumber varchar(255),
    Price       int,
    Rid         varchar(255),
    Name        varchar(255),
    Sale        int,
    Size        varchar(255),
    TotalPrice  int,
    NmId        int,
    Brand       varchar(255),
    Status      int
);

CREATE TABLE MessageDeliveries
(
    id         serial                                                   not null unique,
    MessageId  int references Message (MessageId) on delete cascade     not null,
    DeliveryId int references Deliveries (DeliveryId) on delete cascade not null
);

CREATE TABLE MessagePayments
(
    id        serial                                                not null unique,
    MessageId int references Message (MessageId) on delete cascade  not null,
    PaymentId int references Payments (PaymentId) on delete cascade not null
);

CREATE TABLE MessageItem
(
    id        serial                                               not null unique,
    MessageId int references Message (MessageId) on delete cascade not null,
    ItemId    int references Item (ItemId) on delete cascade       not null
);
