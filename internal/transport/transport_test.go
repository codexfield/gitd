package transport_test

import (
	"fmt"
	"gitd/internal/storage"
	"gitd/internal/transport"
	"github.com/go-git/go-billy/v5/memfs"
	fixtures "github.com/go-git/go-git-fixtures/v4"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/object"
	transport2 "github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/client"
	"github.com/go-git/go-git/v5/plumbing/transport/server"
	"github.com/go-git/go-git/v5/plumbing/transport/test"
	"github.com/go-git/go-git/v5/storage/memory"
	. "gopkg.in/check.v1"
	"os"
	"testing"
	"time"
)

func Test(t *testing.T) { TestingT(t) }

const (
	BucketName = "gitd"
	Endpoint   = "gnfd://gnfd-testnet-fullnode-tendermint-us.bnbchain.org:443/" + BucketName
)

type BaseSuite struct {
	fixtures.Suite
	test.ReceivePackSuite

	loader       server.MapLoader
	client       transport2.Transport
	clientBackup transport2.Transport
	asClient     bool
}

func (s *BaseSuite) SetUpSuite(c *C) {
	s.clientBackup = client.Protocols[transport.GnfdProtocol]
	transport.InstallGreenfieldTransport()

	s.client = client.Protocols[transport.GnfdProtocol]
}

func (s *BaseSuite) TearDownSuite(c *C) {
	if s.clientBackup == nil {
		delete(client.Protocols, transport.GnfdProtocol)
	} else {
		client.Protocols[transport.GnfdProtocol] = s.clientBackup
	}
}

type BasicTestSuite struct {
	BaseSuite
}

var _ = Suite(&BasicTestSuite{})

func (s *BasicTestSuite) SetUpSuite(c *C) {
	s.BaseSuite.SetUpSuite(c)
}

func (s *BasicTestSuite) SetUpTest(c *C) {
}

func (s *BasicTestSuite) TestBasic(c *C) {
	var err error

	endpoint, err := transport2.NewEndpoint(Endpoint)
	if err != nil {
		fmt.Printf("New endpoint error: %s", err)
		return
	}
	fmt.Println("Endpoint: ", endpoint.String())
	c.Assert(endpoint, NotNil)
	c.Assert(endpoint.Protocol, Equals, transport.GnfdProtocol)
	c.Assert(endpoint.Host, Equals, "gnfd-testnet-fullnode-tendermint-us.bnbchain.org")
	c.Assert(endpoint.Port, Equals, int(443))

	// load storage
	session, err := s.client.NewUploadPackSession(endpoint, nil)
	if err != nil {
		return
	}
	c.Assert(session, NotNil)
}

func (s *BasicTestSuite) TestBasicClone(c *C) {
	var err error

	endpoint, err := transport2.NewEndpoint(Endpoint)
	if err != nil {
		fmt.Printf("New endpoint error: %s", err)
		return
	}

	fmt.Println("Endpoint: ", endpoint.String())
	c.Assert(endpoint, NotNil)
	c.Assert(endpoint.Protocol, Equals, transport.GnfdProtocol)
	c.Assert(endpoint.Host, Equals, "gnfd-testnet-fullnode-tendermint-us.bnbchain.org")
	c.Assert(endpoint.Port, Equals, int(443))

	r, err := git.PlainOpen("../gitd/")
	c.Assert(err, IsNil)

	remoteName := "greenfield"
	_, err = r.CreateRemote(&config.RemoteConfig{
		Name: remoteName,
		URLs: []string{"gnfd://gnfd-testnet-fullnode-tendermint-us.bnbchain.org:443/" + BucketName},
	})
	c.Assert(err, IsNil)
	defer func() {
		c.Assert(r.DeleteRemote(remoteName), IsNil)
	}()

	err = r.Push(&git.PushOptions{RemoteName: remoteName, Force: true})
	c.Assert(err, IsNil)

	// clone the empty repo
	_, err = git.Clone(memory.NewStorage(), memfs.New(), &git.CloneOptions{
		URL: endpoint.String(),
	})
	if err != nil {
		return
	}
	// load storage
	session, err := s.client.NewUploadPackSession(endpoint, nil)
	if err != nil {
		return
	}
	c.Assert(session, NotNil)
}

func (s *BasicTestSuite) TestBasicPush(c *C) {
	var err error

	endpoint, err := transport2.NewEndpoint(Endpoint)
	if err != nil {
		fmt.Printf("New endpoint error: %s", err)
		return
	}

	fmt.Println("Endpoint: ", endpoint.String())
	c.Assert(endpoint, NotNil)
	c.Assert(endpoint.Protocol, Equals, transport.GnfdProtocol)
	c.Assert(endpoint.Host, Equals, "gnfd-testnet-fullnode-tendermint-us.bnbchain.org")
	c.Assert(endpoint.Port, Equals, int(443))

	rep, err := git.Init(memory.NewStorage(), memfs.New())
	c.Assert(err, IsNil)
	_, err = rep.CreateRemote(&config.RemoteConfig{
		Name: "greenfield",
		URLs: []string{"gnfd://gnfd-testnet-fullnode-tendermint-us.bnbchain.org:443/" + BucketName},
	})
	c.Assert(err, IsNil)

	wt, err := rep.Worktree()
	c.Assert(err, IsNil)
	err = wt.Pull(&git.PullOptions{RemoteName: "greenfield"})
	c.Assert(err, IsNil)
	createCommit(c, rep)
	err = rep.Push(&git.PushOptions{
		RemoteName: "greenfield",
		Force:      true,
	})
	c.Assert(err, IsNil)
}

func createCommit(c *C, r *git.Repository) {
	// Create a commit so there is a HEAD to check
	wt, err := r.Worktree()
	c.Assert(err, IsNil)

	rm, err := wt.Filesystem.Create("foo.txt")
	c.Assert(err, IsNil)

	_, err = rm.Write([]byte("foo text"))
	c.Assert(err, IsNil)

	_, err = wt.Add("foo.txt")
	c.Assert(err, IsNil)

	author := object.Signature{
		Name:  "go-git",
		Email: "go-git@fake.local",
		When:  time.Now(),
	}
	_, err = wt.Commit("test commit message", &git.CommitOptions{
		All:       true,
		Author:    &author,
		Committer: &author,
	})
	c.Assert(err, IsNil)

}

func newEmptyGreenfieldRepo(c *C) {
	newStorage, err := storage.NewStorage(
		os.Getenv(transport.EnvChainID),
		"https://gnfd-testnet-fullnode-tendermint-us.bnbchain.org:443/",
		os.Getenv(transport.EnvPrivateKey),
		"helloworld",
	)
	c.Assert(err, IsNil)
	_, err = git.Init(newStorage, memfs.New())
	c.Assert(err, IsNil)
}

func (s *BasicTestSuite) TestEmptyRepoInGreenfield(c *C) {
	newEmptyGreenfieldRepo(c)
}
