package forward

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func Run(insecure bool,
	url, listen, remote, user, password string,
) {
	r := bufio.NewReader(os.Stdin)

	if url == `` {
		url = inputURL(r)
	}
	if listen == `` {
		listen = inputAddr(r)
	}

	fmt.Println(`connect`, url)
	fmt.Println(`forward`, listen, `to`, remote)
	fmt.Println(`insecure`, insecure)
}
func readString(r *bufio.Reader) string {
	var result string
	for {
		b, _, e := r.ReadLine()
		if e != nil {
			log.Fatalln(e)
		}
		result = strings.TrimSpace(string(b))
		if result != `` {
			break
		}
	}
	return result
}
func inputURL(r *bufio.Reader) string {
	var v string
	for {
		fmt.Print(`connect websocket url: `)
		v = readString(r)
		if strings.HasPrefix(v, `ws://`) || strings.HasPrefix(v, `wss://`) {
			break
		}
	}
	return v
}
func inputAddr(r *bufio.Reader) string {
	var v string
	for {
		fmt.Print(`local listen address: `)
		v = readString(r)
		v = strings.TrimSpace(v)
		if v != `` {
			break
		}
	}
	return v
}
