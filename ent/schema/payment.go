package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Payment struct {
	ent.Schema
}

func (Payment) Fields() []ent.Field {
	return []ent.Field{
		field.String("provider"),
		field.String("external_id").Unique(),
		field.Int64("amount"),
		field.String("status"),
		field.String("redirect_url"),
		field.Text("callback_data").Optional(),
	}
}

func (Payment) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("transaction", Transaction.Type).Ref("payment").Unique(),
	}
}
