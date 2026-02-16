package idgen

import "errors"

var ErrInvalidWorkerID = errors.New("snowflake workerId must be between 0 and 1023")
