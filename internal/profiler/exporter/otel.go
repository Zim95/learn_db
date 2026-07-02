package exporter

/*
	OpenTelemetry exporter — FUTURE COMPONENT (intentionally not implemented).

	OpenTelemetry shines for distributed tracing: spans propagated across
	service boundaries and shipped to a collector/backend (Jaeger, Tempo, ...).
	A single-process CLI that does one operation and exits has nothing to trace
	across and nowhere to ship to, so implementing this now would add a heavy
	dependency for no signal.

	This file reserves the package boundary. It becomes worthwhile if learn_db
	ever participates in a larger system (e.g. called as part of a request that
	spans multiple services), where each operation would become a span.
*/
