Economy, gambling, waifu and so on features for YAGDPDB

To use:

1. Add 	"github.com/jonas747/yageconomy" to the imports in yags main.go
2. Add "yageconomy.RegisterPlugin()" at the bottom of all the other "RegisterFunctions()" calls in yags main.go

If using docker att the following to the Dockerfile

`COPY --from=builder /go/src/github.com/jonas747/yageconomy/*/assets/*.html templates/plugins/`

on one of the lines under "# Handle templates for plugins automatically"

