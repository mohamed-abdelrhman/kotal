package bitcoin

import (
	"os"

	bitcoinv1alpha1 "github.com/kotalco/kotal/apis/bitcoin/v1alpha1"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Bitcoin core client", func() {

	node := &bitcoinv1alpha1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "bitcoin-node",
			Namespace: "default",
		},
		Spec: bitcoinv1alpha1.NodeSpec{
			Network: "mainnet",
			RPCPort: 7777,
		},
	}

	node.Default()
	client := NewClient(node)

	It("Should get correct image", func() {
		// default image
		img := client.Image()
		Expect(img).To(Equal(DefaultBitcoinCoreImage))
		// after setting custom image
		testImage := "kotalco/bitcoin-core:test"
		os.Setenv(EnvBitcoinCoreImage, testImage)
		img = client.Image()
		Expect(img).To(Equal(testImage))
	})

	It("Should get correct command", func() {
		Expect(client.Command()).To(BeNil())
	})

	It("Should get correct home directory", func() {
		Expect(client.HomeDir()).To(Equal(BitcoinCoreHomeDir))
	})

	It("Should generate correct client arguments", func() {
		Expect(client.Args()).To(ContainElements([]string{
			"-chain=main",
			"-datadir=/home/bitcoin/kotal-data",
			"-rpcport=7777",
		}))
	})

})