module times

go 1.14

replace leapsy.com/databases => ../databases

replace leapsy.com/packages/logings => ../logings

replace leapsy.com/records => ../records

replace leapsy.com/packages/configurations => ../configurations

replace leapsy.com/packages/network => ../network

require (
	github.com/robfig/cron v1.2.0
	leapsy.com/databases v0.0.0-00010101000000-000000000000
	leapsy.com/packages/logings v0.0.0-00010101000000-000000000000
	leapsy.com/records v0.0.0-00010101000000-000000000000
)
