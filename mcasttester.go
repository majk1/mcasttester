package main

import (
	"net"
	"log"
	"os"
	"fmt"
	"os/signal"
	"time"
	"strings"
	"strconv"
)

const (
	DEFAULTADDR    = "224.0.0.1:9999"
	DEFAULTPAYLOAD = "DATA"
)

var gracefulStop = false

func main() {

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		fmt.Println("\nInterrupt signal cought, stopping...")
		gracefulStop = true
	}()
	
	commands := &Commands{}
	
	sendCommand := commands.AddCommand("send", "Use command \"send\" to send data to the specified multicast address")
	sendCommand.flags.StringVar(&sendCommand.addr, "addr", DEFAULTADDR, "Multicast address and port to send data")
	sendCommand.flags.StringVar(&sendCommand.payload, "data", DEFAULTPAYLOAD, "The string to send")
	sendCommand.flags.IntVar(&sendCommand.loopCount, "loop", 0, "Loop count (default is forever)")
	sendCommand.execFunc = func(addr string, payload string, loopCount int) {
		udpAddr := resolveAddr(addr)
		conn, e := net.DialUDP("udp4", nil, udpAddr)
		if e != nil {
			log.Fatal("Cannot connect to UDP address ", addr)
		}
		
		i := 0
		if loopCount > 0 {
			fmt.Printf("Sending %d payloads (\"%s\") to %s\n", loopCount, payload, udpAddr)
		} else {
			fmt.Printf("Sending payloads (\"%s\") to %s\n", payload, udpAddr)
		}
		for !gracefulStop && (loopCount <= 0 || loopCount > i) {
			bytes := []byte(payload + " [" + strconv.Itoa(i) + "]");
			conn.Write(bytes)
			fmt.Printf("[%4d] Sent %d bytes to %s -> \"%s\"\n", i + 1, len(bytes), udpAddr, string(bytes))
			time.Sleep(1 * time.Second)
			i++
		}
		fmt.Printf("Sent %d payloads\n", i)
		fmt.Print("Closing connection... ")
		conn.Close()
		fmt.Println("done")
	}
	
	receiveCommand := commands.AddCommand("receive", "Use command \"receive\" to listen on the specified multicast address")
	receiveCommand.flags.StringVar(&receiveCommand.addr, "addr", DEFAULTADDR, "Multicast address and port to listen on")
	receiveCommand.flags.IntVar(&receiveCommand.loopCount, "loop", 0, "Loop count (default is forever)")
	receiveCommand.execFunc = func(addr string, payload string, loopCount int) {
		udpAddr := resolveAddr(addr)
		conn, e := net.ListenMulticastUDP("udp4", nil, udpAddr)
		if e != nil {
			log.Fatal("Cannot listen on multicast UDP address ", addr)
		}

		i := 0
		buffer := make([]byte, 512)
		conn.SetReadBuffer(512)
		if loopCount > 0 {
			fmt.Printf("Receiving %d payloads from %s\n", loopCount, udpAddr)
		} else {
			fmt.Printf("Receiving payloads from %s\n", udpAddr)
		}
		for !gracefulStop && (loopCount <= 0 || loopCount > i) {
			conn.SetReadDeadline(time.Now().Add(1 * time.Second))
			n, src, err := conn.ReadFromUDP(buffer)
			if err != nil {
				if !strings.Contains(err.Error(), "i/o timeout") {
					log.Fatal("ReadFromUDP failed:", err)
				}
			} else {
				fmt.Printf("[%4d] Read %d bytes from %s <- \"%s\"\n", i + 1, n, src, string(buffer[:n]))
				i++
			}
		}
		fmt.Printf("Received %d payloads\n", i)
		fmt.Print("Stopping listener... ")
		conn.Close()
		fmt.Println("done")
	}
	
	if len(os.Args) < 2 {
		commands.PrintUsage(os.Args[0])
		os.Exit(1)
	}

	commandName := os.Args[1]
	command := commands.GetByName(commandName)
	if command != nil {
		command.ParseArgs(os.Args[2:])
		command.Execute()
	} else {
		fmt.Printf("Unknown command: %s\n\n", commandName)
		commands.PrintUsage(os.Args[0])
		os.Exit(2)
	}
}

func resolveAddr(addr string) *net.UDPAddr {
	udpAddr, e := net.ResolveUDPAddr("udp4", addr)
	if e != nil {
		log.Fatalf("Could not resolve address: %s", addr)
	}
	return udpAddr
}
