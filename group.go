package xfsmiddle

type Group struct {
	Types   string
	Methods []string
}
type Groups struct {
	Rights []*Group
}

func NewGroup(name string, Methods []string) *Group {
	return &Group{
		Types:   name,
		Methods: Methods,
	}
}

// func (g *Groups) PrefixData(types string)
func (g *Groups) Get(Methods []string) []string {
	rights := make([]string, 0)
	for _, val := range Methods {
		for _, gv := range g.Rights {
			for _, gvs := range gv.Methods {
				if gvs == val {
					rights = append(rights, val)
				}
			}
		}
	}
	return rights
}

func (g *Groups) GetAll() Groups {
	return *g
}

func (g *Groups) GetTypes() []string {
	result := make([]string, 0)
	for _, v := range g.Rights {
		result = append(result, v.Types)
	}
	return result
}

func (g *Groups) GetTypesGroup(Types string) *Group {
	for _, g := range g.Rights {
		if g.Types == Types {
			return g
		}
	}
	return nil
}

func (g *Groups) CheckGroup(x string, y []string) bool {

	for _, val := range y {
		if val == x {
			return true
		}
	}
	return false
}
