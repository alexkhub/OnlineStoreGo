create table comment( 
    id SERIAL primary key,
    title varchar(150) NOT NULL, 
    message text NULL, 
    raiting int  NOT NULL DEFAULT 5 check(raiting >= 0 or raiting <= 10), 
    user_id int NOT NULL, 
    product int  NOT NULL REFERENCES product(id) on delete CASCADE,
    create_at timestamptz  DEFAULT now()
);