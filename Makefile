
lib_sources=$(wildcard *.go)

examples= minimal               \
          component             \
          task                  \
          threaded              \
          multiplexer           \
          toplevel              \
          blueberry             \
          flightpath


#-----------------------------------------------------------------------
#-----------------------------------------------------------------------
all: test examples

test:
	go test -v

examples: bin $(addprefix bin/, $(examples))


#-----------------------------------------------------------------------
#-----------------------------------------------------------------------
bin/%: examples/%/build.go examples/%/main.go $(lib_sources)
	cd bin; go build github.com/BellerophonMobile/logberry/examples/$(subst bin/,,$@)


#-----------------------------------------------------------------------
#-----------------------------------------------------------------------
bin:
	mkdir bin

%build.go: 
	./util/build-stmt-go.sh > $@

.PHONY: all test examples
.SECONDARY:
