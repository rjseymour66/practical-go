package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"text/template"
)

var errInvalidPosArgSpecified = errors.New("More than one positional argument specified")

// in-memory representation of data used for runtime behavior
type config struct {
	numTimes       int
	name           string
	outputHtmlPath string
}

var usageString = fmt.Sprintf(`Usage: %s <integer> [-h|--help]

A greeter application that prints the name you entered <integer> number of times.
`, os.Args[0])

// validation funcs only return errors
func validateArgs(c config) error {
	if !(c.numTimes > 0) && len(c.outputHtmlPath) == 0 {
		return errors.New("Must specify a number greater than 0")
	}
	return nil
}

// 1. create empty config
// 2. create FlagSet object that names the command and what to do for errors
// 3. set the writer
// 4. define the command line option (where to save addr, name, default, default msg)
// 5. parse args
// 6. handle positional args
func parseArgs(w io.Writer, args []string) (config, error) {
	// create empty config
	c := config{}
	// create FlagSet object to handle CLI args (command-name, error action)
	fs := flag.NewFlagSet("greeter", flag.ContinueOnError)
	// set the writer
	fs.SetOutput(w)
	fs.Usage = func() {
		var usageString = `
	A greeter application that prints the name you entered a specified number of times.

	Usage of %s: <options> [name]`
		fmt.Fprintf(w, usageString, fs.Name())
		fmt.Fprintln(w)
		fmt.Fprintln(w, "Options: ")
		fs.PrintDefaults()
	}
	fs.StringVar(&c.outputHtmlPath, "o", "", "Path to the HTML page containing the greeting")
	// setting a command line option. IntVar creates int option
	fs.IntVar(&c.numTimes, "n", 0, "Number of times to greet")
	err := fs.Parse(args)
	if err != nil {
		return c, err
	}
	// if err == nil, check to see if there was a positional arg
	// NArg() returns the number of positional args
	if fs.NArg() > 1 {
		return c, errInvalidPosArgSpecified
	}

	if fs.NArg() == 1 {
		c.name = fs.Arg(0)
	}

	return c, nil
}

// prompts the user and accepts input
func getName(r io.Reader, w io.Writer) (string, error) {
	msg := "Your name please? Press the Enter key when done.\n"
	fmt.Fprintf(w, msg)

	// create a scanner to read user input
	scanner := bufio.NewScanner(r)
	scanner.Scan()
	if err := scanner.Err(); err != nil {
		return "", err
	}
	name := scanner.Text()
	if len(name) == 0 {
		return "", errors.New("You didn't enter your name")
	}

	return name, nil
}

func greetUser(c config, w io.Writer) {
	msg := fmt.Sprintf("Nice to meet you %s\n", c.name)
	for i := 0; i < c.numTimes; i++ {
		fmt.Fprintf(w, msg)
	}
}

func createHtmlGreeter(c config, name string) error {
	f, err := os.Create(c.outputHtmlPath)
	if err != nil {
		return err
	}
	defer f.Close()

	tmpl, err := template.New("greeterHtml").Parse("<h1>Hello {{.}}</h1>")
	if err != nil {
		return err
	}
	return tmpl.Execute(f, name)
}

// performs action based on config values
func runCmd(rd io.Reader, w io.Writer, c config) error {
	var err error
	if len(c.name) == 0 {
		c.name, err = getName(rd, w)
		if err != nil {
			return err
		}
	}
	if len(c.outputHtmlPath) != 0 {
		return createHtmlGreeter(c, c.name)
	}
	greetUser(c, w)
	return nil
}

func main() {
	// parse the args
	c, err := parseArgs(os.Stderr, os.Args[1:])
	if err != nil {
		if errors.Is(err, errInvalidPosArgSpecified) {
			fmt.Fprintln(os.Stdout, err)
		}
		os.Exit(1)
	}
	// validate there are args
	err = validateArgs(c)
	if err != nil {
		fmt.Fprintln(os.Stdout, err)
		os.Exit(1)
	}
	// run getUsage or print the name to the console
	err = runCmd(os.Stdin, os.Stdout, c)
	if err != nil {
		fmt.Fprintln(os.Stdout, err)
		os.Exit(1)
	}
}
