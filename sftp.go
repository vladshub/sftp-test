package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

func main() {
	lines := run_sftp_command("ls -1")
	fmt.Println(len(lines))

}

type sftpReader struct {
	buffer          []string
	currentLocation int
}

func (s *sftpReader) Write(p []byte) (n int, err error) {
	line := string(p)
	if !strings.Contains(line, ">") && !strings.Contains(line, "Connected to") {
		s.buffer = append(s.buffer, line)
	}
	log.Println(line)
	return len(p), nil
}

func (s *sftpReader) ReadNext() string {
	if len(s.buffer) > s.currentLocation {
		defer func() { s.currentLocation = s.currentLocation + 1 }()
		return s.buffer[s.currentLocation]
	}
	return ""
}

func run_sftp_command(command string) []string {
	sftpPath, err := exec.LookPath("sftp")
	if err != nil {
		log.Fatal("installing sftp zis in your future\nFor osx run brew install sftp\nFor Ubuntu run apt install openssh-server")
	}
	log.Println("sftp is available at ", sftpPath)
	sshpassPath, err := exec.LookPath("sshpass")
	if err != nil {
		log.Fatal("installing sshpass is in your future\nFor osx run brew install https://raw.githubusercontent.com/kadwanev/bigboybrew/master/Library/Formula/sshpass.rb\nFor Ubuntu run apt install sshpass")
	}
	log.Println("sshpass is available at ", sshpassPath)
	cmd := exec.CommandContext(context.Background(), sshpassPath, "-p", "pass", sftpPath, "-P", "2222", "foo@127.0.0.1:/upload")
	cmd.Env = os.Environ()

	stdOut, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	stdErr, err := cmd.StderrPipe()
	if err != nil {
		log.Fatal(err)
	}
	out := io.MultiReader(stdOut, stdErr)

	writeBuffer, err := cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	wpath := fmt.Sprintf("%s/downloads", dir)
	err = os.MkdirAll(wpath, 755)
	if err != nil {
		log.Fatal(err)
	}
	cmd.Dir = wpath
	err = cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	_, err = writeBuffer.Write(append([]byte(command), '\n'))
	_, err = writeBuffer.Write([]byte("exit\n"))

	response := make([]string, 0)
	scanner := bufio.NewScanner(out)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		txt := scanner.Text()
		if strings.Contains(txt, "Changing to:") || strings.Contains(txt, "Connected to") || strings.Contains(txt, ">") {
			continue
		}
		response = append(response, txt)
	}
	err = cmd.Wait()
	if err != nil {
		log.Fatal(err)
	}
	return response
}

func main_sftp() {
	sftpPath, err := exec.LookPath("sftp")
	if err != nil {
		log.Fatal("installing sftp zis in your future\nFor osx run brew install sftp\nFor Ubuntu run apt install openssh-server")
	}
	log.Println("sftp is available at ", sftpPath)
	sshpassPath, err := exec.LookPath("sshpass")
	if err != nil {
		log.Fatal("installing sshpass is in your future\nFor osx run brew install https://raw.githubusercontent.com/kadwanev/bigboybrew/master/Library/Formula/sshpass.rb\nFor Ubuntu run apt install sshpass")
	}
	log.Println("sshpass is available at ", sshpassPath)
	cmd := exec.CommandContext(context.Background(), sshpassPath, "-p", "pass", sftpPath, "-P", "2222", "foo@127.0.0.1")
	cmd.Env = os.Environ()

	stdOut, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	stdErr, err := cmd.StderrPipe()
	if err != nil {
		log.Fatal(err)
	}
	out := io.MultiReader(stdOut, stdErr)

	writeBuffer, err := cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	wpath := fmt.Sprintf("%s/downloads", dir)
	err = os.MkdirAll(wpath, 755)
	if err != nil {
		log.Fatal(err)
	}
	cmd.Dir = wpath
	err = cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	_, err = writeBuffer.Write([]byte("cd upload\n"))
	if err != nil {
		log.Fatal(err)
	}

	_, err = writeBuffer.Write([]byte("ls\n"))
	if err != nil {
		log.Fatal(err)
	}

	_, err = writeBuffer.Write([]byte("get example\n"))
	if err != nil {
		log.Fatal(err)
	}

	var buf bytes.Buffer
	tee := io.TeeReader(out, &buf)
	data, err := ioutil.ReadAll(tee)
	if err != nil {
		log.Fatal(err)
	}

	//scanner := bufio.NewScanner(out)
	//scanner.Split(bufio.ScanLines)
	//for scanner.Scan() {
	//	txt := scanner.Text()
	//	fmt.Println(txt)
	//}
	_, err = writeBuffer.Write([]byte("exit\n"))
	if err != nil {
		log.Fatal(err)
	}
	//data, err := ioutil.ReadAll(out)
	//if err != nil {
	//	log.Fatal(err)
	//}

	log.Println(string(data))
	err = cmd.Wait()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("done")
}

func native_main() {
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
