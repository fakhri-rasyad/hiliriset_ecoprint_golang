ALTER TABLE fabric_types
DROP CONSTRAINT fabric_types_public_id_unique
,
DROP CONSTRAINT fabric_types_name_unique;

DROP TABLE fabric_types;

DROP TYPE fabric_type_enum;
