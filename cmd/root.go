/*
Copyright The Helm Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"errors"
	"io"
	"os"

	"github.com/spf13/cobra"
)

var (
	settings *EnvSettings
)

func NewRootCmd(out io.Writer, args []string) *cobra.Command {
	cmd := &cobra.Command{
		Use:          "2to3",
		Short:        "Migrate and Cleanup Helm v2 configuration and releases in-place to Helm v3",
		Long:         "Migrate and Cleanup Helm v2 configuration and releases in-place to Helm v3",
		SilenceUsage: true,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				return errors.New("no arguments accepted")
			}
			return nil
		},
	}

	flags := cmd.PersistentFlags()
	flags.Parse(args)
	settings = new(EnvSettings)

	// When run with the Helm plugin framework, Helm plugins are not passed the
	// plugin flags that correspond to Helm global flags e.g. helm 2to3 convert --kube-context ...
	// The flag values are set to corresponding environment variables instead.
	// The flags are passed as expected when run directly using the binary.
	// The below allows to use Helm's --kube-context global flag.
	if ctx := os.Getenv("HELM_KUBECONTEXT"); ctx != "" {
		settings.KubeContext = ctx
	}

	// Note that the plugin's --kubeconfig flag is set by the Helm plugin framework to
	// the KUBECONFIG environment variable instead of being passed into the plugin.
	// That variable is transparently handled by the helm-plugin-utils package so does not
	// need to be explicitely handled here.

	cmd.AddCommand(
		newCleanupCmd(out),
		newConvertCmd(out),
		newMoveConfigCmd(out),
	)

	return cmd
}
