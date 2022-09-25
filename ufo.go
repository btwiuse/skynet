package ufo

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/marten-seemann/webtransport-go"
)

var _ net.Listener = (*listener)(nil)

func Serve(u string, handler http.Handler) error {
	ln, err := Listen(u)
	if err != nil {
		return err
	}
	log.Println("listening on", ln.URL())
	if handler == nil {
		handler = http.DefaultServeMux
	}
	return http.Serve(ln, handler)
}

func Listen(u string) (*listener, error) {
	// localhost:3000 will be parsed by net/url as URL{Scheme: localhost, Port: 3000}
	// hence the hack
	if !strings.Contains(u, "://") {
		u = "http://" + u
	}
	up, err := url.Parse(u)
	if err != nil {
		return nil, err
	}
	ctx, _ := context.WithTimeout(context.TODO(), 3*time.Second)
	session, err := Dial(ctx, up, nil)
	if err != nil {
		return nil, err
	}
	stm0, err := session.AcceptStream(ctx)
	if err != nil {
		return nil, err
	}
	errchan := make(chan string)
	hostchan := make(chan string)
	// go io.Copy(os.Stdout, stm0)
	go func() {
		scanner := bufio.NewScanner(stm0)
		for scanner.Scan() {
			line := scanner.Text()
			// ignore server pings
			if line == "PING" {
				continue
			}
			if strings.HasPrefix(line, "HOST ") {
				hostchan <- strings.TrimPrefix(line, "HOST ")
				continue
			}
			if strings.HasPrefix(line, "ERR ") {
				errchan <- strings.TrimPrefix(line, "ERR ")
				continue
			}
			log.Println("stm0: unknown command:", line)
		}
	}()
	// go io.Copy(stm0, os.Stdin)
	ln := &listener{
		session: session,
		stm0:    stm0,
		scheme:  up.Scheme,
		port:    getport(up),
	}
	select {
	case emsg := <-errchan:
		return nil, fmt.Errorf("server: %s", emsg)
	case ln.host = <-hostchan:
		return ln, nil
	}
}

func getport(u *url.URL) string {
	_, p, ok := strings.Cut(u.Host, ":")
	if ok {
		return ":" + p
	}
	return ""
}

type listener struct {
	session *webtransport.Session
	stm0    webtransport.Stream
	scheme  string
	host    string
	port    string
}

func (l *listener) Accept() (net.Conn, error) {
	stream, err := l.session.AcceptStream(context.Background())
	if err != nil {
		return nil, err
	}
	return &StreamConn{stream, l.session}, nil
}

func (l *listener) Close() error {
	return l.session.Close()
}

// Addr returns listener itself which is an implementor of net.Addr
func (l *listener) Addr() net.Addr {
	return l
}

// Network returns the protocol scheme, either http or https
func (l *listener) Network() string {
	return l.scheme
}

// String returns the host(:port) address of listener
func (l *listener) String() string {
	return l.host + l.port
}

// URL returns the public accessible address of the listener
func (l *listener) URL() string {
	return l.Network() + "://" + l.String()
}
