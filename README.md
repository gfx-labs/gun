# gun


gun is a preconfigured configuration for aconfig https://github.com/cristalhq/aconfig


## Usage

```go
package main

import (
	"fmt"

	"github.com/gfx-labs/gun"
)

var Config struct {
	Host      string `yaml:"host" env:"HOST" json:"host" default:"localhost"`
	Port      int    `yaml:"port" env:"PORT" json:"port" default:"8080"`
	DbUrl     string `yaml:"db_url" env:"DB_URL" json:"db_url"`
	ManyItems []string
}

func init() {
	gun.Load(&Config)
}

func main() {
	fmt.Println(Config.Host, Config.Port)
}
```

With a prefix:

```go
gun.LoadPrefix(&Config, "MYAPP")
// reads env vars like MYAPP_HOST, MYAPP_PORT, etc.
// reads config files named myapp.yml, myapp.yaml, myapp.json
```
