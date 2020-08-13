module github.com/hashlash/descarca/server1

go 1.14

require (
	github.com/google/uuid v1.1.1
	github.com/hashlash/descarca v0.0.0-00010101000000-000000000000
	github.com/streadway/amqp v1.0.0
)

replace github.com/hashlash/descarca => ../
