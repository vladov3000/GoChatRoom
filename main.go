package main

import (
	"os"
	"net"
	"fmt"
	"bufio"
	"time"
	"strings"
)

const SERVER_FIRST byte = 1

type routineRes struct {
    text string
    first byte
}

func serverRountine(reader *bufio.Reader, ch chan<- routineRes) {
	for {
		text, err := reader.ReadString('\n')

		if (err != nil) {
			time.Sleep(5 * time.Second)
			text = fmt.Sprintf("Error reading from connection: %s", err)
			ch <- routineRes{ text, SERVER_FIRST }
			return
		} 
		ch <- routineRes{ text, SERVER_FIRST }
	}
}

func userRoutine(reader *bufio.Reader, ch chan<- routineRes) {
	for {	
		text, err := reader.ReadString('\n')

		if (err != nil) {
			time.Sleep(5 * time.Second)
			text = fmt.Sprintf("Error reading from connection: %s", err)
			ch <- routineRes{ text, 1 - SERVER_FIRST }
			return
		}

		ch <- routineRes{ text, 1 - SERVER_FIRST }
	}
}

func consolePrint(text string) {
	fmt.Print("\x1bM") // move down one line (note: down is up the screen)
	fmt.Print("\x1b[J") // clear everything below cursor
	fmt.Printf("  %s\n> ", text)
}

func main() {
	host := "localhost:8080"
	conn, err := net.Dial("tcp", host)
	if err != nil {
		fmt.Printf("Could not dial %s: %s\n", host, err)
		os.Exit(1)
	}

	fmt.Print("\x1b[2J") // clear screen
	fmt.Print("> ")

	serverReader := bufio.NewReader(conn)
	stdinReader := bufio.NewReader(os.Stdin)

	ch := make(chan routineRes)

	go serverRountine(serverReader, ch)
	go userRoutine(stdinReader, ch)
	
	for {
		res := <-ch
		text := strings.TrimSpace(res.text)

		if (res.first == SERVER_FIRST) { 
			fmt.Print("\x1b[2K") // clear everything right of cursor
			fmt.Print("\x1b[2D") // move cursor two columns left
		} else {
			fmt.Print("\x1bM") // move down one line (note: down is up the screen)
			fmt.Print("\x1b[J") // clear everything below cursor
		}
		fmt.Printf("%s\n> ", text)
	}

}