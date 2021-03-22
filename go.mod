module test

go 1.16

require (
	GoRESTClient/configHandler v0.0.0-00010101000000-000000000000 // indirect
	GoRESTClient/httpClient v0.0.0-00010101000000-000000000000 // indirect
	github.com/aarzilli/nucular v0.0.0-20210224090343-aa83b964abc8
	github.com/go-xmlfmt/xmlfmt v0.0.0-20191208150333-d5b6f63a941b
	github.com/syndtr/goleveldb v1.0.0 // indirect
	github.com/yosssi/gohtml v0.0.0-20201013000340-ee4748c638f4
	golang.org/x/image v0.0.0-20210220032944-ac19c3e999fb
	golang.org/x/mobile v0.0.0-20210220033013-bdb1ca9a1e08
)

replace GoRESTClient/httpClient => ./httpClient

replace GoRESTClient/configHandler => ./configHandler
