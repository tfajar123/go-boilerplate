-- Rename a column from "user_profiles" to "user_id"
ALTER TABLE "public"."profiles" RENAME COLUMN "user_profiles" TO "user_id";
-- Create index "profiles_user_id_key" to table: "profiles"
CREATE UNIQUE INDEX "profiles_user_id_key" ON "public"."profiles" ("user_id");
