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

CREATE OR REPLACE FUNCTION update_order_price()
RETURNS TRIGGER AS $$
DECLARE
    target_order_id INT;
BEGIN
    
    IF TG_TABLE_NAME = 'order_order_point' THEN
        target_order_id := COALESCE(NEW.user_order, OLD.user_order);

    ELSIF TG_TABLE_NAME = 'order_point' THEN
        FOR target_order_id IN
            SELECT user_order
            FROM order_order_point
            WHERE order_point = COALESCE(NEW.id, OLD.id)
        LOOP
            UPDATE user_order u
            SET full_price = (
                SELECT COALESCE(SUM(op.product_price * op.amount), 0)
                FROM order_point op
                JOIN order_order_point oop ON oop.order_point = op.id
                WHERE oop.user_order = u.id
            )
            WHERE u.id = target_order_id;
        END LOOP;
        RETURN NULL;
    END IF;

    UPDATE user_order u
    SET full_price = (
        SELECT COALESCE(SUM(op.product_price * op.amount), 0)
        FROM order_point op
        JOIN order_order_point oop ON oop.order_point = op.id
        WHERE oop.user_order = target_order_id
    )
    WHERE u.id = target_order_id;

    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER trg_update_order_price_link
AFTER INSERT OR UPDATE OR DELETE ON order_order_point
FOR EACH ROW
EXECUTE FUNCTION update_order_price();


CREATE OR REPLACE TRIGGER trg_update_order_price_point
AFTER INSERT OR UPDATE OR DELETE ON order_point
FOR EACH ROW
EXECUTE FUNCTION update_order_price();