GOPATH	= $(CURDIR)
BINDIR	= $(CURDIR)/bin

PROGRAMS = hetzner-cloud-metrics

depend:
	env GOPATH=$(GOPATH) go get github.com/logrusorgru/aurora
	env GOPATH=$(GOPATH) go get github.com/olekukonko/tablewriter

build:
	test -d $(BINDIR) || mkdir -p $(BINDIR)
	env GOPATH=$(GOPATH) go install $(PROGRAMS)

destdirs:
	mkdir -p -m 0755 $(DESTDIR)/usr/bin

strip: build
	strip --strip-all $(BINDIR)/hetzner-cloud-metrics

install: strip destdirs install-bin

install-bin:
	install -m 0755 $(BINDIR)/hetzner-cloud-metrics $(DESTDIR)/usr/bin

clean:
	/bin/rm -f bin/hetzner-cloud-metrics

distclean: clean
	rm -rf src/github.com/
	rm -rf src/gopkg.in/
	rm -rf src/golang.org/
	test -d pkg && chmod -R u+w pkg/ && rm -rf pkg/

uninstall:
	/bin/rm -f $(DESTDIR)/usr/bin

all: depend build strip install

