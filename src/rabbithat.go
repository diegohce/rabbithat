package main

import (
	"os"
	"fmt"
	"flag"
	"rabbit"
)


const (
	VERSION = "0.1.0"
	VERSION_NAME = "Varus"
)


type CmdOptions struct {
	SourceRabbit   string
	SourceUser     string
	SourcePassword string
	SourceVHost    string
	TargetRabbit   string
	TargetUser     string
	TargetPassword string
	TargetVHost    string
	SourceFile     string
	TargetFile     string
	Version        bool
}


func options() (*CmdOptions, error) {

	args := os.Args[1:]

	co := &CmdOptions{}

	fs := flag.NewFlagSet("rabbithat", flag.ExitOnError)

	fs.StringVar(&co.SourceRabbit  , "source-rabbit"  , "", "Source rabbit address:port")
	fs.StringVar(&co.SourceUser    , "source-user"    , "", "Source rabbit username")
	fs.StringVar(&co.SourcePassword, "source-password", "", "Source rabbit password")
	fs.StringVar(&co.SourceVHost   , "source-vhost"   , "", "Source rabbit virtual host")

	fs.StringVar(&co.TargetRabbit  , "target-rabbit"  , "", "Target rabbit address:port")
	fs.StringVar(&co.TargetUser    , "target-user"    , "", "Target rabbit username")
	fs.StringVar(&co.TargetPassword, "target-password", "", "Target rabbit password")
	fs.StringVar(&co.TargetVHost   , "target-vhost"   , "", "Target rabbit virtual host")

	fs.StringVar(&co.TargetFile, "target-file", "", "File to dump source rabbit data (json format)")
	fs.StringVar(&co.SourceFile, "source-file", "", "File to read source rabbit data (json format)")

	fs.BoolVar(&co.Version, "version", false, "Rabbit Hat version")

	if len(args) == 0 {
		args = append(args, "--help")
	}

	err := fs.Parse(args)
	if err != nil {
		return nil, err
	}

	return co, nil

}



func main() {

	var r *rabbit.RabbitMQ


	op, err := options()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if op.Version {
		fmt.Printf("This is Rabbit Hat %s (%s)\n", VERSION, VERSION_NAME)
		os.Exit(0)
	}

	if op.SourceFile != "" && op.SourceRabbit != "" {
		fmt.Println("Ambiguous source")
		os.Exit(1)
	}
	if op.SourceFile == "" && op.SourceRabbit == "" {
		fmt.Println("No source specified")
		os.Exit(1)
	}

	if op.TargetFile != "" && op.TargetRabbit != "" {
		fmt.Println("Ambiguous target")
		os.Exit(1)
	}
	if op.TargetFile == "" && op.TargetRabbit == "" {
		fmt.Println("No target specified")
		os.Exit(1)
	}


	if op.SourceFile != "" {
		var err error

		fmt.Println("Reading from", op.SourceFile)

		r, err = rabbit.CollectFromFile(op.SourceFile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	if op.SourceRabbit != "" {

		fmt.Println("Collecting from", op.SourceRabbit)

		r = &rabbit.RabbitMQ{
			BaseUrl : "http://" + op.SourceRabbit,
			User    : op.SourceUser,
			Password: op.SourcePassword,
			VHost   : op.SourceVHost,
		}

		if err := r.Collect(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}


	if op.TargetRabbit != "" {

		fmt.Println("Cloning to", op.TargetRabbit)

		if err := r.CloneTo(op.TargetRabbit, op.TargetUser, op.TargetPassword, op.TargetVHost); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		os.Exit(0)
	}


	if op.TargetFile != "" {

		fmt.Println("Dumping to", op.TargetFile)

		if err := r.DumpTo(op.TargetFile); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		os.Exit(0)
	}
}

/*

rmq to rmq
rmq to dump file
dump file to rmq

*/


