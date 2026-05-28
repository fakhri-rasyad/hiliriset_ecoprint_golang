ALTER TABLE boiling_sessions
DROP CONSTRAINT boiling_session_public_id_unique,
DROP CONSTRAINT boiling_session_fk_user_id,
DROP CONSTRAINT boiling_session_fk_kompor_id,
DROP CONSTRAINT boiling_session_fk_esp_id,
DROP CONSTRAINT boiling_session_fk_fabric_type_id;

DROP TABLE boiling_sessions;

DROP TYPE boiling_status_enum;
