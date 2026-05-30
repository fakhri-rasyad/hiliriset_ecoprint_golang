CREATE TYPE fabric_type_enum AS ENUM
('sample','katun', 'mori', 'linen', 'sutra');

CREATE TABLE fabric_types
(
  internal_id BIGSERIAL PRIMARY KEY,
  public_id UUID NOT NULL DEFAULT gen_random_uuid(),
  name fabric_type_enum NOT NULL,
  boiling_minutes INT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

  CONSTRAINT fabric_types_public_id_unique UNIQUE (public_id),
  CONSTRAINT fabric_types_name_unique      UNIQUE (name)
);

-- Seed default boiling times
INSERT INTO fabric_types
  (name, boiling_minutes)
VALUES
  ('sample', 3),
  ('katun', 120),
  ('mori', 90),
  ('linen', 30),
  ('sutra', 60);
