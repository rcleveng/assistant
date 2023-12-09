package llm

import (
	"bytes"
	"fmt"
	"log/slog"
	"strings"
	"text/template"
)

// 1st prompt: https://makersuite.google.com/app/prompts/1JtpmT6Efbsg9S-PgxTvAsMbDL_hTEo5F?pli=1
/*
	Example:

You are a helpful assistant please respond  to USERQUESTION with one of the following:

if you can answer the question please respond with:
ANSWER: The answer to the question

If you are asked to remember something, please respond with"
REMEMBER: The text you are asked to remember

If you need more information please respond with:
CALENDAR: I need to look up the calendar on day $DAY
NEEDMORE: I need more information, please ask for what information is needed to answer the question.

Use the following additional information to help answer if needed:
CONTEXT:
Remembered:  My kid's birthday is August 2nd
Calendar: First day of winter break December 19th, 2023

USERQUESTION:
*/
var ChatPromptTemplate string

func init() {
	ChatPromptTemplate = `
You are a helpful assistant please respond to USERQUESTION with one of the following:

if you can answer the question please respond with:
ANSWER: The answer to the question

If you are asked to remember something, please respond with"
REMEMBER: The text you are asked to remember

If you need more information please respond with:
CALENDAR: I need to look up the calendar on day $DAY
NEEDMORE: I need more information, please ask for what information is needed to answer the question.

Use the following additional information to help answer if needed:
CONTEXT:
{{ .Context}}

USERQUESTION: {{ .Query }}
`
}

type PromptId int

const (
	PROMPT_CHAT PromptId = iota
)

func ChatPrompt(query string, context []string) (string, error) {
	c := strings.Join(context, "\n")
	prompt, err := Prompt(PROMPT_CHAT, map[string]string{
		"Query": query, "Context": c})

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
