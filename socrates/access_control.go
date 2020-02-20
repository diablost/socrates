package socrates

import (
	"regexp"
)

type AccessControl struct {
	WhiteList map[string]regexp.Regexp
}

func (a *AccessControl) access(path string) bool {
	if a.WhiteList == nil {
		return true
	}
	for k, v := range a.WhiteList {
		logf("access list key:%v, req:%v, match:%v",k, path, v.MatchString(path))
		return v.MatchString(path)
	}
	return false
}

func (a *AccessControl) trafficControl() {

}
