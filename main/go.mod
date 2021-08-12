module main

go 1.14

require (
	github.com/denisenkom/go-mssqldb v0.10.0
	leapsy.com/databases v0.0.0-00010101000000-000000000000 // indirect
	leapsy.com/packages/configurations v0.0.0-00010101000000-000000000000 // indirect
	leapsy.com/packages/logings v0.0.0-00010101000000-000000000000 // indirect
	leapsy.com/packages/network v0.0.0-00010101000000-000000000000 // indirect
	leapsy.com/records v0.0.0-00010101000000-000000000000 // indirect
	leapsy.com/servers v0.0.0-00010101000000-000000000000
	leapsy.com/times v0.0.0-00010101000000-000000000000 // indirect
)

//servers
replace leapsy.com/packages/logings => ../LeapsyPackages/logings

replace leapsy.com/servers => ../LeapsyPackages/servers

replace leapsy.com/databases => ../LeapsyPackages/databases

replace leapsy.com/packages/configurations => ../LeapsyPackages/configurations

replace leapsy.com/packages/network => ../LeapsyPackages/network

replace leapsy.com/records => ../LeapsyPackages/records

replace leapsy.com/times => ../LeapsyPackages/times
