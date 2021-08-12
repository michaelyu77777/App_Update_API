module servers

go 1.14

require (
	github.com/gin-gonic/gin v1.7.2
	github.com/robfig/cron v1.2.0
	github.com/sirupsen/logrus v1.8.1
	go.mongodb.org/mongo-driver v1.5.4
	leapsy.com/databases v0.0.0-00010101000000-000000000000
	leapsy.com/packages/configurations v0.0.0-00010101000000-000000000000
	leapsy.com/packages/logings v0.0.0-00010101000000-000000000000
	leapsy.com/packages/network v0.0.0-00010101000000-000000000000
	leapsy.com/records v0.0.0-00010101000000-000000000000
	leapsy.com/times v0.0.0-00010101000000-000000000000
)

replace leapsy.com/packages/logings => ../logings

replace leapsy.com/servers => ../servers

replace leapsy.com/databases => ../databases

replace leapsy.com/packages/configurations => ../configurations

replace leapsy.com/packages/network => ../network

replace leapsy.com/records => ../records

replace leapsy.com/times => ../times
