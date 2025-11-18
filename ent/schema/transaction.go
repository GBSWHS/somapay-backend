package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"time"
)

type Transaction struct {
	ent.Schema
}

func (Transaction) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("quantity"),
		field.Int64("amount"),
		field.String("status"),
		field.Time("timestamp").Default(time.Now),
	}
}

func (Transaction) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("transactions").
			Unique().
			Required(),

		edge.From("booth", Booth.Type).
			Ref("transactions").
			Unique().
			Required(),

		edge.From("product", Product.Type).
			Ref("transactions").
			Unique().
			Required(),
	}
}
