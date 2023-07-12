package apiclient

import (
	"context"
	"crypto/tls"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/acme/autocert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"os"
	"strings"
)

type EnvVar string

const (
	LestEncryptEmail            = "admin@planton.cloud" //this is random
	CertCacheLoc                = "/tmp/planton-cloud-kube-agent/cert/cache"
	PlantonBuildEngineDnsDomain = "planton-build-engine"
)

type tokenAuth struct {
	token string
}

// GetRequestMetadata adds oauth token to request headers.
// https://jbrandhorst.com/post/grpc-auth/
func (t tokenAuth) GetRequestMetadata(ctx context.Context, in ...string) (map[string]string, error) {
	return map[string]string{
		"authorization": "Bearer " + t.token,
	}, nil
}

func (tokenAuth) RequireTransportSecurity() bool {
	return true
}

func NewConn(endpoint, token string) (*grpc.ClientConn, error) {
	c, err := creConn(endpoint, token)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get grpc conn to core using endpoint %s", endpoint)
	}
	return c, nil
}

func getHostname(endpoint string) string {
	return strings.Split(endpoint, ":")[0]
}

func creConn(addr, token string) (*grpc.ClientConn, error) {
	log.Debugf("dialling %s to create a grpc conn", addr)
	grpcMaxReceiveMsgSize := 104857600 //100MB

	tlsConfig, err := getTlsConfig(getHostname(addr))
	if err != nil {
		return nil, errors.Wrap(err, "failed to get tls config")
	}
	var dialCredentialOption grpc.DialOption
	if strings.HasSuffix(addr, ":443") {
		dialCredentialOption = grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig))
	} else {
		dialCredentialOption = grpc.WithTransportCredentials(insecure.NewCredentials())
	}
	log.Debugf("using %#v dial option for grpc transport credentials", dialCredentialOption)
	if token == "" {
		conn, err := grpc.Dial(addr,
			dialCredentialOption,
			grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(grpcMaxReceiveMsgSize)),
		)
		if err != nil {
			return nil, errors.Wrap(err, "failed to open grpc connection without token")
		}
		return conn, nil
	} else {
		conn, err := grpc.Dial(addr,
			dialCredentialOption,
			grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(grpcMaxReceiveMsgSize)),
			grpc.WithPerRPCCredentials(tokenAuth{token: token}),
		)
		if err != nil {
			return nil, errors.Wrap(err, "failed to open grpc connection with token")
		}
		return conn, nil
	}
}

func getTlsConfig(host string) (*tls.Config, error) {

	if strings.HasSuffix(host, PlantonBuildEngineDnsDomain) {
		return &tls.Config{
			InsecureSkipVerify: true}, nil
	}

	if err := os.MkdirAll(CertCacheLoc, os.ModePerm); err != nil {
		return nil, errors.Wrapf(err, "failed to ensure cert config dir %s", CertCacheLoc)
	}
	manager := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		Cache:      autocert.DirCache(CertCacheLoc),
		HostPolicy: autocert.HostWhitelist(host),
		Email:      LestEncryptEmail,
	}
	return &tls.Config{GetCertificate: manager.GetCertificate}, nil
}
