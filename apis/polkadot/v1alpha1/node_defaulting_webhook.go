package v1alpha1

import (
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// +kubebuilder:webhook:path=/mutate-polkadot-kotal-io-v1alpha1-node,mutating=true,failurePolicy=fail,groups=polkadot.kotal.io,resources=nodes,verbs=create;update,versions=v1alpha1,name=mutate-polkadot-v1alpha1-node.kb.io,sideEffects=None,admissionReviewVersions=v1

var _ webhook.Defaulter = &Node{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *Node) Default() {
	nodelog.Info("default", "name", r.Name)

}
