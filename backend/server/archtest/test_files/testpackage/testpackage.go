package testpackage

import (
	"crypto" // for test
	"fmt"

	"github.com/notawar/mobius/backend/server/archtest/test_files/dependency"
)

func What(_ crypto.Decrypter) {
	fmt.Println(dependency.Item)
}
