package awsconsole

import (
	"fmt"
	"time"

	"github.com/pkg/browser"
	"github.com/sethvargo/go-password/password"
	"github.com/spf13/cobra"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/kubectl/pkg/util/templates"

	"github.com/cblecker/oc-sre/pkg/options"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/sts"
)

const (
	awsCredsNamespace       = "kube-system"           // #nosec G101
	awsCredsSecretName      = "aws-creds"             // #nosec G101
	awsCredsSecretIDKey     = "aws_access_key_id"     // #nosec G101
	awsCredsSecretAccessKey = "aws_secret_access_key" // #nosec G101
	osdAwsUsername          = "osdManagedAdmin"
	loginURL                = "https://%s.signin.aws.amazon.com/console"
)

var (
	consoleShort = templates.LongDesc(`
		Provide login details for the AWS console.`)

	consoleExample = templates.Examples(`
		# Provide login details for the AWS console
		%[1]s %[2]s

		# Provide login details for the AWS console, and open the URL in your default browser
		%[1]s %[2]s --open`)

	consoleOutput = `Log in here: %s
Username: %s
Password: %s

Press enter when done logging in (timeout in 120 seconds)...
`
)

// ConsoleCmdOptions are options supported by the console command.
type ConsoleCmdOptions struct { //nolint:golint
	rootOptions *options.SRECmdOptions

	// Open is true if the command should also open the URL in the default browser
	Open bool

	// args is the slice of strings containing any arguments passed
	args []string
}

// NewConsoleCmdOptions provides an instance of ConsoleCmdOptions with default values
func NewConsoleCmdOptions(rootOptions *options.SRECmdOptions) *ConsoleCmdOptions {
	return &ConsoleCmdOptions{
		rootOptions: rootOptions,
	}
}

// NewCmdConsoleConfig provides a cobra command wrapping ConsoleCmdOptions
func NewCmdConsoleConfig(rootOptions *options.SRECmdOptions) *cobra.Command {
	o := NewConsoleCmdOptions(rootOptions)

	cmd := &cobra.Command{
		Use:          "awsconsole",
		Short:        consoleShort,
		Example:      fmt.Sprintf(consoleExample, options.RootCmd, "awsconsole"),
		SilenceUsage: true,
		RunE: func(c *cobra.Command, args []string) error {
			if err := o.Complete(args); err != nil {
				return err
			}
			if err := o.Validate(); err != nil {
				return err
			}
			if err := o.Run(); err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().BoolVar(&o.Open, "open", o.Open, "Also open the console URL in your default browser")
	o.rootOptions.ConfigFlags.AddFlags(cmd.Flags())

	return cmd
}

// Complete sets up the KubeClient
func (o *ConsoleCmdOptions) Complete(args []string) error {
	o.args = args

	return nil
}

// Validate ensures that all required arguments and flag values are provided
func (o *ConsoleCmdOptions) Validate() error {
	if len(o.args) > 0 {
		return fmt.Errorf("no arguments are allowed")
	}

	return nil
}

// getAWSCredentials retrieves the AWS credentials from the cluster
func (o *ConsoleCmdOptions) getAWSCredentials() (string, string, error) {
	awsCredentialSecret, err := o.rootOptions.KubeClient.CoreV1().
		Secrets(awsCredsNamespace).
		Get(awsCredsSecretName, metav1.GetOptions{})
	if err != nil {
		return "", "", fmt.Errorf("unable to find the AWS credential secret: %v", err)
	}

	accessKeyID, ok := awsCredentialSecret.Data[awsCredsSecretIDKey]
	if !ok {
		return "", "", fmt.Errorf("AWS credentials secret %v did not contain key %v",
			awsCredsSecretName, awsCredsSecretIDKey)
	}
	secretAccessKey, ok := awsCredentialSecret.Data[awsCredsSecretAccessKey]
	if !ok {
		return "", "", fmt.Errorf("AWS credentials secret %v did not contain key %v",
			awsCredsSecretName, awsCredsSecretAccessKey)
	}

	return string(accessKeyID), string(secretAccessKey), nil
}

// getAWSAccountID retrieves the AWS account ID
func getAWSAccountID(awsSession *session.Session) (string, error) {
	svc := sts.New(awsSession)
	input := &sts.GetCallerIdentityInput{}

	result, err := svc.GetCallerIdentity(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			return "", fmt.Errorf(aerr.Error())
		}
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		return "", fmt.Errorf(err.Error())
	}

	return *result.Account, nil
}

// createAWSLoginProfile creates a login profile
func createAWSLoginProfile(awsSession *session.Session, username, password string) error {
	svc := iam.New(awsSession)
	input := &iam.CreateLoginProfileInput{
		UserName:              aws.String(username),
		Password:              aws.String(password),
		PasswordResetRequired: aws.Bool(false),
	}

	_, err := svc.CreateLoginProfile(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			return fmt.Errorf(aerr.Error())
		}
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		return fmt.Errorf(err.Error())
	}

	return nil
}

// deleteAWSLoginProfile deletes a login profile
func deleteAWSLoginProfile(awsSession *session.Session, username string) error {
	svc := iam.New(awsSession)
	input := &iam.DeleteLoginProfileInput{
		UserName: aws.String(username),
	}

	_, err := svc.DeleteLoginProfile(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			return fmt.Errorf(aerr.Error())
		}
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		return fmt.Errorf(err.Error())
	}

	return nil
}

// Run grabs the console URL, and either prints it to the terminal or opens it
// in your default web browser
func (o *ConsoleCmdOptions) Run() error {
	var err error

	awsAccessKeyID, awsSecretAccessKey, err := o.getAWSCredentials()
	if err != nil {
		return err
	}

	awsConfig := &aws.Config{
		Credentials: credentials.NewStaticCredentials(
			awsAccessKeyID, awsSecretAccessKey, ""),
	}
	awsSession, err := session.NewSession(awsConfig)
	if err != nil {
		return err
	}

	account, err := getAWSAccountID(awsSession)
	if err != nil {
		return err
	}

	url := fmt.Sprintf(loginURL, account)
	pwd, err := password.Generate(32, 10, 0, false, false)
	if err != nil {
		return err
	}

	err = createAWSLoginProfile(awsSession, osdAwsUsername, pwd)
	if err != nil {
		if err = deleteAWSLoginProfile(awsSession, osdAwsUsername); err != nil {
			return err
		}
		return err
	}

	ch := make(chan int)

	if o.Open {
		if err = browser.OpenURL(url); err != nil {
			return err
		}
	}
	fmt.Printf(consoleOutput, url, osdAwsUsername, pwd)

	time.Sleep(15 * time.Second)

	go func() {
		fmt.Scanf("\n")
		ch <- 1
	}()

	select {
	case <-ch:
		fmt.Println("Exiting.")
	case <-time.After(120 * time.Second):
		fmt.Println("Timed out, exiting.")
	}

	if err = deleteAWSLoginProfile(awsSession, osdAwsUsername); err != nil {
		return err
	}

	return nil
}
