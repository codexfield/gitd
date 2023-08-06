package transport_test

import (
	"fmt"
	"gitd/internal/transport"
	"github.com/go-git/go-billy/v5/memfs"
	fixtures "github.com/go-git/go-git-fixtures/v4"
	"github.com/go-git/go-git/v5"
	transport2 "github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/client"
	"github.com/go-git/go-git/v5/plumbing/transport/server"
	"github.com/go-git/go-git/v5/plumbing/transport/test"
	"github.com/go-git/go-git/v5/storage/memory"
	. "gopkg.in/check.v1"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

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

	endpoint, err := transport2.NewEndpoint("gnfd://gnfd-testnet-fullnode-tendermint-us.bnbchain.org:443/test-bucket")
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

	endpoint, err := transport2.NewEndpoint("gnfd://gnfd-testnet-fullnode-tendermint-us.bnbchain.org:443/test-bucket")
	if err != nil {
		fmt.Printf("New endpoint error: %s", err)
		return
	}

	fmt.Println("Endpoint: ", endpoint.String())
	c.Assert(endpoint, NotNil)
	c.Assert(endpoint.Protocol, Equals, transport.GnfdProtocol)
	c.Assert(endpoint.Host, Equals, "gnfd-testnet-fullnode-tendermint-us.bnbchain.org")
	c.Assert(endpoint.Port, Equals, int(443))

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
