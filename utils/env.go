package utils

import (
	"os"
	"strconv"
)

type EnvKey string

type EnvGetter struct {
	Key          EnvKey
	DefaultValue any
}

var envMap map[EnvKey]any

func init() {
	envMap = make(map[EnvKey]any, 0)
}

func GetEnvMap() map[EnvKey]any {
	return envMap
}

func GetterDefault[T int | uint | int64 | uint64 | bool | string](key EnvKey, defaultValue T) *EnvGetter {
	if _, ok := envMap[key]; !ok {
		envMap[key] = defaultValue
	}
	return &EnvGetter{
		Key:          key,
		DefaultValue: defaultValue,
	}
}

func Getter(key EnvKey) *EnvGetter {
	return &EnvGetter{
		Key:          key,
		DefaultValue: nil,
	}
}

func (getter *EnvGetter) KeyName() string {
	return string(getter.Key)
}

func (getter *EnvGetter) String() string {
	if getter.DefaultValue == nil {
		return GetEnv(getter.Key)
	} else {
		return GetEnvDefault(getter.Key, getter.DefaultValue.(string))
	}
}

func (getter *EnvGetter) Bool() bool {
	if getter.DefaultValue == nil {
		return GetEnvBool(getter.Key)
	} else {
		return GetEnvBoolDefault(getter.Key, getter.DefaultValue.(bool))
	}
}

func (getter *EnvGetter) Int() (r int) {
	if getter.DefaultValue == nil {
		return GetEnvIntDefault(getter.Key, r)
	} else {
		return GetEnvIntDefault(getter.Key, getter.DefaultValue.(int))
	}
}

func (getter *EnvGetter) Int64() (r int64) {
	if getter.DefaultValue == nil {
		return GetEnvIntDefault(getter.Key, r)
	} else {
		return GetEnvIntDefault(getter.Key, getter.DefaultValue.(int64))
	}
}

func (getter *EnvGetter) UInt() (r uint) {
	if getter.DefaultValue == nil {
		return GetEnvIntDefault(getter.Key, r)
	} else {
		return GetEnvIntDefault(getter.Key, getter.DefaultValue.(uint))
	}
}

func (getter *EnvGetter) UInt64() (r uint64) {
	if getter.DefaultValue == nil {
		return GetEnvIntDefault(getter.Key, r)
	} else {
		return GetEnvIntDefault(getter.Key, getter.DefaultValue.(uint64))
	}
}

func GetEnv(name EnvKey) string {
	return os.Getenv(string(name))
}

func GetEnvDefault(name EnvKey, defaultValue string) string {
	r := os.Getenv(string(name))
	if r == "" {
		return defaultValue
	}
	return r
}

func GetEnvBool(name EnvKey) bool {
	b, _ := strconv.ParseBool(GetEnv(name))
	return b
}

func GetEnvBoolDefault(name EnvKey, defaultValue bool) bool {
	r := os.Getenv(string(name))
	if r == "" {
		return defaultValue
	}
	b, _ := strconv.ParseBool(r)
	return b
}

func GetEnvInt[T int | uint | int64 | uint64](name EnvKey) T {
	r, _ := strconv.ParseInt(GetEnv(name), 10, 0)
	return T(r)
}

func GetEnvIntDefault[T int | uint | int64 | uint64](name EnvKey, defaultValue T) T {
	r := os.Getenv(string(name))
	if r == "" {
		return defaultValue
	}
	v, _ := strconv.ParseInt(r, 10, 0)
	return T(v)
}
