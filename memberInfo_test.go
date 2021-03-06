package cluster

// xlCluster_go/member_info_test.go

import (
	"fmt"
	xr "github.com/jddixon/rnglib_go"
	. "gopkg.in/check.v1"
)

func (s *XLSuite) TestMISerialization(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("\nTEST_MI_SERIALIZATION")
	}
	rng := xr.MakeSimpleRNG()

	// Generate a random cluster member
	cm := s.makeAMemberInfo(c, rng)

	// Serialize it
	serialized := cm.String()

	// DEBUG
	fmt.Printf("SERIALIZED MemberInfo:\n%s\n", serialized)
	// END

	// Reverse the serialization
	deserialized, rest, err := ParseMemberInfo(serialized)
	c.Assert(err, IsNil)
	c.Assert(len(rest), Equals, 0)

	// Verify that the deserialized member is identical to the original
	c.Assert(deserialized.Equal(cm), Equals, true)
}

//func (s *XLSuite) TestMemberInfoAndTokens(c *C) {
//	if VERBOSITY > 0 {
//		fmt.Println("\nTEST_MEMBER_INFO_AND_TOKENS")
//	}
//	rng := xr.MakeSimpleRNG()
//
//	// Generate a random cluster member
//	cm := s.makeAMemberInfo(c, rng)
//
//	token, err := cm.Token()
//	c.Assert(err, IsNil)
//
//	cm2, err := NewMemberInfoFromToken(token)
//	c.Assert(err, IsNil)
//	c.Assert(cm.Equal(cm2), Equals, true)
//}
