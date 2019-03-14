package main

import (
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/iam"
	survey "gopkg.in/AlecAivazis/survey.v1"
)

func (client *iamClient) prompt() string {
	users := client.getUserNames()

	answer := promptForSelection("Which user's password needs to be reset?", users, "IAM username")
	return answer
}

func (client *iamClient) getUserNames() []string {
	// if client.svc == nil {
	// 	client.svc = client.newIamClient()
	// }

	input := &iam.ListUsersInput{}
	usernames := []string{}
	ctx := aws.BackgroundContext()
	err := client.svc.ListUsersPagesWithContext(ctx, input,
		func(page *iam.ListUsersOutput, lastPage bool) bool {
			for _, user := range page.Users {
				usernames = append(usernames, *user.UserName)
			}
			return true
		}, func(r *request.Request) {

		})
	if err != nil {
		log.Panic(err.Error())
		return nil
	}
	return usernames
}

func promptForSelection(message string, options []string, help string) string {
	var surveyTemplateOriginal = survey.SelectQuestionTemplate
	survey.SelectQuestionTemplate = `
		{{- if .ShowHelp }}{{- color "225"}}{{ HelpIcon }} {{ .Help }}{{color "reset"}}{{"\n"}}{{end}}
		{{- color "yellow+hb"}}{{ QuestionIcon }} {{color "reset"}}
		{{- color "default+h"}}{{ .Message }}{{ .FilterMessage }}{{color "reset"}}
		{{- if .ShowAnswer}}{{color "159"}} {{.Answer}}{{color "reset"}}{{"\n"}}
		{{- else}}
		{{- "  "}}{{- color "159"}}[Use arrows to move, type to filter{{- if and .Help (not .ShowHelp)}}, {{ HelpInputRune }} for more help{{end}}]{{color "reset"}}
		{{- "\n"}}
		{{- range $ix, $choice := .PageEntries}}
			{{- if eq $ix $.SelectedIndex}}{{color "156"}}{{ SelectFocusIcon }} {{else}}{{color "default+h"}}  {{end}}
			{{- $choice}}
			{{- color "reset"}}{{"\n"}}
		{{- end}}
		{{- end}}`
	defer func() {
		survey.SelectQuestionTemplate = surveyTemplateOriginal
	}()

	s := &survey.Select{
		Message: message,
		Options: options,
		Help:    help,
	}

	answer := ""
	err := survey.AskOne(s, &answer, nil)
	if err != nil {
		panic(err)
	}
	if !confirmSelection(answer) {
		os.Exit(1)
	}

	return answer
}

func confirmSelection(user string) bool {
	var surveyConfirmQuestionOriginal = survey.ConfirmQuestionTemplate
	survey.ConfirmQuestionTemplate = `
	{{- if .ShowHelp }}{{- color "225"}}{{ HelpIcon }} {{ .Help }}{{color "reset"}}{{"\n"}}{{end}}
	{{- color "yellow+hb"}}{{ QuestionIcon }} {{color "reset"}}
	{{- color "default+hb"}}{{ .Message }} {{color "reset"}}
	{{- if .Answer}}
	{{- color "159"}}{{.Answer}}{{color "reset"}}{{"\n"}}
	{{- else }}
	{{- if and .Help (not .ShowHelp)}}{{color "225"}}[{{ HelpInputRune }} for help]{{color "reset"}} {{end}}
	{{- color "156"}}{{if .Default}}(Y/n) {{else}}(y/N) {{end}}{{color "reset"}}
	{{- end}}`
	defer func() {
		survey.ConfirmQuestionTemplate = surveyConfirmQuestionOriginal
	}()

	message := fmt.Sprintf("Are you sure you want to reset %s's password?", user)

	b := &survey.Confirm{
		Message: message,
		Default: false,
	}

	var validate bool
	err := survey.AskOne(b, &validate, nil)
	if err != nil {
		panic(err)
	}
	return validate
}
