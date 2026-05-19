env "local" {
  src = "ent://ent/schema"

  dev = getenv("LOCAL_DEV_DB")
  url = getenv("DATABASE_URL")

  migration {
    dir = "file://migrations"
  }
}

env "staging" {
  src = "ent://ent/schema"
  url = getenv("DATABASE_URL")
}

env "production" {
  src = "ent://ent/schema"
  url = getenv("DATABASE_URL")
}