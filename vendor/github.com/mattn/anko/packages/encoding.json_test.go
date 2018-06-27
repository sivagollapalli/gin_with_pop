package packages

import (
	"os"
	"testing"

	"github.com/mattn/anko/internal/testlib"
)

func TestJson(t *testing.T) {
	os.Setenv("ANKO_DEBUG", "1")
	var toByteSlice = func(s string) []byte { return []byte(s) }
	tests := []testlib.Test{
		{Script: `json = import("encoding/json"); a = make(mapStringInterface); a["b"] = "b"; c, err = json.Marshal(a); err`, Types: map[string]interface{}{"mapStringInterface": map[string]interface{}{}}, Output: map[string]interface{}{"a": map[string]interface{}{"b": "b"}, "c": []byte(`{"b":"b"}`)}},
		{Script: `json = import("encoding/json"); b = 1; err = json.Unmarshal(a, &b); err`, Input: map[string]interface{}{"a": []byte(`{"b": "b"}`)}, Output: map[string]interface{}{"a": []byte(`{"b": "b"}`), "b": map[string]interface{}{"b": "b"}}},
		{Script: `json = import("encoding/json"); b = 1; err = json.Unmarshal(toByteSlice(a), &b); err`, Input: map[string]interface{}{"a": `{"b": "b"}`, "toByteSlice": toByteSlice}, Output: map[string]interface{}{"a": `{"b": "b"}`, "b": map[string]interface{}{"b": "b"}}},
		{Script: `json = import("encoding/json"); b = 1; err = json.Unmarshal(toByteSlice(a), &b); err`, Input: map[string]interface{}{"a": `[["1", "2"],["3", "4"]]`, "toByteSlice": toByteSlice}, Output: map[string]interface{}{"a": `[["1", "2"],["3", "4"]]`, "b": []interface{}{[]interface{}{"1", "2"}, []interface{}{"3", "4"}}}},
	}
	testlib.Run(t, tests, &testlib.Options{EnvSetupFunc: &testPackagesEnvSetupFunc})
}
