package main

func main() {
    broker, err := NewBroker("amqp://guest:guest@localhost:5672/")
    failOnError(err, "Failed to connect")

    defer broker.Conn.Close()
    defer broker.Ch.Close()

    err = broker.StartConsumer("hello")
    failOnError(err, "Consumer failed")
}
