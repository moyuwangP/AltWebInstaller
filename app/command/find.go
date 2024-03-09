package command

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
)

func FindBinary(ctx context.Context, args ...string) (string, error) {
	var out bytes.Buffer

	c := exec.CommandContext(ctx, "which", args...)
	c.Stdout = &out

	err := c.Run()
	if err != nil {
		return "", fmt.Errorf("%s: %s", c.String(), err)
	}

	return strings.TrimSpace(out.String()), nil
}
