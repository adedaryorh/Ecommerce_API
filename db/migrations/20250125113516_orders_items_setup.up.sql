CREATE TABLE "order_items" (
                               "id" BIGSERIAL PRIMARY KEY,
                               "order_id" BIGINT NOT NULL,
                               "product_id" BIGINT NOT NULL,
                               "quantity" integer NOT NULL,
                               "price" numeric(10,2) NOT NULL,
                               "created_at" timestamptz NOT NULL DEFAULT NOW()
);
