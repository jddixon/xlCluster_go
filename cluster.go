package cluster

// xlReg_go/testCluster.go

// This file contains functions and structures used to create
// and manage and manage clusters of ClusterMembers.

import (
	"bytes"
	"encoding/hex"
	"fmt"
	ha "github.com/jddixon/hamt_go"
	xi "github.com/jddixon/xlNodeID_go"
	xn "github.com/jddixon/xlNode_go"
	//xo "github.com/jddixon/xlOverlay_go"
	xm "github.com/jddixon/xlUtil_go/math"
	"strconv"
	"strings"
	"sync"
)

var _ = fmt.Print

type Cluster struct {
	Name            string // must be unique within the registry
	ID              []byte // must be globally unique
	Attrs           uint64 // a field of bit flags
	curSize         uint32 // current size, may not exceed maxSize
	maxSize         uint32 // a maximum; must be > 0
	EPCount         uint32 // a positive integer, for now usually 1 or 2
	ClMembers       []*ClusterMember
	ClMembersByName map[string]*ClusterMember
	ClMembersByID   ha.HAMT
	Mu              sync.RWMutex
}

func NewCluster(name string, id *xi.NodeID, attrs uint64,
	maxSize, epCount uint32) (tc *Cluster, err error) {

	var m ha.HAMT

	if name == "" {
		name = "xlCluster"
	}
	nameMap := make(map[string]*ClusterMember)
	if epCount < 1 {
		err = ClusterMembersMustHaveEndPoint
	}
	if err == nil && maxSize < 1 {
		err = ClusterMustHaveMember
	} else {
		t := uint(xm.NextExp2_32(maxSize))
		m, err = ha.NewHAMT(DEFAULT_W, t)
	}
	if err == nil {
		tc = &Cluster{
			Attrs:           attrs,
			Name:            name,
			ID:              id.Value(),
			EPCount:         epCount,
			maxSize:         maxSize,
			ClMembersByName: nameMap,
			ClMembersByID:   m,
		}
	}
	return
}

// ATTRIBUTES ///////////////////////////////////////////////

func (tc *Cluster) GetName() string {
	return tc.Name
}
func (tc *Cluster) GetNodeID() (id *xi.NodeID) {
	id, _ = xi.New(tc.ID)
	return
}

func (tc *Cluster) GetCurSize() uint32 {
	var curSize uint32
	tc.Mu.RLock() // <-------------------------------------
	curSize = uint32(len(tc.ClMembers))
	tc.Mu.RUnlock() // <-----------------------------------
	return curSize
}
func (tc *Cluster) GetMaxSize() uint32 {
	return tc.maxSize
}

// UTILITY FUNCTIONS ////////////////////////////////////////////////

//
func (tc *Cluster) AddToCluster(node *xn.Node, attrs uint64) (
	member *ClusterMember, err error) {

	nodeID, err := xi.New(tc.ID)
	if err == nil {
		member = &ClusterMember{
			Attrs:          attrs,
			ClusterName:    tc.Name,
			ClusterID:      nodeID,
			ClusterAttrs:   tc.Attrs,
			ClusterMaxSize: tc.maxSize,
			EPCount:        tc.EPCount, // need to check
			SelfIndex:      uint32(len(tc.ClMembers)),
			// ClMembers not set
			Node: *node,
		}
		if err == nil {
			err = tc.AddMember(member)
		}
	}
	return
}

func (tc *Cluster) AddMember(member *ClusterMember) (err error) {

	// verify no existing member has the same name
	name := member.GetName()

	tc.Mu.RLock() // <------------------------------------
	_, ok := tc.ClMembersByName[name]
	tc.Mu.RUnlock() // <------------------------------------

	if ok {
		// DEBUG
		fmt.Printf("AddMember: ATTEMPT TO ADD EXISTING MEMBER %s\n", name)
		// END
		err = ClusterMemberNameInUse
	}
	if err == nil {
		var (
			entry interface{}
			bKey  ha.BytesKey
		)
		// check for entry in HAMT
		tc.Mu.RLock() // <---------------------------------
		bKey, err = ha.NewBytesKey(tc.ID)
		entry, err = tc.ClMembersByID.Find(bKey)
		tc.Mu.RUnlock() // <-------------------------------
		if err == nil {
			if entry != nil {
				err = ClusterMemberIDInUse
			}
		}
		if err == nil {
			tc.Mu.Lock()               // <------------------
			index := len(tc.ClMembers) // DEBUG
			_ = index                  // we might want to use this
			tc.ClMembers = append(tc.ClMembers, member)
			tc.ClMembersByName[name] = member
			bKey, err = ha.NewBytesKey(member.GetNodeID().Value())
			if err == nil {
				err = tc.ClMembersByID.Insert(bKey, member)
			}
			tc.Mu.Unlock() // <----------------------------
		}
	}
	return
}

///**
// * XXX Locking occurs at a lower level, making deadlocks possible.
// */
//func (tc *Cluster) Start() (err error) {
//	members := tc.ClMembers //  []*ClusterMember)
//	if members != nil {
//		for i := 0; err == nil && i < len(members); i++ {
//			err = members[i].Node.Run()
//		}
//	}
//	return
//}
//
///**
// * XXX Locking occurs at a lower level, making deadlocks possible.
// */
//func (tc *Cluster) Stop() (err error) {
//	members := tc.ClMembers //  []*ClusterMember)
//	if members != nil {
//		for i := 0; err == nil && i < len(members); i++ {
//			err = members[i].Node.Close()
//		}
//	}
//	return
//}

// EQUAL ////////////////////////////////////////////////////////////
func (tc *Cluster) Equal(any interface{}) bool {

	if any == tc {
		return true
	}
	if any == nil {
		return false
	}
	switch v := any.(type) {
	case *Cluster:
		_ = v
	default:
		return false
	}
	other := any.(*Cluster) // type assertion
	if tc.Attrs != other.Attrs {
		// DEBUG
		fmt.Printf("tc.Equal: ATTRS DIFFER %s vs %s\n", tc.Attrs, other.Attrs)
		// END
		return false
	}
	if tc.Name != other.Name {
		// DEBUG
		fmt.Printf("tc.Equal: NAMES DIFFER %s vs %s\n", tc.Name, other.Name)
		// END
		return false
	}
	if !bytes.Equal(tc.ID, other.ID) {
		// DEBUG
		tcHexID := hex.EncodeToString(tc.ID)
		otherHexID := hex.EncodeToString(other.ID)
		fmt.Printf("tc.Equal: IDs DIFFER %s vs %s\n", tcHexID, otherHexID)
		// END
		return false
	}
	if tc.EPCount != other.EPCount {
		// DEBUG
		fmt.Printf("tc.Equal: EPCOUNTS DIFFER %d vs %d\n",
			tc.EPCount, other.EPCount)
		// END
		return false
	}
	if tc.maxSize != other.maxSize {
		// DEBUG
		fmt.Printf("tc.Equal: MAX SIZES DIFFER %d vs %d\n",
			tc.maxSize, other.maxSize)
		// END
		return false
	}
	if tc.GetCurSize() != other.GetCurSize() {
		// DEBUG
		fmt.Printf("tc.Equal:ACTUAL SIZES DIFFER %d vs %d\n",
			tc.GetCurSize(), other.GetCurSize())
		// END
		return false
	}
	// Members			[]*ClientInfo
	for i := uint32(0); i < tc.GetCurSize(); i++ {
		rcMember := tc.ClMembers[i]
		otherMember := other.ClMembers[i]
		if !rcMember.Equal(otherMember) {
			return false
		}
	}
	return true
}

// SERIALIZATION ////////////////////////////////////////////////////

func (tc *Cluster) Strings() (ss []string) {

	ss = []string{"testCluster {"}

	ss = append(ss, fmt.Sprintf("    Attrs: 0x%016x", tc.Attrs))
	ss = append(ss, "    Name: "+tc.Name)
	ss = append(ss, "    ID: "+hex.EncodeToString(tc.ID))
	ss = append(ss, fmt.Sprintf("    EPCount: %d", tc.EPCount))
	ss = append(ss, fmt.Sprintf("    curSize: %d", tc.curSize))
	ss = append(ss, fmt.Sprintf("    maxSize: %d", tc.maxSize))

	ss = append(ss, "    Members {")
	for i := 0; i < len(tc.ClMembers); i++ {
		mem := tc.ClMembers[i].Strings()
		for i := 0; i < len(mem); i++ {
			ss = append(ss, "        "+mem[i])
		}
	}
	ss = append(ss, "    }")
	ss = append(ss, "}")

	return
}

func (tc *Cluster) String() string {
	return strings.Join(tc.Strings(), "\n")
}
func ParseCluster(s string) (tc *Cluster, rest []string, err error) {
	ss := strings.Split(s, "\n")
	return ParseClusterFromStrings(ss)
}
func ParseClusterFromStrings(ss []string) (
	tc *Cluster, rest []string, err error) {

	var (
		attrs   uint64
		name    string
		id      *xi.NodeID
		EPCount uint32
		curSize uint32
		maxSize uint32
	)
	rest = ss

	line := xn.NextNBLine(&rest) // the line is trimmed
	if line != "testCluster {" {
		fmt.Println("MISSING regCluster {")
		err = IllFormedCluster
	} else {
		line = xn.NextNBLine(&rest)
		if strings.HasPrefix(line, "Attrs: ") {
			var i int64
			i, err = strconv.ParseInt(line[7:], 0, 64)
			if err == nil {
				attrs = uint64(i)
			}
		} else {
			fmt.Printf("BAD ATTRS in line '%s'", line)
			err = IllFormedCluster
		}
	}
	if err == nil {
		line = xn.NextNBLine(&rest)
		if strings.HasPrefix(line, "Name: ") {
			name = line[6:]
		} else {
			fmt.Printf("BAD NAME in line '%s'", line)
			err = IllFormedCluster
		}
	}
	if err == nil {
		// collect ID
		line = xn.NextNBLine(&rest)
		if strings.HasPrefix(line, "ID: ") {
			var val []byte
			val, err = hex.DecodeString(line[4:])
			if err == nil {
				id, err = xi.New(val)
			}
		} else {
			fmt.Println("BAD ID")
			err = IllFormedCluster
		}
	}
	if err == nil {
		line = xn.NextNBLine(&rest)
		if strings.HasPrefix(line, "EPCount: ") {
			var count int
			count, err = strconv.Atoi(line[9:])
			if err == nil {
				EPCount = uint32(count)
			}
		} else {
			fmt.Println("BAD END POINT COUNT")
			err = IllFormedCluster
		}
	}
	if err == nil {
		line = xn.NextNBLine(&rest)
		if strings.HasPrefix(line, "curSize: ") {
			var size int
			size, err = strconv.Atoi(line[9:])
			if err == nil {
				curSize = uint32(size)
			}
		} else {
			fmt.Println("BAD MAX_SIZE")
			err = IllFormedCluster
		}
	}
	if err == nil {
		line = xn.NextNBLine(&rest)
		if strings.HasPrefix(line, "maxSize: ") {
			var size int
			size, err = strconv.Atoi(line[9:])
			if err == nil {
				maxSize = uint32(size)
			}
		} else {
			fmt.Println("BAD MAX_SIZE")
			err = IllFormedCluster
		}
	}
	if err == nil {
		tc, err = NewCluster(name, id, attrs, maxSize, EPCount)
	}
	if err == nil {
		tc.curSize = curSize
		line = xn.NextNBLine(&rest)
		if line == "Members {" {
			for {
				line = strings.TrimSpace(rest[0]) // peek
				if line == "}" {
					break
				}
				var member *ClusterMember
				member, rest, err = ParseClusterMemberFromStrings(rest)
				if err != nil {
					break
				}
				err = tc.AddMember(member)
				if err != nil {
					break
				}
			}
		} else {
			err = MissingMembersList
		}
	}

	// expect closing brace for Members list
	if err == nil {
		line = xn.NextNBLine(&rest)
		if line != "}" {
			err = MissingClosingBrace
		}
	}
	// expect closing brace  for cluster
	if err == nil {
		line = xn.NextNBLine(&rest)
		if line != "}" {
			err = MissingClosingBrace
		}
	}

	return
}
