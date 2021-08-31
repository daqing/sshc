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
	Options  map[string]interface{}
}

// 远程执行命令
func run(client *goph.Client, cmd string) {
	out, err := client.Run(cmd)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s", string(out))
}

func cp(client *goph.Client, src, dest string) {
	dest = fill_path(src, dest)

	err := client.Upload(src, dest)
	if err != nil {
		log.Fatal(err)
	}
}

// 如果`dest`是一个目录
// 把 src 的basename 添加到 dest 上面
// 形成一个文件路径
func fill_path(src, dest string) string {
	dest_name := filepath.Base(dest)
	if !strings.Contains(dest_name, ".") {
		src_filename := filepath.Base(src)

		dest = strings.TrimRight(dest, "/")
		dest = strings.Join([]string{dest, src_filename}, "/")
	}

	return dest
}

func down(client *goph.Client, dest, src string) {
	src = fill_path(dest, src)

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
			fmt.Printf("usage: %s %s [src] [dest]\n", os.Args[0], action)
			return
		}

		cp(client, os.Args[5], os.Args[6])
	case "dl", "down":
		if len(os.Args) == 5 {
			fmt.Printf("usage: %s %s [dest] [src]\n", os.Args[0], action)
			return
		}

		down(client, os.Args[5], os.Args[6])
	default:
		if v, ok := hostInfo.Options[action]; ok {
			if strings.HasPrefix(action, "is_") {
				fmt.Println(v.(int64) > 0)
			} else {
				fmt.Println(v.(string))
			}

			return
		}

		log.Fatal("unknown action: " + action)
	}
}
