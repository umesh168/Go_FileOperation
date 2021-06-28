package graphdb

import "TestFileApplication/model"

type GraphDB interface {
	SetGraph(databaseName string)
	CreateNode(node model.Node) (string, error)
	DeleteNodeById(nodeId string) error
	GetNodeById(nodeId string) (model.Node, error)
	GetNodeByIdWithProperties(nodeId string) (model.Node, error)
	UpdateNode(node model.Node) error
	UpdateNodeOnlyProperties(node model.Node) string
}
