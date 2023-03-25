BUILDER=./.builder
RULES=go
include $(BUILDER)/rules.mk
$(BUILDER)/rules.mk:
	-go run github.com/endobit/builder@latest init

build::
	$(GO_BUILD) .

.PHONY: cobra.go
cobra.go: cobra.yaml
	go run . --import="github.com/endobit/cgen/internal/gen" > cobra
	mv cobra cobra.go

generate:: cobra.go


