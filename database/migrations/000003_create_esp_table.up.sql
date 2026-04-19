CREATE TABLE esps
(
    internal_id BIGSERIAL PRIMARY KEY,
    public_id UUID NOT NULL DEFAULT gen_random_uuid(),
    mac_address VARCHAR(126) NOT NULL UNIQUE,
    user_id INT,
    device_status VARCHAR(16) NOT NULL DEFAULT 'offline',
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ NULL,
    CONSTRAINT mac_address_unique UNIQUE(mac_address),
    CONSTRAINT esps_public_id_unique UNIQUE(public_id),
    CONSTRAINT esps_fk_user FOREIGN KEY (user_id) REFERENCES users(internal_id) ON DELETE SET NULL
);