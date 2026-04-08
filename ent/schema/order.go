package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Order struct {
	ent.Schema
}

func (Order) Fields() []ent.Field {
	return []ent.Field{
		field.Enum("status").Values("created", "on_hold", "in_progress", "review", "completed", "dispute", "cancelled"),
		field.Int64("amount"),
		field.Text("delivery_notes").Optional(),
		field.Text("dispute_reason").Optional(),
	}
}

func (Order) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("card", Card.Type).Ref("orders").Unique(),
		edge.From("bid", Bid.Type).Ref("order").Unique(),
		edge.From("customer", User.Type).Ref("customer_orders").Unique().Required(),
		edge.From("engineer", User.Type).Ref("engineer_orders").Unique().Required(),
		edge.To("transactions", Transaction.Type),
		edge.To("chat_room", ChatRoom.Type).Unique(),
		edge.To("disputes", Dispute.Type),
	}
}
