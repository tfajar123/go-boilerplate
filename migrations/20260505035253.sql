-- Modify "users" table
ALTER TABLE "public"."users" ADD COLUMN "email_verified" boolean NOT NULL DEFAULT false;
