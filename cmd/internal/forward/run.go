package forward

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/powerpuffpenguin/webpc/cmd/internal/client"
	"golang.org/x/term"
)

func Run(insecure bool,
	url, listen, remote,
	user, password string,
	heart int,
) {
	r := bufio.NewReader(os.Stdin)
	if !strings.HasPrefix(url, `ws://`) && !strings.HasPrefix(url, `wss://`) {
		url = inputURL(r, `connect websocket url: `)
	}
	if listen == `` {
		listen = inputString(r, `local listen address: `)
	}
	if remote == `` {
		remote = inputString(r, `remote connect address: `)
	}
	if user == `` {
		user = inputString(r, `user name: `)
	}
	if password == `` {
		password = inputPassword(r, `user password: `)
	}

	dialer, e := client.NewDialer(`forward`, url, insecure, user, password, heart)
	if e != nil {
		log.Fatalln(e)
	}

	l, e := net.Listen(`tcp`, listen)
	if e != nil {
		log.Fatalln(e)
	}
	fmt.Println(`listen on `, listen)

	newWorker(dialer).Serve(l, remote)
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
func inputURL(r *bufio.Reader, placeholder string) string {
	var v string
	for {
		fmt.Print(placeholder)
		v = readString(r)
		if strings.HasPrefix(v, `ws://`) || strings.HasPrefix(v, `wss://`) {
			break
		}
	}
	return v
}
func inputString(r *bufio.Reader, placeholder string) string {
	var v string
	for {
		fmt.Print(placeholder)
		v = readString(r)
		v = strings.TrimSpace(v)
		if v != `` {
			break
		}
	}
	return v
}
func inputPassword(r *bufio.Reader, placeholder string) string {
	fmt.Print(placeholder)
	b, e := term.ReadPassword(int(os.Stdin.Fd()))
	if e != nil {
		log.Fatalln(e)
	}
	fmt.Println()
	return string(b)
}
