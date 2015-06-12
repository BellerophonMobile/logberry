
lib_sources=$(wildcard *.go)

examples= minimal                \
          small
#          component             \
#          task                  \
3          threaded              \
#          fanout                \
#          toplevel              \
#          blueberry             \
#          flightpath

gopath=${subst /src/github.com/BellerophonMobile/logberry,,${shell pwd}}

#-----------------------------------------------------------------------
#-----------------------------------------------------------------------
all: test examples

test:
	env GOPATH=$(gopath) go test -v

examples: bin $(addprefix bin/, $(examples))


#-----------------------------------------------------------------------
#-----------------------------------------------------------------------
bin/%: examples/%/build.go examples/%/main.go $(lib_sources)
	cd bin; env GOPATH=$(gopath) go build github.com/BellerophonMobile/logberry/examples/$(subst bin/,,$@)


#-----------------------------------------------------------------------
#-----------------------------------------------------------------------
bin:
	mkdir bin

%build.go: 
	./util/build-metadata-go.sh > $@

.PHONY: all test examples
.SECONDARY:
