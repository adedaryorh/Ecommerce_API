CREATE TABLE "sessions" (
                            "id" BIGSERIAL PRIMARY KEY,
                            "user_id" BIGINT NOT NULL,
                            "token" text NOT NULL,
                            "created_at" timestamptz NOT NULL DEFAULT NOW(),
                            "expires_at" timestamptz NOT NULL
);