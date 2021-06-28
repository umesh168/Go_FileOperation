package graphdb

import (
	"model"
	"testing"
)

var graph = GetGraphDBProvider("fileoperationgraph")

func TestCreateNode(t *testing.T) {
	var nodeObj model.Node
	nodeObj.SetName("TestFile")
	nodeObj.SetId("graphdb.com%2Ffileoperationgraph%2Fnodetype%2FTestFile")
	nodeId, GraphDBErr := graph.CreateNode(nodeObj)

	if GraphDBErr != nil {
		t.Errorf("TestCreateNode failed  got %v", GraphDBErr)
	} else {
		t.Logf("TestCreateNode succeed  got %v", nodeId)
	}
}
