package tree

type Element struct {
	Parent      *Element
	ID          int64
	Name        string
	Description string
	Children    []*Element
}

func NewElement(id int64, name, description string) *Element {
	return &Element{
		ID:          id,
		Name:        name,
		Description: description,
	}
}
func (ele *Element) IsLeaf() bool {
	return len(ele.Children) == 0
}
func (ele *Element) AddChild(child *Element) {
	parent := child.Parent
	if parent == ele {
		return
	}

	ele.Children = append(ele.Children, child)
	if parent != nil {
		parent.RemoveChild(child)
	}

	child.Parent = ele
}
func (ele *Element) RemoveChild(child *Element) {
	parent := child.Parent
	if ele != parent {
		return
	}

	child.Parent = nil
	for i, node := range ele.Children {
		if node == child {
			copy(ele.Children[i:], ele.Children[i+1:])
			ele.Children = ele.Children[:len(ele.Children)-1]
			break
		}
	}
}
