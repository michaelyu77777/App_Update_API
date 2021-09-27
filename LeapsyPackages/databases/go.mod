module databases

go 1.14

replace leapsy.com/packages/logings => ../logings

replace leapsy.com/model => ../model

replace leapsy.com/packages/network => ../network

replace leapsy.com/records => ../records

replace leapsy.com/packages/configurations => ../configurations

require (
	github.com/sirupsen/logrus v1.8.1
	go.mongodb.org/mongo-driver v1.7.2
	leapsy.com/model v0.0.0-00010101000000-000000000000
	leapsy.com/packages/configurations v0.0.0-00010101000000-000000000000
	leapsy.com/packages/logings v0.0.0-00010101000000-000000000000
	leapsy.com/packages/network v0.0.0-00010101000000-000000000000
	leapsy.com/records v0.0.0-00010101000000-000000000000
)
