package module

import (
	"context"
	"fmt"

	"github.com/maxlandon/wiregost/client/connection"
	cliCtx "github.com/maxlandon/wiregost/client/context"
	pb "github.com/maxlandon/wiregost/proto/v1/gen/go/module"
)

type Use struct {
}

// Execute - Run
func (c *Use) Execute(args []string) (err error) {

	// Set request
	in := &pb.UseRequest{
		Path:   "Test",
		Client: cliCtx.Context.Client,
	}

	res, err := connection.ModuleRPC.UseModule(context.Background(), in)
	if err != nil {
		fmt.Println(err)
	}

	if res.Err != "" {
		fmt.Println(res.Err)
	}

	if res.Loaded != nil {
		cliCtx.Context.Menu = cliCtx.ModuleMenu
		cliCtx.Context.Module = *res.Loaded
		cliCtx.Context.NeedsCommandRefresh = true
	}

	return
}
