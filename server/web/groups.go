package web

import (
	"context"
	"xfsmiddle"
)

type Groups struct {
	Groups *xfsmiddle.Groups
}

func (g *Groups) GetGroups(ctx context.Context, _ *Empty, reply *[]xfsmiddle.Group) error {
	result := make([]xfsmiddle.Group, 0)
	for _, v := range g.Groups.GetAll().Rights {
		result = append(result, *v)
	}
	*reply = result
	return nil
}
