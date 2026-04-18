ALTER TABLE esps 
DROP CONSTRAINT mac_address_unique,
DROP CONSTRAINT esps_public_id_unique,
DROP CONSTRAINT esps_fk_user,
DROP CONSTRAINT esps_fk_kompor;
DROP TABLE esps;