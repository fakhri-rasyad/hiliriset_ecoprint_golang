ALTER TABLE kompors 
DROP CONSTRAINT kompors_public_id_unique,
DROP CONSTRAINT kompors_fk_user_id;
DROP TABLE kompors;