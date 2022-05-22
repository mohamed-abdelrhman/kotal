package controllers

import (
	"context"

	aptosv1alpha1 "github.com/kotalco/kotal/apis/aptos/v1alpha1"
	"github.com/kotalco/kotal/controllers/shared"
	"gopkg.in/yaml.v2"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Waypoint struct {
	FromConfig string `yaml:"from_config"`
}

type Execution struct {
	GenesisFileLocation string `yaml:"genesis_file_location"`
}

type Base struct {
	Role     string   `yaml:"role"`
	DataDir  string   `yaml:"data_dir"`
	Waypoint Waypoint `yaml:"waypoint"`
}

type Identity struct {
	Type   string `yaml:"type"`
	Key    string `yaml:"key"`
	PeerId string `yaml:"peer_id"`
}

type Network struct {
	NetworkId       string   `yaml:"network_id"`
	DiscoveryMethod string   `yaml:"discovery_method"`
	Identity        Identity `yaml:"identity,omitempty"`
}

type Config struct {
	Base             Base      `yaml:"base"`
	Execution        Execution `yaml:"execution"`
	FullNodeNetworks []Network `yaml:"full_node_networks,omitempty"`
}

// ConfigFromSpec generates config.toml file from node spec
func ConfigFromSpec(node *aptosv1alpha1.Node, client client.Client) (config string, err error) {
	var role string
	if node.Spec.Validator {
		role = "validator"
	} else {
		role = "full_node"
	}

	var nodePrivateKey string
	var identity Identity
	if node.Spec.NodePrivateKeySecretName != "" {
		key := types.NamespacedName{
			Name:      node.Spec.NodePrivateKeySecretName,
			Namespace: node.Namespace,
		}

		if nodePrivateKey, err = shared.GetSecret(context.Background(), client, key, "key"); err != nil {
			return
		}

		identity = Identity{
			Type: "from_config",
			Key:  nodePrivateKey,
			// TODO: update with peer ID
		}

	}

	c := Config{
		Base: Base{
			Role:    role,
			DataDir: "/opt/aptos/data",
			Waypoint: Waypoint{
				FromConfig: node.Spec.Waypoint,
			},
		},
		Execution: Execution{
			GenesisFileLocation: "/opt/aptos/config/genesis.blob",
		},
		FullNodeNetworks: []Network{
			{
				NetworkId:       "public",
				DiscoveryMethod: "onchain",
				Identity:        identity,
			},
		},
	}

	data, err := yaml.Marshal(&c)
	if err != nil {
		return
	}

	config = string(data)
	return
}
