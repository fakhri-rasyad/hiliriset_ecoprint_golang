ALTER TABLE session_records
DROP CONSTRAINT session_records_public_id_unique
,
DROP CONSTRAINT session_records_fk_session_id;

DROP TABLE session_records;
