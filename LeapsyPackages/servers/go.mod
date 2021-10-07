module servers

go 1.14

replace leapsy.com/packages/logings => ../logings

replace leapsy.com/databases => ../databases

replace leapsy.com/packages/configurations => ../configurations

replace leapsy.com/packages/network => ../network

replace leapsy.com/records => ../records

replace leapsy.com/times => ../times

replace leapsy.com/packages/paths => ../paths

require (
	github.com/gin-gonic/gin v1.7.4
	github.com/robfig/cron v1.2.0
	github.com/shogo82148/androidbinary v1.0.2
	github.com/sirupsen/logrus v1.8.1
	go.mongodb.org/mongo-driver v1.7.2
	gopkg.in/mgo.v2 v2.0.0-20190816093944-a6b53ec6cb22
	leapsy.com/databases v0.0.0-00010101000000-000000000000
	leapsy.com/packages/configurations v0.0.0-00010101000000-000000000000
	leapsy.com/packages/logings v0.0.0-00010101000000-000000000000
	leapsy.com/packages/network v0.0.0-00010101000000-000000000000
	leapsy.com/packages/paths v0.0.0-00010101000000-000000000000
	leapsy.com/records v0.0.0-00010101000000-000000000000
	leapsy.com/times v0.0.0-00010101000000-000000000000
)
