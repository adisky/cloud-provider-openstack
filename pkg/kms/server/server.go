package server

import (
	"fmt"
	"net"
	"os"

	"github.com/golang/glog"
	"golang.org/x/net/context"
	"golang.org/x/sys/unix"
	"google.golang.org/grpc"
	gcfg "gopkg.in/gcfg.v1"
	pb "k8s.io/apiserver/pkg/storage/value/encrypt/envelope/v1beta1"
	"k8s.io/cloud-provider-openstack/pkg/kms/barbican"
	"k8s.io/cloud-provider-openstack/pkg/kms/encryption/aescbc"
)

const (
	netProtocol    = "unix"
	version        = "v1beta1"
	runtimename    = "barbican"
	runtimeversion = "0.0.1"
)

type server struct{}

var cfg barbican.Config

func initConfig(configFilePath string) error {

	config, err := os.Open(configFilePath)
	defer config.Close()
	if err != nil {
		return err
	}
	err = gcfg.FatalOnly(gcfg.ReadInto(&cfg, config))
	if err != nil {
		return err
	}
	return nil
}

// Run Grpc server for barbican KMS
func Run(configFilePath string, socketpath string, sigchan <-chan os.Signal) error {

	glog.Infof("Barbican KMS Plugin Starting Version: %s, RunTimeVersion: %s", version, runtimeversion)

	if err := initConfig(configFilePath); err != nil {
		glog.V(4).Infof("Error in Getting Config File: %v", err)
		return err
	}

	// unlink the unix socket
	if err := unix.Unlink(socketpath); err != nil {
		glog.V(4).Infof("Error to unlink unix socket: %v", err)
	}

	listener, err := net.Listen(netProtocol, socketpath)
	if err != nil {
		glog.Fatalf("Failed to Listen: %v", err)
		return err
	}

	s := grpc.NewServer()
	pb.RegisterKeyManagementServiceServer(s, &server{})

	go s.Serve(listener)

	for {
		sig := <-sigchan
		if sig == unix.SIGINT || sig == unix.SIGTERM {
			fmt.Println("force stop, shutting down grpc server")
			s.GracefulStop()
			return nil
		}
	}
}

func (s *server) Version(ctx context.Context, req *pb.VersionRequest) (*pb.VersionResponse, error) {

	glog.V(4).Infof("Version Information Requested by Kubernetes api server")

	res := &pb.VersionResponse{
		Version:        version,
		RuntimeName:    runtimename,
		RuntimeVersion: runtimeversion,
	}

	return res, nil
}

func (s *server) Decrypt(ctx context.Context, req *pb.DecryptRequest) (*pb.DecryptResponse, error) {

	glog.V(4).Infof("Decrypt Request by Kubernetes api server")

	barbicanClient, err := barbican.NewBarbicanClient(&cfg)
	if err != nil {
		glog.V(4).Infof("Failed to get Barbican client %v: ", err)
		return nil, err
	}

	keyID := cfg.KeyManager.KeyID

	key, err := barbicanClient.GetSecret(keyID)
	if err != nil {
		glog.V(4).Infof("Failed to get key %v: ", err)
		return nil, err
	}

	plain, err := aescbc.Decrypt(req.Cipher, key)
	if err != nil {
		glog.V(4).Infof("Failed to decrypt data %v: ", err)
		return nil, err
	}

	return &pb.DecryptResponse{Plain: plain}, nil
}

func (s *server) Encrypt(ctx context.Context, req *pb.EncryptRequest) (*pb.EncryptResponse, error) {

	glog.V(4).Infof("Encrypt Request by Kubernetes api server")

	barbicanClient, err := barbican.NewBarbicanClient(&cfg)
	if err != nil {
		glog.V(4).Infof("Failed to get Barbican client %v: ", err)
		return nil, err
	}

	keyID := cfg.KeyManager.KeyID

	key, err := barbicanClient.GetSecret(keyID)

	if err != nil {
		glog.V(4).Infof("Failed to get key %v: ", err)
		return nil, err
	}

	cipher, err := aescbc.Encrypt(req.Plain, key)

	if err != nil {
		glog.V(4).Infof("Failed to encrypt data %v: ", err)
		return nil, err
	}
	return &pb.EncryptResponse{Cipher: cipher}, nil
}
