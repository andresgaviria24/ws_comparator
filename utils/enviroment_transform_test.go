package utils

/*
import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetStrEnv(t *testing.T) {

	key := "TEST_KEY"
	value := "test_value"
	os.Setenv(key, value)
	defer os.Unsetenv(key)

	result := GetStrEnv(key)
	assert.Equal(t, value, result)
}

func TestGetIntEnv(t *testing.T) {

	key := "TEST_INT_KEY"
	value := "42"
	os.Setenv(key, value)
	defer os.Unsetenv(key)

	result := GetIntEnv(key)
	assert.Equal(t, 42, result)
}

func TestGetDoubleEnv(t *testing.T) {

	key := "TEST_DOUBLE_KEY"
	value := "3.14"
	os.Setenv(key, value)
	defer os.Unsetenv(key)

	result := GetDoubleEnv(key)
	assert.Equal(t, 3.14, result)
}

func TestGetBoolEnv(t *testing.T) {

	key := "TEST_BOOL_KEY"
	value := "true"
	os.Setenv(key, value)
	defer os.Unsetenv(key)

	result := GetBoolEnv(key)
	assert.True(t, result)
}

func TestGetStrEnv_Panic(t *testing.T) {
	assert.Panics(t, func() {
		GetStrEnv("NonExistentKey")
	})
}

func TestGetIntEnv_Panic(t *testing.T) {

	key := "INVALID_INT_KEY"
	value := "not_an_integer"
	os.Setenv(key, value)
	defer os.Unsetenv(key)

	assert.Panics(t, func() {
		GetIntEnv(key)
	})
}

func TestGetDoubleEnv_Panic(t *testing.T) {

	key := "INVALID_DOUBLE_KEY"
	value := "not_a_double"
	os.Setenv(key, value)
	defer os.Unsetenv(key)

	assert.Panics(t, func() {
		GetDoubleEnv(key)
	})
}

func TestGetBoolEnv_Panic(t *testing.T) {

	key := "INVALID_BOOL_KEY"
	value := "not_a_boolean"
	os.Setenv(key, value)
	defer os.Unsetenv(key)

	assert.Panics(t, func() {
		GetBoolEnv(key)
	})
}
*/
