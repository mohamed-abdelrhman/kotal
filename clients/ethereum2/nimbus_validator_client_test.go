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

var _ = Describe("Nimbus validator client", func() {
	validator := &ethereum2v1alpha1.Validator{
		Spec: ethereum2v1alpha1.ValidatorSpec{
			Client:          ethereum2v1alpha1.NimbusClient,
			Network:         "mainnet",
			BeaconEndpoints: []string{"http://nimbus-beacon-node"},
			Graffiti:        "Validated by Kotal",
			Keystores: []ethereum2v1alpha1.Keystore{
				{
					SecretName: "my-validator",
				},
			},
			Logging: sharedAPI.FatalLogs,
		},
	}

	validator.Default()
	client, _ := NewClient(validator)

	It("Should get correct image", func() {
		// default image
		img := client.Image()
		Expect(img).To(Equal(DefaultNimbusValidatorImage))
		// after changing .spec.image
		testImage := "kotalco/nimbus:spec"
		validator.Spec.Image = &testImage
		img = client.Image()
		Expect(img).To(Equal(testImage))
		// after setting custom image
		testImage = "kotalco/nimbus:test"
		os.Setenv(EnvNimbusValidatorImage, testImage)
		img = client.Image()
		Expect(img).To(Equal(testImage))
	})

	It("Should get correct command", func() {
		Expect(client.Command()).To(ConsistOf("nimbus_validator_client"))
	})

	It("Should get correct env", func() {
		Expect(client.Env()).To(BeNil())
	})

	It("Should get correct home dir", func() {
		Expect(client.HomeDir()).To(Equal(NimbusHomeDir))
	})

	It("Should generate correct client arguments", func() {

		args := client.Args()

		Expect(args).To(ContainElements([]string{
			NimbusNonInteractive,
			argWithVal(NimbusLogging, string(validator.Spec.Logging)),
			argWithVal(NimbusDataDir, shared.PathData(client.HomeDir())),
			argWithVal(NimbusBeaconNodes, "http://nimbus-beacon-node"),
			argWithVal(NimbusGraffiti, "Validated by Kotal"),
			argWithVal(NimbusValidatorsDir, fmt.Sprintf("%s/kotal-validators/validator-keys", shared.PathData(client.HomeDir()))),
			argWithVal(NimbusSecretsDir, fmt.Sprintf("%s/kotal-validators/validator-secrets", shared.PathData(client.HomeDir()))),
		}))

	})

})
