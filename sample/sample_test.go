package sample

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"
	"strings"
	"testing"
)

func TestSample(t *testing.T) {

	reqPath := "testdata/request"
	expPath := "testdata/expected"

	fis, err := ioutil.ReadDir(reqPath)
	if err != nil {
		t.Fatalf("cannot read testdata dir (err=%s)", err.Error())
	}

	for _, fi := range fis {
		if fi.IsDir() || strings.Index(fi.Name(), ".json") < 0 {
			continue
		}

		jsonb, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", reqPath, fi.Name()))
		if err != nil {
			t.Fatalf("can not read request file (file=%s, err=%s)", fi.Name(), err.Error())
		}

		req := &Request{}
		if json.Unmarshal(jsonb, req) != nil {
			t.Fatalf("can not unmarshal request json (file=%s, err=%s)", fi.Name(), err.Error())
		}

		expJSONb, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", expPath, fi.Name()))
		if err != nil {
			t.Fatalf("can not read expected file (file=%s, err=%s)", fi.Name(), err.Error())
		}

		exp := &Expected{}
		if err := json.Unmarshal(expJSONb, &exp); err != nil {
			t.Fatalf("can not unmarshal expected json (file=%s, err=%s)", fi.Name(), err.Error())
		}

		r := Sample(req)

		if !reflect.DeepEqual(r, exp) {
			t.Errorf("result is not match (file=%s, expected=%v, actual=%v)", fi.Name(), exp, r)
		}
	}
}
