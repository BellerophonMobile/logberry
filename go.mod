module "logberry"

require (
  "github.com/jroimartin/gocui" v0.3.0
)

replace "github.com/jroimartin/gocui" v0.3.0 => "github.com/tjkopena/gocui" v0.3.1
