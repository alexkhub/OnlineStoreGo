create table cart (
    id SERIAL primary key,
    product_id int, 
    user_id int,
    amount int DEFAULT 1, 
    CONSTRAINT product_user UNIQUE (product_id, user_id)  
);

create  table payment_method (
    id SERIAL primary key,
    Name varchar(50),
    Description text
);

create table user_order ( 
    id SERIAL primary key,
    user_id int,
    full_price int NULL,
    payment_method int NULL REFERENCES payment_method(id) ON DELETE SET NULL,
    status varchar(50) DEFAULT 'Collect',
    delivery_method varchar(50) DEFAULT 'Delivery',
    address varchar(200),
    create_at timestamptz DEFAULT now(),
    delivery_date timestamptz NULL, 
    CHECK (status IN ('Collect', 'Ready to ship', 'At the pick-up point', 'Delivered', 'Canceled')),
    CHECK (delivery_method IN ('Delivery', 'Self pickup'))
);

create table order_point (
    id SERIAL primary key,
    product_id int, 
    product_price int,
    amount int
);

create table order_order_point (
    user_order int REFERENCES user_order(id) ON DELETE CASCADE,
    order_point int REFERENCES order_point(id) ON DELETE CASCADE,
    CONSTRAINT id primary key (user_order, order_point) 
);



