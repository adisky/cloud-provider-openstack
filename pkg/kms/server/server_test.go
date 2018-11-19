package server

import (
	"testing"

	"k8s.io/cloud-provider-openstack/pkg/kms/barbican"

	"golang.org/x/net/context"
	pb "k8s.io/apiserver/pkg/storage/value/encrypt/envelope/v1beta1"
)

func TestInitConfig(t *testing.T) {

}

/*
func TestVersion(t *testing.T) {
	req := &pb.VersionRequest{Version: "v1beta1"}
	resp, err := s.Version(context.TODO(), req)
	t.Log(resp)
	t.Log(err)

}
*/

func TestEncryptDecrypt(t *testing.T) {
	s := new(KMSserver)
	s.barbican = &barbican.FakeBarbican{}
	fakeData := []byte("fakedata")
	req := &pb.EncryptRequest{Version: "v1beta1", Plain: fakeData}
	resp, err := s.Encrypt(context.TODO(), req)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	t.Log(resp)

}
