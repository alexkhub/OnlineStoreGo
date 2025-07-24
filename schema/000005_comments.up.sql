create table comment( 
    id SERIAL primary key,
    title varchar(150) NOT NULL, 
    message text NULL, 
    rating int  NOT NULL DEFAULT 5 check(rating >= 0 or rating <= 10), 
    user_id int NOT NULL, 
    product int  NOT NULL REFERENCES product(id) on delete CASCADE,
    create_at timestamptz  DEFAULT now()
);