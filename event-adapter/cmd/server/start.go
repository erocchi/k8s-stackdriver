/*
Copyright 2017 The Kubernetes Authors.
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

package server

import (
	"io"
	"fmt"

	"github.com/spf13/cobra"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	coreclient "k8s.io/client-go/kubernetes/typed/core/v1"
	stackdriver "google.golang.org/api/monitoring/v3"

	"github.com/GoogleCloudPlatform/k8s-stackdriver/event-adapter/pkg/cmd/server"
	"github.com/GoogleCloudPlatform/k8s-stackdriver/event-adapter/cmd/provider"
	"time"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// NewCommandStartMaster provides a CLI handler for 'start master' command
func NewCommandStartSampleAdapterServer(out, errOut io.Writer, stopCh <-chan struct{}) *cobra.Command {
	baseOpts := server.NewEventsAdapterServerOptions(out, errOut)
	o := SampleAdapterServerOptions{
		EventsAdapterServerOptions: baseOpts,
	}

	cmd := &cobra.Command{
		Short: "Launch the events API adapter server",
		Long:  "Launch the events API adapter server",
		RunE: func(c *cobra.Command, args []string) error {
			if err := o.Complete(); err != nil {
				return err
			}
			if err := o.Validate(args); err != nil {
				return err
			}
			if err := o.RunEventsAdapterServer(stopCh); err != nil {
				return err
			}
			return nil
		},
	}

	flags := cmd.Flags()
	o.SecureServing.AddFlags(flags)
	o.Authentication.AddFlags(flags)
	o.Authorization.AddFlags(flags)
	o.Features.AddFlags(flags)
	fmt.Println(o.RemoteKubeConfigFile)
	flags.StringVar(&o.RemoteKubeConfigFile, "lister-kubeconfig", o.RemoteKubeConfigFile, ""+
		"kubeconfig file pointing at the 'core' kubernetes server with enough rights to list "+
		"any described objets")

	return cmd
}

func (o SampleAdapterServerOptions) RunEventsAdapterServer(stopCh <-chan struct{}) error {
	config, err := o.Config()
	if err != nil {
		return err
	}

	var clientConfig *rest.Config
	if len(o.RemoteKubeConfigFile) > 0 {
		loadingRules := &clientcmd.ClientConfigLoadingRules{ExplicitPath: o.RemoteKubeConfigFile}
		loader := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, &clientcmd.ConfigOverrides{})

		clientConfig, err = loader.ClientConfig()
	} else {
		clientConfig, err = rest.InClusterConfig()
	}
	if err != nil {
		return fmt.Errorf("unable to construct lister client config to initialize provider: %v", err)
	}

	client, err := coreclient.NewForConfig(clientConfig)
	if err != nil {
		return fmt.Errorf("unable to construct lister client to initialize provider: %v", err)
	}
	oauthClient := oauth2.NewClient(oauth2.NoContext, google.ComputeTokenSource(""))
	stackdriverService, err := stackdriver.New(oauthClient)
	if err != nil {
		return fmt.Errorf("Failed to create Stackdriver client: %v", err)
	}
	evProvider := provider.NewStackdriverProvider(client.RESTClient(), stackdriverService, 5 * time.Minute)

	server, err := config.Complete().New(evProvider)
	if err != nil {
		return err
	}
	return server.GenericAPIServer.PrepareRun().Run(stopCh)
}

type SampleAdapterServerOptions struct {
	*server.EventsAdapterServerOptions

	// RemoteKubeConfigFile is the config used to list pods from the master API server
	RemoteKubeConfigFile string
}