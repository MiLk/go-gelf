package writer

import (
	"strings"
	"testing"
)

func TestGetCaller(t *testing.T) {
	file, line := getCallerIgnoringLogMulti(1000)
	if line != 0 || file != "???" {
		t.Errorf("didn't fail 1 %s %d", file, line)
		return
	}

	file, _ = getCaller(0)
	if !strings.HasSuffix(file, "/gelf/writer/utils_test.go") {
		t.Errorf("not utils_test.go 1? %s", file)
	}

	file, _ = getCallerIgnoringLogMulti(0)
	if !strings.HasSuffix(file, "/gelf/writer/utils_test.go") {
		t.Errorf("not utils_test.go 2? %s", file)
	}
}
