package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Transaction struct {
	ent.Schema
}

func (Transaction) Fields() []ent.Field {
	return []ent.Field{
		field.Enum("type").Values("deposit", "hold", "release", "refund", "partial_refund"),
		field.Int64("amount"),
		field.String("external_id").Optional(),
		field.Time("created_at"),
	}
}

func (Transaction) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("order", Order.Type).Ref("transactions").Unique(),
		edge.To("payment", Payment.Type).Unique(),
	}
}
