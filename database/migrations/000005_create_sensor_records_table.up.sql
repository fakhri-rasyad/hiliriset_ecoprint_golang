CREATE TABLE session_records
(
  internal_id BIGSERIAL PRIMARY KEY,
  public_id UUID NOT NULL DEFAULT gen_random_uuid(),
  session_id BIGINT NOT NULL,
  air_temp FLOAT NOT NULL,
  water_temp FLOAT NOT NULL,
  humidity FLOAT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMPTZ NULL,

  CONSTRAINT    session_records_public_id_unique   UNIQUE(public_id),
  CONSTRAINT    session_records_fk_session_id      FOREIGN KEY (session_id)  REFERENCES boiling_sessions(internal_id) ON DELETE SET NULL
);
