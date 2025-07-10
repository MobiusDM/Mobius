package dep

import (
	"testing"

	"github.com/notawar/mobius/backend/server/archtest/test_files/testfiledeps/testonlydependency"
)

func TestDoIBreakYou(t *testing.T) {
	testonlydependency.OohNoBadCode()
}
