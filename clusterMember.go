package cluster

import (
	"encoding/hex"
	"fmt"
	xc "github.com/jddixon/xlCrypto_go"
	xi "github.com/jddixon/xlNodeID_go"
	xn "github.com/jddixon/xlNode_go"
	"strconv"
	"strings"
)

var (
	INDENT = xn.INDENT // in a better go, this would be const
)

type ClusterMember struct {
	Attrs          uint64 // possibly negotiated with/decreed by a reg server
	ClusterName    string
	ClusterID      *xi.NodeID
	ClusterAttrs   uint64
	ClusterMaxSize uint32 // this is a FIXED size, aka MaxSize, including self

	// EPCount is the number of endPoints dedicated to use for cluster-
	// related purposes.  By convention endPoints[0] is used for
	// member-member communications and [1] for comms with cluster clients,
	// should they exist. The first EPCount endPoints are passed
	// to other cluster members via the registry.
	EPCount uint32

	SelfIndex uint32        // which member we are in the Members slice
	Members   []*MemberInfo // information on (other) cluster members

	xn.Node
}

// Add the other members in the Members list to this member's Node as
// Peers.  This function must be called after the struct is created but
// BEFORE it is serialized.
func (cm *ClusterMember) AddPeers() (err error) {
	node := &cm.Node
	if node == nil {
		err = NilNode
	} else {
		for i := uint32(0); i < cm.ClusterMaxSize; i++ {
			if i == cm.SelfIndex {
				continue
			}
			_, err = node.AddPeer(cm.Members[i].Peer)
			if err != nil {
				break
			}
		}
	}
	return
}

// EQUAL ////////////////////////////////////////////////////////////

func (cm *ClusterMember) Equal(any interface{}) bool {

	if any == cm {
		return true
	}
	if any == nil {
		return false
	}
	switch v := any.(type) {
	case *ClusterMember:
		_ = v
	default:
		return false
	}
	other := any.(*ClusterMember) // type assertion

	if cm.Attrs != other.Attrs || cm.ClusterName != other.ClusterName ||
		cm.ClusterAttrs != other.ClusterAttrs ||
		cm.ClusterMaxSize != other.ClusterMaxSize || cm.EPCount != other.EPCount {
		return false
	}
	if !cm.ClusterID.Equal(other.ClusterID) {
		return false
	}
	for i := 0; i < len(cm.Members); i++ {
		if !cm.Members[i].Equal(other.Members[i]) {
			return false
		}
	}
	return true
}

// SERIALIZATION ////////////////////////////////////////////////////

func (cm *ClusterMember) Strings() (ss []string) {
	ss = []string{"clusterMember {"}
	ns := cm.Node.Strings()
	for i := 0; i < len(ns); i++ {
		ss = append(ss, INDENT+ns[i])
	}
	ss = append(ss, fmt.Sprintf("%sattrs: %d", INDENT, cm.Attrs))

	ss = append(ss, fmt.Sprintf("%sclusterName: %s", INDENT, cm.ClusterName))
	ss = append(ss, fmt.Sprintf("%sclusterID: %s", INDENT,
		hex.EncodeToString(cm.ClusterID.Value())))
	ss = append(ss, fmt.Sprintf("%sclusterAttrs: %d", INDENT, cm.ClusterAttrs))
	ss = append(ss, fmt.Sprintf("%sclusterMaxSize: %d", INDENT, cm.ClusterMaxSize))
	ss = append(ss, fmt.Sprintf("%sepCount: %d", INDENT, cm.EPCount))

	ss = append(ss, fmt.Sprintf("%sselfIndex: %d", INDENT, cm.SelfIndex))
	ss = append(ss, fmt.Sprintf("%smembers {", INDENT))
	for i := 0; i < len(cm.Members); i++ {
		// DEBUG
		//fmt.Printf("serializing member %d\n", i)
		// END
		miss := cm.Members[i].Strings()
		for j := 0; j < len(miss); j++ {
			ss = append(ss, fmt.Sprintf("%s%s%s", INDENT, INDENT, miss[j]))
		}
	}
	ss = append(ss, fmt.Sprintf("%s}", INDENT))
	ss = append(ss, "}")
	return
}
func (cm *ClusterMember) String() string {
	return strings.Join(cm.Strings(), "\n")
}

func ParseClusterMember(s string) (
	cm *ClusterMember, rest []string, err error) {

	ss := strings.Split(s, "\n")
	return ParseClusterMemberFromStrings(ss)
}

func ParseClusterMemberFromStrings(ss []string) (
	cm *ClusterMember, rest []string, err error) {

	var (
		node           *xn.Node
		attrs          uint64
		clusterName    string
		clusterAttrs   uint64
		clusterID      *xi.NodeID
		clusterMaxSize uint32
		selfIndex      uint32
		memberInfos    []*MemberInfo
		epCount        uint32
	)
	line, err := xc.NextNBLine(&ss)
	if err == nil {
		if line != "clusterMember {" {
			err = IllFormedClusterMember
		} else {
			node, rest, err = xn.ParseFromStrings(ss)
		}
	}
	if err == nil {
		line, err = xc.NextNBLine(&rest)
		if err == nil {
			parts := strings.Split(line, ": ")
			if len(parts) == 2 && parts[0] == "attrs" {
				var n int
				raw := strings.TrimSpace(parts[1])
				n, err = strconv.Atoi(raw)
				if err == nil {
					attrs = uint64(n)
				}
			} else {
				err = IllFormedClusterMember
			}
		}
	}
	if err == nil {
		line, err = xc.NextNBLine(&rest)
		if err == nil {
			parts := strings.Split(line, ": ")
			if len(parts) == 2 && parts[0] == "clusterName" {
				clusterName = strings.TrimLeft(parts[1], " \t")
			} else {
				err = IllFormedClusterMember
			}
		}
	}
	if err == nil {
		line, err = xc.NextNBLine(&rest)
		if err == nil {
			parts := strings.Split(line, ": ")
			if len(parts) == 2 && parts[0] == "clusterID" {
				var h []byte
				raw := strings.TrimSpace(parts[1])
				h, err = hex.DecodeString(raw)
				if err == nil {
					clusterID, err = xi.New(h)
				}
			} else {
				err = IllFormedClusterMember
			}
		}
	}
	if err == nil {
		line, err = xc.NextNBLine(&rest)
		if err == nil {
			parts := strings.Split(line, ": ")
			if len(parts) == 2 && parts[0] == "clusterAttrs" {
				var n int
				raw := strings.TrimSpace(parts[1])
				n, err = strconv.Atoi(raw)
				if err == nil {
					clusterAttrs = uint64(n)
				}
			} else {
				err = IllFormedClusterMember
			}
		}
	}
	if err == nil {
		line, err = xc.NextNBLine(&rest)
		if err == nil {
			parts := strings.Split(line, ": ")
			if len(parts) == 2 && parts[0] == "clusterMaxSize" {
				var n int
				raw := strings.TrimSpace(parts[1])
				n, err = strconv.Atoi(raw)
				if err == nil {
					clusterMaxSize = uint32(n)
				}
			} else {
				err = IllFormedClusterMember
			}
		}
	}
	if err == nil {
		line, err = xc.NextNBLine(&rest)
		if err == nil {
			parts := strings.Split(line, ": ")
			if len(parts) == 2 && parts[0] == "epCount" {
				var n int
				raw := strings.TrimSpace(parts[1])
				n, err = strconv.Atoi(raw)
				if err == nil {
					epCount = uint32(n)
				}
			} else {
				err = IllFormedClusterMember
			}
		}
	}
	if err == nil {
		line, err = xc.NextNBLine(&rest)
		if err == nil {
			parts := strings.Split(line, ": ")
			if len(parts) == 2 && parts[0] == "selfIndex" {
				var n int
				raw := strings.TrimSpace(parts[1])
				n, err = strconv.Atoi(raw)
				if err == nil {
					selfIndex = uint32(n)
				}
			} else {
				err = IllFormedClusterMember
			}
		}
	}

	if err == nil {
		line, err = xc.NextNBLine(&rest)
		if err == nil {
			if line == "members {" {
				line = strings.TrimSpace(rest[0]) // a peek
				for line == "memberInfo {" {
					var mi *MemberInfo
					mi, rest, err = ParseMemberInfoFromStrings(rest)
					if err != nil {
						break
					} else {
						memberInfos = append(memberInfos, mi)
						line = strings.TrimSpace(rest[0]) // a peek
					}
				}
				// we need a closing brace at this point
				line, err = xc.NextNBLine(&rest)
				if err == nil {
					if line != "}" {
						err = IllFormedClusterMember
					}
				}
			} else {
				err = IllFormedClusterMember
			}
		}
	}
	if err == nil {
		line, err = xc.NextNBLine(&rest)
		if err == nil {
			if line != "}" {
				err = IllFormedClusterMember
			}
		}
	}
	if err == nil {
		cm = &ClusterMember{
			Attrs:          attrs,
			ClusterName:    clusterName,
			ClusterAttrs:   clusterAttrs,
			ClusterID:      clusterID,
			ClusterMaxSize: clusterMaxSize,
			SelfIndex:      selfIndex,
			Members:        memberInfos,
			EPCount:        epCount,
			Node:           *node,
		}
	}
	return
}
