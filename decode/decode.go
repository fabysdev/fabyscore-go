package decode

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

// contentTypeKey is the "Content-Type" header key.
var contentTypeKey = http.CanonicalHeaderKey("Content-Type")

// Request decodes the request data based on the Content-Type header (query + content-type).
func Request(r *http.Request, v interface{}) error {
	if v == nil {
		return nil
	}

	contentType := r.Header.Get(contentTypeKey)

	// query / form
	decodeQuery := r.URL.RawQuery != ""
	decodeForm := strings.HasPrefix(contentType, "multipart/form-data") || contentType == "application/x-www-form-urlencoded"

	if decodeQuery || decodeForm {
		queryValues := r.URL.Query()

		tType := reflect.TypeOf(v).Elem()
		if tType.Kind() != reflect.Struct {
			return errors.New("value must be a struct")
		}

		vType := reflect.ValueOf(v).Elem()

		for i := 0; i < tType.NumField(); i++ {
			vField := vType.Field(i)
			if !vField.CanSet() {
				continue
			}

			tField := tType.Field(i)

			// resolve value
			value := ""

			// form
			if decodeForm {
				if r.PostForm == nil {
					r.ParseMultipartForm(32 << 20)
				}

				value = resolveValue(r.PostForm, "form", tField)
			}

			// query
			if decodeQuery && value == "" {
				if r.Method != "GET" && r.Form == nil {
					r.ParseForm()
					queryValues = r.Form
				}

				value = resolveValue(queryValues, "query", tField)
			}

			err := setValue(vField, value, tField.Type.Kind())
			if err != nil {
				return err
			}
		}

		if decodeForm {
			return nil
		}
	}

	// json
	if r.ContentLength == 0 {
		return nil
	}

	if contentType == "application/json" {
		if err := json.NewDecoder(r.Body).Decode(v); err != nil {
			return err
		}

		return nil
	}

	return fmt.Errorf("content type '%v' can not be decoded", contentType)
}

// Query decodes the request data based on the query.
func Query(r *http.Request, v interface{}) error {
	if v == nil {
		return nil
	}

	if r.URL.RawQuery == "" {
		return nil
	}

	// resolve query data
	queryValues := r.URL.Query()
	if r.Method != "GET" && r.Form == nil {
		r.ParseForm()
		queryValues = r.Form
	}

	// resolve fields from type
	tType := reflect.TypeOf(v).Elem()
	if tType.Kind() != reflect.Struct {
		return errors.New("value must be a struct")
	}

	vType := reflect.ValueOf(v).Elem()

	for i := 0; i < tType.NumField(); i++ {
		vField := vType.Field(i)
		if !vField.CanSet() {
			continue
		}

		tField := tType.Field(i)

		// set field value
		value := resolveValue(queryValues, "query", tField)
		err := setValue(vField, value, tField.Type.Kind())
		if err != nil {
			return err
		}
	}

	return nil
}

// Form decodes the request data based on the form.
func Form(r *http.Request, v interface{}) error {
	if v == nil {
		return nil
	}

	// parse form
	r.ParseMultipartForm(32 << 20)

	// resolve fields from type
	tType := reflect.TypeOf(v).Elem()
	if tType.Kind() != reflect.Struct {
		return errors.New("value must be a struct")
	}

	vType := reflect.ValueOf(v).Elem()

	for i := 0; i < tType.NumField(); i++ {
		vField := vType.Field(i)
		if !vField.CanSet() {
			continue
		}

		tField := tType.Field(i)

		// set field value
		value := resolveValue(r.PostForm, "form", tField)
		err := setValue(vField, value, tField.Type.Kind())
		if err != nil {
			return err
		}
	}

	return nil
}

// JSON decodes the request body as json.
func JSON(r *http.Request, v interface{}) error {
	if v == nil {
		return nil
	}

	if r.ContentLength == 0 {
		return nil
	}

	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return err
	}

	return nil
}

// resolveValue returns the value for the field as string.
func resolveValue(data url.Values, tag string, tField reflect.StructField) string {
	key := tField.Tag.Get(tag)
	if key == "" {
		return ""
	}

	val := data.Get(key)
	if val == "" {
		for vk, vs := range data {
			if strings.EqualFold(vk, key) && len(vs) > 0 {
				val = vs[0]
				break
			}
		}
	}

	return val
}

// setValue sets the field value with the correct type.
func setValue(field reflect.Value, value string, kind reflect.Kind) error {
	switch kind {
	case reflect.Bool:
		return setBool(field, value)
	case reflect.Float32:
		return setFloat(field, value, 32)
	case reflect.Float64:
		return setFloat(field, value, 64)
	case reflect.Int:
		return setInt(field, value, 0)
	case reflect.Int8:
		return setInt(field, value, 8)
	case reflect.Int16:
		return setInt(field, value, 16)
	case reflect.Int32:
		return setInt(field, value, 32)
	case reflect.Int64:
		return setInt(field, value, 64)
	case reflect.Ptr:
		if field.IsNil() {
			field.Set(reflect.New(field.Type().Elem()))
		}

		return setValue(field.Elem(), value, field.Elem().Kind())
	case reflect.String:
		field.SetString(value)
	case reflect.Uint:
		return setUint(field, value, 0)
	case reflect.Uint8:
		return setUint(field, value, 8)
	case reflect.Uint16:
		return setUint(field, value, 16)
	case reflect.Uint32:
		return setUint(field, value, 32)
	case reflect.Uint64:
		return setUint(field, value, 64)
	default:
		return fmt.Errorf("unknown type: %v", kind.String())
	}

	return nil
}

// setInt sets the value as int.
func setInt(field reflect.Value, value string, bitSize int) error {
	if value == "" {
		value = "0"
	}

	val, err := strconv.ParseInt(value, 10, bitSize)
	if err != nil {
		return err
	}

	field.SetInt(val)

	return nil
}

// setUint sets the value as uint.
func setUint(field reflect.Value, value string, bitSize int) error {
	if value == "" {
		value = "0"
	}

	val, err := strconv.ParseUint(value, 10, bitSize)
	if err != nil {
		return err
	}

	field.SetUint(val)

	return nil
}

// setBool sets the value as bool.
func setBool(field reflect.Value, value string) error {
	if value == "" {
		value = "false"
	}

	val, err := strconv.ParseBool(value)
	if err != nil {
		return err
	}

	field.SetBool(val)

	return nil
}

// setFloat sets the value as float.
func setFloat(field reflect.Value, value string, bitSize int) error {
	if value == "" {
		value = "0.0"
	}

	val, err := strconv.ParseFloat(value, bitSize)
	if err != nil {
		return err
	}

	field.SetFloat(val)

	return nil
}
