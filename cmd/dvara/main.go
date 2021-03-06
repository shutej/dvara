package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/facebookgo/gangliamr"
	"github.com/facebookgo/inject"
	"github.com/facebookgo/startstop"
	"github.com/facebookgo/stats"
	"github.com/shutej/dvara"
)

func main() {
	if err := Main(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func Main() error {
	addrs := flag.String("addrs", "localhost:27017", "comma separated list of mongo addresses")
	certFile := flag.String("cert_file", "", "path to the certificate file")
	clientIdleTimeout := flag.Duration("client_idle_timeout", 60*time.Minute, "idle timeout for client connections")
	getLastErrorTimeout := flag.Duration("get_last_error_timeout", time.Minute, "timeout for getLastError pinning")
	keyFile := flag.String("key_file", "", "path to the key file")
	maxConnections := flag.Uint("max_connections", 100, "maximum number of connections per mongo")
	maxPerClientConnections := flag.Uint("max_per_client_connections", 1, "maximum number of connections from a single client")
	messageTimeout := flag.Duration("message_timeout", 2*time.Minute, "timeout for one message to be proxied")
	password := flag.String("password", "", "")
	portEnd := flag.Int("port_end", 6010, "end of port range")
	portStart := flag.Int("port_start", 6000, "start of port range")
	serverClosePoolSize := flag.Uint("server_close_pool_size", 1, "number of goroutines that will handle closing server connections.")
	serverIdleTimeout := flag.Duration("server_idle_timeout", 60*time.Minute, "duration after which a server connection will be considered idle")
	username := flag.String("username", "", "")

	flag.Parse()

	replicaSet := dvara.ReplicaSet{
		Addrs:                   *addrs,
		CertFile:                *certFile,
		ClientIdleTimeout:       *clientIdleTimeout,
		GetLastErrorTimeout:     *getLastErrorTimeout,
		KeyFile:                 *keyFile,
		MaxConnections:          *maxConnections,
		MaxPerClientConnections: *maxPerClientConnections,
		MessageTimeout:          *messageTimeout,
		Password:                *password,
		PortEnd:                 *portEnd,
		PortStart:               *portStart,
		ServerClosePoolSize:     *serverClosePoolSize,
		ServerIdleTimeout:       *serverIdleTimeout,
		Username:                *username,
	}

	var log stdLogger
	var graph inject.Graph
	err := graph.Provide(
		&inject.Object{Value: &log},
		&inject.Object{Value: &replicaSet},
		&inject.Object{Value: &stats.HookClient{}},
	)
	if err != nil {
		return err
	}
	if err := graph.Populate(); err != nil {
		return err
	}
	objects := graph.Objects()

	// Temporarily setup the metrics against a test registry.
	gregistry := gangliamr.NewTestRegistry()
	for _, o := range objects {
		if rmO, ok := o.Value.(registerMetrics); ok {
			rmO.RegisterMetrics(gregistry)
		}
	}
	if err := startstop.Start(objects, &log); err != nil {
		return err
	}
	defer startstop.Stop(objects, &log)

	ch := make(chan os.Signal, 2)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)
	<-ch
	signal.Stop(ch)
	return nil
}

type registerMetrics interface {
	RegisterMetrics(r *gangliamr.Registry)
}
