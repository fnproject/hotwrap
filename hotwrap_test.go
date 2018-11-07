package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/fnproject/fdk-go"
)

type testContext struct {}

func (c testContext) Config() map[string]string { return nil }
func (c testContext) Header() http.Header {
	hs := http.Header{}
	hs.Set("Fn-call-id", "some-blah")

	return hs
}
func (c testContext) AppID() string { return "blah-app" }
func (c testContext) CallID() string { return "blah-app" }
func (c testContext) FnID() string { return "blah-app" }
func (c testContext) ContentType() string { return "application/json" }

type Person struct {
	Name string `json:"name"`
}

func TestHotWrap(t *testing.T) {
	var in, out bytes.Buffer
	cmd := "echo"
	expectedPerson := Person{"John"}
	json.NewEncoder(&in).Encode(expectedPerson)

	ctx := fdk.WithContext(context.Background(), testContext{})
	err := runExec(ctx, cmd, []string{in.String(),}, nil, &out)
	if err != nil {
		t.Fatalf(err.Error())
	}

	var actualPerson Person
	err = json.NewDecoder(&out).Decode(&actualPerson)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if expectedPerson.Name != actualPerson.Name {
		t.Fatalf("Output content mismatch!"+
			"\n\tExpected: %v"+
			"\n\tActual: %v", expectedPerson.Name, actualPerson.Name)
	}
}
