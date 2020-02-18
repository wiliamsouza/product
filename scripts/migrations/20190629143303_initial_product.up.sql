CREATE TABLE "product" (
    "id" uuid   NOT NULL,
    "title" varchar   NOT NULL,
    "description" text   NOT NULL,
    "price_in_cents" integer   NOT NULL
);

ALTER TABLE ONLY "product" ADD CONSTRAINT "pk_product" PRIMARY KEY (id);
