module github.com/alexkreidler/wiz

go 1.13

require (
	github.com/StackExchange/wmi v0.0.0-20190523213315-cbe66965904d // indirect
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751
	github.com/alexkreidler/deepcopy v0.0.0-20191229055930-5d6df38e9bd6
	github.com/alexkreidler/jsonscrape v0.0.0-20200109233942-83dfc73db91b
	github.com/davecgh/go-spew v1.1.1
	github.com/emicklei/go-restful/v3 v3.0.0
	github.com/gin-gonic/gin v1.5.0
	github.com/go-ole/go-ole v1.2.4 // indirect
	github.com/hashicorp/go-getter v1.4.0
	github.com/iancoleman/strcase v0.0.0-20191112232945-16388991a334
	github.com/imdario/mergo v0.3.8
	github.com/joeybloggs/go-download v2.1.0+incompatible
	github.com/json-iterator/go v1.1.9
	github.com/mitchellh/go-homedir v1.1.0
	github.com/mitchellh/mapstructure v1.1.2
	github.com/segmentio/ksuid v1.0.2
	github.com/shirou/gopsutil v2.19.11+incompatible
	github.com/shirou/w32 v0.0.0-20160930032740-bb4de0191aa4 // indirect
	github.com/spf13/cobra v0.0.5
	github.com/swaggo/swag v1.6.4
	golang.org/x/net v0.0.0-20191209160850-c0dbc17a3553 // indirect
	gonum.org/v1/gonum v0.6.1
	gopkg.in/src-d/go-git.v4 v4.13.1
	gopkg.in/yaml.v2 v2.2.7
	gotest.tools v2.2.0+incompatible
)

replace github.com/alexkreidler/jsonscrape => ../jsonscrape
