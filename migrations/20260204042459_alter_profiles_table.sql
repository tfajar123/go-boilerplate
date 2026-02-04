-- Modify "profiles" table
ALTER TABLE "public"."profiles" ALTER COLUMN "image_url" DROP NOT NULL, ALTER COLUMN "birth_date" DROP NOT NULL, ALTER COLUMN "address" DROP NOT NULL;
