package main

import (
	"fmt"
	"regexp"
)

type AccessControl struct {
	whiteList map[string]bool
}

func (a *AccessControl) access(path string) bool {
	for k, _ := range a.whiteList {
		regStr := fmt.Sprintf(`%s`, k)
		reg, _ := regexp.Compile(regStr)
		fmt.Println(k, path, reg.MatchString(path))
		return reg.MatchString(path)
	}
	return false
}

func (a *AccessControl) trafficControl() {

}
