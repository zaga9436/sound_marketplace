package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Bid struct {
	ent.Schema
}

func (Bid) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("price"),
		field.Text("message"),
	}
}

func (Bid) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("request", Card.Type).Ref("bids").Unique().Required(),
		edge.From("engineer", User.Type).Ref("sent_bids").Unique().Required(),
		edge.To("order", Order.Type).Unique(),
	}
}
