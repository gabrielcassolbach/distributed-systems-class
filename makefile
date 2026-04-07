GO_FILES=broker.go enveloper.go

gateway:
	go run $(GO_FILES) gateway.go

promotions:
	go run $(GO_FILES) promotions.go

notification:
	go run $(GO_FILES) notification.go

client:
	go run broker.go client.go

clean:
	go clean