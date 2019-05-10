package sample

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestSample(t *testing.T) {

	reqPath := "testdata/request"
	expPath := "testdata/expected"

	fis, err := ioutil.ReadDir(reqPath)
	if err != nil {
		t.Fatalf("cannot read testdata dir (err=%s)", err.Error())
	}

	for _, fi := range fis {
		idxJSON := strings.Index(fi.Name(), ".json")
		if fi.IsDir() || idxJSON < 0 {
			continue
		}

		testname := fi.Name()[:idxJSON]
		t.Run(testname, func(t *testing.T) {

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
			diff := cmp.Diff(exp, r)
			if diff != "" {
				t.Errorf("result is not match (diff=%s)", diff)
			}
		})
	}
}
