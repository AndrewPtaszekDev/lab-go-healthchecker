package internal

// example data
var seedServices = []service{
	service{}.init("calibre", healthcheck{}.init("calibre.ptaszek.studio")),
	service{}.init("grafana", healthcheck{}.init("grafana.ptaszek.studio")),
	service{}.init("notes", healthcheck{}.init("notes.ptaszek.studio")),
	service{}.init("gitea", healthcheck{}.init("gitea.ptaszek.studio")),
	service{}.init("studio", healthcheck{}.init("ptaszek.studio")),
	service{}.init("cloudflare", healthcheck{}.init("1.1.1.1")),
}
