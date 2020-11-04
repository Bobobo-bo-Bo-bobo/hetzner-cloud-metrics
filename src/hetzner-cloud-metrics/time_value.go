package main

import (
	"fmt"
	"strconv"
)

func extractTimeValue(tv HetznerTimeValue) (float64, float64, error) {
	var t float64
	var v float64
	var err error

	if getInterfaceType(tv[0]) != TypeFloat {
		return t, v, fmt.Errorf("First element of time-value array is of type %s instead of float", TypeNameMap[getInterfaceType(tv[0])])
	}
	t = tv[0].(float64)

	if getInterfaceType(tv[1]) != TypeString {
		return t, v, fmt.Errorf("Second element of time-value array is of type %s instead of float", TypeNameMap[getInterfaceType(tv[1])])
	}
	v, err = strconv.ParseFloat(tv[1].(string), 64)
	if err != nil {
		return t, v, err
	}

	return t, v, nil
}
