module servers

go 1.14

require (
	github.com/smartystreets/goconvey v1.6.4 // indirect
	gopkg.in/ini.v1 v1.62.0
	leapsy.com/packages/logings v0.0.0-00010101000000-000000000000
)

replace leapsy.com/packages/logings => ../logings

replace leapsy.com/servers => ../servers

replace leapsy.com/databases => ../databases

replace leapsy.com/packages/configurations => ../configurations

replace leapsy.com/packages/network => ../network

replace leapsy.com/records => ../records

replace leapsy.com/times => ../times
