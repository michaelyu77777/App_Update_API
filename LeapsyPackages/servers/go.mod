module servers

go 1.14

replace leapsy.com/packages/logings => ../logings

replace leapsy.com/servers => ../servers

replace leapsy.com/databases => ../databases

replace leapsy.com/packages/configurations => ../configurations

replace leapsy.com/packages/network => ../network

replace leapsy.com/records => ../records

replace leapsy.com/times => ../times

replace leapsy.com/packages/model => ../model

require (
	github.com/gin-gonic/gin v1.7.3
	github.com/robfig/cron v1.2.0
	github.com/sirupsen/logrus v1.8.1
	leapsy.com/databases v0.0.0-00010101000000-000000000000
	leapsy.com/packages/configurations v0.0.0-00010101000000-000000000000
	leapsy.com/packages/logings v0.0.0-00010101000000-000000000000
	leapsy.com/packages/model v0.0.0-00010101000000-000000000000
	leapsy.com/packages/network v0.0.0-00010101000000-000000000000
	leapsy.com/records v0.0.0-00010101000000-000000000000
	leapsy.com/times v0.0.0-00010101000000-000000000000
)
