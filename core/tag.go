package core

type LatticeTag struct {
	tag     string
	hashset map[string]bool
}

func (t *LatticeTag) Add(tag string) {
	if t.hashset == nil{
		t.hashset = make(map[string]bool)
	}
	if _,ok := t.hashset[tag];ok{
		return
	}

	t.tag += tag + " | "
	t.hashset[tag] = true
}

func (t *LatticeTag) Contains(tag string) bool {
	_,ok :=t.hashset[tag]
	return ok
}

func (t *LatticeTag) Delete(tag string) {
	delete(t.hashset,tag)
}

func (t *LatticeTag) String() string {
	return t.tag
}
