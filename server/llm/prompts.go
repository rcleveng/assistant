package llm

import (
	"bytes"
	"fmt"
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
{{ CONTEXT }}

USERQUESTION: {{ QUERY }}
`
}

type PromptId int

const (
	PROMPT_CHAT PromptId = iota
)

func ChatPrompt(query string, context []string) (string, error) {
	c := strings.Join(context, "\n")
	return Prompt(PROMPT_CHAT, map[string]string{
		"QUERY": query, "CONTEXT": c})
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

		return buf.String(), nil
	default:
		return "", fmt.Errorf("unknown prompt id")
	}
}
