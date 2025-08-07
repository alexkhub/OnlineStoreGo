create table verify_order_email(
    id uuid primary key DEFAULT gen_random_uuid(),
    order_id int,
    user_id int, 
    datetime_create timestamptz DEFAULT now()

);

create table confirm_order_email(
    id SERIAL primary key,
    order_id int,
    confirm_code int DEFAULT FLOOR(random() * (999999 - 100000 + 1) + 100000)::int,
    datetime_create timestamptz DEFAULT now()

)