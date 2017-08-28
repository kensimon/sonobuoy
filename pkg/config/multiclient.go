/*
Copyright 2017 Heptio Inc.

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

package config

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type MultiClient struct {
	CoreClient       kubernetes.Interface
	ExtensionsClient apiextensionsclient.Interface
}

// LoadMultiClient multiple kubernetes clientsets using given sonobuoy configuration
func LoadMultiClient(cfg *Config) (ret MultiClient, err error) {
	var config *rest.Config

	// 1 - gather config information used to initialize
	kubeconfig := viper.GetString("kubeconfig")
	if len(kubeconfig) > 0 {
		cfg.Kubeconfig = kubeconfig
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	} else {
		config, err = rest.InClusterConfig()
	}
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	// 2 - creates the clientset from kubeconfig
	ret.CoreClient, err = kubernetes.NewForConfig(config)
	if err != nil {
		return
	}

	// 3 - creates the extensions client from kubeconfig
	ret.ExtensionsClient, err = apiextensionsclient.NewForConfig(config)
	if err != nil {
		return
	}

	return
}
