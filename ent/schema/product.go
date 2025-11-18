package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Product struct {
	ent.Schema
}

func (Product) Fields() []ent.Field {
	return []ent.Field{
		field.String("name"),
		field.String("description").Optional(),
		field.Int64("price"),
		field.Int64("quantity"),
	}
}

func (Product) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("booth", Booth.Type).
			Ref("products").
			Required(),
		edge.To("transactions", Transaction.Type),
	}
}
