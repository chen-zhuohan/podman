package generate

import (
	"net/url"
	"reflect"
	"strconv"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
)

/*
This file is generated automatically by go generate.  Do not edit.

Created 2020-12-18 15:58:20.522950566 -0600 CST m=+0.000154384
*/

// Changed
func (o *KubeOptions) Changed(fieldName string) bool {
	r := reflect.ValueOf(o)
	value := reflect.Indirect(r).FieldByName(fieldName)
	return !value.IsNil()
}

// ToParams
func (o *KubeOptions) ToParams() (url.Values, error) {
	params := url.Values{}
	if o == nil {
		return params, nil
	}
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	s := reflect.ValueOf(o)
	if reflect.Ptr == s.Kind() {
		s = s.Elem()
	}
	sType := s.Type()
	for i := 0; i < s.NumField(); i++ {
		fieldName := sType.Field(i).Name
		if !o.Changed(fieldName) {
			continue
		}
		f := s.Field(i)
		if reflect.Ptr == f.Kind() {
			f = f.Elem()
		}
		switch f.Kind() {
		case reflect.Bool:
			params.Set(fieldName, strconv.FormatBool(f.Bool()))
		case reflect.String:
			params.Set(fieldName, f.String())
		case reflect.Int, reflect.Int64:
			// f.Int() is always an int64
			params.Set(fieldName, strconv.FormatInt(f.Int(), 10))
		case reflect.Uint, reflect.Uint64:
			// f.Uint() is always an uint64
			params.Set(fieldName, strconv.FormatUint(f.Uint(), 10))
		case reflect.Slice:
			typ := reflect.TypeOf(f.Interface()).Elem()
			switch typ.Kind() {
			case reflect.String:
				sl := f.Slice(0, f.Len())
				s, ok := sl.Interface().([]string)
				if !ok {
					return nil, errors.New("failed to convert to string slice")
				}
				for _, val := range s {
					params.Add(fieldName, val)
				}
			default:
				return nil, errors.Errorf("unknown slice type %s", f.Kind().String())
			}
		case reflect.Map:
			lowerCaseKeys := make(map[string][]string)
			iter := f.MapRange()
			for iter.Next() {
				lowerCaseKeys[iter.Key().Interface().(string)] = iter.Value().Interface().([]string)

			}
			s, err := json.MarshalToString(lowerCaseKeys)
			if err != nil {
				return nil, err
			}

			params.Set(fieldName, s)
		}
	}
	return params, nil
}

// WithService
func (o *KubeOptions) WithService(value bool) *KubeOptions {
	v := &value
	o.Service = v
	return o
}

// GetService
func (o *KubeOptions) GetService() bool {
	var service bool
	if o.Service == nil {
		return service
	}
	return *o.Service
}
