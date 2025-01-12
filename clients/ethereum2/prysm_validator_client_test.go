package ethereum2

import (
	"fmt"
	"os"

	ethereum2v1alpha1 "github.com/kotalco/kotal/apis/ethereum2/v1alpha1"
	sharedAPI "github.com/kotalco/kotal/apis/shared"
	"github.com/kotalco/kotal/controllers/shared"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Prysm validator client", func() {

	validator := &ethereum2v1alpha1.Validator{
		Spec: ethereum2v1alpha1.ValidatorSpec{
			Client:          ethereum2v1alpha1.PrysmClient,
			Network:         "mainnet",
			BeaconEndpoints: []string{"http://localhost:8899"},
			Graffiti:        "Validated by Kotal",
			Keystores: []ethereum2v1alpha1.Keystore{
				{
					SecretName: "my-validator",
				},
			},
			WalletPasswordSecret: "wallet-password",
			CertSecretName:       "my-cert",
			Logging:              sharedAPI.ErrorLogs,
		},
	}

	validator.Default()
	client, _ := NewClient(validator)

	It("Should get correct image", func() {
		// default image
		img := client.Image()
		Expect(img).To(Equal(DefaultPrysmValidatorImage))
		// after changing .spec.image
		testImage := "kotalco/prysm:spec"
		validator.Spec.Image = &testImage
		img = client.Image()
		Expect(img).To(Equal(testImage))
		// after setting custom image
		testImage = "kotalco/prysm:test"
		os.Setenv(EnvPrysmValidatorImage, testImage)
		img = client.Image()
		Expect(img).To(Equal(testImage))
	})

	It("Should get correct command", func() {
		Expect(client.Command()).To(ConsistOf("validator"))
	})

	It("Should get correct env", func() {
		Expect(client.Env()).To(BeNil())
	})

	It("Should get correct home dir", func() {
		Expect(client.HomeDir()).To(Equal(PrysmHomeDir))
	})

	It("Should generate correct client arguments", func() {
		args := client.Args()

		Expect(args).To(ContainElements([]string{
			PrysmAcceptTermsOfUse,
			PrysmDataDir,
			shared.PathData(client.HomeDir()),
			"--mainnet",
			PrysmBeaconRPCProvider,
			"http://localhost:8899",
			PrysmGraffiti,
			"Validated by Kotal",
			PrysmWalletDir,
			fmt.Sprintf("%s/prysm-wallet", shared.PathData(client.HomeDir())),
			PrysmWalletPasswordFile,
			fmt.Sprintf("%s/prysm-wallet/prysm-wallet-password.txt", shared.PathSecrets(client.HomeDir())),
			PrysmLogging,
			string(sharedAPI.ErrorLogs),
			PrysmTLSCert,
			fmt.Sprintf("%s/cert/tls.crt", shared.PathSecrets(client.HomeDir())),
		}))

	})

})
