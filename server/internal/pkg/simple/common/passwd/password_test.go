package passwd_test

import (
	"fmt"
	"testing"

	"bbs-go/internal/pkg/simple/common/passwd"
)

func TestPassword(t *testing.T) {
	fmt.Println(passwd.GenerateRandomPassword(16))
}
