
sources=$(wildcard *.go)

all: bin bin/minimal bin/component bin/task bin/threaded bin/multiplexer


#-----------------------------------------------------------------------
#-----------------------------------------------------------------------
bin/minimal: examples/minimal/build.go examples/minimal/main.go $(sources)
	cd bin; go build github.com/BellerophonMobile/logberry/examples/minimal

bin/component: examples/component/main.go $(sources)
	cd bin; go build github.com/BellerophonMobile/logberry/examples/component

bin/task: examples/task/main.go $(sources)
	cd bin; go build github.com/BellerophonMobile/logberry/examples/task

bin/threaded: examples/threaded/main.go $(sources)
	cd bin; go build github.com/BellerophonMobile/logberry/examples/threaded

bin/multiplexer: examples/multiplexer/main.go $(sources)
	cd bin; go build github.com/BellerophonMobile/logberry/examples/multiplexer


#-----------------------------------------------------------------------
#-----------------------------------------------------------------------
bin:
	mkdir bin

%build.go: 
	./util/build-stmt-go.sh > $@

.PHONY: all
