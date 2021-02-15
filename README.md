Meta-Viper is a wrapper over [viper](https://github.com/spf13/viper), it uses a [Go tagged struct](https://golang.org/ref/spec#Tag) to:

1.  Initialize [viper](https://github.com/spf13/viper)'s configuration.

2.  Load application configuration on that same struct.

You can find examples at [./examples](./examples) or in this article: https://carlosvin.github.io/posts/create-cmd-tool-golang.

![meta viper](https://pkg.go.dev/badge/github.com/carlosvin/meta-viper)

Usage {#_usage}
====================

Define the struct holding your application config
-------------------------------------------------

Meta-Viper will try to load the configuration in that struct from configuration files, environment variables or flags.

**Example.**

~~~ go
package main

import (
    "fmt"
    "os"

    config "github.com/carlosvin/meta-viper"
)

type appConfig struct {
    Host      string `cfg_name:"host" cfg_desc:"Server host"`
    Port      int    `cfg_name:"port" cfg_desc:"Server port"`
    SearchAPI string `cfg_name:"apis.search" cfg_desc:"Search API endpoint"`
}

func main() {
    cfg := &appConfig{ 
        Host: "localhost",
        Port: 6000,
        SearchAPI: "https://google.es",
    }

    _, err := config.New(cfg, os.Args) 
    if err != nil {
        panic(err)
    }
    log.Printf("Loaded Configuration %v...", cfg)
}
~~~

-   We instantiate the declared struct. As you can see you can optionally specify default values.

-   It loads the configuration in the passed struct `cfg`. The `os.Args` are required to parse the application flags.

Let’s focus on one application configuration attribute to explain the example. Meta-Viper will allow you to load the config into `Host` structure attribute in 3 different ways:

**Using flags.**

~~~ bash
./your-program --host=my.host
~~~

**Using environment variables.**

~~~ bash
HOST=my.host ./your-program
~~~

**Loading the data from a file in json, yaml or toml format.**

~~~ bash
./your-program --config=qa 
~~~

-   Following the qa configuration file content

**qa.json.**

~~~ json
{
    "host": "qa.host",
    "port": 8000,
    "apis": {
        "search": "https://search.api/"
    }
}
~~~

<div class="note" markdown="1">

You can combine flags, environment variables and configuration files.

</div>

Tags {#_tags}
----

### cfg_name {#_cfg_name}

Required tag to specify the configuration parameter name.

### cfg_description {#_cfg_description}

Optional tag to describe how to use the configuration parameter.
