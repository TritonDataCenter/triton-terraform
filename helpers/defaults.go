package helpers

var (
	// TestAccount is a fake SDC account or the value of "SDC_ACCOUNT"
	TestAccount = EnvOrElse("SDC_ACCOUNT", "test")

	// TestKeyFile is the relative path to the test key file from the root of the
	// repo or the value of "SDC_KEY")
	TestKeyFile = EnvOrElse("SDC_KEY", "fixtures/id_rsa.test")

	// TestKeyID is the key ID of the test key file, or the value of "SDC_KEY_ID"
	TestKeyID = EnvOrElse("SDC_KEY_ID", "SHA256:x9NvQdjk+aBVcckyJnAibVbnyt/UyghtJUJxr41KgTI")

	// TestPublicKeyFile is the relative path to the test public key file from the
	// root of the repo or the value of "SDC_KEY_PUBLIC"
	TestPublicKeyFile = EnvOrElse("SDC_KEY_PUBLIC", "fixtures/id_rsa.test.pub")
)
