package cognito

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
)

var (
	clientID = "5l74ttc4m9etagg1jh8n5b8vic"
	region   = "us-west-1"
)

type CognitoTokenManager struct {
	username string
	password string
	idToken  string
}

func NewCognitoTokenManager(username, password string) *CognitoTokenManager {
	return &CognitoTokenManager{
		username: username,
		password: password,
	}
}

// Authenticate and get a new token
func (ctm *CognitoTokenManager) Authenticate(ctx context.Context) error {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return fmt.Errorf("unable to load AWS SDK config: %v", err)
	}

	svc := cognitoidentityprovider.NewFromConfig(cfg)

	authParams := &cognitoidentityprovider.InitiateAuthInput{
		AuthFlow: "USER_PASSWORD_AUTH",
		AuthParameters: map[string]string{
			"USERNAME": ctm.username,
			"PASSWORD": ctm.password,
		},
		ClientId: aws.String(clientID),
	}

	resp, err := svc.InitiateAuth(ctx, authParams)
	if err != nil {
		return fmt.Errorf("failed to authenticate user: %v", err)
	}

	ctm.idToken = *resp.AuthenticationResult.IdToken
	return nil
}

// GetIdToken returns the current ID token
func (ctm *CognitoTokenManager) GetIdToken() string {
	return ctm.idToken
}
