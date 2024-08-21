package infra

import (
	"context"
	"encoding/json"
	"net"
	"strconv"
	"time"

	"github.com/outofforest/logger"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"go.uber.org/zap"

	"github.com/sei-protocol/build/pkg/retry"
	"github.com/sei-protocol/build/pkg/tools/docker"
)

// AlpineImage returns the default docker image used to run containers.
func AlpineImage() string {
	return "alpine:" + docker.AlpineVersion
}

// WaitUntilHealthy waits until app is healthy or context expires.
func WaitUntilHealthy(ctx context.Context, hFuncs map[string]func(ctx context.Context) error) error {
	log := logger.Get(ctx)
	for name, hFunc := range hFuncs {
		ctx := logger.With(ctx, zap.String("app", name))
		log.Info("Waiting for app to start.")
		if err := retry.Do(ctx, time.Second, func() error {
			return hFunc(ctx)
		}); err != nil {
			return err
		}
	}
	return nil
}

// IsRunning returns a health check which succeeds if application is running.
func IsRunning(app *App) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		if app.info.Status == AppStatusRunning {
			return nil
		}
		return retry.Retryable(errors.New("application hasn't been started yet"))
	}
}

// JoinNetAddr joins protocol, hostname and port.
func JoinNetAddr(proto, hostname string, port int) string {
	if proto != "" {
		proto += "://"
	}
	return proto + net.JoinHostPort(hostname, strconv.Itoa(port))
}

// JoinNetAddrIP joins protocol, IP and port.
func JoinNetAddrIP(proto string, ip net.IP, port int) string {
	return JoinNetAddr(proto, ip.String(), port)
}

// PortsToMap converts structure containing port numbers to a map.
func PortsToMap(ports interface{}) map[string]int {
	unmarshaled := map[string]interface{}{}
	lo.Must0(json.Unmarshal(lo.Must(json.Marshal(ports)), &unmarshaled))

	res := map[string]int{}
	for k, v := range unmarshaled {
		res[k] = int(v.(float64))
	}
	return res
}

// AppSetToHealthChecks builds a list of healthcheck functions from the given AppSet.
func AppSetToHealthChecks(appSet AppSet) map[string]func(ctx context.Context) error {
	hFuncs := map[string]func(ctx context.Context) error{}
	for _, app := range appSet {
		if app.HealthCheckFunc == nil {
			hFuncs[app.Name] = IsRunning(app)
		} else {
			hFuncs[app.Name] = app.HealthCheckFunc
		}
	}
	return hFuncs
}
