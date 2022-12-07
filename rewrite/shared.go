package track

import (
	"fmt"
)

func fmtMsg(msg string, args ...any) string {
	return fmt.Sprintf(msg, args...)
}
