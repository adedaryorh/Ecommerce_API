CREATE TABLE "orders" (
                          "id" BIGSERIAL PRIMARY KEY,
                          "user_id" BIGINT NOT NULL,
                          "status" varchar(50) NOT NULL DEFAULT 'Pending',
                          "total_amount" numeric(10,2) NOT NULL DEFAULT 0,
                          "created_at" timestamptz NOT NULL DEFAULT NOW(),
                          "updated_at" timestamptz NOT NULL DEFAULT NOW()
);