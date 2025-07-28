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