package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Profile struct {
	ent.Schema
}

func (Profile) Fields() []ent.Field {
	return []ent.Field{
		field.String("display_name"),
		field.Text("bio").Optional(),
		field.Float("rating").Default(0),
		field.Time("created_at"),
		field.Time("updated_at"),
	}
}

func (Profile) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).Ref("profile").Unique().Required(),
	}
}
