package v3

import (
	"fmt"
	"strings"
)

func (a ApiErrorResponseBody) String() string {
	var msg strings.Builder
	msg.WriteString("--- error ---\n")
	msg.WriteString(a.Error.String())
	return msg.String()
}

func (a ApiError) String() string {
	var msg strings.Builder
	msg.WriteString(fmt.Sprintf("Request: %s\n", a.RequestId))
	msg.WriteString(fmt.Sprintf("Code:    %s\n", a.Code))
	msg.WriteString(fmt.Sprintf("Msg:     %s\n", a.Message))
	if a.DebugMessage != nil && *a.DebugMessage != a.Message {
		msg.WriteString(fmt.Sprintf("Debug:   %s\n", *a.DebugMessage))
	}
	if a.Details != nil {
		msg.WriteString("Details:\n")
		for _, detail := range *a.Details {
			msg.WriteString(fmt.Sprintf("  - Code:    %s\n", detail.Code))
			msg.WriteString(fmt.Sprintf("    Message: %s\n", detail.Message))
			if detail.DebugMessage != nil && *detail.DebugMessage != detail.Message {
				msg.WriteString(fmt.Sprintf("    Debug:   %s\n", *detail.DebugMessage))
			}
		}
	}
	return msg.String()
}
