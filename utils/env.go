package utils

import (
	"os"
	"reflect"
	"strconv"
	"sync"
)

type EnvKey string

type EnvGetter struct {
	Key          EnvKey
	defaultValue any
	kind         reflect.Kind
}

var envMap sync.Map

func init() {
	envMap = sync.Map{}
}

func getterOf(key EnvKey, defaultValue any, kind reflect.Kind) *EnvGetter {
	ref, ok := envMap.Load(key)
	if !ok {
		ref = &EnvGetter{
			Key:          key,
			defaultValue: defaultValue,
			kind:         kind,
		}
		envMap.Store(key, ref)
	}
	return ref.(*EnvGetter)
}

func Getter(key EnvKey) *EnvGetter {
	if ref, ok := envMap.Load(key); ok {
		return ref.(*EnvGetter)
	}
	return nil
}

func GetterDefault[T int | int8 | int16 | int32 | int64 |
	uint | uint8 | uint16 | uint32 | uint64 |
	string | bool](key EnvKey, defaultValue T) *EnvGetter {
	return getterOf(key, defaultValue, reflect.TypeOf(defaultValue).Kind())
}

func (getter *EnvGetter) WithDefault(value any) *EnvGetter {
	getter.defaultValue = value
	return getter
}

func (getter *EnvGetter) KeyName() string {
	return string(getter.Key)
}

func (getter *EnvGetter) Kind() reflect.Kind {
	return getter.kind
}

func (getter *EnvGetter) Value() any {
	if getter.kind == reflect.Int || getter.kind == reflect.Int8 || getter.kind == reflect.Int16 || getter.kind == reflect.Int32 {
		return getter.Int()
	} else if getter.kind == reflect.Uint || getter.kind == reflect.Uint8 || getter.kind == reflect.Uint16 || getter.kind == reflect.Uint32 {
		return getter.UInt()
	} else if getter.kind == reflect.Int64 {
		return getter.Int64()
	} else if getter.kind == reflect.Uint64 {
		return getter.UInt64()
	} else if getter.kind == reflect.String {
		return getter.String()
	} else if getter.kind == reflect.Bool {
		return getter.Bool()
	}
	return nil
}

func (getter *EnvGetter) String() string {
	return GetEnvDefault(getter.Key, getter.defaultValue.(string))
}

func (getter *EnvGetter) Bool() bool {
	return GetEnvBoolDefault(getter.Key, getter.defaultValue.(bool))
}

func (getter *EnvGetter) Int() (r int) {
	return GetEnvIntDefault(getter.Key, getter.defaultValue.(int))
}

func (getter *EnvGetter) Int64() (r int64) {
	return GetEnvIntDefault(getter.Key, getter.defaultValue.(int64))
}

func (getter *EnvGetter) UInt() (r uint) {
	return GetEnvIntDefault(getter.Key, getter.defaultValue.(uint))
}

func (getter *EnvGetter) UInt64() (r uint64) {
	return GetEnvIntDefault(getter.Key, getter.defaultValue.(uint64))
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
