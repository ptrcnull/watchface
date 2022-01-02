package system

import (
	"io/ioutil"
	"strings"
)

func GetBattery() string {
	res, _ := ioutil.ReadFile("/sys/class/power_supply/battery/capacity")
	value := strings.Trim(string(res), "\n")
	if len(value) == 1 {
		value = "0" + value
	}
	if value == "100" {
		return "uwu"
	}
	return value + "%"
}
