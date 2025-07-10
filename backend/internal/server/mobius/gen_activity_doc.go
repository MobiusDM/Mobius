//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/notawar/mobius/internal/server/mobius"
)

func main() {
	var b strings.Builder

	b.WriteString(`<!-- DO NOT EDIT. This document is automatically generated. -->
# Audit logs

Mobius logs activities.

To see activities in Mobius, select the Mobius icon in the top navigation and see the **Activity** section.

This page includes a list of activities.

`)

	activityMap := map[string]struct{}{}
	for _, activity := range mobius.ActivityDetailsList {
		if _, ok := activityMap[activity.ActivityName()]; ok {
			panic(fmt.Sprintf("type %s already used", activity.ActivityName()))
		}
		activityMap[activity.ActivityName()] = struct{}{}

		fmt.Fprintf(&b, "## %s\n\n", activity.ActivityName())
		activityTypeDoc, detailsDoc, detailsExampleDoc := activity.Documentation()
		fmt.Fprintf(&b, activityTypeDoc+"\n\n"+detailsDoc+"\n\n")
		if detailsExampleDoc != "" {
			fmt.Fprintf(&b, "#### Example\n\n```json\n%s\n```\n\n", detailsExampleDoc)
		}
	}
	b.WriteString(`
<meta name="title" value="Audit logs">
<meta name="pageOrderInSection" value="1400">
<meta name="description" value="Learn how Mobius logs administrative actions in JSON format.">
<meta name="navSection" value="Dig deeper">
`)

	if err := os.WriteFile(os.Args[1], []byte(b.String()), 0600); err != nil {
		panic(err)
	}
}
