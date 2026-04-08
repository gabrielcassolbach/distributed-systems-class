GO_FILES=broker.go enveloper.go

gateway:
	go run $(GO_FILES) gateway.go

promotions:
	go run $(GO_FILES) promotions.go

notification:
	go run $(GO_FILES) notification.go

ranking:
	go run $(GO_FILES) ranking.go

client:
	go run $(GO_FILES) client.go

clean:
	go clean