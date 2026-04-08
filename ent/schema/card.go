package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Card struct {
	ent.Schema
}

func (Card) Fields() []ent.Field {
	return []ent.Field{
		field.Enum("card_type").Values("offer", "request"),
		field.Enum("kind").Values("product", "service"),
		field.String("title"),
		field.Text("description"),
		field.Int64("price"),
		field.Strings("tags").Optional(),
		field.Bool("is_published").Default(true),
	}
}

func (Card) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("author", User.Type).Ref("cards").Unique().Required(),
		edge.To("bids", Bid.Type),
		edge.To("orders", Order.Type),
		edge.To("media_files", MediaFile.Type),
		edge.To("reviews", Review.Type),
	}
}
