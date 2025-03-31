package config

import (
	"fmt"
	"strconv"
	"strings"
)

type UintSlice []uint

func (u *UintSlice) String() string {
	return fmt.Sprintf("%v", *u)
}

func (u *UintSlice) Set(value string) error {
	values := strings.Split(value, ",")
	for _, v := range values {
		num, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid uint value: %s", v)
		}
		*u = append(*u, uint(num))
	}
	return nil
}
