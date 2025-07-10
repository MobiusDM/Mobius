package dep_test

import (
	"testing"

	"github.com/notawar/mobius/backend/server/archtest/test_files/testfiledeps/testpkgdependency"
)

func Test(t *testing.T) {
	testpkgdependency.OohNoBadCode()
}
