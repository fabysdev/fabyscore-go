package decode

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type testTypesReq struct {
	Bool       bool     `form:"bool" json:"bool"`
	PtrBool    *bool    `form:"ptrbool" json:"ptrbool"`
	Float32    float32  `form:"float32" json:"float32"`
	PtrFloat32 *float32 `form:"ptrfloat32" json:"ptrfloat32"`
	Float64    float64  `form:"float64" json:"float64"`
	PtrFloat64 *float64 `form:"ptrfloat64" json:"ptrfloat64"`
	Int        int      `form:"int" json:"int"`
	PtrInt     *int     `form:"ptrint" json:"ptrint"`
	Int8       int8     `form:"int8" json:"int8"`
	PtrInt8    *int8    `form:"ptrint8" json:"ptrint8"`
	Int16      int16    `form:"int16" json:"int16"`
	PtrInt16   *int16   `form:"ptrint16" json:"ptrint16"`
	Int32      int32    `form:"int32" json:"int32"`
	PtrInt32   *int32   `form:"ptrint32" json:"ptrint32"`
	Int64      int64    `form:"int64" json:"int64"`
	PtrInt64   *int64   `form:"ptrint64" json:"ptrint64"`
	String     string   `form:"string" json:"string"`
	PtrString  *string  `form:"ptrstring" json:"ptrstring"`
	Uint       uint     `form:"uint" json:"uint"`
	PtrUint    *uint    `form:"ptruint" json:"ptruint"`
	Uint8      uint8    `form:"uint8" json:"uint8"`
	PtrUint8   *uint8   `form:"ptruint8" json:"ptruint8"`
	Uint16     uint16   `form:"uint16" json:"uint16"`
	PtrUint16  *uint16  `form:"ptruint16" json:"ptruint16"`
	Uint32     uint32   `form:"uint32" json:"uint32"`
	PtrUint32  *uint32  `form:"ptruint32" json:"ptruint32"`
	Uint64     uint64   `form:"uint64" json:"uint64"`
	PtrUint64  *uint64  `form:"ptruint64" json:"ptruint64"`
	intern     string
}

type testBaseReq struct {
	Name  string `query:"name" form:"name" json:"name"`
	Email string `query:"email" form:"email" json:"email"`
}

func TestRequestQuery(t *testing.T) {
	r, _ := http.NewRequest("GET", "/?name=test&email=test@test", nil)

	tr := new(testBaseReq)

	err := Request(r, tr)
	assert.Nil(t, err)
	assert.Equal(t, "test", tr.Name)
	assert.Equal(t, "test@test", tr.Email)
}

func TestRequestQueryCaseInsensitive(t *testing.T) {
	r, _ := http.NewRequest("GET", "/?name=test&NAME=example&EMAIL=test@test", nil)

	tr := new(testBaseReq)

	err := Request(r, tr)
	assert.Nil(t, err)
	assert.Equal(t, "test", tr.Name)
	assert.Equal(t, "test@test", tr.Email)
}

func TestQuery(t *testing.T) {
	r, _ := http.NewRequest("GET", "/?name=test&email=test@test", nil)

	type testBaseInternReq struct {
		Name   string `query:"name" form:"name" json:"name"`
		Email  string `query:"email" form:"email" json:"email"`
		intern string
	}

	tr := new(testBaseInternReq)

	err := Query(r, tr)
	assert.Nil(t, err)
	assert.Equal(t, "test", tr.Name)
	assert.Equal(t, "test@test", tr.Email)
}

func TestQueryCaseInsensitive(t *testing.T) {
	r, _ := http.NewRequest("GET", "/?name=test&NAME=example&EMAIL=test@test", nil)

	tr := new(testBaseReq)

	err := Query(r, tr)
	assert.Nil(t, err)
	assert.Equal(t, "test", tr.Name)
	assert.Equal(t, "test@test", tr.Email)
}

func TestQueryEmpty(t *testing.T) {
	r, _ := http.NewRequest("GET", "/", nil)

	tr := new(testBaseReq)

	err := Query(r, tr)
	assert.Nil(t, err)
	assert.Equal(t, "", tr.Name)
	assert.Equal(t, "", tr.Email)
}

func TestRequestFormUrlencoded(t *testing.T) {
	r, _ := http.NewRequest("POST", "/", strings.NewReader("name=test&email=test@test"))
	r.Header.Set("content-type", "application/x-www-form-urlencoded")

	tr := new(testBaseReq)

	err := Request(r, tr)
	assert.Nil(t, err)
	assert.Equal(t, "test", tr.Name)
	assert.Equal(t, "test@test", tr.Email)
}

func TestRequestFormUrlencodedCaseInsensitive(t *testing.T) {
	r, _ := http.NewRequest("POST", "/", strings.NewReader("name=test&NAME=example&EMAIL=test@test"))
	r.Header.Set("content-type", "application/x-www-form-urlencoded")

	tr := new(testBaseReq)

	err := Request(r, tr)
	assert.Nil(t, err)
	assert.Equal(t, "test", tr.Name)
	assert.Equal(t, "test@test", tr.Email)
}

func TestFormUrlencoded(t *testing.T) {
	r, _ := http.NewRequest("POST", "/", strings.NewReader("name=test&email=test@test"))
	r.Header.Set("content-type", "application/x-www-form-urlencoded")

	type testBaseInternReq struct {
		Name   string `query:"name" form:"name" json:"name"`
		Email  string `query:"email" form:"email" json:"email"`
		intern string
	}

	tr := new(testBaseInternReq)

	err := Form(r, tr)
	assert.Nil(t, err)
	assert.Equal(t, "test", tr.Name)
	assert.Equal(t, "test@test", tr.Email)
}

func TestFormUrlencodedCaseInsensitive(t *testing.T) {
	r, _ := http.NewRequest("POST", "/", strings.NewReader("name=test&NAME=example&EMAIL=test@test"))
	r.Header.Set("content-type", "application/x-www-form-urlencoded")

	tr := new(testBaseReq)

	err := Form(r, tr)
	assert.Nil(t, err)
	assert.Equal(t, "test", tr.Name)
	assert.Equal(t, "test@test", tr.Email)
}

func TestRequestFormUrlencodedAndQuery(t *testing.T) {
	r, _ := http.NewRequest("POST", "/?name=test", strings.NewReader("email=test@test"))
	r.Header.Set("content-type", "application/x-www-form-urlencoded")

	tr := new(testBaseReq)

	err := Request(r, tr)
	assert.Nil(t, err)
	assert.Equal(t, "test", tr.Name)
	assert.Equal(t, "test@test", tr.Email)
}

func TestRequestFormUrlencodedOnlyTypeAndQuery(t *testing.T) {
	r, _ := http.NewRequest("POST", "/?name=test", strings.NewReader("email=test@test"))
	r.Header.Set("content-type", "application/x-www-form-urlencoded")

	type testBaseReqForm struct {
		Name  string `form:"name"`
		Email string `form:"email"`
	}

	tr := new(testBaseReqForm)

	err := Request(r, tr)
	assert.Nil(t, err)
	assert.Equal(t, "", tr.Name)
	assert.Equal(t, "test@test", tr.Email)
}

func TestRequestPostQuery(t *testing.T) {
	r, _ := http.NewRequest("POST", "/?name=test", nil)

	type testBaseReqForm struct {
		Name  string `query:"name"`
		Email string `form:"email"`
	}

	tr := new(testBaseReqForm)

	err := Request(r, tr)
	assert.Nil(t, err)
	assert.Equal(t, "test", tr.Name)
	assert.Equal(t, "", tr.Email)
}

func TestPostQuery(t *testing.T) {
	r, _ := http.NewRequest("POST", "/?name=test", nil)

	type testBaseReqForm struct {
		Name  string `query:"name"`
		Email string `form:"email"`
	}

	tr := new(testBaseReqForm)

	err := Query(r, tr)
	assert.Nil(t, err)
	assert.Equal(t, "test", tr.Name)
	assert.Equal(t, "", tr.Email)
}

func TestRequestFormMultipart(t *testing.T) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	writer.WriteField("name", "test")
	writer.WriteField("email", "test@test")

	writer.Close()

	r, _ := http.NewRequest("POST", "/", body)
	r.Header.Set("content-type", writer.FormDataContentType())

	tr := new(testBaseReq)

	err := Request(r, tr)
	assert.Nil(t, err)
	assert.Equal(t, "test", tr.Name)
	assert.Equal(t, "test@test", tr.Email)
}

func TestRequestFormMultipartCaseInsensitive(t *testing.T) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	writer.WriteField("name", "test")
	writer.WriteField("NAME", "example")
	writer.WriteField("EMAIL", "test@test")

	writer.Close()

	r, _ := http.NewRequest("POST", "/", body)
	r.Header.Set("content-type", writer.FormDataContentType())

	tr := new(testBaseReq)

	err := Request(r, tr)
	assert.Nil(t, err)
	assert.Equal(t, "test", tr.Name)
	assert.Equal(t, "test@test", tr.Email)
}

func TestRequestJSON(t *testing.T) {
	r, _ := http.NewRequest("POST", "/", strings.NewReader(`{"name":"test","email":"test@test"}`))
	r.Header.Set("content-type", "application/json")

	tr := new(testBaseReq)

	err := Request(r, tr)
	assert.Nil(t, err)
	assert.Equal(t, "test", tr.Name)
	assert.Equal(t, "test@test", tr.Email)
}

func TestRequestJSONInvalid(t *testing.T) {
	r, _ := http.NewRequest("POST", "/", strings.NewReader("no json"))
	r.Header.Set("content-type", "application/json")

	tr := new(testBaseReq)

	err := Request(r, tr)
	assert.Error(t, err)
}

func TestRequestJSONInvalidType(t *testing.T) {
	r, _ := http.NewRequest("POST", "/", strings.NewReader(`{"name":123,"email":"test@test"}`))
	r.Header.Set("content-type", "application/json")

	tr := new(testBaseReq)

	err := Request(r, tr)
	assert.Error(t, err)
}

func TestJSON(t *testing.T) {
	r, _ := http.NewRequest("POST", "/", strings.NewReader(`{"name":"test","email":"test@test"}`))
	r.Header.Set("content-type", "application/json")

	tr := new(testBaseReq)

	err := JSON(r, tr)
	assert.Nil(t, err)
	assert.Equal(t, "test", tr.Name)
	assert.Equal(t, "test@test", tr.Email)
}

func TesJSONInvalid(t *testing.T) {
	r, _ := http.NewRequest("POST", "/", strings.NewReader("no json"))
	r.Header.Set("content-type", "application/json")

	tr := new(testBaseReq)

	err := JSON(r, tr)
	assert.Error(t, err)
}

func TestJSONInvalidType(t *testing.T) {
	r, _ := http.NewRequest("POST", "/", strings.NewReader(`{"name":123,"email":"test@test"}`))
	r.Header.Set("content-type", "application/json")

	tr := new(testBaseReq)

	err := JSON(r, tr)
	assert.Error(t, err)
}

func TestRequestAllTypes(t *testing.T) {
	r, _ := http.NewRequest("POST", "/", strings.NewReader(`bool=true&ptrbool=false&float32=0.1&ptrfloat32=0.11&float64=1.2&ptrfloat64=1.22&int=-1&ptrint=-10&int8=2&ptrint8=20&int16=-3&ptrint16=-30&int32=4&ptrint32=40&int64=5&ptrint64=50&string=str&ptrstring=string&uint=1&ptruint=10&uint8=2&ptruint8=20&uint16=3&ptruint16=30&uint32=4&ptruint32=40&uint64=5&ptruint64=50`))
	r.Header.Set("content-type", "application/x-www-form-urlencoded")

	tr := new(testTypesReq)

	err := Request(r, tr)
	assert.NoError(t, err)

	res, _ := json.Marshal(tr)
	assert.Equal(t, `{"bool":true,"ptrbool":false,"float32":0.1,"ptrfloat32":0.11,"float64":1.2,"ptrfloat64":1.22,"int":-1,"ptrint":-10,"int8":2,"ptrint8":20,"int16":-3,"ptrint16":-30,"int32":4,"ptrint32":40,"int64":5,"ptrint64":50,"string":"str","ptrstring":"string","uint":1,"ptruint":10,"uint8":2,"ptruint8":20,"uint16":3,"ptruint16":30,"uint32":4,"ptruint32":40,"uint64":5,"ptruint64":50}`, string(res))
}

func TestRequestAllTypesWithEmpty(t *testing.T) {
	r, _ := http.NewRequest("POST", "/", strings.NewReader(`bool=&ptrbool=false&float32=&ptrfloat32=0.11&float64=1.2&ptrfloat64=1.22&int=-1&ptrint=-10&int8=&ptrint8=20&int16=-3&ptrint16=-30&int32=4&ptrint32=40&int64=5&ptrint64=50&string=str&ptrstring=string&uint=1&ptruint=10&uint8=2&ptruint8=20&uint16=&ptruint16=30&uint32=4&ptruint32=40&uint64=5&ptruint64=50`))
	r.Header.Set("content-type", "application/x-www-form-urlencoded")

	tr := new(testTypesReq)

	err := Request(r, tr)
	assert.NoError(t, err)

	res, _ := json.Marshal(tr)
	assert.Equal(t, `{"bool":false,"ptrbool":false,"float32":0,"ptrfloat32":0.11,"float64":1.2,"ptrfloat64":1.22,"int":-1,"ptrint":-10,"int8":0,"ptrint8":20,"int16":-3,"ptrint16":-30,"int32":4,"ptrint32":40,"int64":5,"ptrint64":50,"string":"str","ptrstring":"string","uint":1,"ptruint":10,"uint8":2,"ptruint8":20,"uint16":0,"ptruint16":30,"uint32":4,"ptruint32":40,"uint64":5,"ptruint64":50}`, string(res))
}

func TestRequestAllTypesWithIntError(t *testing.T) {
	r, _ := http.NewRequest("POST", "/", strings.NewReader(`int=-*--1`))
	r.Header.Set("content-type", "application/x-www-form-urlencoded")

	tr := new(testTypesReq)

	err := Request(r, tr)
	assert.Error(t, err)
}

func TestRequestAllTypesWithFloatError(t *testing.T) {
	r, _ := http.NewRequest("POST", "/", strings.NewReader(`float32=-*--1`))
	r.Header.Set("content-type", "application/x-www-form-urlencoded")

	tr := new(testTypesReq)

	err := Request(r, tr)
	assert.Error(t, err)
}

func TestRequestAllTypesWithUintError(t *testing.T) {
	r, _ := http.NewRequest("POST", "/", strings.NewReader(`uint=-*--1`))
	r.Header.Set("content-type", "application/x-www-form-urlencoded")

	tr := new(testTypesReq)

	err := Request(r, tr)
	assert.Error(t, err)
}

func TestRequestAllTypesWithBoolError(t *testing.T) {
	r, _ := http.NewRequest("POST", "/", strings.NewReader(`bool=-*--`))
	r.Header.Set("content-type", "application/x-www-form-urlencoded")

	tr := new(testTypesReq)

	err := Request(r, tr)
	assert.Error(t, err)
}

func TestRequestNoValue(t *testing.T) {
	r, _ := http.NewRequest("POST", "/", strings.NewReader(`{"name":123,"email":"test@test"}`))
	r.Header.Set("content-type", "application/json")

	err := Request(r, nil)
	assert.NoError(t, err)
}

func TestRequestInvalidValue(t *testing.T) {
	r, _ := http.NewRequest("GET", "/?name=test&email=test@test", nil)

	invalid := "invalid"
	err := Request(r, &invalid)
	assert.Error(t, err)
}

func TestRequestInvalidContentType(t *testing.T) {
	r, _ := http.NewRequest("POST", "/", strings.NewReader(`{"name":123,"email":"test@test"}`))
	r.Header.Set("content-type", "invalid")

	tr := new(testBaseReq)

	err := Request(r, tr)
	assert.Error(t, err)
}

func TestRequestUnknownStructType(t *testing.T) {
	r, _ := http.NewRequest("POST", "/?name=test", strings.NewReader("email=test@test"))
	r.Header.Set("content-type", "application/x-www-form-urlencoded")

	type testBaseReqForm struct {
		Name  string    `form:"name"`
		Email time.Time `form:"email"`
	}

	tr := new(testBaseReqForm)

	err := Request(r, tr)
	assert.Error(t, err)
}

func TestQueryInvalidValue(t *testing.T) {
	r, _ := http.NewRequest("GET", "/?name=test&email=test@test", nil)

	invalid := "invalid"
	err := Query(r, &invalid)
	assert.Error(t, err)
}

func TestFormInvalidValue(t *testing.T) {
	r, _ := http.NewRequest("GET", "/?name=test&email=test@test", nil)

	invalid := "invalid"
	err := Form(r, &invalid)
	assert.Error(t, err)
}

func TestQueryNoValue(t *testing.T) {
	r, _ := http.NewRequest("POST", "/", strings.NewReader(`{"name":123,"email":"test@test"}`))
	r.Header.Set("content-type", "application/json")

	err := Query(r, nil)
	assert.NoError(t, err)
}

func TestFormNoValue(t *testing.T) {
	r, _ := http.NewRequest("POST", "/", strings.NewReader(`{"name":123,"email":"test@test"}`))
	r.Header.Set("content-type", "application/json")

	err := Form(r, nil)
	assert.NoError(t, err)
}

func TestJSONNoValue(t *testing.T) {
	r, _ := http.NewRequest("POST", "/", strings.NewReader(`{"name":123,"email":"test@test"}`))
	r.Header.Set("content-type", "application/json")

	err := JSON(r, nil)
	assert.NoError(t, err)
}

func TestJSONNoBody(t *testing.T) {
	r, _ := http.NewRequest("POST", "/", strings.NewReader(""))
	r.Header.Set("content-type", "application/json")

	type testBaseReqForm struct {
		Name  string    `form:"name"`
		Email time.Time `form:"email"`
	}

	tr := new(testBaseReqForm)

	err := JSON(r, tr)
	assert.NoError(t, err)
}

func TestFormAllTypesWithIntError(t *testing.T) {
	r, _ := http.NewRequest("POST", "/", strings.NewReader(`int=-*--1`))
	r.Header.Set("content-type", "application/x-www-form-urlencoded")

	tr := new(testTypesReq)

	err := Form(r, tr)
	assert.Error(t, err)
}

func TestQueryTypesWithIntError(t *testing.T) {
	r, _ := http.NewRequest("POST", "/?int=-*--1", nil)

	type testQueryTypesReq struct {
		Int    int `query:"int"`
		intern string
	}

	tr := new(testQueryTypesReq)

	err := Query(r, tr)
	assert.Error(t, err)
}

func TestJsonAllTypesWithIntError(t *testing.T) {
	r, _ := http.NewRequest("POST", "/", strings.NewReader(`{"int":"-*--1"}`))

	tr := new(testTypesReq)

	err := JSON(r, tr)
	assert.Error(t, err)
}

func TestQueryUnknownStructType(t *testing.T) {
	r, _ := http.NewRequest("POST", "/?name=test", strings.NewReader("email=test@test"))
	r.Header.Set("content-type", "application/x-www-form-urlencoded")

	type testBaseReqForm struct {
		Name  string    `query:"name"`
		Email time.Time `query:"email"`
	}

	tr := new(testBaseReqForm)

	err := Query(r, tr)
	assert.Error(t, err)
}

func TestFormUnknownStructType(t *testing.T) {
	r, _ := http.NewRequest("POST", "/?name=test", strings.NewReader("email=test@test"))
	r.Header.Set("content-type", "application/x-www-form-urlencoded")

	type testBaseReqForm struct {
		Name  string    `form:"name"`
		Email time.Time `form:"email"`
	}

	tr := new(testBaseReqForm)

	err := Form(r, tr)
	assert.Error(t, err)
}
