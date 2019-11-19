package main

import (
	"regexp"
)

type AccessControl struct {
	whiteList map[string]regexp.Regexp
}

func (a *AccessControl) access(path string) bool {
	if a.whiteList == nil {
		return true
	}
	for k, v := range a.whiteList {
		logf("access list key:%v, req:%v, match:%v",k, path, v.MatchString(path))
		return v.MatchString(path)
	}
	return false
}

func (a *AccessControl) trafficControl() {

}
