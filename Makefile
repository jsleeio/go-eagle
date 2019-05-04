binaries: schroff panelgen

schroff:
	go build ./cmd/schroff

panelgen:
	go build ./cmd/panelgen

clean:
	$(RM) schroff panelgen
