package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type MediaFile struct {
	ent.Schema
}

func (MediaFile) Fields() []ent.Field {
	return []ent.Field{
		field.String("uploaded_by"),
		field.String("storage_key"),
		field.String("mime_type"),
		field.String("visibility").Default("private"),
		field.String("purpose"),
		field.Bool("is_processed").Default(false),
		field.Time("created_at"),
	}
}

func (MediaFile) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("card", Card.Type).Ref("media_files").Unique(),
		edge.To("order", Order.Type).Unique(),
	}
}
