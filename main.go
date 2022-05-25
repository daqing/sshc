package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/daqing/goph"
)

type HostInfo struct {
	Key      string
	IP       string
	User     string
	Port     uint
	Password string
	Options  map[string]interface{}
}

// 远程执行命令
func run(client *goph.Client, cmd string) {
	fmt.Printf("\033[0;33m---> running command: [%s]\033[0m\n", cmd)

	out, err := client.Run(cmd)
	if err != nil {
		log.Fatalf("run remote cmd error: %s\n", err)
	}

	fmt.Printf("%s", string(out))
}

func cp(client *goph.Client, src, dest string) {
	dest = fill_path(src, dest)

	err := client.Upload(src, dest)
	if err != nil {
		log.Fatalf("cp to remote error: src=%s, dest=%s, err=%s\n", src, dest, err)
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
		log.Fatalf("download error: dest=%s, src=%s, error=%s\n", dest, src, err)
	}
}

func docker_exec(client *goph.Client, image string, cmd string) {
	fullCmd := fmt.Sprintf("docker exec -t %s sh -c '%s'", image, cmd)

	run(client, fullCmd)
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
		log.Fatalf("read toml error: %s\n", err)
	}

	var hostInfo HostInfo

	if _, err := toml.Decode(string(content), &hostInfo); err != nil {
		log.Fatalf("toml deocode error: %s\n", err)
	}

	if hostInfo.Port == 0 {
		hostInfo.Port = 22
	}

	action := os.Args[4]
	if action == "pwd" {
		fmt.Printf("%s\n", hostInfo.Password)
		return
	}

	var client *goph.Client
	var auth goph.Auth

	if len(hostInfo.Key) > 0 {
		auth, err = goph.Key(hostInfo.Key, "")
		if err != nil {
			log.Fatal("goph.Key error:", err)
		}
	} else {
		auth = goph.Password(hostInfo.Password)
	}

	client, err = goph.New(hostInfo.User, hostInfo.IP, hostInfo.Port, auth)
	if err != nil {
		log.Fatal("goph.New error:", err)
	}
	defer client.Close()

	switch action {
	case "run":
		if len(os.Args) == 5 {
			fmt.Println("please provide command to run")
			return
		}

		run(client, strings.Join(os.Args[5:], " "))
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
	case "dk", "docker":
		if len(os.Args) == 5 {
			fmt.Printf("usage: %s %s [cmd]\n", os.Args[0], action)
			return
		}

		image, valid := hostInfo.Options["image"]
		if !valid {
			fmt.Println("no image name configured in host info Options")
			return
		}

		docker_exec(client, image.(string), strings.Join(os.Args[5:], " "))
	default:
		if v, ok := hostInfo.Options[action]; ok {
			if strings.HasPrefix(action, "is_") {
				fmt.Println(v.(int64) > 0)
			} else {
				fmt.Println(v.(string))
			}

			return
		} else {
			// fallback to run command
			run(client, strings.Join(os.Args[4:], " "))
		}
	}
}
