// +build go1.5

package main

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/codegangsta/cli"
	"github.com/sirupsen/logrus"
	"github.com/soopsio/gotty-client"
)

var VERSION string

func main() {
	app := cli.NewApp()
	app.EnableBashCompletion = true
	app.Name = path.Base(os.Args[0])
	app.Author = "Manfred Touron"
	app.Email = "https://github.com/moul/gotty-client"
	app.Version = VERSION
	app.Usage = "GoTTY client for your terminal"
	app.ArgsUsage = "GOTTY_URL"
	app.BashComplete = func(c *cli.Context) {
		for _, command := range []string{
			"--debug", "--skip-tls-verify", "--use-proxy-from-env",
			"--v2", "--detach-keys", "--ws-origin", "--help",
			"--generate-bash-completion", "--version",
			"http://user:pass@host:1234/path/\\\\?arg=abcdef\\\\&arg=ghijkl",
			"https://user:pass@host:1234/path/\\\\?arg=abcdef\\\\&arg=ghijkl",
			"http://localhost:8000",
		} {
			fmt.Println(command)
		}
	}

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:   "debug, D",
			Usage:  "Enable debug mode",
			EnvVar: "GOTTY_CLIENT_DEBUG",
		},
		cli.BoolFlag{
			Name:   "skip-tls-verify",
			Usage:  "Skip TLS verify",
			EnvVar: "SKIP_TLS_VERIFY",
		},
		cli.BoolFlag{
			Name:   "use-proxy-from-env",
			Usage:  "Use Proxy from environment",
			EnvVar: "USE_PROXY_FROM_ENV",
		},
		cli.StringFlag{
			Name:  "detach-keys",
			Usage: "Key sequence for detaching gotty-client",
			Value: "ctrl-p,ctrl-q",
		},
		cli.BoolFlag{
			Name:   "v2",
			Usage:  "For Gotty 2.0",
			EnvVar: "GOTTY_CLIENT_GOTTY2",
		},
		cli.StringFlag{
			Name:   "ws-origin, w",
			Usage:  "WebSocket Origin URL",
			EnvVar: "GOTTY_CLIENT_WS_ORIGIN",
		},
		cli.BoolFlag{
			Name:   "webconsole",
			Usage:  "Connect webconsole server",
			EnvVar: "WEBCONSOLE",
		},
		cli.StringFlag{
			Name:  "server-addr,s",
			Usage: "SSHD server addr",
		},
		cli.StringFlag{
			Name:   "username,u",
			Usage:  "SSHD username",
			EnvVar: "WEBCONSOLE_SSH_USERNAME",
		},
		cli.StringFlag{
			Name:   "password,p",
			Usage:  "SSH password",
			EnvVar: "WEBCONSOLE_SSH_PASSWORD",
		},
	}

	app.Action = action

	app.Run(os.Args)
}

func action(c *cli.Context) error {
	if len(c.Args()) != 1 {
		logrus.Fatalf("usage: gotty-client [GoTTY URL]")
	}

	// setting up logrus
	logrus.SetOutput(os.Stderr)
	if c.Bool("debug") {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}

	// create Client
	url := c.Args()[0]
	client, err := gottyclient.NewClient(url)
	if err != nil {
		logrus.Fatalf("Cannot create client: %v", err)
	}

	if c.Bool("skip-tls-verify") {
		client.SkipTLSVerify = true
	}

	if c.Bool("use-proxy-from-env") {
		client.UseProxyFromEnv = true
	}

	if c.Bool("v2") {
		client.V2 = true
	}

	if wsOrigin := c.String("ws-origin"); wsOrigin != "" {
		client.WSOrigin = wsOrigin
	}

	if detachKey := c.String("detach-keys"); detachKey != "" {
		escapeKeys, err := gottyclient.ToBytes(detachKey)
		if err != nil {
			logrus.Warnf("Invalid escape key: %v", err)
		} else {
			client.EscapeKeys = escapeKeys
		}
	}

	// webconsole 模式参数
	if c.Bool("webconsole") {
		client.Webconsole = true
	}

	if serverAddr := c.String("server-addr"); serverAddr != "" {
		client.ServerAddr = serverAddr
	}

	if username := c.String("username"); username != "" {
		client.Username = username
	}

	if password := c.String("password"); password != "" {
		client.Password = password
	}

	// loop
	defer log.Println("9999999999")
	defer os.Exit(0)
	if err = client.Loop(); err != nil {
		logrus.Fatalf("Communication error: %v", err)
	}
	return nil
}
