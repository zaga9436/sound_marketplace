package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Message struct {
	ent.Schema
}

func (Message) Fields() []ent.Field {
	return []ent.Field{
		field.Text("body"),
		field.String("sender_id"),
	}
}

func (Message) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("chat_room", ChatRoom.Type).Ref("messages").Unique().Required(),
	}
}
