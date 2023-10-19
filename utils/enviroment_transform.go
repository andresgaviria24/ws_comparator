package utils

import (
	"fmt"
	"os"
	"strconv"
)

const SOME_ERROR = "some error"

type EnvGetter interface {
	GetStrEnv(key string) string
	GetIntEnv(key string) int
	GetDoubleEnv(key string) float64
	GetBoolEnv(key string) bool
}

type RealUtils struct{}

func (ru *RealUtils) GetStrEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		panic(fmt.Sprintf("some error msg"))
	}
	return val
}

func (ru *RealUtils) GetIntEnv(key string) int {
	val := ru.GetStrEnv(key)
	ret, err := strconv.Atoi(val)
	if err != nil {
		panic(fmt.Sprintf(SOME_ERROR))
	}
	return ret
}

func (ru *RealUtils) GetDoubleEnv(key string) float64 {
	val := ru.GetStrEnv(key)
	ret, err := strconv.ParseFloat(val, 64)
	if err != nil {
		panic(fmt.Sprintf(SOME_ERROR))
	}
	return ret
}

func (ru *RealUtils) GetBoolEnv(key string) bool {
	val := ru.GetStrEnv(key)
	ret, err := strconv.ParseBool(val)
	if err != nil {
		panic(fmt.Sprintf(SOME_ERROR))
	}
	return ret
}
