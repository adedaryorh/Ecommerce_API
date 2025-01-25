CREATE TABLE "users" (
                         "id" BIGSERIAL PRIMARY KEY,
                         "email" varchar(256) UNIQUE NOT NULL,
                         "hashed_password" varchar(256) NOT NULL,
                         "username" varchar(256) UNIQUE NOT NULL,
                         "role" TEXT NOT NULL DEFAULT 'user',
                         "created_at" timestamptz NOT NULL DEFAULT NOW(),
                         "updated_at" timestamptz NOT NULL DEFAULT NOW()
);