package usecase

import (
	"fmt"

	"github.com/flatmix/final-otus-project/internal/migrator/out"
)

func TerminalStatusOut(statusStructs *Outs) error {
	terminalOut := out.NewTabWriter()
	defer terminalOut.Flush()

	for _, status := range *statusStructs {
		err := terminalOut.
			Printf("%s \t| status: \t %s \t| Vesion: \t %s \t| Date: \t %s \t| \n",
				status.Name, status.Status, status.Version, status.TimeMigrate)
		if err != nil {
			return fmt.Errorf("terminalOut: %w", err)
		}
	}

	return nil
}

func TerminalOut(outs *Outs) error {
	terminalOut := out.NewTabWriter()
	defer terminalOut.Flush()

	for _, out := range *outs {
		err := terminalOut.Printf("%s \t| status: %s \t| \n", out.Name, out.Status)
		if err != nil {
			return fmt.Errorf("terminalOut: %w", err)
		}
	}

	return nil
}
