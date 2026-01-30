-- Create "comments" table
CREATE TABLE "public"."comments" (
  "id" uuid NOT NULL,
  "contents" character varying NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "user_comments" uuid NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "comments_users_comments" FOREIGN KEY ("user_comments") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
