package schema

import (
    "entgo.io/ent"
    "entgo.io/ent/schema/edge"
    "entgo.io/ent/schema/field"
)

type ChargeRequest struct {
    ent.Schema
}

func (ChargeRequest) Fields() []ent.Field {
    return []ent.Field{
        field.Int64("amount"),
        field.String("status").Default("PENDING"),
    }
}

func (ChargeRequest) Edges() []ent.Edge {
    return []ent.Edge{
        edge.To("user", User.Type).
            Unique().
            Required(),
    }
}
