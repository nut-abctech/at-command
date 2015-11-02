package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/nut-abctech/at-command/Godeps/_workspace/src/github.com/barnybug/gogsmmodem"
	"github.com/nut-abctech/at-command/Godeps/_workspace/src/github.com/codegangsta/cli"
	"github.com/nut-abctech/at-command/Godeps/_workspace/src/github.com/tarm/goserial"
)

type server struct {
	net.Conn
	Reader *bufio.Reader
	Writer *bufio.Writer
}

func (s *server) SendCmd(cmd string) {
	fmt.Fprint(s.Writer, cmd)
}

func (s *server) Receiver() {
	bytes, err := s.Reader.ReadByte()
	if err != nil {
		log.Println(err)
	}
	log.Println(string(bytes))
}

func main() {
	app := cli.NewApp()
	app.Name = "AT Command tools"
	app.Usage = "Test sim server"
	app.Version = "0.0.1"
	app.Commands = []cli.Command{
		{
			Name:  "run-tcp",
			Usage: "Listen SMS from specific server via tcp protocal",
			Action: func(c *cli.Context) {
				var host string = c.String("host")
				if c.String("host") == "" {
					log.Println("Unknow host name. Please specific host")
					os.Exit(1)
				}
				fmt.Printf("Listening SMS from : %s \n", host)
				for {
					conn, err := net.Dial("tcp", host)
					if err != nil {
						log.Panicf("Error connecting %s %s", host, err)
					}
					handleConn(conn)
				}
			},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "host, H",
					Usage: "Specific module ip address. e.g. 127.0.0.1:8080",
				},
			},
		},
		{
			Name:  "run-serial",
			Usage: "Listen SMS from specific device via serial port",
			Action: func(c *cli.Context) {
				var name = c.String("name")
				var baud = c.Int("baud")
				log.Printf("usb: %s, Baud: %d", name, baud)
				conf := serial.Config{
					Name: name,
					Baud: baud,
				}
				modem, err := gogsmmodem.Open(&conf, true)
				if err != nil {
					log.Panic(err)
				}
				defer modem.Close()
				for packet := range modem.OOB {
					log.Printf("Received : %#v\n", packet)
					switch p := packet.(type) {
					case gogsmmodem.MessageNotification:
						log.Println("Message notification:", p)
						msg, err := modem.GetMessage(p.Index)
						if err == nil {
							fmt.Printf("Message from %s: %s\n", msg.Telephone, msg.Body)
							modem.DeleteMessage(p.Index)
						}
					}
				}
			},
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:  "baud, b",
					Usage: "Baud number e.g 115200",
				},
				cli.StringFlag{
					Name:  "name, n",
					Usage: "e.g. /dev/ttyUSB1",
				},
			},
		},
	}
	app.RunAndExitOnError()
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	// for {
	// 	server := &server{
	// 		Conn:   conn,
	// 		Reader: bufio.NewReader(conn),
	// 		Writer: bufio.NewWriter(conn),
	// 	}
	// }
}
