package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Booth struct {
	ent.Schema
}

func (Booth) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").Unique(),
	}
}

func (Booth) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).Ref("booth").Unique().Required(),
		edge.To("products", Product.Type),
		edge.To("transactions", Transaction.Type),
	}
}
