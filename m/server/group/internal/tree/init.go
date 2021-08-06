package tree

import (
	signal_group "github.com/powerpuffpenguin/webpc/signal/group"
)

func init() {
	signal_group.ConnectIDS(soltIDS)
	signal_group.ConnectExists(soltExists)
}

func soltIDS(req *signal_group.IDSRequest, resp *signal_group.IDSResponse) (e error) {
	t := defaultTree
	t.rw.RLock()
	defer t.rw.RUnlock()
	root := t.keys[req.ID]
	if root == nil {
		return
	}

	t.foreach(root, func(ele *Element) (e error) {
		if req.Args {
			resp.Args = append(resp.Args, ele.ID)
		} else {
			resp.ID = append(resp.ID, ele.ID)
		}
		return nil
	})
	return
}

func soltExists(req *signal_group.ExistsRequest, resp *signal_group.ExistsResponse) (e error) {
	t := defaultTree
	t.rw.RLock()
	_, resp.Exists = t.keys[req.ID]
	t.rw.RUnlock()
	return
}
