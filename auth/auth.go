package auth

import (
	"context"
	"encoding/base64"
	"strconv"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
)

// GetClientWithDefaultAzureCredential creates a new Cosmos DB client using DefaultAzureCredential authentication strategy.
// Recommended way to authenticate in production environments.
func GetClientWithDefaultAzureCredential(endpoint string, opts *azcosmos.ClientOptions) (*azcosmos.Client, error) {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, err
	}
	return azcosmos.NewClient(endpoint, cred, opts)
}

// emulatorTokenCredential implements azcore.TokenCredential for Cosmos DB Emulator
// It always returns the static emulator token

type emulatorTokenCredential struct {
	token azcore.AccessToken
}

func (e *emulatorTokenCredential) GetToken(ctx context.Context, opts policy.TokenRequestOptions) (azcore.AccessToken, error) {
	return e.token, nil
}

// GetEmulatorClientWithAzureADAuth creates a Cosmos DB client for the local emulator using a static Azure AD token.
// This enables local development and testing with Cosmos DB Emulator using Azure AD-like authentication.
func GetEmulatorClientWithAzureADAuth(endpoint string, opts *azcosmos.ClientOptions) (*azcosmos.Client, error) {
	token, err := getADTokenForEmulator()
	if err != nil {
		return nil, err
	}
	cred := &emulatorTokenCredential{token: token}
	return azcosmos.NewClient(endpoint, cred, opts)
}

// taken from https://github.com/Azure/azure-sdk-for-go/blob/main/sdk/data/azcosmos/emulator_tests.go
func getADTokenForEmulator() (azcore.AccessToken, error) {
	header := `{"typ":"JWT","alg":"RS256","x5t":"CosmosEmulatorPrimaryMaster","kid":"CosmosEmulatorPrimaryMaster"}`
	unixNow := time.Now().Unix()
	expiration := unixNow + 7200
	payload := `{ 
		"appid":"localhost", 
		"aio":"",
		"appidacr":"1",
		"idp": "https://localhost:8081/",
		"oid": "96313034-4739-43cb-93cd-74193adbe5b6",
		"rh": "",
		"sub": "localhost",
		"tid": "EmulatorFederation",
		"uti": "",
		"ver": "1.0",
		"scp": "user_impersonation",
		"groups":[ 
			"7ce1d003-4cb3-4879-b7c5-74062a35c66e",
			"e99ff30c-c229-4c67-ab29-30a6aebc3e58",
			"5549bb62-c77b-4305-bda9-9ec66b85d9e4",
			"c44fd685-5c58-452c-aaf7-13ce75184f65",
			"be895215-eab5-43b7-9536-9ef8fe130330"], 
		"nbf":` + strconv.FormatInt(unixNow, 10) + `, 
		"exp":` + strconv.FormatInt(expiration, 10) + `, 
		"iat":` + strconv.FormatInt(unixNow, 10) + `,
		"iss":"https://sts.fake-issuer.net/7b1999a1-dfd7-440e-8204-00170979b984",
		"aud":"https://localhost.localhost" 
	}`

	headerBase64 := base64.RawURLEncoding.EncodeToString([]byte(header))
	payloadBase64 := base64.RawURLEncoding.EncodeToString([]byte(payload))
	masterKeyBase64 := base64.RawURLEncoding.EncodeToString([]byte("C2y6yDjf5/R+ob0N8A7Cgv30VRDJIWEHLM+4QDU5DE2nQ9nDuVTqobD4b8mGGyPMbIZnqyMsEcaGQy67XIw/Jw=="))

	token := headerBase64 + "." + payloadBase64 + "." + masterKeyBase64

	return azcore.AccessToken{
		Token:     token,
		ExpiresOn: time.Unix(expiration, 0),
	}, nil
}
