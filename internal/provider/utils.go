package provider

import (
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func convertStringSet(set *schema.Set) []string {
	s := make([]string, set.Len())
	for i, v := range set.List() {
		s[i] = v.(string)
	}

	return s
}

func convertInt32Set(set *schema.Set) []int32 {
	s := make([]int32, set.Len())
	for i, v := range set.List() {
		s[i] = int32(v.(int))
	}

	return s
}

func convertStringList(list []interface{}) []string {
	s := make([]string, len(list))
	for i, v := range list {
		s[i] = v.(string)
	}

	return s
}

func stringElem(v interface{}) string {
	val := reflect.Indirect(reflect.ValueOf(v))
	if v == nil || isEmptyValue(val) || val.IsZero() {
		return ""
	}

	return val.Interface().(string)
}

func isValid(v interface{}) bool {
	return !isEmptyValue(reflect.ValueOf(v))
}

// https://github.com/hashicorp/terraform-provider-google/commit/9900e6a4c70294db07dec023c9da7e27a12ee464#diff-16751bdc2e307bd9d601bfeb3d2e62100c18679a6f7633cc5afca8d783a1dd4cR61
func isEmptyValue(v reflect.Value) bool {
	if !v.IsValid() {
		return true
	}

	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return false
}

// merge returns a new map with all the keys specified within each arguement.
// Keys will never be overriden once set.
func merge(maps ...map[string]interface{}) map[string]interface{} {
	merged := make(map[string]interface{})
	for _, m := range maps {
		for k, v := range m {
			if _, ok := merged[k]; ok {
				continue
			}
			merged[k] = v
		}
	}
	return merged
}
