package sseoutput

import (
	"net/http"
	"github.com/BellerophonMobile/logberry"
	"github.com/BellerophonMobile/sse"
)

type Options struct {
	HistoryLimit int
}

type SSEOutput struct {

	internalroot *logberry.Root
	log *logberry.Task

	receiveroot *logberry.Root

	server *sse.EventServer
	
}

func New(options *Options) (*SSEOutput,error) {

	if options == nil {
		options = &Options{}
	}

	sseopts := sse.EventServerOptions{
		HistoryLimit: options.HistoryLimit,
	}
	
	sse := &SSEOutput{
		internalroot: logberry.NewRoot(11),
		server: sse.NewEventServer(&sseopts),
	}

	sse.internalroot.AddOutputDriver(logberry.NewStdOutput("sseoutput"))
	sse.log = sse.internalroot.Component("SSEOutput")
	
	return sse,nil
}

func (x *SSEOutput) Attach(root *logberry.Root) {
	x.log.Info("Attached")
	x.receiveroot = root
}

func (x *SSEOutput) Detach() {
	x.log.Info("Detached")	
	x.receiveroot = nil
}

func (x *SSEOutput) Event (event *logberry.Event) {

	x.server.JSONMessage(event)
	
}

func (x *SSEOutput) Handler() func(w http.ResponseWriter, r *http.Request) {
	return x.server.Handle
}
