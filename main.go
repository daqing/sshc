package main

import (
	"fmt"
	"log"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/melbahja/goph"
)

type HostInfo struct {
	IP       string
	User     string
	Password string
}

func run(client *goph.Client, cmd string) {
	out, err := client.Run(cmd)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s", string(out))
}

func cp(client *goph.Client, src, dest string) {
	err := client.Upload(src, dest)
	if err != nil {
		panic(err)
	}
}

func down(client *goph.Client, dest, src string) {
	err := client.Download(dest, src)
	if err != nil {
		panic(err)
	}
}

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("usage: %s [host] [action] [arg1 arg2 arg3 ...]\n", os.Args[0])
		return
	}

	host := os.Args[1]

	content, err := os.ReadFile(fmt.Sprintf("config/%s.toml", host))
	if err != nil {
		log.Fatal(err)
	}

	var hostInfo HostInfo

	if _, err := toml.Decode(string(content), &hostInfo); err != nil {
		log.Fatal(err)
	}

	client, err := goph.New(hostInfo.User, hostInfo.IP, goph.Password(hostInfo.Password))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	action := os.Args[2]
	switch action {
	case "run":
		if len(os.Args) == 3 {
			fmt.Println("please provide command to run")
			return
		}

		run(client, os.Args[3])
	case "cp":
		if len(os.Args) < 5 {
			fmt.Printf("usage: %s cp [src] [dest]\n", os.Args[0])
			return
		}

		cp(client, os.Args[3], os.Args[4])
	case "down":
		if len(os.Args) < 5 {
			fmt.Printf("usage: %s down [dest] [src]\n", os.Args[0])
			return
		}

		down(client, os.Args[3], os.Args[4])
	}
}
