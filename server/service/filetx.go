package service

import (
	"bytes"
	"fmt"
	"rabbitnft/server/config"
	"strings"
	"time"

	ipfsapi "github.com/ipfs/go-ipfs-api"
)

type IPFSNode struct {
	IP      string
	Port    int
	Gateway int
}

type FileTxService struct {
}

func NewFileTxService() *FileTxService {
	return &FileTxService{}
}

// upload file and return hash
func (f *FileTxService) UploadFile(buf *bytes.Buffer) (string, error) {
	node, err := f.getFirstNode()
	if err != nil {
		return "", err
	}
	ip := node.IP
	port := node.Port
	shell := ipfsapi.NewShell(fmt.Sprintf("%s:%d", ip, port))
	return shell.Add(buf, ipfsapi.OnlyHash(false))
}

func (f *FileTxService) VerifyHash(hash string) error {
	node, err := f.getFirstNode()
	if err != nil {
		return err
	}
	ip := node.IP
	port := node.Port
	shell := ipfsapi.NewShell(fmt.Sprintf("%s:%d", ip, port))
	shell.SetTimeout(3 * time.Second)
	_, err = shell.Cat(hash)
	if err != nil && strings.Contains(err.Error(), "context deadline exceeded") {
		return fmt.Errorf("hash query timeout.")
	}
	return err
}

func (f *FileTxService) getFirstNode() (*IPFSNode, error) {
	node := &IPFSNode{
		IP:      config.C.GetString("ipfs.url"),
		Port:    config.C.GetInt("ipfs.port"),
		Gateway: config.C.GetInt("ipfs.gateway"),
	}
	return node, nil
}

func CollectibleUrl(hash string) string {
	return fmt.Sprintf(
		"%s:%d/ipfs/%s",
		config.C.GetString("ipfs.url"),
		config.C.GetInt("ipfs.gateway"),
		hash,
	)
}
