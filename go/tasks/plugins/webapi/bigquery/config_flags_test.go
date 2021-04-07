// Code generated by go generate; DO NOT EDIT.
// This file was generated by robots.

package bigquery

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
)

var dereferencableKindsConfig = map[reflect.Kind]struct{}{
	reflect.Array: {}, reflect.Chan: {}, reflect.Map: {}, reflect.Ptr: {}, reflect.Slice: {},
}

// Checks if t is a kind that can be dereferenced to get its underlying type.
func canGetElementConfig(t reflect.Kind) bool {
	_, exists := dereferencableKindsConfig[t]
	return exists
}

// This decoder hook tests types for json unmarshaling capability. If implemented, it uses json unmarshal to build the
// object. Otherwise, it'll just pass on the original data.
func jsonUnmarshalerHookConfig(_, to reflect.Type, data interface{}) (interface{}, error) {
	unmarshalerType := reflect.TypeOf((*json.Unmarshaler)(nil)).Elem()
	if to.Implements(unmarshalerType) || reflect.PtrTo(to).Implements(unmarshalerType) ||
		(canGetElementConfig(to.Kind()) && to.Elem().Implements(unmarshalerType)) {

		raw, err := json.Marshal(data)
		if err != nil {
			fmt.Printf("Failed to marshal Data: %v. Error: %v. Skipping jsonUnmarshalHook", data, err)
			return data, nil
		}

		res := reflect.New(to).Interface()
		err = json.Unmarshal(raw, &res)
		if err != nil {
			fmt.Printf("Failed to umarshal Data: %v. Error: %v. Skipping jsonUnmarshalHook", data, err)
			return data, nil
		}

		return res, nil
	}

	return data, nil
}

func decode_Config(input, result interface{}) error {
	config := &mapstructure.DecoderConfig{
		TagName:          "json",
		WeaklyTypedInput: true,
		Result:           result,
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToTimeDurationHookFunc(),
			mapstructure.StringToSliceHookFunc(","),
			jsonUnmarshalerHookConfig,
		),
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}

	return decoder.Decode(input)
}

func join_Config(arr interface{}, sep string) string {
	listValue := reflect.ValueOf(arr)
	strs := make([]string, 0, listValue.Len())
	for i := 0; i < listValue.Len(); i++ {
		strs = append(strs, fmt.Sprintf("%v", listValue.Index(i)))
	}

	return strings.Join(strs, sep)
}

func testDecodeJson_Config(t *testing.T, val, result interface{}) {
	assert.NoError(t, decode_Config(val, result))
}

func testDecodeSlice_Config(t *testing.T, vStringSlice, result interface{}) {
	assert.NoError(t, decode_Config(vStringSlice, result))
}

func TestConfig_GetPFlagSet(t *testing.T) {
	val := Config{}
	cmdFlags := val.GetPFlagSet("")
	assert.True(t, cmdFlags.HasFlags())
}

func TestConfig_SetFlags(t *testing.T) {
	actual := Config{}
	cmdFlags := actual.GetPFlagSet("")
	assert.True(t, cmdFlags.HasFlags())

	t.Run("Test_webApi.readRateLimiter.qps", func(t *testing.T) {
		t.Run("DefaultValue", func(t *testing.T) {
			// Test that default value is set properly
			if vInt, err := cmdFlags.GetInt("webApi.readRateLimiter.qps"); err == nil {
				assert.Equal(t, int(defaultConfig.WebAPI.ReadRateLimiter.QPS), vInt)
			} else {
				assert.FailNow(t, err.Error())
			}
		})

		t.Run("Override", func(t *testing.T) {
			testValue := "1"

			cmdFlags.Set("webApi.readRateLimiter.qps", testValue)
			if vInt, err := cmdFlags.GetInt("webApi.readRateLimiter.qps"); err == nil {
				testDecodeJson_Config(t, fmt.Sprintf("%v", vInt), &actual.WebAPI.ReadRateLimiter.QPS)

			} else {
				assert.FailNow(t, err.Error())
			}
		})
	})
	t.Run("Test_webApi.readRateLimiter.burst", func(t *testing.T) {
		t.Run("DefaultValue", func(t *testing.T) {
			// Test that default value is set properly
			if vInt, err := cmdFlags.GetInt("webApi.readRateLimiter.burst"); err == nil {
				assert.Equal(t, int(defaultConfig.WebAPI.ReadRateLimiter.Burst), vInt)
			} else {
				assert.FailNow(t, err.Error())
			}
		})

		t.Run("Override", func(t *testing.T) {
			testValue := "1"

			cmdFlags.Set("webApi.readRateLimiter.burst", testValue)
			if vInt, err := cmdFlags.GetInt("webApi.readRateLimiter.burst"); err == nil {
				testDecodeJson_Config(t, fmt.Sprintf("%v", vInt), &actual.WebAPI.ReadRateLimiter.Burst)

			} else {
				assert.FailNow(t, err.Error())
			}
		})
	})
	t.Run("Test_webApi.writeRateLimiter.qps", func(t *testing.T) {
		t.Run("DefaultValue", func(t *testing.T) {
			// Test that default value is set properly
			if vInt, err := cmdFlags.GetInt("webApi.writeRateLimiter.qps"); err == nil {
				assert.Equal(t, int(defaultConfig.WebAPI.WriteRateLimiter.QPS), vInt)
			} else {
				assert.FailNow(t, err.Error())
			}
		})

		t.Run("Override", func(t *testing.T) {
			testValue := "1"

			cmdFlags.Set("webApi.writeRateLimiter.qps", testValue)
			if vInt, err := cmdFlags.GetInt("webApi.writeRateLimiter.qps"); err == nil {
				testDecodeJson_Config(t, fmt.Sprintf("%v", vInt), &actual.WebAPI.WriteRateLimiter.QPS)

			} else {
				assert.FailNow(t, err.Error())
			}
		})
	})
	t.Run("Test_webApi.writeRateLimiter.burst", func(t *testing.T) {
		t.Run("DefaultValue", func(t *testing.T) {
			// Test that default value is set properly
			if vInt, err := cmdFlags.GetInt("webApi.writeRateLimiter.burst"); err == nil {
				assert.Equal(t, int(defaultConfig.WebAPI.WriteRateLimiter.Burst), vInt)
			} else {
				assert.FailNow(t, err.Error())
			}
		})

		t.Run("Override", func(t *testing.T) {
			testValue := "1"

			cmdFlags.Set("webApi.writeRateLimiter.burst", testValue)
			if vInt, err := cmdFlags.GetInt("webApi.writeRateLimiter.burst"); err == nil {
				testDecodeJson_Config(t, fmt.Sprintf("%v", vInt), &actual.WebAPI.WriteRateLimiter.Burst)

			} else {
				assert.FailNow(t, err.Error())
			}
		})
	})
	t.Run("Test_webApi.caching.size", func(t *testing.T) {
		t.Run("DefaultValue", func(t *testing.T) {
			// Test that default value is set properly
			if vInt, err := cmdFlags.GetInt("webApi.caching.size"); err == nil {
				assert.Equal(t, int(defaultConfig.WebAPI.Caching.Size), vInt)
			} else {
				assert.FailNow(t, err.Error())
			}
		})

		t.Run("Override", func(t *testing.T) {
			testValue := "1"

			cmdFlags.Set("webApi.caching.size", testValue)
			if vInt, err := cmdFlags.GetInt("webApi.caching.size"); err == nil {
				testDecodeJson_Config(t, fmt.Sprintf("%v", vInt), &actual.WebAPI.Caching.Size)

			} else {
				assert.FailNow(t, err.Error())
			}
		})
	})
	t.Run("Test_webApi.caching.resyncInterval", func(t *testing.T) {
		t.Run("DefaultValue", func(t *testing.T) {
			// Test that default value is set properly
			if vString, err := cmdFlags.GetString("webApi.caching.resyncInterval"); err == nil {
				assert.Equal(t, string(defaultConfig.WebAPI.Caching.ResyncInterval.String()), vString)
			} else {
				assert.FailNow(t, err.Error())
			}
		})

		t.Run("Override", func(t *testing.T) {
			testValue := defaultConfig.WebAPI.Caching.ResyncInterval.String()

			cmdFlags.Set("webApi.caching.resyncInterval", testValue)
			if vString, err := cmdFlags.GetString("webApi.caching.resyncInterval"); err == nil {
				testDecodeJson_Config(t, fmt.Sprintf("%v", vString), &actual.WebAPI.Caching.ResyncInterval)

			} else {
				assert.FailNow(t, err.Error())
			}
		})
	})
	t.Run("Test_webApi.caching.workers", func(t *testing.T) {
		t.Run("DefaultValue", func(t *testing.T) {
			// Test that default value is set properly
			if vInt, err := cmdFlags.GetInt("webApi.caching.workers"); err == nil {
				assert.Equal(t, int(defaultConfig.WebAPI.Caching.Workers), vInt)
			} else {
				assert.FailNow(t, err.Error())
			}
		})

		t.Run("Override", func(t *testing.T) {
			testValue := "1"

			cmdFlags.Set("webApi.caching.workers", testValue)
			if vInt, err := cmdFlags.GetInt("webApi.caching.workers"); err == nil {
				testDecodeJson_Config(t, fmt.Sprintf("%v", vInt), &actual.WebAPI.Caching.Workers)

			} else {
				assert.FailNow(t, err.Error())
			}
		})
	})
	t.Run("Test_webApi.caching.maxSystemFailures", func(t *testing.T) {
		t.Run("DefaultValue", func(t *testing.T) {
			// Test that default value is set properly
			if vInt, err := cmdFlags.GetInt("webApi.caching.maxSystemFailures"); err == nil {
				assert.Equal(t, int(defaultConfig.WebAPI.Caching.MaxSystemFailures), vInt)
			} else {
				assert.FailNow(t, err.Error())
			}
		})

		t.Run("Override", func(t *testing.T) {
			testValue := "1"

			cmdFlags.Set("webApi.caching.maxSystemFailures", testValue)
			if vInt, err := cmdFlags.GetInt("webApi.caching.maxSystemFailures"); err == nil {
				testDecodeJson_Config(t, fmt.Sprintf("%v", vInt), &actual.WebAPI.Caching.MaxSystemFailures)

			} else {
				assert.FailNow(t, err.Error())
			}
		})
	})
	t.Run("Test_googleTokenSource.type", func(t *testing.T) {
		t.Run("DefaultValue", func(t *testing.T) {
			// Test that default value is set properly
			if vString, err := cmdFlags.GetString("googleTokenSource.type"); err == nil {
				assert.Equal(t, string(defaultConfig.GoogleTokenSource.Type), vString)
			} else {
				assert.FailNow(t, err.Error())
			}
		})

		t.Run("Override", func(t *testing.T) {
			testValue := "1"

			cmdFlags.Set("googleTokenSource.type", testValue)
			if vString, err := cmdFlags.GetString("googleTokenSource.type"); err == nil {
				testDecodeJson_Config(t, fmt.Sprintf("%v", vString), &actual.GoogleTokenSource.Type)

			} else {
				assert.FailNow(t, err.Error())
			}
		})
	})
	t.Run("Test_googleTokenSource.identityNamespace", func(t *testing.T) {
		t.Run("DefaultValue", func(t *testing.T) {
			// Test that default value is set properly
			if vString, err := cmdFlags.GetString("googleTokenSource.identityNamespace"); err == nil {
				assert.Equal(t, string(defaultConfig.GoogleTokenSource.IdentityNamespace), vString)
			} else {
				assert.FailNow(t, err.Error())
			}
		})

		t.Run("Override", func(t *testing.T) {
			testValue := "1"

			cmdFlags.Set("googleTokenSource.identityNamespace", testValue)
			if vString, err := cmdFlags.GetString("googleTokenSource.identityNamespace"); err == nil {
				testDecodeJson_Config(t, fmt.Sprintf("%v", vString), &actual.GoogleTokenSource.IdentityNamespace)

			} else {
				assert.FailNow(t, err.Error())
			}
		})
	})
	t.Run("Test_googleTokenSource.scope", func(t *testing.T) {
		t.Run("DefaultValue", func(t *testing.T) {
			// Test that default value is set properly
			if vStringSlice, err := cmdFlags.GetStringSlice("googleTokenSource.scope"); err == nil {
				assert.Equal(t, []string([]string{}), vStringSlice)
			} else {
				assert.FailNow(t, err.Error())
			}
		})

		t.Run("Override", func(t *testing.T) {
			testValue := join_Config("1,1", ",")

			cmdFlags.Set("googleTokenSource.scope", testValue)
			if vStringSlice, err := cmdFlags.GetStringSlice("googleTokenSource.scope"); err == nil {
				testDecodeSlice_Config(t, join_Config(vStringSlice, ","), &actual.GoogleTokenSource.Scope)

			} else {
				assert.FailNow(t, err.Error())
			}
		})
	})
	t.Run("Test_googleTokenSource.kubeConfig", func(t *testing.T) {
		t.Run("DefaultValue", func(t *testing.T) {
			// Test that default value is set properly
			if vString, err := cmdFlags.GetString("googleTokenSource.kubeConfig"); err == nil {
				assert.Equal(t, string(defaultConfig.GoogleTokenSource.KubeConfigPath), vString)
			} else {
				assert.FailNow(t, err.Error())
			}
		})

		t.Run("Override", func(t *testing.T) {
			testValue := "1"

			cmdFlags.Set("googleTokenSource.kubeConfig", testValue)
			if vString, err := cmdFlags.GetString("googleTokenSource.kubeConfig"); err == nil {
				testDecodeJson_Config(t, fmt.Sprintf("%v", vString), &actual.GoogleTokenSource.KubeConfigPath)

			} else {
				assert.FailNow(t, err.Error())
			}
		})
	})
	t.Run("Test_googleTokenSource.kubeClientConfig.timeout", func(t *testing.T) {
		t.Run("DefaultValue", func(t *testing.T) {
			// Test that default value is set properly
			if vString, err := cmdFlags.GetString("googleTokenSource.kubeClientConfig.timeout"); err == nil {
				assert.Equal(t, string(defaultConfig.GoogleTokenSource.KubeClientConfig.Timeout.String()), vString)
			} else {
				assert.FailNow(t, err.Error())
			}
		})

		t.Run("Override", func(t *testing.T) {
			testValue := defaultConfig.GoogleTokenSource.KubeClientConfig.Timeout.String()

			cmdFlags.Set("googleTokenSource.kubeClientConfig.timeout", testValue)
			if vString, err := cmdFlags.GetString("googleTokenSource.kubeClientConfig.timeout"); err == nil {
				testDecodeJson_Config(t, fmt.Sprintf("%v", vString), &actual.GoogleTokenSource.KubeClientConfig.Timeout)

			} else {
				assert.FailNow(t, err.Error())
			}
		})
	})
	t.Run("Test_bigQueryEndpoint", func(t *testing.T) {
		t.Run("DefaultValue", func(t *testing.T) {
			// Test that default value is set properly
			if vString, err := cmdFlags.GetString("bigQueryEndpoint"); err == nil {
				assert.Equal(t, string(defaultConfig.bigQueryEndpoint), vString)
			} else {
				assert.FailNow(t, err.Error())
			}
		})

		t.Run("Override", func(t *testing.T) {
			testValue := "1"

			cmdFlags.Set("bigQueryEndpoint", testValue)
			if vString, err := cmdFlags.GetString("bigQueryEndpoint"); err == nil {
				testDecodeJson_Config(t, fmt.Sprintf("%v", vString), &actual.bigQueryEndpoint)

			} else {
				assert.FailNow(t, err.Error())
			}
		})
	})
}
