
sources=$(wildcard *.go)

all: bin bin/minimal bin/component bin/task


#-----------------------------------------------------------------------
#-----------------------------------------------------------------------
bin/minimal: examples/minimal/build.go examples/minimal/main.go $(sources)
	cd bin; go build github.com/BellerophonMobile/logberry/examples/minimal

bin/component: examples/component/build.go examples/component/main.go $(sources)
	cd bin; go build github.com/BellerophonMobile/logberry/examples/component

bin/task: examples/task/build.go examples/task/main.go $(sources)
	cd bin; go build github.com/BellerophonMobile/logberry/examples/task


#-----------------------------------------------------------------------
#-----------------------------------------------------------------------
bin:
	mkdir bin

%build.go: 
	./util/build-stmt-go.sh > $@

.PHONY: all
