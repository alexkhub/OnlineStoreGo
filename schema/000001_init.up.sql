Create Table Roles(
    id SERIAL primary key,
    role_name varchar(100),
    UNIQUE(role_name)
);

Insert into Roles (role_name) VALUES ('client'), ('seller'), ('carrier'), ('admin');

Create Table Users(
    id SERIAL primary key,
    username varchar(150) UNIQUE,
    first_name varchar(30) NULL,
    last_name varchar(50) NULL,
    email varchar(150) UNIQUE,
    hash_password varchar(150),
    activate BOOLEAN Default FALSE, 
    block BOOLEAN Default FALSE, 
    role_id int REFERENCES Roles(id) DEFAULT 1,
    datetime_create timestamptz  DEFAULT now(),
    image varchar(250) NULL

);

create table Refresh(
    id SERIAL primary key ,
    user_id  int REFERENCES Users(id),
    refresh_token varchar(250),
    expiration_time timestamptz DEFAULT now() 
);
