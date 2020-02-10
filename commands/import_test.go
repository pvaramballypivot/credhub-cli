package commands_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
)

var _ = Describe("Import", func() {
	BeforeEach(func() {
		login()
	})

	ItRequiresAuthentication("import", "-f", "../test/test_import_file.yml")
	ItRequiresAnAPIToBeSet("import", "-f", "../test/test_import_file.yml")

	Describe("importing a file with mixed credentials", func() {
		Describe("when importing yaml", func() {
			It("sets all the credentials", func() {
				setUpImportRequests()

				session := runCommand("import", "-f", "../test/test_import_file.yml")

				Eventually(session).Should(Exit(0))

				Expect(string(session.Out.Contents())).To(Equal(`Import complete.
Successfully set: 7
Failed to set: 0
`))
			})
		})
		Describe("when importing json", func() {
			It("sets all the credentials", func() {
				setUpImportRequests()

				session := runCommand("import", "-f", "../test/test_import_file.json", "-j")

				Eventually(session).Should(Exit(0))

				Expect(string(session.Out.Contents())).To(Equal(`Import complete.
Successfully set: 7
Failed to set: 0
`))
			})
		})
	})

	Describe("when the yaml file starts with ---", func() {
		It("sets all the credentials", func() {
			setUpImportRequests()

			session := runCommand("import", "-f", "../test/test_import_file_with_document_end.yml")

			Eventually(session).Should(Exit(0))

			Expect(string(session.Out.Contents())).To(Equal(`Import complete.
Successfully set: 7
Failed to set: 0
`))
		})
	})

	Describe("when the yaml file starts with --- and spaces for Iryna", func() {
		It("sets all the credentials", func() {
			setUpImportRequests()

			session := runCommand("import", "-f", "../test/test_import_file_with_document_end_and_spaces.yml")

			Eventually(session).Should(Exit(0))

			Expect(string(session.Out.Contents())).To(Equal(`Import complete.
Successfully set: 7
Failed to set: 0
`))
		})
	})

	Describe("when importing file with no name specified", func() {
		Describe("when importing yaml", func() {
			It("passes through the server error", func() {
				jsonBody := `{"name":"","type":"password","value":"test-password"}`
				SetupPutBadRequestServer(jsonBody)

				session := runCommand("import", "-f", "../test/test_import_missing_name.yml")

				Eventually(session.Out).Should(Say(`test error`))
			})
		})
		Describe("when importing json", func() {
			It("passes through the server error", func() {
				jsonBody := `{"name":"","type":"password","value":"test-password"}`
				SetupPutBadRequestServer(jsonBody)

				session := runCommand("import", "-f", "../test/test_import_missing_name.json", "-j")

				Eventually(session.Out).Should(Say(`test error`))
			})
		})
	})

	Describe("when importing file with incorrect structure", func() {
		Describe("when importing yaml", func() {
			It("returns an error message", func() {
				errorMessage := "The referenced file does not contain valid yaml structure. Please update and retry your request."

				session := runCommand("import", "-f", "../test/test_import_incorrect_yaml.yml")

				Eventually(session.Err).Should(Say(errorMessage))
			})
		})
		Describe("when importing json", func() {
			It("returns an error message", func() {
				errorMessage := `The referenced file does not contain valid json structure. Please update and retry your request.`

				session := runCommand("import", "-f", "../test/test_import_incorrect_json.json", "-j")

				Eventually(session.Err).Should(Say(errorMessage))
			})
		})
	})

	Describe("when some credentials fail to set it prints errors in summary", func() {
		Describe("when importing yaml", func() {
			It("should display error message", func() {
				request := `{"type":"invalid_type","name":"/test/invalid_type","value":"some string"}`
				request1 := `{"type":"invalid_type","name":"/test/invalid_type1","value":"some string"}`
				SetupPutBadRequestServer(request)
				SetupPutBadRequestServer(request1)
				SetupPutUserServer("/test/user", `{"username": "covfefe", "password": "test-user-password"}`, "covfefe", "test-user-password", "P455W0rd-H45H")

				session := runCommand("import", "-f", "../test/test_import_partial_fail_set.yml")
				summaryMessage := `Import complete.
Successfully set: 1
Failed to set: 2
`
				Eventually(session.Out).Should(Say(`Credential '/test/invalid_type' at index 0 could not be set: test error`))
				Eventually(session.Out).Should(Say(`Credential '/test/invalid_type1' at index 1 could not be set: test error`))
				Expect(session.ExitCode()).ToNot(Equal(0))
				Expect(session.Out.Contents()).NotTo(ContainSubstring(`id: 5a2edd4f-1686-4c8d-80eb-5daa866f9f86`))
				Expect(session.Err.Contents()).To(ContainSubstring(`One or more credentials failed to import`))
				Eventually(session.Out).Should(Say(summaryMessage))
			})
		})
		Describe("when importing json", func() {
			It("should display error message", func() {
				request := `{"type":"invalid_type","name":"/test/invalid_type","value":"some string"}`
				request1 := `{"type":"invalid_type","name":"/test/invalid_type1","value":"some string"}`
				SetupPutBadRequestServer(request)
				SetupPutBadRequestServer(request1)
				SetupPutUserServer("/test/user", `{"username": "covfefe", "password": "test-user-password"}`, "covfefe", "test-user-password", "P455W0rd-H45H")

				session := runCommand("import", "-f", "../test/test_import_partial_fail_set.json", "-j")
				summaryMessage := `Import complete.
Successfully set: 1
Failed to set: 2
`
				Eventually(session.Out).Should(Say(`Credential '/test/invalid_type' at index 0 could not be set: test error`))
				Eventually(session.Out).Should(Say(`Credential '/test/invalid_type1' at index 1 could not be set: test error`))
				Expect(session.ExitCode()).ToNot(Equal(0))
				Expect(session.Out.Contents()).NotTo(ContainSubstring(`id: 5a2edd4f-1686-4c8d-80eb-5daa866f9f86`))
				Expect(session.Err.Contents()).To(ContainSubstring(`One or more credentials failed to import`))
				Eventually(session.Out).Should(Say(summaryMessage))
			})
		})
	})

	Describe("when no credential tag present in import file", func() {
		Describe("when importing yaml", func() {
			It("prints correct error message", func() {
				session := runCommand("import", "-f", "../test/test_import_incorrect_format.yml")

				noCredentialTagError := "The referenced file does not contain valid yaml structure. Please update and retry your request."
				Eventually(session.Err).Should(Say(noCredentialTagError))
			})
		})
		Describe("when importing json", func() {
			It("prints correct error message", func() {
				session := runCommand("import", "-f", "../test/test_import_incorrect_format.json", "-j")

				noCredentialTagError := `The referenced file does not contain valid json structure. Please update and retry your request.`
				Eventually(session.Err).Should(Say(noCredentialTagError))
			})
		})
	})

	Describe("when importing an ssh type with key public_key_fingerprint", func() {
		Describe("when importing yaml", func() {
			It("ignore public_key_fingerprint", func() {
				SetupPutSshServer("/test/sshCred", "ssh", "some-key", "some-private-key")

				session := runCommand("import", "-f", "../test/test_import_ssh_type_with_public_key_fingerprint.yml")
				Eventually(session).Should(Exit(0))
				Expect(string(session.Out.Contents())).To(Equal(`Import complete.
Successfully set: 1
Failed to set: 0
`))
			})
		})
		Describe("when importing json", func() {
			It("ignore public_key_fingerprint", func() {
				SetupPutSshServer("/test/sshCred", "ssh", "some-key", "some-private-key")

				session := runCommand("import", "-f", "../test/test_import_ssh_type_with_public_key_fingerprint.json", "-j")
				Eventually(session).Should(Exit(0))
				Expect(string(session.Out.Contents())).To(Equal(`Import complete.
Successfully set: 1
Failed to set: 0
`))
			})
		})

	})

	Describe("when importing a user type with password_hash", func() {
		Describe("when importing yaml", func() {
			It("ignore password_hash", func() {
				SetupPutUserServer("/test/userCred", `{"username": "sample-username", "password": "test-user-password"}`, "sample-username", "test-user-password", "P455W0rd-H45H")

				session := runCommand("import", "-f", "../test/test_import_user_type_with_password_hash.yml")
				Eventually(session).Should(Exit(0))
				Expect(string(session.Out.Contents())).To(Equal(`Import complete.
Successfully set: 1
Failed to set: 0
`))
			})
		})
		Describe("when importing json", func() {
			It("ignore password_hash", func() {
				SetupPutUserServer("/test/userCred", `{"username": "sample-username", "password": "test-user-password"}`, "sample-username", "test-user-password", "P455W0rd-H45H")

				session := runCommand("import", "-f", "../test/test_import_user_type_with_password_hash.json", "-j")
				Eventually(session).Should(Exit(0))
				Expect(string(session.Out.Contents())).To(Equal(`Import complete.
Successfully set: 1
Failed to set: 0
`))
			})
		})

	})

	Describe("when importing a value as a integer", func() {
		Describe("when importing yaml", func() {
			It("casts int to string and successfully imports", func() {
				SetupPutValueServer("/test/intStringValue", "value", "123")

				session := runCommand("import", "-f", "../test/test_import_with_int_for_value.yml")

				Eventually(session).Should(Exit(0))

				Expect(string(session.Out.Contents())).To(Equal(`Import complete.
Successfully set: 1
Failed to set: 0
`))
			})
		})
		Describe("when importing json", func() {
			It("casts int to string and successfully imports", func() {
				SetupPutValueServer("/test/intStringValue", "value", "123")

				session := runCommand("import", "-f", "../test/test_import_with_int_for_value.json", "-j")

				Eventually(session).Should(Exit(0))

				Expect(string(session.Out.Contents())).To(Equal(`Import complete.
Successfully set: 1
Failed to set: 0
`))
			})
		})

	})

	Describe("when importing certificate chain", func() {
		Context("and leaf comes after signing CA", func() {
			Describe("when importing yaml", func() {
				It("imports the signing CA first", func() {
					SetupPutCertificateServer("/root_ca",
						"root ca",
						"root certificate",
						"root private key")
					SetupPutCertificateWithCaNameForImportServer("/intermediate_ca",
						"/root_ca",
						"intermediate certificate",
						"intermediate private key")
					SetupPutCertificateWithCaNameForImportServer("/leaf_cert",
						"/intermediate_ca",
						"leaf certificate",
						"leaf private key")

					session := runCommand("import", "-f", "../test/certificate-chain.yml")
					Eventually(session).Should(Exit(0))
					Expect(string(session.Out.Contents())).To(Equal(`Import complete.
Successfully set: 3
Failed to set: 0
`))

				})
			})
			Describe("when importing json", func() {
				It("imports the signing CA first", func() {
					SetupPutCertificateServer("/root_ca",
						"root ca",
						"root certificate",
						"root private key")
					SetupPutCertificateWithCaNameForImportServer("/intermediate_ca",
						"/root_ca",
						"intermediate certificate",
						"intermediate private key")
					SetupPutCertificateWithCaNameForImportServer("/leaf_cert",
						"/intermediate_ca",
						"leaf certificate",
						"leaf private key")

					session := runCommand("import", "-f", "../test/certificate-chain.json", "-j")
					Eventually(session).Should(Exit(0))
					Expect(string(session.Out.Contents())).To(Equal(`Import complete.
Successfully set: 3
Failed to set: 0
`))

				})
			})

		})
	})
})

func setUpImportRequests() {
	SetupPutValueServer("/test/password", "password", "test-password-value")
	SetupPutValueServer("/test/value", "value", "test-value")
	SetupPutCertificateServer("/test/certificate",
		`ca-certificate`,
		`certificate`,
		`private-key`)
	SetupPutRsaServer("/test/rsa", "rsa", "public-key", "private-key")
	SetupPutRsaServer("/test/ssh", "ssh", "ssh-public-key", "private-key")
	SetupPutUserServer("/test/user", `{"username": "covfefe", "password": "test-user-password"}`, "covfefe", "test-user-password", "P455W0rd-H45H")
	setupPutJsonServer("/test/json", `{"1":"key is not a string","3.14":"pi","true":"key is a bool","arbitrary_object":{"nested_array":["array_val1",{"array_object_subvalue":"covfefe"}]}}`)
}
