CREATE TYPE boiling_status_enum AS ENUM
('boiling', 'finished', 'cancelled');
CREATE TYPE fabric_type_enum AS ENUM
('katun', 'polyester', 'linen', 'sutra');

CREATE TABLE boiling_sessions
(
    internal_id     BIGSERIAL               PRIMARY KEY,
    public_id       UUID                    NOT NULL DEFAULT gen_random_uuid(),
    boiling_status  boiling_status_enum     NOT NULL DEFAULT 'boiling',
    fabric_type     fabric_type_enum        NOT NULL DEFAULT 'katun',
    user_id         BIGINT,
    kompor_id       BIGINT,
    esp_id          BIGINT,
    created_at      TIMESTAMPTZ             NOT NULL DEFAULT CURRENT_TIMESTAMP,
    finished_at     TIMESTAMPTZ             NULL,
    updated_at      TIMESTAMPTZ             NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at      TIMESTAMPTZ             NULL,

    CONSTRAINT boiling_session_public_id_unique UNIQUE (public_id),
    CONSTRAINT boiling_session_fk_user_id   FOREIGN KEY (user_id)   REFERENCES users(internal_id)   ON DELETE SET NULL,
    CONSTRAINT boiling_session_fk_kompor_id FOREIGN KEY (kompor_id) REFERENCES kompors(internal_id) ON DELETE SET NULL,
    CONSTRAINT boiling_session_fk_esp_id    FOREIGN KEY (esp_id)    REFERENCES esps(internal_id)    ON DELETE SET NULL
);