create table category (
    id SERIAL primary key,
    name varchar(150),
    UNIQUE(name)
);

create table product (
    id SERIAL primary key,
    name varchar(150),
    first_price int, 
    description text NULL,
    discount int Default 0,
    price int ,
    category int REFERENCES category(id)  on delete set NULL,
    UNIQUE(name)

);
create table image (
    id SERIAL primary key,
    image_uuid varchar(250) NULL
);

create table product_image(
    id SERIAL primary key,
    product int REFERENCES product(id),
    image int REFERENCES image(id) ON DELETE CASCADE
);

CREATE OR REPLACE FUNCTION price_calculation()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.discount >= 0 and NEW.discount <= 100 THEN
        NEW.price = NEW.first_price - (NEW.first_price * NEW.discount / 100);
    ELSE 
        RAISE EXCEPTION 'discount is not in range from 0 to 100';
    
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER product_update
BEFORE UPDATE OR INSERT ON product
FOR EACH ROW
EXECUTE FUNCTION price_calculation();

