CREATE OR REPLACE FUNCTION del_cart_point()
RETURNS TRIGGER AS $$

BEGIN 
    IF NEW.amount <= 0 THEN
        delete from cart where id = NEW.id;
        RETURN NULL;
    END IF;
    RETURN NEW;

END;
$$ LANGUAGE plpgsql;


CREATE OR REPLACE TRIGGER trigger_del_cart_proint
AFTER UPDATE ON cart
FOR EACH ROW
EXECUTE FUNCTION del_cart_point();

CREATE OR REPLACE FUNCTION del_order_point()
RETURNS TRIGGER AS $$
BEGIN
    delete from order_point where id in (select order_point from order_order_point where user_order=OLD.id);
    RETURN OLD;
END;
$$ LANGUAGE plpgsql;


CREATE OR REPLACE TRIGGER trigger_del_order_point
BEFORE DELETE ON user_order
FOR EACH ROW
EXECUTE FUNCTION del_order_point();