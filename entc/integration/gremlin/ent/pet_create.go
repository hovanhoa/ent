// Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"

	"github.com/facebookincubator/ent/dialect/gremlin"
	"github.com/facebookincubator/ent/dialect/gremlin/graph/dsl"
	"github.com/facebookincubator/ent/dialect/gremlin/graph/dsl/__"
	"github.com/facebookincubator/ent/dialect/gremlin/graph/dsl/g"
	"github.com/facebookincubator/ent/dialect/gremlin/graph/dsl/p"
	"github.com/facebookincubator/ent/entc/integration/gremlin/ent/pet"
	"github.com/facebookincubator/ent/entc/integration/gremlin/ent/user"
)

// PetCreate is the builder for creating a Pet entity.
type PetCreate struct {
	config
	name  *string
	team  map[string]struct{}
	owner map[string]struct{}
}

// SetName sets the name field.
func (pc *PetCreate) SetName(s string) *PetCreate {
	pc.name = &s
	return pc
}

// SetTeamID sets the team edge to User by id.
func (pc *PetCreate) SetTeamID(id string) *PetCreate {
	if pc.team == nil {
		pc.team = make(map[string]struct{})
	}
	pc.team[id] = struct{}{}
	return pc
}

// SetNillableTeamID sets the team edge to User by id if the given value is not nil.
func (pc *PetCreate) SetNillableTeamID(id *string) *PetCreate {
	if id != nil {
		pc = pc.SetTeamID(*id)
	}
	return pc
}

// SetTeam sets the team edge to User.
func (pc *PetCreate) SetTeam(u *User) *PetCreate {
	return pc.SetTeamID(u.ID)
}

// SetOwnerID sets the owner edge to User by id.
func (pc *PetCreate) SetOwnerID(id string) *PetCreate {
	if pc.owner == nil {
		pc.owner = make(map[string]struct{})
	}
	pc.owner[id] = struct{}{}
	return pc
}

// SetNillableOwnerID sets the owner edge to User by id if the given value is not nil.
func (pc *PetCreate) SetNillableOwnerID(id *string) *PetCreate {
	if id != nil {
		pc = pc.SetOwnerID(*id)
	}
	return pc
}

// SetOwner sets the owner edge to User.
func (pc *PetCreate) SetOwner(u *User) *PetCreate {
	return pc.SetOwnerID(u.ID)
}

// Save creates the Pet in the database.
func (pc *PetCreate) Save(ctx context.Context) (*Pet, error) {
	if pc.name == nil {
		return nil, errors.New("ent: missing required field \"name\"")
	}
	if len(pc.team) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"team\"")
	}
	if len(pc.owner) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"owner\"")
	}
	return pc.gremlinSave(ctx)
}

// SaveX calls Save and panics if Save returns an error.
func (pc *PetCreate) SaveX(ctx context.Context) *Pet {
	v, err := pc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (pc *PetCreate) gremlinSave(ctx context.Context) (*Pet, error) {
	res := &gremlin.Response{}
	query, bindings := pc.gremlin().Query()
	if err := pc.driver.Exec(ctx, query, bindings, res); err != nil {
		return nil, err
	}
	if err, ok := isConstantError(res); ok {
		return nil, err
	}
	pe := &Pet{config: pc.config}
	if err := pe.FromResponse(res); err != nil {
		return nil, err
	}
	return pe, nil
}

func (pc *PetCreate) gremlin() *dsl.Traversal {
	type constraint struct {
		pred *dsl.Traversal // constraint predicate.
		test *dsl.Traversal // test matches and its constant.
	}
	constraints := make([]*constraint, 0, 1)
	v := g.AddV(pet.Label)
	if pc.name != nil {
		v.Property(dsl.Single, pet.FieldName, *pc.name)
	}
	for id := range pc.team {
		v.AddE(user.TeamLabel).From(g.V(id)).InV()
		constraints = append(constraints, &constraint{
			pred: g.E().HasLabel(user.TeamLabel).OutV().HasID(id).Count(),
			test: __.Is(p.NEQ(0)).Constant(NewErrUniqueEdge(pet.Label, user.TeamLabel, id)),
		})
	}
	for id := range pc.owner {
		v.AddE(user.PetsLabel).From(g.V(id)).InV()
	}
	if len(constraints) == 0 {
		return v.ValueMap(true)
	}
	tr := constraints[0].pred.Coalesce(constraints[0].test, v.ValueMap(true))
	for _, cr := range constraints[1:] {
		tr = cr.pred.Coalesce(cr.test, tr)
	}
	return tr
}