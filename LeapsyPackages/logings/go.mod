module servers

go 1.14

require (
	github.com/jonboulle/clockwork v0.2.2 // indirect
	github.com/keepeye/logrus-filename v0.0.0-20190711075016-ce01a4391dd1
	github.com/lestrrat-go/file-rotatelogs v2.4.0+incompatible
	github.com/lestrrat-go/strftime v1.0.4 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/rifflock/lfshook v0.0.0-20180920164130-b9218ef580f5
	github.com/sirupsen/logrus v1.8.1
)

replace leapsy.com/packages/logings => ../logings

replace leapsy.com/servers => ../servers

replace leapsy.com/databases => ../databases

replace leapsy.com/packages/configurations => ../configurations

replace leapsy.com/packages/network => ../network

replace leapsy.com/records => ../records

replace leapsy.com/times => ../times