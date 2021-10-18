package core

type LatticeTag struct {
	tag    string
	msgSet map[string]bool
}

func (t *LatticeTag) Add(tag string) {
	if t.msgSet == nil{
		t.msgSet = make(map[string]bool)
	}
	if _,ok := t.msgSet[tag];ok{
		return
	}

	t.tag += tag + " | "
	t.msgSet[tag] = true
}

func (t *LatticeTag) Contains(tag string) bool {
	_,ok :=t.msgSet[tag]
	return ok
}

func (t *LatticeTag) Delete(tag string) {
	delete(t.msgSet,tag)
}

func (t *LatticeTag) String() string {
	return t.tag
}
