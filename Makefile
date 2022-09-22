# tracklist
# See LICENSE for copyright and license details.
.POSIX:

include config.mk

all: tracklist

tracklist:
	$(GO) build $(GOFLAGS)

clean:
	$(RM) tracklist

install: all
	mkdir -p $(DESTDIR)$(PREFIX)/bin
	cp -f tracklist $(DESTDIR)$(PREFIX)/bin
	chmod 755 $(DESTDIR)$(PREFIX)/bin/tracklist

uninstall:
	$(RM) $(DESTDIR)$(PREFIX)/bin/tracklist

.DEFAULT_GOAL := all

.PHONY: all tracklist clean install uninstall
