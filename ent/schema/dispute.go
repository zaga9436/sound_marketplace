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
		field.String("opened_by"),
		field.Text("reason"),
		field.String("status").Default("open"),
		field.Text("resolution").Default(""),
		field.Time("created_at"),
		field.Time("resolved_at").Optional().Nillable(),
	}
}

func (Dispute) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("order", Order.Type).Ref("disputes").Unique().Required(),
	}
}
