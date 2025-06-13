package runner

import (
	"fmt"
	"net/http"

	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gagliardetto/solana-go/rpc"
	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"

	tonrunner "github.com/RWAs-labs/muse/e2e/runner/ton"
	btcclient "github.com/RWAs-labs/muse/museclient/chains/bitcoin/client"
	musecore_rpc "github.com/RWAs-labs/muse/pkg/rpc"
)

// Clients contains all the RPC clients and gRPC clients for E2E tests
type Clients struct {
	Musecore musecore_rpc.Clients

	// the RPC clients for external chains in the localnet
	BtcRPC  *btcclient.Client
	Solana  *rpc.Client
	Evm     *ethclient.Client
	EvmAuth *bind.TransactOpts
	TON     *tonrunner.Client
	Sui     sui.ISuiAPI

	// the RPC clients for MuseChain
	Mevm     *ethclient.Client
	MevmAuth *bind.TransactOpts

	MuseclientMetrics *MetricsClient
}

type MetricsClient struct {
	URL string
}

// Fetch retrieves and parses the prometheus metrics from the provided URL
func (m *MetricsClient) Fetch() (map[string]*dto.MetricFamily, error) {
	// Fetch metrics
	resp, err := http.Get(m.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch metrics: %w", err)
	}
	defer resp.Body.Close()

	// Parse metrics
	parser := expfmt.TextParser{}
	metricFamilies, err := parser.TextToMetricFamilies(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse metrics: %w", err)
	}

	return metricFamilies, nil
}

// FetchGauge fetches and individual gauge metric by it's name
func (m *MetricsClient) FetchGauge(name string) (float64, error) {
	metrics, err := m.Fetch()
	if err != nil {
		return 0, err
	}
	metric, ok := metrics[name]
	if !ok {
		return 0, fmt.Errorf("%s metric is not found", name)
	}
	gauge := metric.Metric[0].Gauge
	if gauge == nil {
		return 0, fmt.Errorf("%s metric is not a gauge", name)
	}
	return *gauge.Value, nil
}
