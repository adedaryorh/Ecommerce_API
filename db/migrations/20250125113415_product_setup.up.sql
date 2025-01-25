CREATE TABLE "products" (
                            "id" BIGSERIAL PRIMARY KEY,
                            "name" varchar(256) NOT NULL,
                            "description" text,
                            "price" numeric(10,2) NOT NULL,
                            "stock" integer NOT NULL DEFAULT 0,
                            "created_at" timestamptz NOT NULL DEFAULT NOW(),
                            "updated_at" timestamptz NOT NULL DEFAULT NOW()
);