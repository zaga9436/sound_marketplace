package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type User struct {
	ent.Schema
}

func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("email").Unique(),
		field.String("password_hash"),
		field.Enum("role").Values("customer", "engineer", "admin"),
		field.Time("created_at"),
	}
}

func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("profile", Profile.Type).Unique(),
		edge.To("cards", Card.Type),
		edge.To("sent_bids", Bid.Type),
		edge.To("customer_orders", Order.Type),
		edge.To("engineer_orders", Order.Type),
	}
}
