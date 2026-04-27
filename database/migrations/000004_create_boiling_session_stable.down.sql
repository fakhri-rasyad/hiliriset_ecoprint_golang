ALTER TABLE boiling_sessions
DROP CONSTRAINT boiling_session_public_id_unique,
DROP CONSTRAINT boiling_session_fk_user_id,
DROP CONSTRAINT boiling_session_fk_kompor_id,
DROP CONSTRAINT boiling_session_fk_esp_id;
DROP TABLE boiling_sessions;

DROP TYPE boiling_status_enum;
DROP TYPE fabric_type_enum;
