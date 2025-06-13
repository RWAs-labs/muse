package musecore

import (
	"context"

	observertypes "github.com/RWAs-labs/muse/x/observer/types"
)

func (c *Client) GetOperationalFlags(ctx context.Context) (observertypes.OperationalFlags, error) {
	res, err := c.Observer.OperationalFlags(ctx, &observertypes.QueryOperationalFlagsRequest{})
	return res.OperationalFlags, err
}
