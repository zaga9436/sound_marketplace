package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
)

type ChatRoom struct {
	ent.Schema
}

func (ChatRoom) Fields() []ent.Field {
	return nil
}

func (ChatRoom) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("order", Order.Type).Ref("chat_room").Unique().Required(),
		edge.To("messages", Message.Type),
	}
}
