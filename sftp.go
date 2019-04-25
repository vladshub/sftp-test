package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

func main() {
	//hostKey := getHostKey("127.0.0.1")
	hostKey1, _, _, _, err := ssh.ParseAuthorizedKey([]byte("[127.0.0.1]:2222 ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQCaHiEzMIM/I3aCgmWk2KYK+T8FpRW1sl2a08qHTWd8sh/HNixa75z4hr5/to+8vaHyxuvkj47XYOYYZB2d7RN5AjcmGD4FMh0zFcYOHXklWLQMyaZLXMiFR+Ee/AilDomCHHA4ZVMmZV0Me6uhcRf5gzrvPsNwh6+185IhG4pv+9Al3ocBS5itViweRxYmKjGoAo4ulzqxK68HV0RyVYtV+hYvZwl6oUw/6V5DKpFDoVwgVhPdY4+Fc81QYG9Da7Bg6nC+nDCLvN0bmMKeK81YytynTrszNWAHrHFPDZAM7OgfUIZsGCoBKM3IE7x0oxv4oo6vtsHMpFAS6n7pRS1h1gLd2PwFhtvKWYwdC9re3FVY9ouvvI81aqFnhmXZKA5nrRHx01S9PqEVl6ddYKitAEzXkfQCzpFrBHoHpq3f88yzUHn2Mq62TbWz1EhONCvp7xn33MzZCzysvGLwEUBInhTTcJPzNRbJRvKb96E+iThly4Yfz/xq452l9LZFkZsBUt43BzSI2XIS9b7IxSrfW8gc8gPaoFCMJ6iKeK9s0NLWPGY5xEhZvTHwHPRMa4lC7EBsynRPnDIhgb8wRqSxA++IQnzQDW2F7u0ak/3Y8IowTBMlhl6+gPM42L2mXVlhaU5GCs11+xov2zych0ImmcydoRT9JYqlDg1dXQOnPQ=="))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(readFromSftp("foo", "pass", "127.0.0.1", "2222", "upload/example", hostKey1))
}

func readFromSftp(user, pass, host, port, path string, hostKey ssh.PublicKey) string {

	// user := ""
	// pass := ""
	// remote := ""
	// port := ":22"

	// get host public key
	//hostKey := getHostKey(host)

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(pass),
		},
		// HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		HostKeyCallback: ssh.FixedHostKey(hostKey),
	}

	// connect
	conn, err := ssh.Dial("tcp", host+":"+port, config)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// create new SFTP client
	client, err := sftp.NewClient(conn)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// create destination file
	dstFile, err := os.Create("./file.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		dstFile.Close()
		os.Remove("./file.txt")
	}()

	// open source file
	//dir, file := sftp.Split(path)

	srcFile, err := client.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	// copy source file to destination file
	bytes, err := io.Copy(dstFile, srcFile)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%d bytes copied\n", bytes)

	// flush in-memory copy
	err = dstFile.Sync()
	if err != nil {
		log.Fatal(err)
	}
	dat, err := ioutil.ReadFile("./file.txt")
	if err != nil {
		log.Fatal(err)
	}
	return string(dat)
}

func getHostKey(host string) ssh.PublicKey {
	// parse OpenSSH known_hosts file
	// ssh or use ssh-keyscan to get initial key
	file, err := os.Open(filepath.Join(os.Getenv("HOME"), ".ssh", "known_hosts"))
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var hostKey ssh.PublicKey
	for scanner.Scan() {
		fields := strings.Split(scanner.Text(), " ")
		if len(fields) != 3 {
			continue
		}
		if strings.Contains(fields[0], host) {
			var err error
			hostKey, _, _, _, err = ssh.ParseAuthorizedKey(scanner.Bytes())
			if err != nil {
				log.Fatalf("error parsing %q: %v", fields[2], err)
			}
			break
		}
	}

	if hostKey == nil {
		log.Fatalf("no hostkey found for %s", host)
	}

	return hostKey
}