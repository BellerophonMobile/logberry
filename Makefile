
sources=$(wildcard *.go)

all: test \
     bin \
     bin/minimal \
     bin/component \
     bin/task \
     bin/threaded \
     bin/multiplexer \
     bin/toplevel \
     bin/blueberry \
     bin/flightpath

#-----------------------------------------------------------------------
#-----------------------------------------------------------------------
test:
	go test -v

#-----------------------------------------------------------------------
#-----------------------------------------------------------------------
bin/%: examples/%/build.go examples/%/main.go $(sources)
	cd bin; go build github.com/BellerophonMobile/logberry/examples/$(subst bin/,,$@)


#-----------------------------------------------------------------------
#-----------------------------------------------------------------------
bin:
	mkdir bin

%build.go: 
	./util/build-stmt-go.sh > $@

.PHONY: all test
.SECONDARY:
