package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/melbahja/goph"
)

type HostInfo struct {
	Key      string
	IP       string
	User     string
	Password string
}

func run(client *goph.Client, cmd string) {
	out, err := client.Run(cmd)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s", string(out))
}

func cp(client *goph.Client, src, dest string) {
	dest_name := filepath.Base(dest)
	if !strings.Contains(dest_name, ".") {
		src_filename := filepath.Base(src)

		dest = strings.TrimRight(dest, "/")
		dest = strings.Join([]string{dest, src_filename}, "/")
	}

	err := client.Upload(src, dest)
	if err != nil {
		log.Fatal(err)
	}
}

func down(client *goph.Client, dest, src string) {
	err := client.Download(dest, src)
	if err != nil {
		panic(err)
	}
}

// sshc -c $HOME/hosts v8os run ls
// sshc -c $HOME/hosts v8os cp "/tmp/abc.txt" "/tmp/foo.txt"
// sshc -c $HOME/hosts v8os down "/tmp/foo.txt" "/tmp/demo.txt"
func main() {
	if len(os.Args) < 5 {
		fmt.Printf("usage: %s [host] [action] [arg1 arg2 arg3 ...]\n", os.Args[0])
		return
	}

	confDir := flag.String("c", "", "config dir to load [host].toml")
	flag.Parse()

	if len(*confDir) == 0 {
		log.Fatal("Please provide conf dir via -c [path/to/dir]")
	}

	host := os.Args[3]

	content, err := os.ReadFile(fmt.Sprintf("%s/%s.toml", *confDir, host))
	if err != nil {
		log.Fatal(err)
	}

	var hostInfo HostInfo

	if _, err := toml.Decode(string(content), &hostInfo); err != nil {
		log.Fatal(err)
	}

	var client *goph.Client
	var auth goph.Auth

	if len(hostInfo.Key) > 0 {
		auth, err = goph.Key(hostInfo.Key, "")
		if err != nil {
			log.Fatal(err)
		}
	} else {
		auth = goph.Password(hostInfo.Password)
	}

	client, err = goph.New(hostInfo.User, hostInfo.IP, auth)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	action := os.Args[4]
	switch action {
	case "run":
		if len(os.Args) == 5 {
			fmt.Println("please provide command to run")
			return
		}

		run(client, os.Args[5])
	case "scp", "cp":
		if len(os.Args) == 5 {
			fmt.Printf("usage: %s cp [src] [dest]\n", os.Args[0])
			return
		}

		cp(client, os.Args[5], os.Args[6])
	case "down":
		if len(os.Args) == 5 {
			fmt.Printf("usage: %s down [dest] [src]\n", os.Args[0])
			return
		}

		down(client, os.Args[5], os.Args[6])
	default:
		log.Fatal("unknown action: " + action)
	}
}
