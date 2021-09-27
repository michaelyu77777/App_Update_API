module tlss

go 1.16

replace leapsy.com/packages/configurations => ../configurations

replace leapsy.com/packages/logings => ../logings

replace leapsy.com/packages/network => ../network

require (
	github.com/sirupsen/logrus v1.8.1
	leapsy.com/packages/logings v0.0.0-00010101000000-000000000000
	leapsy.com/packages/network v0.0.0-00010101000000-000000000000
)
