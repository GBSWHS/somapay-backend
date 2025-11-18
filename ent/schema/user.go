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
		field.String("username").Unique(),
		field.String("password"),
		field.Int64("point"),
		field.String("pin"),
		field.String("role").Default("USER"),
	}
}

func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("booth", Booth.Type).Unique(),
		edge.To("transactions", Transaction.Type),
		edge.To("charge_requests", ChargeRequest.Type),
	}
}
