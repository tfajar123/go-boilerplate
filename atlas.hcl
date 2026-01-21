env "local" {
  src = "ent://ent/schema"

  dev = getenv("LOCAL_DEV_DB")
  url = getenv("LOCAL_DB")

  migration {
    dir = "file://migrations"
  }
}

env "staging" {
  src = "ent://ent/schema"
  url = getenv("STAGING_DB")
}

env "production" {
  src = "ent://ent/schema"
  url = getenv("PROD_DB")
}