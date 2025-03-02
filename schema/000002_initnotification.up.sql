create table verifyemail(
    id SERIAL primary key,
    user_id int ,
    verify_uuid uuid,
    datetime_create timestamptz DEFAULT now() ,
    UNIQUE(verify_uuid)
)