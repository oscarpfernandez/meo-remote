package main

import (
	"fmt"
	"log"
	"net"

	term "github.com/nsf/termbox-go"
)

const (
	boxAddress = "192.168.1.64:8082"
)

func main() {
	if err := term.Init(); err != nil {
		log.Fatalf("Failed to start terminal hanndler: %v", err)
	}
	defer term.Close()

	fmt.Println("--- Meo Remote")
	fmt.Println("--- Connecting to box: " + boxAddress)

	conn, err := net.Dial("tcp", boxAddress)
	if err != nil {
		log.Fatalf("Failed to connect to box: %v", err)
	}
	defer conn.Close()

	fmt.Println("--- Connected!")

	sendCmdKey := func(key string) {
		fmt.Printf("<<< Command Key=%s\n", key)
		if _, err := conn.Write([]byte("key=" + key + "\n")); err != nil {
			log.Fatal(err)
		}
	}

	go func() {
		for {
			buff := make([]byte, 10)
			if _, err := conn.Read(buff); err != nil {
				return
			}
			fmt.Printf(">>> Response: %s", string(buff))
		}
	}()

	standardKeys := map[string]string{
		"0": "48", "1": "49", "2": "50", "3": "51", "4": "52", "5": "53", "6": "54", "7": "55", "8": "56", "9": "57",
		"a": "97", "b": "98", "c": "99", "e": "100", "d": "101", "f": "102", "g": "103", "h": "104", "i": "105",
		"j": "106", "k": "107", "l": "108", "m": "109", "n": "110", "o": "111", "p": "112", "q": "113", "r": "114",
		"s": "115", "t": "116", "u": "117", "v": "118", "w": "119", "x": "120", "y": "121", "z": "122",
	}

	specialKey := map[term.Key]string{
		term.KeyEsc:        "0",   // special exit code.
		term.KeySpace:      "32",  // space.
		term.KeyEnter:      "13",  // enter.
		term.KeyArrowUp:    "38",  // up.
		term.KeyArrowDown:  "40",  // down.
		term.KeyArrowLeft:  "37",  // left.
		term.KeyArrowRight: "39",  // right.
		term.KeyTab:        "36",  // menu.
		term.KeyPgup:       "33",  // program up
		term.KeyPgdn:       "34",  // program down.
		term.KeyDelete:     "46",  // delete.
		term.KeyCtrlR:      "140", // red button.
		term.KeyCtrlG:      "141", // green button.
		term.KeyCtrlY:      "142", // yellow button.
		term.KeyCtrlB:      "143", // blue button.
		term.KeyBackspace:  "166", // browser back.
		term.KeyHome:       "8",   // back.
		term.KeyCtrlQ:      "175", // volume up.
		term.KeyCtrlA:      "174", // volume down.
		term.KeyEnd:        "233", // power.
	}

	for {
		switch ev := term.PollEvent(); ev.Type {
		case term.EventKey:
			key, ok := specialKey[ev.Key]
			if !ok {
				key, ok := standardKeys[string(ev.Ch)]
				if !ok {
					continue
				}
				sendCmdKey(key)
				continue
			}
			if key == "0" {
				return
			}
			sendCmdKey(key)
		}
	}
}
