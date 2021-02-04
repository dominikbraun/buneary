package buneary

import "github.com/spf13/cobra"

// globalOptions defines global command line options available for all commands.
// They're read by the top-level command and passed to the sub-command factories.
type globalOptions struct {
	user     string
	password string
}

// RootCommand creates the top-level `buneary` command without any functionality.
func RootCommand() *cobra.Command {
	var options globalOptions

	root := &cobra.Command{
		Use:   "buneary",
		Short: "An easy-to-use CLI client for RabbitMQ.",
		Long: `buneary, pronounced bun-ear-y, is an easy-to-use RabbitMQ command line client
for managing exchanges, managing queues and publishing messages to exchanges.`,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	root.AddCommand(publishCommand(&options))

	root.PersistentFlags().
		StringVarP(&options.user, "user", "u", "", "the username to connect with")
	root.PersistentFlags().
		StringVarP(&options.password, "password", "p", "", "the password to authenticate with")

	return root
}

// publishCommand creates the `buneary publish` command, making sure that exactly
// four command arguments are passed.
func publishCommand(options *globalOptions) *cobra.Command {
	publish := &cobra.Command{
		Use:   "publish <ADDRESS> <EXCHANGE> <ROUTING KEY> <BODY>",
		Short: "Publish a message to an exchange",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runPublishCommand(options, args)
		},
	}

	return publish
}

// runPublishCommand publishes a command by reading the command line data, setting
// the configuration and calling the PublishMessage function. In case the password
// or both the user and password aren't provided, it will go into interactive mode.
func runPublishCommand(options *globalOptions, args []string) error {
	var (
		address    = args[0]
		exchange   = args[1]
		routingKey = args[2]
		body       = args[3]
	)

	buneary := buneary{
		config: &AMQPConfig{
			Address:  address,
			User:     options.user,
			Password: options.password,
		},
	}

	message := Message{
		Target:     Exchange{Name: exchange},
		RoutingKey: routingKey,
		Body:       []byte(body),
	}

	if err := buneary.PublishMessage(message); err != nil {
		return err
	}

	return nil
}
