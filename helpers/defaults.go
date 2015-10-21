package helpers

import (
	"io/ioutil"
)

var (
	// TestAccount is a fake SDC account or the value of "SDC_ACCOUNT"
	TestAccount = EnvOrElse("SDC_ACCOUNT", "test")

	// TestKeyFile is the relative path to the test key file from the root of the
	// repo or the value of "SDC_KEY")
	TestKeyFile = EnvOrElse("SDC_KEY", "fixtures/id_rsa.test")

	// TestKeyDataBytes is the data for the private test key, in bytes
	TestKeyDataBytes, _ = ioutil.ReadFile(TestKeyFile)

	// TestKeyData is the data for the private test key, as a string
	TestKeyData = string(TestKeyDataBytes)

	// TestKeyID is the key ID of the test key file, or the value of "SDC_KEY_ID"
	TestKeyID = EnvOrElse("SDC_KEY_ID", "7f:65:5c:73:b2:9e:b5:7b:68:21:4a:ea:18:26:ac:1f")

	// TestPublicKeyFile is the relative path to the test public key file from the
	// root of the repo or the value of "SDC_KEY_PUBLIC"
	TestPublicKeyFile = EnvOrElse("SDC_KEY_PUBLIC", "fixtures/id_rsa.test.pub")

	// TestPublicKeyDataBytes is the data for the test public key, in bytes
	TestPublicKeyDataBytes, _ = ioutil.ReadFile(TestPublicKeyFile)

	// TestPublicKeyData is the data for the test public key, as a string
	TestPublicKeyData = string(TestPublicKeyDataBytes)
)
