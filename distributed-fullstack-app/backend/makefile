GO_FILES=broker.go sse.go enveloper.go 

gateway:
	clear && go run $(GO_FILES) routes.go gateway.go

promotions:
	clear && go run $(GO_FILES) promotions.go

notification:
	clear && go run $(GO_FILES) notification.go

ranking:
	clear && go run $(GO_FILES) ranking.go

clean:
	go clean