CREATE TABLE kompors
(
    internal_id BIGSERIAL PRIMARY KEY,
    public_id UUID NOT NULL DEFAULT gen_random_uuid(),
    user_id INT,
    kompor_name VARCHAR(32) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ NULL,
    CONSTRAINT kompors_public_id_unique UNIQUE(public_id),
    CONSTRAINT kompors_fk_user_id FOREIGN KEY (user_id) REFERENCES users(internal_id) ON DELETE SET NULL
);