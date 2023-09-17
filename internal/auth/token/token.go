package token

import (
	iamv1authnmachinegrpc "buf.build/gen/go/plantoncloud/planton-cloud-apis/grpc/go/cloud/planton/apis/v1/iam/authn/machine/rpc/rpcgrpc"
	iamv1authnmachinepb "buf.build/gen/go/plantoncloud/planton-cloud-apis/protocolbuffers/go/cloud/planton/apis/v1/iam/authn/machine/rpc"
	"context"
	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
	"github.com/plantoncloud-inc/planton-cloud-kube-agent/internal/apiclient"
	"github.com/plantoncloud-inc/planton-cloud-kube-agent/internal/config"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"time"
)

var (
	Token = ""
)

// Initialize sets the initial value of the token by sending a rpc call to planton cloud service
// the token will be updated by
func Initialize(ctx context.Context, c *config.Config) error {
	conn, err := apiclient.NewConn(c.PlantonCloudServiceApiEndpoint, "")
	if err != nil {
		log.Fatalf("failed to create api client conn with error %v", err)
	}
	accessToken, err := getToken(ctx, c, conn)
	if err != nil {
		return errors.Wrap(err, "failed to get token")
	}
	Token = accessToken
	return nil
}

// StartRotator the periodic scheduler to check and rotate access token
func StartRotator(ctx context.Context, c *config.Config) {

	ticker := time.NewTicker(time.Duration(c.TokenExpirationCheckIntervalSeconds) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := checkAndRotateToken(ctx, c); err != nil {
				log.Fatalf("failed to check and rotate token")
			}
		}
	}
}

// getToken by sending a rpc call to planton cloud service
func getToken(ctx context.Context, c *config.Config, conn *grpc.ClientConn) (string, error) {
	machineAuthnQueryClient := iamv1authnmachinegrpc.NewMachineAuthenticationQueryControllerClient(conn)
	token, err := machineAuthnQueryClient.GetAccessToken(ctx, &iamv1authnmachinepb.GetMachineAccessTokenQueryInput{
		MachineAccountEmail: c.PlantonCloudKubeAgentMachineAccountEmail,
		ClientSecret:        c.PlantonCloudKubeAgentClientSecret,
	})
	if err != nil {
		return "", errors.Wrap(err, "failed to get token")
	}
	return token.Value, nil
}

func checkAndRotateToken(ctx context.Context, c *config.Config) error {
	log.Info("checking access-token expiration")
	isAboutToExpire, err := isTokenAboutToExpire(c.TokenExpirationBufferMinutes, Token)
	if err != nil {
		return errors.Wrapf(err, "failed to check if token is about to expire")
	}
	if !isAboutToExpire {
		return nil
	}
	log.Infof("jwt token is about to expire in the next %d minutes... rotating...", c.TokenExpirationBufferMinutes)
	if err := Initialize(ctx, c); err != nil {
		return errors.Wrap(err, "failed to rotate token")
	}
	log.Infof("access token has been rotated successfully")
	return nil
}

func isTokenAboutToExpire(tokenExpirationBufferMinutes int, jwtToken string) (isAboutToExpire bool, err error) {
	// Parse the JWT token without validating the signature
	token, _, err := new(jwt.Parser).ParseUnverified(jwtToken, jwt.MapClaims{})
	if err != nil {
		return true, errors.Wrap(err, "failed to parse token")
	}

	// Check if the token contains claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return true, errors.Wrap(err, "failed retrieving claims from jwt token")
	}

	// Get the expiration time from the claims
	exp, ok := claims["exp"].(float64)
	if !ok {
		return true, errors.Wrap(err, "failed retrieving expiration from jwt claims")
	}

	// Convert the expiration time to a time.Time value
	expirationTime := time.Unix(int64(exp), 0)

	// Check if the token is about to expire in the next 5 minutes
	isAboutToExpire = time.Now().Add(time.Duration(tokenExpirationBufferMinutes) * time.Minute).After(expirationTime)
	return
}
