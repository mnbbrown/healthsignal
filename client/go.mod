module github.com/mnbbrown/healthsignal/client

require (
	github.com/mnbbrown/healthsignal/healthsignal v0.0.0
	golang.org/x/net v0.0.0-20180906233101-161cd47e91fd
	golang.org/x/sys v0.0.0-20180909124046-d0be0721c37e // indirect
	google.golang.org/grpc v1.14.0
)

replace github.com/mnbbrown/healthsignal/healthsignal => ../healthsignal
