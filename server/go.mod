module github.com/mnbbrown/healthsignal/server

require (
	github.com/go-chi/chi v3.3.3+incompatible
	github.com/influxdata/influxdb v1.6.2
	github.com/jmoiron/sqlx v0.0.0-20180614180643-0dae4fefe7c0
	github.com/lib/pq v1.0.0
	github.com/mnbbrown/healthsignal/healthsignal v0.0.0 //indirect
	github.com/rs/cors v1.5.0
	golang.org/x/net v0.0.0-20180906233101-161cd47e91fd
	golang.org/x/sys v0.0.0-20180909124046-d0be0721c37e // indirect
	google.golang.org/grpc v1.14.0
)

replace github.com/mnbbrown/healthsignal/healthsignal => ../healthsignal
