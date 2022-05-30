package main

import (
	"github.com/roy2220/virtual-kubelet/cmd/virtual-kubelet/internal/provider"
	"github.com/roy2220/virtual-kubelet/cmd/virtual-kubelet/internal/provider/mock"
)

func registerMock(s *provider.Store) {
	s.Register("mock", func(cfg provider.InitConfig) (provider.Provider, error) { //nolint:errcheck
		return mock.NewMockProvider(
			cfg.ConfigPath,
			cfg.NodeName,
			cfg.OperatingSystem,
			cfg.InternalIP,
			cfg.DaemonPort,
		)
	})
}
