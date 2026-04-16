CREATE TABLE kompors
(
    internal_id BIGSERIAL PRIMARY KEY,
    public_id UUID NOT NULL DEFAULT gen_random_uuid(),
    user_id INT,
    CONSTRAINT kompors_public_id_unique UNIQUE(public_id),
    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(internal_id) ON DELETE SET NULL
)