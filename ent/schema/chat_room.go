package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type ChatRoom struct {
	ent.Schema
}

func (ChatRoom) Fields() []ent.Field {
	return []ent.Field{
		field.Time("created_at"),
	}
}

func (ChatRoom) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("order", Order.Type).Ref("chat_room").Unique().Required(),
		edge.To("messages", Message.Type),
	}
}
