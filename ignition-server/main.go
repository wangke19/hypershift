package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"go.uber.org/zap/zapcore"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	"github.com/openshift/hypershift/ignition-server/cmd"
)

func main() {
	ctrl.SetLogger(zap.New(zap.UseDevMode(true), zap.JSONEncoder(func(o *zapcore.EncoderConfig) {
		o.EncodeTime = zapcore.RFC3339TimeEncoder
	})))

	root := &cobra.Command{
		Use:   "debug",
		Short: "A program for debugging the ignition server",
	}

	root.AddCommand(cmd.NewStartCommand())
	root.AddCommand(cmd.NewRunLocalIgnitionProviderCommand())

	if err := root.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
