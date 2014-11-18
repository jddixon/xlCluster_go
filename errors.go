package cluster

import (
	e "errors"
)

var (
	BadAttrsLine                   = e.New("badly formed attrs line")
	ClusterMembersMustHaveEndPoint = e.New("cluster members must have at least one endPoint")
	ClusterMemberIDInUse           = e.New("cluster member id already in use")
	ClusterMemberNameInUse         = e.New("cluster member name already in use")
	ClusterMustHaveMember          = e.New("cluster must have a member!")
	ClusterMustHaveTwo             = e.New("cluster must have at least two members")
	IllFormedCluster               = e.New("ill-formed cluster serialization")
	IllFormedClusterMember         = e.New("ill-formed cluster member serialization")
	MemberMustHaveEndPoint         = e.New("member must have at least one endPoint")
	MissingClosingBrace            = e.New("missing closing brace")
	MissingMembersList             = e.New("missing members list")
	NilNode                        = e.New("nil node argument")
	WrongNumberOfBytesInAttrs      = e.New("wrong number of bytes in attrs")
)
