-- Create "profiles" table
CREATE TABLE "public"."profiles" (
  "id" uuid NOT NULL,
  "name" character varying NOT NULL,
  "image_url" character varying NOT NULL,
  "birth_date" character varying NOT NULL,
  "address" character varying NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "user_profiles" uuid NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "profiles_users_profiles" FOREIGN KEY ("user_profiles") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Drop "comments" table
DROP TABLE "public"."comments";
