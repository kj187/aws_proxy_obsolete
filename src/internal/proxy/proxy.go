package proxy

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"syscall"

	"github.com/fatih/color"
	"github.com/kj187/sshtunnel"
	"go.kj187.de/aws_proxy/src/internal/aws"
	"go.kj187.de/aws_proxy/src/internal/configuration"
	"golang.org/x/crypto/ssh/terminal"
)

// Proxy exported
type Proxy struct {
	Bastion       struct{ PublicIP string }
	Label         string
	Verbose       bool
	SSHPassphrase bool
	Destination   string
	LocalPort     int64
	Tunnel        *sshtunnel.SSHTunnel
}

// Start exported
func (p *Proxy) Start(additionalInformation []string) {
	err := p.loadBastionPublicIP()
	if err != nil {
		color.Red("Cant fetch bastion public ip: %s", err)
		return
	}

	config := configuration.Get()

	p.printHeader()

	var sshKeyPassphrase []byte
	if p.SSHPassphrase {
		fmt.Print("Enter SSH Key password: ")
		bytePassword, _ := terminal.ReadPassword(int(syscall.Stdin))
		sshKeyPassphrase = []byte(string(bytePassword))
	}

	key, err := sshtunnel.PrivateKeyFile(config.Bastion.SSH.IdentityFile, sshKeyPassphrase)
	if err != nil {
		color.Red("Cant read SSH key: %s", err)
		color.Red("If it is passphrase protected just execute the same command with the -p flag")
		return
	}

	tunnel := sshtunnel.NewSSHTunnel(
		fmt.Sprintf("%s@%s:%d", config.Bastion.SSH.User, p.Bastion.PublicIP, config.Bastion.SSH.Port),
		key,
		p.Destination,
		strconv.FormatInt(p.LocalPort, 10),
	)

	p.printInfo(additionalInformation)

	if p.Verbose {
		tunnel.Log = log.New(os.Stdout, "", log.Ldate|log.Lmicroseconds)
	}

	p.Tunnel = tunnel

	go tunnel.Start()
	fmt.Scanln()
}

func (p *Proxy) printHeader() {
	localhost := fmt.Sprintf("127.0.0.1:%v", p.LocalPort)
	localhostLabelLine := strings.Repeat("-", len(localhost))
	bastionLabelLine := strings.Repeat("-", len(p.Bastion.PublicIP))
	targetLabelLine := strings.Repeat("-", len(p.Destination))
	color.Green("+-%s-+      +-%s-+      +-%s-+", localhostLabelLine, bastionLabelLine, targetLabelLine)
	color.Green("| Localhost%v |      | Bastion%v |      | %v%v |", strings.Repeat(" ", (len(localhost)-9)), strings.Repeat(" ", (len(p.Bastion.PublicIP)-7)), p.Label, strings.Repeat(" ", (len(p.Destination)-len(p.Label))))
	color.Green("| %v | <--> | %v | <--> | %v |", localhost, p.Bastion.PublicIP, p.Destination)
	color.Green("+-%s-+      +-%s-+      +-%s-+", localhostLabelLine, bastionLabelLine, targetLabelLine)
	fmt.Println()
}

func (p *Proxy) printInfo(additionalInformation []string) {
	fmt.Println()
	config := configuration.Get()
	if p.Verbose {
		color.Blue("Proxy Host: 127.0.0.1 (over bastion '%s@%s:%s')", config.Bastion.SSH.User, p.Bastion.PublicIP, config.Bastion.SSH.Port)
		color.Blue("Proxy Port: %v", p.LocalPort)
	}
	color.Green("Proxy is enabled now")
	fmt.Println()

	if len(additionalInformation) > 0 {
		for _, line := range additionalInformation {
			color.Green(line)
		}
	}

	fmt.Println()
}

func (p *Proxy) loadBastionPublicIP() error {
	config := configuration.Get()
	p.Bastion.PublicIP = config.Bastion.PublicIP
	if config.Bastion.PublicIP == "" {
		bastionPublicIP, err := aws.GetBastionPublicIP(p.Verbose)
		if err != nil {
			return err
		}
		p.Bastion.PublicIP = bastionPublicIP
	}

	return nil
}
