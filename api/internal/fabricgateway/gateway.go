package fabricgateway

import (
	"crypto/x509"
	"log"
	"os"

	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-gateway/pkg/identity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type Gateway struct {
	clientConn *grpc.ClientConn
	gateway    *client.Gateway
}

func NewGateway() *Gateway {
	return &Gateway{}
}

func (g *Gateway) GetContract() (*client.Contract, error) {
	g.clientConn = newGrpcConnection()

	id := newIdentity()
	sign := newSign()

	gw, err := client.Connect(
		id,
		client.WithSign(sign),
		client.WithClientConnection(g.clientConn),
	)
	if err != nil {
		log.Fatalf("failed to connect gateway: %v", err)
	}
	g.gateway = gw

	channelName := os.Getenv("CHANNEL_NAME")
	if channelName == "" {
		channelName = "channel"
	}
	network := gw.GetNetwork(channelName)

	chaincodeName := os.Getenv("CHAINCODE_NAME")
	if chaincodeName == "" {
		chaincodeName = "chaincode"
	}
	contract := network.GetContract(chaincodeName)

	return contract, nil
}

func (g *Gateway) Close() {
	if g.gateway != nil {
		g.gateway.Close()
	}
	if g.clientConn != nil {
		g.clientConn.Close()
	}
}

// newGrpcConnection creates a gRPC connection to the Fabric Gateway
func newGrpcConnection() *grpc.ClientConn {
	// Get peer endpoint from environment or use default
	peerEndpoint := os.Getenv("PEER_ENDPOINT")
	if peerEndpoint == "" {
		peerEndpoint = "localhost:7051"
	}

	// Path to peer TLS certificate
	tlsCertPath := os.Getenv("PEER_TLS_CERT")
	if tlsCertPath == "" {
		tlsCertPath = "./crypto-config/peers/peer.example.com/tls/ica.tls.organisation.cert"
	}

	certificate, err := os.ReadFile(tlsCertPath)
	if err != nil {
		log.Fatalf("failed to read TLS certificate: %v", err)
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(certificate) {
		log.Fatalf("failed to add certificate to pool")
	}

	// Override server name for TLS verification to match certificate
	// This allows connecting to localhost while using peer's certificate
	tlsServerName := os.Getenv("TLS_SERVER_NAME")

	transportCredentials := credentials.NewClientTLSFromCert(certPool, tlsServerName)

	connection, err := grpc.NewClient(peerEndpoint, grpc.WithTransportCredentials(transportCredentials))
	if err != nil {
		log.Fatalf("failed to create gRPC connection: %v", err)
	}

	return connection
}

// newIdentity creates a client identity from certificate and MSP ID
func newIdentity() *identity.X509Identity {
	// Get MSP ID from environment or use default
	mspID := os.Getenv("MSP_ID")
	if mspID == "" {
		mspID = "orgMSP"
	}

	// Path to user certificate
	certPath := os.Getenv("CERT_PATH")
	if certPath == "" {
		certPath = "./crypto-config/users/user/msp/signcerts/cert.pem"
	}

	certificatePEM, err := os.ReadFile(certPath)
	if err != nil {
		log.Fatalf("failed to read certificate: %v", err)
	}

	certificate, err := identity.CertificateFromPEM(certificatePEM)
	if err != nil {
		log.Fatalf("failed to parse certificate: %v", err)
	}

	id, err := identity.NewX509Identity(mspID, certificate)
	if err != nil {
		log.Fatalf("failed to create identity: %v", err)
	}

	return id
}

// newSign creates a signing function using the private key
func newSign() identity.Sign {
	// Path to user private key
	keyPath := os.Getenv("KEY_PATH")
	if keyPath == "" {
		keyPath = "./crypto-config/users/user/msp/keystore/key_sk"
	}

	privateKeyPEM, err := os.ReadFile(keyPath)
	if err != nil {
		log.Fatalf("failed to read private key: %v", err)
	}

	privateKey, err := identity.PrivateKeyFromPEM(privateKeyPEM)
	if err != nil {
		log.Fatalf("failed to parse private key: %v", err)
	}

	sign, err := identity.NewPrivateKeySign(privateKey)
	if err != nil {
		log.Fatalf("failed to create sign function: %v", err)
	}

	return sign
}
