package main

import (
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
)

type model struct {
	config   Config
	err      error
	textarea textarea.Model
}

func New(cfg Config) model {
	ta := textarea.New()
	ta.Placeholder = "Enter your message here..."
	ta.Focus()

	ta.Prompt = "| "
	ta.CharLimit = int(cfg.MaxLength)

	ta.SetHeight(3)

	return model{
		config:   cfg,
		err: nil,
		textarea: ta,
	}
}

