package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Dispute struct {
	ent.Schema
}

func (Dispute) Fields() []ent.Field {
	return []ent.Field{
		field.Text("reason"),
		field.String("status").Default("open"),
	}
}

func (Dispute) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("order", Order.Type).Ref("disputes").Unique().Required(),
	}
}
