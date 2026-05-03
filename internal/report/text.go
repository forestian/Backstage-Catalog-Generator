package report

import (
	"fmt"
	"io"

	"github.com/example/backstage-catalog-generator/internal/validate"
)

// PrintResult writes validation errors and warnings to w.
func PrintResult(w io.Writer, res *validate.Result) {
	for _, warn := range res.Warnings {
		fmt.Fprintf(w, "WARNING: %s\n", warn)
	}
	for _, err := range res.Errors {
		fmt.Fprintf(w, "ERROR:   %s\n", err)
	}
}
