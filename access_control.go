package main

import (
	"fmt"
	"regexp"
)

type AccessControl struct {
	whiteList map[string]bool
}

func (a *AccessControl) access(path string) bool {
	for k, v := range a.whiteList {
		fmt.Println(k, v)

		regStr := fmt.Sprintf("`(?i:^%s).*`",k)
		reg := regexp.MustCompile(regStr)
		fmt.Printf("%q\n", reg.FindAllString(path, -1))
	}
	return true
}

func (a *AccessControl) trafficControl() {

}
