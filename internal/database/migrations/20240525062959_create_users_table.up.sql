CREATE TABLE IF NOT EXISTS users (
    id varchar PRIMARY KEY,
    created_at bigint NOT NULL,
    username varchar NOT NULL,
    email varchar NOT NULL,
    password varchar NOT NULL,
    role varchar NOT NULL,

    CONSTRAINT username_unique UNIQUE (username),
    CONSTRAINT email_role_unique UNIQUE (email, role),
    CONSTRAINT role_check CHECK (role IN ('customer', 'admin'))
)
