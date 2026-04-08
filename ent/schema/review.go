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
		field.String("author_id"),
		field.String("target_user_id"),
		field.Time("created_at"),
	}
}

func (Review) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("card", Card.Type).Ref("reviews").Unique(),
		edge.To("order", Order.Type).Unique(),
	}
}
