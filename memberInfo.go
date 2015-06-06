package cluster

// xlCluster_go/member_info.go

// This file contains functions and structures used by registry clients
// to manage information about clusters and their members.

import (
	//"crypto/rsa"
	"encoding/hex"
	"fmt"
	xn "github.com/jddixon/xlNode_go"
	"strings"
)

var _ = fmt.Print

type MemberInfo struct {
	Attrs uint64 //  bit flags are defined in const.go
	Peer  *xn.Peer
}

//func NewMemberInfo(name string, id *xi.NodeID,
//	commsPubKey, sigPubKey *rsa.PublicKey, attrs uint64, myEnds []string) (

func NewMemberInfo(attrs uint64, peer *xn.Peer) (
	mi *MemberInfo, err error) {

	// all attrs bits are zero by default

	mi = &MemberInfo{
		Attrs: attrs,
		Peer:  peer,
	}
	return
}

// EQUAL ////////////////////////////////////////////////////////////

func (mi *MemberInfo) Equal(any interface{}) bool {

	if any == mi {
		return true
	}
	if any == nil {
		return false
	}
	switch v := any.(type) {
	case *MemberInfo:
		_ = v
	default:
		return false
	}
	other := any.(*MemberInfo) // type assertion
	if mi.Attrs != other.Attrs {
		return false
	}
	if mi.Peer.SizeConnectors() != other.Peer.SizeConnectors() {
		return false
	} else {
		count := mi.Peer.SizeConnectors()
		for i := 0; i < count; i++ {
			if mi.Peer.GetConnector(i).String() !=
				other.Peer.GetConnector(i).String() {
				return false
			}
		}
	}
	// WARNING: panics without the ampersand !
	return mi.Peer.BaseNode.Equal(&other.Peer.BaseNode)
}

// SERIALIZATION ////////////////////////////////////////////////////

func (mi *MemberInfo) Strings() (ss []string) {
	ss = []string{"memberInfo {"}
	bns := mi.Peer.BaseNode.Strings()
	for i := 0; i < len(bns); i++ {
		ss = append(ss, "    "+bns[i])
	}
	ss = append(ss, "    connectors {")
	for i := 0; i < mi.Peer.SizeConnectors(); i++ {
		ss = append(ss, "        "+mi.Peer.GetConnector(i).String())
	}
	ss = append(ss, "    }")
	ss = append(ss, fmt.Sprintf("    attrs: 0x%016x", mi.Attrs))
	ss = append(ss, "}")
	return
}

func (mi *MemberInfo) String() string {
	return strings.Join(mi.Strings(), "\n")
}
func collectAttrs(mi *MemberInfo, ss []string) (rest []string, err error) {
	rest = ss
	line, err := xn.NextNBLine(&rest) // trims
	if err == nil {
		// attrs line looks like "attrs: 0xHHHH..." where H is a hex digit
		if strings.HasPrefix(line, "attrs: 0x") {
			var val []byte
			var attrs uint64
			line := line[9:]
			val, err = hex.DecodeString(line)
			if err == nil {
				if len(val) != 8 {
					err = WrongNumberOfBytesInAttrs
				} else {
					for i := 0; i < 8; i++ {
						// assume little-endian ; but printf has put
						// high order bytes first - ie, it's big-endian
						attrs |= uint64(val[i]) << uint(8*(7-i))
					}
					mi.Attrs = attrs
				}
			}
		} else {
			err = BadAttrsLine
		}
	}
	return
}

func ParseMemberInfo(s string) (
	mi *MemberInfo, rest []string, err error) {

	ss := strings.Split(s, "\n")
	return ParseMemberInfoFromStrings(ss)
}

func ParseMemberInfoFromStrings(ss []string) (
	mi *MemberInfo, rest []string, err error) {

	var line string
	bn, rest, err := xn.ParseBNFromStrings(ss, "memberInfo")
	if err == nil {
		peerPart := &xn.Peer{
			BaseNode: *bn,
		}
		mi = &MemberInfo{
			Peer: peerPart,
		}
		// expect and consume a closing brace
		rest, err = xn.CollectConnectors(peerPart, rest)
		if err == nil {
			rest, err = collectAttrs(mi, rest)
			if err == nil {
				line, err = xn.NextNBLine(&rest)
				if err == nil {
					if line == "" {
						if line != "}" {
							err = MissingClosingBrace
						}
					}
				}
			}
		}
	}
	return
}
