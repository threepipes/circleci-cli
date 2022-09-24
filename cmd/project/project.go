package project

import (
	projectapi "github.com/CircleCI-Public/circleci-cli/api/project"
	"github.com/CircleCI-Public/circleci-cli/cmd/validator"
	"github.com/CircleCI-Public/circleci-cli/prompt"

	"github.com/CircleCI-Public/circleci-cli/settings"
	"github.com/spf13/cobra"
)

type Reader interface {
	ReadSecretString(msg string) (string, error)
	AskConfirm(msg string) bool
}

type projectOpts struct {
	client projectapi.ProjectClient
	reader Reader
}

type ProjectOption interface {
	apply(*projectOpts)
}

type PromptReader struct{}

func (p PromptReader) ReadSecretString(msg string) (string, error) {
	return prompt.ReadSecretStringFromUser(msg)
}

func (p PromptReader) AskConfirm(msg string) bool {
	return prompt.AskUserToConfirm(msg)
}

// NewProjectCommand generates a cobra command for managing projects
func NewProjectCommand(config *settings.Config, preRunE validator.Validator, opts ...ProjectOption) *cobra.Command {
	pos := projectOpts{
		reader: &PromptReader{},
	}
	for _, o := range opts {
		o.apply(&pos)
	}
	command := &cobra.Command{
		Use:   "project",
		Short: "Operate on projects",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			client, err := projectapi.NewProjectRestClient(*config)
			if err != nil {
				return err
			}
			pos.client = client
			return nil
		},
	}

	command.AddCommand(newProjectEnvironmentVariableCommand(&pos, preRunE))

	return command
}

type CustomReaderProjectOption struct {
	r Reader
}

func (c CustomReaderProjectOption) apply(opts *projectOpts) {
	opts.reader = c.r
}

func CustomReader(r Reader) ProjectOption {
	return CustomReaderProjectOption{r}
}
