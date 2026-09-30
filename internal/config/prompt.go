package config

import (
	"bytes"
	"text/template"
)

type PromptContext struct {
	PetName   string
	Character string
	UserEvent string
}

type PromptGenerator struct {
	systemPrompt string
	tmpl		 *template.Template
}

const DEFAULT_SHIMEJI_PROMPT = `You are {{.PetName}}, a {{.Character}}.
Respond in character to the following event or poke.
Always reply with strict JSON using this schema:
{
  "subtitle_en": "English translation or companion thought",
  "animation": "idle" | "pout" | "shocked" | "happy"
}

Event: {{.UserEvent}}`

func NewPromptGenerator(systemPrompt string) (*PromptGenerator) {
	base := systemPrompt
	if base == "" {
		base = DEFAULT_SHIMEJI_PROMPT
	}
	
	tmpl, err := template.New("prompt").Parse(base)
	if err != nil {
		tmpl = template.Must(template.New("prompt").Parse(DEFAULT_SHIMEJI_PROMPT))
	}

	return &PromptGenerator{
		systemPrompt: base,
		tmpl:         tmpl,
	}
}

func (pg *PromptGenerator) BuildPrompt(ctx PromptContext) (string, error) {
	var buf bytes.Buffer
	if err := pg.tmpl.Execute(&buf, ctx); err != nil {
		return "", err
	}
	
	return buf.String(), nil
}
