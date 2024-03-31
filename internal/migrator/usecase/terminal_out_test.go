package usecase_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/flatmix/final-otus-project/internal/migrator/usecase"
)

func Test_TerminalOutOK(t *testing.T) {
	outs := usecase.Outs{
		&usecase.Out{
			Name:   "test1",
			Status: "migrate ok",
		}, &usecase.Out{
			Name:   "test2",
			Status: "migrate ok",
		},
	}

	err := usecase.TerminalOut(&outs)
	assert.NoError(t, err)
}

func Test_TerminalStatusOutOK(t *testing.T) {
	now := time.Now().Format(time.RFC3339)

	outs := usecase.Outs{
		&usecase.Out{
			Name:        "test3",
			Status:      "migrate ok",
			Version:     "1",
			TimeMigrate: now,
		}, &usecase.Out{
			Name:        "test4",
			Status:      "migrate ok",
			Version:     "2",
			TimeMigrate: now,
		},
	}

	err := usecase.TerminalStatusOut(&outs)
	assert.NoError(t, err)
}
