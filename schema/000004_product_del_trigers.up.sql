CREATE OR REPLACE FUNCTION del_images()
RETURNS TRIGGER AS $$
BEGIN
    delete from image where id in (select image from product_image where product=OLD.id);
    RETURN OLD;
END;
$$ LANGUAGE plpgsql;


CREATE  OR REPLACE TRIGGER trigger_del_images
BEFORE DELETE ON product
FOR EACH ROW
EXECUTE FUNCTION del_images();