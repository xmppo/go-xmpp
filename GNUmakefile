#
#	Makefile for hookAPI
#
# switches:
#	define the ones you want in the CFLAGS definition...
#
#	TRACE		- turn on tracing/debugging code
#
#
#
#

# Version for distribution
VER=1_0r1

MAKEFILE=GNUmakefile

# We Use Compact Memory Model

all: bin/example
	@[ -d bin ] || exit

bin/example: _example/example.go xmpp.go
	@[ -d bin ] || mkdir bin
	go build -o $@ _example/example.go
	@strip $@ || echo "example OK"

clean:

distclean: clean
	@rm -rf bin
