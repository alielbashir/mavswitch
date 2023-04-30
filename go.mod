module github.com/alielbashir/mavswitch

go 1.20

require github.com/bluenviron/gomavlib/v2 v2.0.1

require (
	github.com/tarm/serial v0.0.0-20180830185346-98f6abe2eb07 // indirect
	golang.org/x/sys v0.1.0 // indirect
)

replace github.com/bluenviron/gomavlib/v2 v2.0.1 => github.com/alielbashir/gomavlib/v2 v2.1.1
