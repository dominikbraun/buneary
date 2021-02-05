package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

var version = "UNDEFINED"

// globalOptions defines global command line options available for all commands.
// They're read by the top-level command and passed to the sub-command factories.
type globalOptions struct {
	user     string
	password string
}

// rootCommand creates the top-level `buneary` command without any functionality.
func rootCommand() *cobra.Command {
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

	root.AddCommand(createCommand(&options))
	root.AddCommand(publishCommand(&options))
	root.AddCommand(deleteCommand(&options))
	root.AddCommand(versionCommand())

	root.PersistentFlags().
		StringVarP(&options.user, "user", "u", "", "the username to connect with")
	root.PersistentFlags().
		StringVarP(&options.password, "password", "p", "", "the password to authenticate with")

	return root
}

// createCommand creates the `buneary create` command without any functionality.
func createCommand(options *globalOptions) *cobra.Command {
	create := &cobra.Command{
		Use:   "create <COMMAND>",
		Short: "Create a resource",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	create.AddCommand(createExchangeCommand(options))
	create.AddCommand(createQueueCommand(options))
	create.AddCommand(createBindingCommand(options))

	return create
}

// createExchangeOptions defines options for creating a new exchange.
type createExchangeOptions struct {
	*globalOptions
	durable    bool
	autoDelete bool
	internal   bool
	noWait     bool
}

// createExchangeCommand creates the `buneary create exchange` command, making sure
// that exactly three arguments are passed.
//
// At the moment, there is no support for setting Exchange.NoWait via this command.
func createExchangeCommand(options *globalOptions) *cobra.Command {
	createExchangeOptions := &createExchangeOptions{
		globalOptions: options,
	}

	createExchange := &cobra.Command{
		Use:   "exchange <ADDRESS> <NAME> <TYPE>",
		Short: "Create a new exchange",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCreateExchange(createExchangeOptions, args)
		},
	}

	createExchange.Flags().
		BoolVar(&createExchangeOptions.durable, "durable", false, "make the exchange durable")
	createExchange.Flags().
		BoolVar(&createExchangeOptions.autoDelete, "auto-delete", false, "make the exchange auto-deleted")
	createExchange.Flags().
		BoolVar(&createExchangeOptions.internal, "internal", false, "make the exchange internal")

	return createExchange
}

// runCreateExchange creates a new exchange by reading the command line data, setting
// the configuration and calling the runCreateExchange function. In case the password
// or both the user and password aren't provided, it will go into interactive mode.
//
// ToDo: Move the logic for parsing the exchange type into Exchange.
func runCreateExchange(options *createExchangeOptions, args []string) error {
	var (
		address      = args[0]
		name         = args[1]
		exchangeType = args[2]
	)

	user, password := getOrReadInCredentials(options.globalOptions)

	buneary := buneary{
		config: &AMQPConfig{
			Address:  address,
			User:     user,
			Password: password,
		},
	}

	exchange := Exchange{
		Name:       name,
		Durable:    options.durable,
		AutoDelete: options.autoDelete,
		Internal:   options.internal,
		NoWait:     options.noWait,
	}

	switch exchangeType {
	case "direct":
		exchange.Type = Direct
	case "headers":
		exchange.Type = Headers
	case "fanout":
		exchange.Type = Fanout
	case "topic":
		exchange.Type = Topic
	}

	if err := buneary.CreateExchange(exchange); err != nil {
		return err
	}

	return nil
}

// createQueueOptions defines options for creating a new queue.
type createQueueOptions struct {
	*globalOptions
	durable    bool
	autoDelete bool
}

// createQueueCommand creates the `buneary create queue` command, making sure that
// exactly three arguments are passed.
//
// The <TYPE> argument may become optional for convenience in the future. In this
// case, it should default to the classic queue type.
func createQueueCommand(options *globalOptions) *cobra.Command {
	createQueueOptions := &createQueueOptions{
		globalOptions: options,
	}

	createQueue := &cobra.Command{
		Use:   "queue <ADDRESS> <NAME> <TYPE>",
		Short: "Create a new queue",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCreateQueue(createQueueOptions, args)
		},
	}

	createQueue.Flags().
		BoolVar(&createQueueOptions.durable, "durable", false, "make the queue durable")
	createQueue.Flags().
		BoolVar(&createQueueOptions.autoDelete, "auto-delete", false, "make the queue auto-deleted")

	return createQueue
}

// runCreateQueue creates a new queue by reading the command line data, setting the
// configuration and calling the CreateQueue function. In case the password or both
// the user and password aren't provided, it will go into interactive mode.
//
// If the queue type is empty or invalid, the queue type defaults to Classic.
func runCreateQueue(options *createQueueOptions, args []string) error {
	var (
		address   = args[0]
		name      = args[1]
		queueType = args[2]
	)

	user, password := getOrReadInCredentials(options.globalOptions)

	buneary := buneary{
		config: &AMQPConfig{
			Address:  address,
			User:     user,
			Password: password,
		},
	}

	queue := Queue{
		Name:       name,
		Durable:    options.durable,
		AutoDelete: options.autoDelete,
	}

	switch queueType {
	case "quorum":
		queue.Type = Quorum
	case "classic":
		fallthrough
	default:
		queue.Type = Classic
	}

	_, err := buneary.CreateQueue(queue)
	if err != nil {
		return err
	}

	return nil
}

// createBindingOptions defines options for creating a new binding.
type createBindingOptions struct {
	*globalOptions
	toExchange bool
}

// createBindingCommand creates the `buneary create binding` command, making sure
// that exactly four arguments are passed.
func createBindingCommand(options *globalOptions) *cobra.Command {
	createBindingOptions := &createBindingOptions{
		globalOptions: options,
	}

	createQueue := &cobra.Command{
		Use:   "binding <ADDRESS> <NAME> <TARGET> <BINDING KEY>",
		Short: "Create a new queue",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCreateBinding(createBindingOptions, args)
		},
	}

	createQueue.Flags().
		BoolVar(&createBindingOptions.toExchange, "to-exchange", false, "the target is another exchange")

	return createQueue
}

// runCreateBinding creates a new binding by reading the command line data, setting
// the configuration and calling the CreateQueue function. In case the password or
// both the user and password aren't provided, it will go into interactive mode.
//
// The binding type defaults to ToQueue. To create a binding to another exchange, the
// --to-exchange flag has to be used.
func runCreateBinding(options *createBindingOptions, args []string) error {
	var (
		address    = args[0]
		name       = args[1]
		target     = args[2]
		bindingKey = args[3]
	)

	user, password := getOrReadInCredentials(options.globalOptions)

	buneary := buneary{
		config: &AMQPConfig{
			Address:  address,
			User:     user,
			Password: password,
		},
	}

	binding := Binding{
		From:       Exchange{Name: name},
		TargetName: target,
		Key:        bindingKey,
	}

	switch options.toExchange {
	case true:
		binding.Type = ToExchange
	default:
		binding.Type = ToQueue
	}

	if err := buneary.CreateBinding(binding); err != nil {
		return err
	}

	return nil
}

// publishCommand creates the `buneary publish` command, making sure that exactly
// four command arguments are passed.
func publishCommand(options *globalOptions) *cobra.Command {
	publish := &cobra.Command{
		Use:   "publish <ADDRESS> <EXCHANGE> <ROUTING KEY> <BODY>",
		Short: "Publish a message to an exchange",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runPublish(options, args)
		},
	}

	return publish
}

// runPublish publishes a message by reading the command line data, setting the
// configuration and calling the PublishMessage function. In case the password or
// both the user and password aren't provided, it will go into interactive mode.
func runPublish(options *globalOptions, args []string) error {
	var (
		address    = args[0]
		exchange   = args[1]
		routingKey = args[2]
		body       = args[3]
	)

	user, password := getOrReadInCredentials(options)

	buneary := buneary{
		config: &AMQPConfig{
			Address:  address,
			User:     user,
			Password: password,
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

// deleteCommand creates the `buneary delete` command without any functionality.
func deleteCommand(options *globalOptions) *cobra.Command {
	delete := &cobra.Command{
		Use:   "delete <COMMAND>",
		Short: "Delete a resource",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	delete.AddCommand(deleteExchangeCommand(options))
	delete.AddCommand(deleteQueueCommand(options))

	return delete
}

// deleteExchangeCommand creates the `buneary delete exchange` command, making sure
// that exactly two arguments are passed.
func deleteExchangeCommand(options *globalOptions) *cobra.Command {
	deleteExchange := &cobra.Command{
		Use:   "exchange <ADDRESS> <NAME>",
		Short: "Delete an exchange",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDeleteExchange(options, args)
		},
	}

	return deleteExchange
}

// runDeleteExchange deletes an exchange by reading the command line data, setting the
// configuration and calling the DeleteExchange function. In case the password or
// both the user and password aren't provided, it will go into interactive mode.
func runDeleteExchange(options *globalOptions, args []string) error {
	var (
		address = args[0]
		name    = args[1]
	)

	user, password := getOrReadInCredentials(options)

	buneary := buneary{
		config: &AMQPConfig{
			Address:  address,
			User:     user,
			Password: password,
		},
	}

	exchange := Exchange{
		Name: name,
	}

	if err := buneary.DeleteExchange(exchange); err != nil {
		return err
	}

	return nil
}

// deleteQueueCommand creates the `buneary delete queue` command, making sure
// that exactly two arguments are passed.
func deleteQueueCommand(options *globalOptions) *cobra.Command {
	deleteExchange := &cobra.Command{
		Use:   "queue <ADDRESS> <NAME>",
		Short: "Delete a queue",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDeleteQueue(options, args)
		},
	}

	return deleteExchange
}

// runDeleteQueue deletes a queue by reading the command line data, setting the
// configuration and calling the DeleteQueue function. In case the password or
// both the user and password aren't provided, it will go into interactive mode.
func runDeleteQueue(options *globalOptions, args []string) error {
	var (
		address = args[0]
		name    = args[1]
	)

	user, password := getOrReadInCredentials(options)

	buneary := buneary{
		config: &AMQPConfig{
			Address:  address,
			User:     user,
			Password: password,
		},
	}

	queue := Queue{
		Name: name,
	}

	_, err := buneary.DeleteQueue(queue)
	if err != nil {
		return err
	}

	return nil
}

// versionCommand creates the `buneary version` command for printing release
// information. This data is injected by the CI pipeline.
func versionCommand() *cobra.Command {
	version := &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("buneary %s", version)
			return nil
		},
	}

	return version
}

// getOrReadInCredentials either returns the credentials directly from the global
// options or prompts the user to type them in.
//
// If both user and password have been set using the --user and --password flags,
// those values will be used. Otherwise, the user will be asked to type in both.
//
// Another option might be to only ask the user for the password in case the --user
// flag has been specified, but this is not implemented at the moment.
func getOrReadInCredentials(options *globalOptions) (string, string) {
	user := options.user
	password := options.password

	if user != "" && password != "" {
		return user, password
	}

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("User: ")
	user, _ = reader.ReadString('\n')
	user = strings.TrimSpace(user)

	signalCh := make(chan os.Signal)
	signal.Notify(signalCh, os.Interrupt)

	go func() {
		<-signalCh
		os.Exit(0)
	}()

	fmt.Print("Password: ")

	p, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		fmt.Println("error reading password from stdin")
		os.Exit(1)
	}

	password = string(p)

	return user, password
}
