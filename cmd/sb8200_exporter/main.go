package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	sb8200exporter "github.com/nickvanw/sb8200_exporter"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	metricsPath = flag.String("metrics.path", "/metrics", "path to fetch metrics")
	metricsAddr = flag.String("metrics.addr", ":9292", "address to listen")

	modemAddr = flag.String("modem.host", "http://192.168.100.1/cmconnectionstatus.html", "url to fetch modem")
)

func main() {
	flag.Parse()

	c, err := createExporter(*modemAddr)
	if err != nil {
		log.Fatalf("unable to create client: %s", err)
	}

	prometheus.MustRegister(c)

	mux := http.NewServeMux()
	mux.Handle(*metricsPath, promhttp.Handler())
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>Arris SB8200 Exporter</title></head>
			<body>
			<h1>Arris SB8200 Exporter</h1>
			<p><a href='` + *metricsPath + `'>Metrics</a></p>
			</body>
			</html>`))
	})
	loggedMux := handlers.LoggingHandler(os.Stdout, mux)
	if err := http.ListenAndServe(*metricsAddr, loggedMux); err != nil {
		log.Fatalf("unable to start metrics server: %s", err)
	}
}

func createExporter(modem string) (*sb8200exporter.Exporter, error) {
	client := http.Client{}
	return sb8200exporter.New(&client, modem)
}
