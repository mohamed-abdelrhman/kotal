package v1alpha1

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Chainlink node defaulting", func() {
	It("Should default node", func() {

		node := Node{
			ObjectMeta: metav1.ObjectMeta{
				Name: "my-node",
			},
			Spec: NodeSpec{
				CertSecretName: "my-certificate",
			},
		}

		node.Default()

		Expect(node.Spec.TLSPort).To(Equal(DefaultTLSPort))
		// resources
		Expect(node.Spec.Resources.CPU).To(Equal(DefaultNodeCPURequest))
		Expect(node.Spec.Resources.CPULimit).To(Equal(DefaultNodeCPULimit))
		Expect(node.Spec.Resources.Memory).To(Equal(DefaultNodeMemoryRequest))
		Expect(node.Spec.Resources.MemoryLimit).To(Equal(DefaultNodeMemoryLimit))
		Expect(node.Spec.Resources.Storage).To(Equal(DefaultNodeStorageRequest))

	})
})
