package main

import (
	"flag"
	"fmt"
  "net"
  "bufio"
  "os"
  "strings"
)

var reader *bufio.Reader
var writer *bufio.Writer

type User struct {
  name string
  id string
}

func main() {
	var (
		ip = flag.String("b", "127.0.0.1", "Base IP Address")
		port = flag.String("p", "8080", "Port")
	)
	flag.Parse()

  reader = bufio.NewReader(os.Stdin)
  writer = bufio.NewWriter(os.Stdin)

  var full_address = fmt.Sprintf("%s:%s", *ip, *port)
	fmt.Printf("Attempting to connect to server on %s...\n", full_address)
	fmt.Println("Hello, world!")

  connectToServer(*ip, *port)
}

func connectToServer(ip string, port string) {
  conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", ip, port))
  defer conn.Close()

  if err != nil {
    fmt.Printf("Error: %s\n", err)
    return
  }

  logIn(conn)
  go handleMessages(conn)

  for {
    printMessage("")
    text, _ := reader.ReadString('\n')
    sendMessage(conn, text)
    if(strings.Trim(text, " \r\n") == "quit") {
      fmt.Printf("Closing connection...")
      return
    }
  }
}

func logIn(conn net.Conn) {
  id, _ := bufio.NewReader(conn).ReadBytes('\n')
  fmt.Printf("Your id is %s", id)
}

func sendMessage(conn net.Conn, text string) {
  message := strings.Trim(text, " \r\n")
  fmt.Fprintf(conn, "%s\n", message)
}

func handleMessages(conn net.Conn) {
  for {
    response, err := bufio.NewReader(conn).ReadString('\n')

    printMessage(response)
    if err != nil {
      if err.Error() == "Error: EOF" {
        fmt.Printf("Error: Disconnected from Host")
      }
      fmt.Printf("Error: %s", err)
      return
    }
  }
}

func handleConnection(conn net.Conn) {
  fmt.Printf("Looks like someone is trying to connect...\n Connection details: %s", conn)
}

func printMessage(message string) {
  var tmp []byte
  reader.Read(tmp)
  fmt.Printf("\r%s", message)
  fmt.Printf("Me: ")
  writer.Write(tmp)
}
