package pktComms

// xlCluster_go/pktComms/lib_test.go

import (
	"bytes"
	"fmt"
	xr "github.com/jddixon/rnglib_go"
	xcl "github.com/jddixon/xlCluster_go"
	xi "github.com/jddixon/xlNodeID_go"
	xn "github.com/jddixon/xlNode_go"
	//xt "github.com/jddixon/xlTransport_go"
	. "gopkg.in/check.v1"
	"os"
	"path"
	"strings"
)

const (
	VERBOSITY = 1
)

func (s *XLSuite) TestMisc(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("\nTEST_MISC")
	}

	rng := xr.MakeSimpleRNG()
	maxSize := uint32(5)
	kluster := s.makeSimpleCluster(c, rng, maxSize)

	_ = kluster
}

//////////////////////////////////////////////////////////////////////////
// UTILITY FUNCTIONS FOR USE IN TESTING.  NONE SHOULD HAVE AN ERROR RETURN
//////////////////////////////////////////////////////////////////////////

func (s *XLSuite) makeNode(c *C, rng *xr.PRNG, epCount uint32) (
	node *xn.Node) {

	lfs := s.makeUniqueLFS(c, rng)
	c.Assert(strings.HasPrefix(lfs, "tmp/"), Equals, true)
	nodeName := lfs[4:]
	c.Assert(len(nodeName) > 1, Equals, true)

	id := rng.SomeBytes(20) // size of SHA1
	nodeID, err := xi.New(id)
	c.Assert(err, IsNil)

	node, err = xn.NewNew(nodeName, nodeID, lfs)
	c.Assert(err, IsNil)
	return
}
func (s *XLSuite) makeSimpleCluster(c *C, rng *xr.PRNG, maxSize uint32) (
	kluster *xcl.Cluster) {

	c.Assert(maxSize > 1, Equals, true)

	// cluster metadata ---------------------------------------------
	dummyLFS := s.makeUniqueLFS(c, rng)
	c.Assert(strings.HasPrefix(dummyLFS, "tmp/"), Equals, true)
	clusterName := dummyLFS[4:]

	id := rng.SomeBytes(20) // size of SHA1
	nodeID, err := xi.New(id)
	c.Assert(err, IsNil)
	attrs := uint64(0)
	epCount := uint32(2)
	kluster, err = xcl.NewCluster(clusterName, nodeID, attrs, maxSize, epCount)
	c.Assert(err, IsNil)
	c.Assert(kluster, NotNil)

	c.Assert(kluster.GetName(), Equals, clusterName)
	c.Assert(bytes.Equal(kluster.GetNodeID().Value(), nodeID.Value()),
		Equals, true)
	c.Assert(kluster.GetCurSize(), Equals, uint32(0))
	c.Assert(kluster.GetMaxSize(), Equals, maxSize)
	c.Assert(kluster.GetAttrs(), Equals, attrs)
	c.Assert(kluster.GetEPCount(), Equals, epCount)

	// create enough Nodes to fully populate the cluster ------------
	var nodes []*xn.Node
	for i := uint32(0); i < maxSize; i++ {
		node := s.makeNode(c, rng, uint32(2))
		err = node.OpenAcc()
		c.Assert(err, IsNil)
		nodes = append(nodes, node)

	}
	defer func() {
		for i := uint32(0); i < maxSize; i++ {
			if nodes[i] != nil {
				nodes[i].CloseAcc()
			}
		}
	}()
	// create K=maxSize MemberInfo ----------------------------------
	// XXX STUB XXX

	// create K ClusterMember ---------------------------------------
	// XXX STUB XXX

	// copy full MemberInfo list to each ClusterMember, setting self
	// XXX STUB XXX

	// copy full MemberInfo list to each ClusterMember, setting self
	// XXX STUB XXX

	return
}

// Create a uniquely named subdirectory of tmp/
func (s *XLSuite) makeUniqueLFS(c *C, rng *xr.PRNG) (lfs string) {

	var err error

	// first make sure that tmp/ exists and is a directory
	tmpDir := "tmp"
	if _, err := os.Stat(tmpDir); os.IsNotExist(err) {
		err = os.Mkdir(tmpDir, 0775)
	}
	if err == nil {
		for {
			lfs = path.Join(tmpDir, rng.NextFileName(8))
			if _, err = os.Stat(lfs); os.IsNotExist(err) {
				// lfs does not exist, so we are done
				err = os.Mkdir(lfs, 0775)
				break
			}
		}
	}
	c.Assert(err, IsNil)
	return
}

//////////////////////////////////////////////////////////////////////////
//// XXX ALL OF WHAT FOLLOWS SHOULD BE DROPPED XXX
//////////////////////////////////////////////////////////////////////////

//// Create and start an ephemeral xlReg server -----------------------
//// Calling routine should call defer eph.Close()
//func (s *XLSuite) launchEphServer(c *C) (eph *xg.EphServer,
//	reg *xg.Registry, regID *xi.NodeID, server *xg.RegServer) {
//
//	eph, err := xg.NewEphServer()
//	c.Assert(eph, NotNil)
//	c.Assert(err, IsNil)
//
//	server = eph.Server
//
//	// start the ephemeral server -------------------------
//	err = eph.Start()
//	c.Assert(err, IsNil)
//
//	// verify Bloom filter is running
//	reg = &eph.Server.Registry
//	c.Assert(reg, NotNil)
//	regID = reg.GetNodeID()
//	c.Assert(reg.IDCount(), Equals, uint(1)) // the registry's own ID
//	found, err := reg.ContainsID(regID)
//	c.Assert(err, IsNil)
//	c.Assert(found, Equals, true)
//
//	return
//}
//
//// Create and register a solo cluster -------------------------------
//func (s *XLSuite) createAndRegSoloCluster(c *C, rng *xr.PRNG,
//	reg *xg.Registry, regID *xi.NodeID, server *xg.RegServer) (
//	clusterName string, clusterAttrs uint64, clusterID *xi.NodeID, K uint32) {
//
//	serverName := server.GetName()
//	serverID := server.GetNodeID()
//	serverEnd := server.GetEndPoint(0)
//	serverCK := server.GetCommsPublicKey()
//	serverSK := server.GetSigPublicKey()
//	c.Assert(serverEnd, NotNil)
//
//	clusterName = rng.NextFileName(8)
//	clusterAttrs = uint64(rng.Int63())
//	K = uint32(2 + rng.Intn(6)) // so the size is 2 .. 7
//
//	// create an AdminClient, use it to get the clusterID --------
//	an, err := xg.NewAdminClient(serverName, serverID, serverEnd,
//		serverCK, serverSK, clusterName, clusterAttrs, K, uint32(3), nil)
//	c.Assert(err, IsNil)
//
//	an.Start()
//	<-an.DoneCh
//
//	c.Assert(an.ClusterID, NotNil)          // the purpose of the exercise
//	c.Assert(an.EPCount, Equals, uint32(3)) // NEED >= 2
//	c.Assert(an.ClusterMaxSize, Equals, K)
//
//	anID := an.GetNodeID()
//	c.Assert(reg.IDCount(), Equals, uint(3)) // regID + anID + clusterID
//	clusterID = an.ClusterID
//
//	// DEBUG
//	fmt.Printf("regID     %s\n", regID.String())
//	fmt.Printf("anID      %s\n", anID.String())
//	fmt.Printf("clusterID %s\n", an.ClusterID.String())
//	fmt.Printf("  size    %d\n", an.ClusterMaxSize)
//	fmt.Printf("  name    %s\n", an.ClusterName)
//	// END
//
//	found, err := reg.ContainsID(regID)
//	c.Assert(err, IsNil)
//	c.Assert(found, Equals, true)
//
//	found, err = reg.ContainsID(anID)
//	c.Assert(err, IsNil)
//	c.Assert(found, Equals, true)
//
//	found, err = reg.ContainsID(an.ClusterID)
//	c.Assert(err, IsNil)
//	c.Assert(found, Equals, true)
//
//	return
//}
//
//// Create PktComm layers for K members
////
/////////////////////////////////////////////////////////////////////////
//// XXX THIS IS NOW CREATING Bootstrappers !!
/////////////////////////////////////////////////////////////////////////
//func (s *XLSuite) createKMemberPktLayers(c *C, rng *xr.PRNG,
//	server *xg.RegServer,
//	clusterName string, clusterAttrs uint64, clusterID *xi.NodeID,
//	K uint32) (bs []*xg.Bootstrapper, bsNames []string) {
//
//	serverName := server.GetName()
//	serverID := server.GetNodeID()
//	serverEnd := server.GetEndPoint(0)
//	serverCK := server.GetCommsPublicKey()
//	serverSK := server.GetSigPublicKey()
//	c.Assert(serverEnd, NotNil)
//
//	var err error
//	bs = make([]*xg.Bootstrapper, K)
//	bsNames = make([]string, K)
//	namesInUse := make(map[string]bool)
//	for i := uint32(0); i < K; i++ {
//		var ep *xt.TcpEndPoint
//		ep, err = xt.NewTcpEndPoint("127.0.0.1:0")
//		c.Assert(err, IsNil)
//		e := []xt.EndPointI{ep}
//		newName := rng.NextFileName(8)
//		_, ok := namesInUse[newName]
//		for ok {
//			newName = rng.NextFileName(8)
//			_, ok = namesInUse[newName]
//		}
//		namesInUse[newName] = true
//		bsNames[i] = newName // guaranteed to be LOCALLY unique
//		attrs := uint64(rng.Int63())
//
//		// XXX the bs[i] are actually NOT PktLayer
//		bs[i], err = xg.NewBootstrapper(bsNames[i], "",
//			nil, nil, // private RSA keys are generated if nil
//			attrs,
//			serverName, serverID, serverEnd, serverCK, serverSK,
//			clusterName, clusterAttrs, clusterID,
//			K, uint32(2), e) // 2 is endPoint count
//		c.Assert(err, IsNil)
//		c.Assert(bs[i], NotNil)
//		c.Assert(bs[i].ClusterID, NotNil)
//	}
//	return
//}
//
//// THIS IS INCOMPLETE
