package llm

import (
	"bytes"
	_ "embed"
	"fmt"
	"log/slog"
	"strings"
	"text/template"
	"time"
)

// 1st prompt: https://makersuite.google.com/app/prompts/1JtpmT6Efbsg9S-PgxTvAsMbDL_hTEo5F?pli=1

//go:embed prompts/chat.prompt
var ChatPromptTemplate string

type PromptId int

const (
	PROMPT_CHAT PromptId = iota
)

func ChatPrompt(query string, context []string) (string, error) {
	c := strings.Join(context, "\n")
	now := time.Now()
	todaysDate := now.Format("Monday January 2, 2006")
	prompt, err := Prompt(PROMPT_CHAT, map[string]string{
		"Query":      query,
		"Context":    c,
		"TodaysDate": todaysDate})

	if err != nil {
		fmt.Printf("error '%s' creating prompt for: '%s", err.Error(), query)
	}
	return prompt, err
}

func Prompt(id PromptId, data map[string]string) (string, error) {
	switch id {
	case PROMPT_CHAT:
		tmpl, err := template.New("chat").Parse(ChatPromptTemplate)
		if err != nil {
			return "", err
		}
		buf := new(bytes.Buffer)
		if err = tmpl.Execute(buf, data); err != nil {
			return "", err
		}

		result := buf.String()
		slog.Info(fmt.Sprintf("using chatprompt: \n%s\n===============================", result))
		return result, nil
	default:
		return "", fmt.Errorf("unknown prompt id")
	}
}
