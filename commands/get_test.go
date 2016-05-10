package commands_test

import (
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
	. "github.com/onsi/gomega/ghttp"
)

var _ = Describe("Get", func() {
	It("displays help", func() {
		session := runCommand("get", "-h")

		Eventually(session).Should(Exit(1))
		Expect(session.Err).To(Say("get"))
		Expect(session.Err).To(Say("--name"))
	})

	It("gets a secret", func() {
		responseJson := `{"potatoes":"delicious"}`

		server.AppendHandlers(
			CombineHandlers(
				VerifyRequest("GET", "/api/v1/secret/my-secret"),
				RespondWith(http.StatusOK, responseJson),
			),
		)

		session := runCommand("get", "-n", "my-secret")

		Eventually(session).Should(Exit(0))
		Eventually(session.Out).Should(Say(`"potatoes": "delicious"`))
	})
})
