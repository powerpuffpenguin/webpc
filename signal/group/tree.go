package group

import (
	"context"
	"strconv"
	"sync"
	"time"

	"github.com/powerpuffpenguin/webpc/signal/group/internal/db"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const rootID = 1

var defaultTree = newTree()

func DefaultTree() *Tree {
	return defaultTree
}

type Tree struct {
	modtime time.Time
	keys    map[int64]*Element
	root    *Element
	rw      sync.RWMutex
}

func newTree() *Tree {
	return &Tree{
		keys: make(map[int64]*Element),
	}
}
func (t *Tree) Init() {
	items := db.Init()
	keys := make(map[int64]*Element)
	for _, item := range items {
		ele := NewElement(item.ID, item.Name, item.Description)
		keys[ele.ID] = ele
	}
	for _, item := range items {
		parent := keys[item.Parent]
		if parent == nil {
			return
		}
		ele := keys[item.ID]
		ele.Parent = parent
		parent.Children = append(parent.Children, ele)
	}
	t.root = keys[db.RootID]
	t.modtime = db.LastModified()
}

func (t *Tree) Move(ctx context.Context, id, pid int64) (changed bool, e error) {
	if id == rootID {
		e = status.Error(codes.PermissionDenied, `root not supported move`)
		return
	}

	t.rw.Lock()
	defer t.rw.Unlock()

	ele, exists := t.keys[id]
	if !exists {
		e = status.Error(codes.NotFound, `id not exists: `+strconv.FormatInt(id, 10))
		return
	} else if ele.Parent.ID == pid {
		return
	}
	parent, exists := t.keys[pid]
	if !exists {
		e = status.Error(codes.NotFound, `parent not exists: `+strconv.FormatInt(pid, 10))
		return
	}
	// update db
	at, e := db.Move(ctx, id, pid)
	if e != nil {
		return
	}

	// update memory
	parent.AddChild(ele)
	changed = true
	t.modtime = at
	return
}
func (t *Tree) Add(ctx context.Context, pid int64, name, description string) (id int64, e error) {
	t.rw.Lock()
	defer t.rw.Unlock()
	parent, exists := t.keys[pid]
	if !exists {
		e = status.Error(codes.NotFound, `parent not exists: `+strconv.FormatInt(pid, 10))
		return
	}
	// update db
	id, at, e := db.Add(ctx, pid, name, description)
	if e != nil {
		return
	}
	// update memory
	parent.AddChild(NewElement(id, name, description))

	t.modtime = at
	return
}
func (t *Tree) Len() int {
	t.rw.RLock()
	result := len(t.keys)
	t.rw.RUnlock()
	return result
}
func (t *Tree) Foreach(callback func(ele *Element) (e error)) (e error) {
	t.rw.RLock()
	e = t.foreach(t.root, callback)
	t.rw.RUnlock()
	return
}
func (t *Tree) foreach(ele *Element, callback func(ele *Element) (e error)) (e error) {
	e = callback(ele)
	if e != nil {
		return
	}
	for _, child := range ele.Children {
		e = t.foreach(child, callback)
		if e != nil {
			break
		}
	}
	return
}
func (t *Tree) LastModified() (modtime time.Time) {
	return t.modtime
}
func (t *Tree) Change(ctx context.Context, id int64, name, description string) (changed bool, e error) {
	t.rw.Lock()
	defer t.rw.Unlock()
	ele, exists := t.keys[id]
	if !exists {
		e = status.Error(codes.NotFound, `id not exists: `+strconv.FormatInt(id, 10))
		return
	} else if ele.Name == name && ele.Description == description {
		return
	}
	// update db
	at, e := db.Change(ctx, id, name, description)
	if e != nil {
		return
	}
	// update memory
	ele.Name = name
	ele.Description = description

	changed = true
	t.modtime = at
	return
}
func (t *Tree) Remove(ctx context.Context, id int64) (rowsAffected int, e error) {
	if id == rootID {
		e = status.Error(codes.PermissionDenied, `root not supported move`)
		return
	}

	t.rw.Lock()
	defer t.rw.Unlock()
	ele, exists := t.keys[id]
	if !exists {
		e = status.Error(codes.NotFound, `id not exists: `+strconv.FormatInt(id, 10))
		return
	}

	var args []interface{}
	t.foreach(ele, func(current *Element) (e error) {
		args = append(args, current.ID)
		return nil
	})
	// update db
	at, e := db.Remove(ctx, args)
	if e != nil {
		return
	}

	// update memory
	ele.Parent.RemoveChild(ele)

	rowsAffected = len(args)
	t.modtime = at
	return
}
