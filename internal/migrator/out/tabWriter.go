package out

import (
	"fmt"
	"log/slog"
	"os"
	"text/tabwriter"
)

type TabWriter struct {
	tw *tabwriter.Writer
}

func NewTabWriter() *TabWriter {
	stdOutWriter := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
	tabwriterStruct := TabWriter{
		tw: stdOutWriter,
	}

	return &tabwriterStruct
}

func (tw *TabWriter) Printf(format string, a ...any) error {
	_, err := fmt.Fprintf(tw.tw, format, a...)
	if err != nil {
		return err
	}
	return nil
}

func (tw *TabWriter) Flush() {
	err := tw.tw.Flush()
	if err != nil {
		slog.Error("%s", err)
	}
}
