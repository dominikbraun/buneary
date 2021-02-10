package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/olekukonko/tablewriter"
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
		Version:       version,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	root.AddCommand(createCommand(&options))
	root.AddCommand(getCommand(&options))
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

	provider := NewProvider(&RabbitMQConfig{
		Address:  address,
		User:     user,
		Password: password,
	})

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

	if err := provider.CreateExchange(exchange); err != nil {
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

	provider := NewProvider(&RabbitMQConfig{
		Address:  address,
		User:     user,
		Password: password,
	})

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

	_, err := provider.CreateQueue(queue)
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
		Short: "Create a new binding",
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

	provider := NewProvider(&RabbitMQConfig{
		Address:  address,
		User:     user,
		Password: password,
	})

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

	if err := provider.CreateBinding(binding); err != nil {
		return err
	}

	return nil
}

// getCommand creates the `buneary get` command without any functionality.
func getCommand(options *globalOptions) *cobra.Command {
	get := &cobra.Command{
		Use:   "get <COMMAND>",
		Short: "Create a resource",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	get.AddCommand(getExchangesCommand(options))
	get.AddCommand(getExchangeCommand(options))
	get.AddCommand(getQueuesCommand(options))
	get.AddCommand(getQueueCommand(options))
	get.AddCommand(getBindingsCommand(options))
	get.AddCommand(getBindingCommand(options))

	return get
}

// getExchangesCommand creates the `buneary get exchanges` command, making sure that
// exactly one argument is passed.
func getExchangesCommand(options *globalOptions) *cobra.Command {
	getExchanges := &cobra.Command{
		Use:   "exchanges <ADDRESS>",
		Short: "Get all available exchanges",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGetExchanges(options, args)
		},
	}

	return getExchanges
}

// getExchangeCommand creates the `buneary get exchange` command, making sure that exactly
// two arguments are passed.
func getExchangeCommand(options *globalOptions) *cobra.Command {
	getExchange := &cobra.Command{
		Use:   "exchange <ADDRESS> <NAME>",
		Short: "Get a single exchange",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGetExchanges(options, args)
		},
	}

	return getExchange
}

// runGetExchanges either returns all exchanges or - if an exchange name has been
// specified as second argument - a single exchange. In case the password or both
// the user and password aren't provided, it will go into interactive mode.
//
// This flexibility allows runGetExchanges to be used by both `buneary get exchanges`
// as well as `buneary get exchange`.
func runGetExchanges(options *globalOptions, args []string) error {
	var (
		address = args[0]
	)

	user, password := getOrReadInCredentials(options)

	provider := NewProvider(&RabbitMQConfig{
		Address:  address,
		User:     user,
		Password: password,
	})

	// The default filter will let pass all exchanges regardless of their names.
	filter := func(_ Exchange) bool {
		return true
	}

	// However, if an exchange name has been specified as second argument, only
	// that particular exchange should be returned.
	if len(args) > 1 {
		filter = func(exchange Exchange) bool {
			return exchange.Name == args[1]
		}
	}

	exchanges, err := provider.GetExchanges(filter)
	if err != nil {
		return err
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Type", "Durable", "Auto-Delete", "Internal"})

	for _, exchange := range exchanges {
		row := make([]string, 5)
		row[0] = exchange.Name
		row[1] = string(exchange.Type)
		row[2] = boolToString(exchange.Durable)
		row[3] = boolToString(exchange.AutoDelete)
		row[4] = boolToString(exchange.Internal)
		table.Append(row)
	}

	table.Render()

	return nil
}

// getQueuesCommand creates the `buneary get queues` command, making sure that
// exactly one argument is passed.
func getQueuesCommand(options *globalOptions) *cobra.Command {
	getQueues := &cobra.Command{
		Use:   "queues <ADDRESS>",
		Short: "Get all available queues",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGetQueues(options, args)
		},
	}

	return getQueues
}

// getQueueCommand creates the `buneary get queue` command, making sure that exactly two
// arguments are passed.
func getQueueCommand(options *globalOptions) *cobra.Command {
	getQueue := &cobra.Command{
		Use:   "queue <ADDRESS> <NAME>",
		Short: "Get a single queue",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGetQueues(options, args)
		},
	}

	return getQueue
}

// runGetQueues either returns all queues or - if a queue name has been specified as second
// argument - a single queue. In case the password or both the user and password aren't
// provided, it will go into interactive mode.
//
// This flexibility allows runGetQueues to be used by both `buneary get queues` as well as
// `buneary get queue`.
func runGetQueues(options *globalOptions, args []string) error {
	var (
		address = args[0]
	)

	user, password := getOrReadInCredentials(options)

	provider := NewProvider(&RabbitMQConfig{
		Address:  address,
		User:     user,
		Password: password,
	})

	// The default filter will let pass all queues regardless of their names.
	filter := func(_ Queue) bool {
		return true
	}

	// However, if a queue name has been specified as second argument, only that
	// particular queue should be returned.
	if len(args) > 1 {
		filter = func(queue Queue) bool {
			return queue.Name == args[1]
		}
	}

	queues, err := provider.GetQueues(filter)
	if err != nil {
		return err
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Durable", "Auto-Delete"})

	for _, queue := range queues {
		row := make([]string, 3)
		row[0] = queue.Name
		row[1] = boolToString(queue.Durable)
		row[2] = boolToString(queue.AutoDelete)
		table.Append(row)
	}

	table.Render()

	return nil
}

// getBindingsCommand creates the `buneary get bindings` command, making sure that
// exactly one argument is passed.
func getBindingsCommand(options *globalOptions) *cobra.Command {
	getQueues := &cobra.Command{
		Use:   "bindings <ADDRESS>",
		Short: "Get all available bindings",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGetBindings(options, args)
		},
	}

	return getQueues
}

// getBindingCommand creates the `buneary get binding` command, making sure that exactly
// three arguments are passed.
func getBindingCommand(options *globalOptions) *cobra.Command {
	getQueue := &cobra.Command{
		Use:   "binding <ADDRESS> <EXCHANGE NAME> <TARGET NAME>",
		Short: "Get the binding or bindings between two resources",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGetBindings(options, args)
		},
	}

	return getQueue
}

// runGetBindings either returns all bindings  or - if a queue name has been specified as second
// argument - a single binding. In case the password or both the user and password aren't
// provided, it will go into interactive mode.
//
// This flexibility allows runGetBindings to be used by both `buneary get bindings` as well as
// `buneary get binding`.
func runGetBindings(options *globalOptions, args []string) error {
	var (
		address = args[0]
	)

	user, password := getOrReadInCredentials(options)

	provider := NewProvider(&RabbitMQConfig{
		Address:  address,
		User:     user,
		Password: password,
	})

	// The default filter will let pass all bindings regardless of their names.
	filter := func(_ Binding) bool {
		return true
	}

	// However, if a source exchange and a binding target have been specified as
	// second argument, only that particular binding should be returned.
	if len(args) > 2 {
		filter = func(binding Binding) bool {
			return binding.From.Name == args[1] &&
				binding.TargetName == args[2]
		}
	}

	bindings, err := provider.GetBindings(filter)
	if err != nil {
		return err
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"From", "Target", "Type", "Binding Key"})

	for _, binding := range bindings {
		row := make([]string, 4)
		row[0] = binding.From.Name
		row[1] = binding.TargetName
		row[2] = string(binding.Type)
		row[3] = binding.Key
		table.Append(row)
	}

	table.Render()

	return nil
}

// publishOptions defines options for publishing a message.
type publishOptions struct {
	*globalOptions
	headers string
}

// publishCommand creates the `buneary publish` command, making sure that exactly
// four command arguments are passed.
func publishCommand(options *globalOptions) *cobra.Command {
	publishOptions := &publishOptions{
		globalOptions: options,
	}

	publish := &cobra.Command{
		Use:   "publish <ADDRESS> <EXCHANGE> <ROUTING KEY> <BODY>",
		Short: "Publish a message to an exchange",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runPublish(publishOptions, args)
		},
	}

	publish.Flags().
		StringVar(&publishOptions.headers, "headers", "", "headers as comma-separated key-value pairs")

	return publish
}

// runPublish publishes a message by reading the command line data, setting the
// configuration and calling the PublishMessage function. In case the password or
// both the user and password aren't provided, it will go into interactive mode.
func runPublish(options *publishOptions, args []string) error {
	var (
		address    = args[0]
		exchange   = args[1]
		routingKey = args[2]
		body       = args[3]
	)

	user, password := getOrReadInCredentials(options.globalOptions)

	provider := NewProvider(&RabbitMQConfig{
		Address:  address,
		User:     user,
		Password: password,
	})

	message := Message{
		Target:     Exchange{Name: exchange},
		Headers:    make(map[string]interface{}),
		RoutingKey: routingKey,
		Body:       []byte(body),
	}

	if options.headers != "" {
		// Parse the message headers in the form key1=val1,key2=val2. If the headers
		// do not adhere to this syntax, an error is returned. In case the same key
		// exists multiple times, the last one wins.
		for _, header := range strings.Split(options.headers, ",") {
			tokens := strings.Split(strings.TrimSpace(header), "=")

			if len(tokens) != 2 {
				return errors.New("expected header in form key=value")
			}

			key := tokens[0]
			value := tokens[1]

			message.Headers[key] = value
		}
	}

	if err := provider.PublishMessage(message); err != nil {
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

	provider := NewProvider(&RabbitMQConfig{
		Address:  address,
		User:     user,
		Password: password,
	})

	exchange := Exchange{
		Name: name,
	}

	if err := provider.DeleteExchange(exchange); err != nil {
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

	provider := NewProvider(&RabbitMQConfig{
		Address:  address,
		User:     user,
		Password: password,
	})

	queue := Queue{
		Name: name,
	}

	if err := provider.DeleteQueue(queue); err != nil {
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
			fmt.Printf("buneary version %s", version)
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

	_, _ = os.Stdout.Write([]byte{'\n'})

	password = string(p)

	return user, password
}

// boolToString returns "yes" if the given bool is true and "no" if it is false.
func boolToString(source bool) string {
	if source {
		return "yes"
	}
	return "no"
}
