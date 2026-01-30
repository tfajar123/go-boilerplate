package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

type Comments struct {
	ent.Schema
}

func (Comments) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Unique().Default(uuid.New),
		field.String("contents"),
		field.Time("created_at").Default(time.Now()),
		field.Time("updated_at").Default(time.Now()),
	}
}

func (Comments) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("comments").
			Unique().
			Required(),
	}
}
