package dates_test

import (
	"fmt"
	"testing"
	"time"

	"bbs-go/internal/pkg/simple/common/dates"
)

func TestWithTimeAsEndOfDay(t *testing.T) {
	fmt.Println(dates.Timestamp(dates.WithTimeAsEndOfDay(time.Now())))
}
