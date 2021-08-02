package sessions

import "fmt"

const (
	KeyUserdata = `userdata`
	KeyModtime  = `modtime`
)

type Userdata struct {
	ID            int64
	Name          string
	Nickname      string
	Authorization []int32
}

func (d *Userdata) Who() string {
	return fmt.Sprintf("id=%d name=%s nickname=%s", d.ID, d.Name, d.Nickname)
}

// Test if has all authorization return true
func (d *Userdata) AuthTest(vals ...int32) bool {
	if d.ID == 0 {
		return false
	}
	var found bool
	for _, val := range vals {
		found = false
		for _, auth := range d.Authorization {
			if val == auth {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

// AuthAny if has any authorization return true
func (d *Userdata) AuthAny(vals ...int32) bool {
	if d.ID == 0 {
		return false
	}
	for _, val := range vals {
		for _, auth := range d.Authorization {
			if val == auth {
				return true
			}
		}
	}
	return false
}

// AuthNone if not has any authorization return true
func (d *Userdata) AuthNone(vals ...int32) bool {
	if d.ID == 0 {
		return false
	}
	for _, val := range vals {
		for _, auth := range d.Authorization {
			if val == auth {
				return false
			}
		}
	}
	return true
}
