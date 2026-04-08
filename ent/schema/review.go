package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Review struct {
	ent.Schema
}

func (Review) Fields() []ent.Field {
	return []ent.Field{
		field.Int("rating"),
		field.Text("comment"),
	}
}

func (Review) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("card", Card.Type).Ref("reviews").Unique(),
	}
}
