module main

go 1.14

//servers
replace leapsy.com/packages/logings => ../LeapsyPackages/logings

replace leapsy.com/servers => ../LeapsyPackages/servers

replace leapsy.com/databases => ../LeapsyPackages/databases

replace leapsy.com/packages/configurations => ../LeapsyPackages/configurations

replace leapsy.com/packages/network => ../LeapsyPackages/network

replace leapsy.com/records => ../LeapsyPackages/records

replace leapsy.com/times => ../LeapsyPackages/times

replace leapsy.com/packages/model => ../LeapsyPackages/model

require (
	github.com/denisenkom/go-mssqldb v0.10.0
	leapsy.com/servers v0.0.0-00010101000000-000000000000
)
