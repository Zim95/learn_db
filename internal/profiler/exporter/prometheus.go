package exporter

/*
	Prometheus exporter — FUTURE COMPONENT (intentionally not implemented).

	Prometheus works by scraping a long-running process's /metrics endpoint on
	an interval. This project is a one-shot CLI: it runs a single command and
	exits in milliseconds, so there is no process for Prometheus to scrape and
	no time window over which metrics accumulate.

	This file exists to reserve the package boundary in the architecture. It
	becomes real only if learn_db grows a long-running server mode (e.g. a
	socket/HTTP daemon that serves many operations). At that point this would
	register counters/histograms (op latency, bytes allocated, ...) and expose
	them via an HTTP handler.
*/
