package dbcallhelper

import (
	"constants"
	"TestFileApplication/graphdb"
	"log"
	"TestFileApplication/model"
	"net/url"
	"strings"
	"github.com/golang/glog"
)

var prefix = "http://www.graphdb.com/"

func GetNodeIdAsPerProductConvention(graphDB string, nodeName string, typeOfNode string) string {
	//temp code to retrive orgname
	orgName := strings.TrimSuffix(graphDB, "graph")
	nodeName = strings.ReplaceAll(strings.ToLower(nodeName), " ", "_")
	return prefix + "graphs/" + orgName + "/" + GetFormatedNodeTypeForNode(typeOfNode) + "/" + nodeName
}

func GetNodeTypeIdAsPerProductConvention(graphDB string, nodeTypeName string) string {
	return prefix + "ontologies/" + getIndustryCode(graphDB) + "#" + nodeTypeName
}

func GetFormatedNodeTypeForNode(nodeTypeName string) string {
	return strings.ReplaceAll(strings.ToLower(nodeTypeName), " ", "_")
}
func GetNodeIdAsPerProductConvention_old(graphDB string, nodeName string, typeOfNode string) string {
	return strings.ToLower(constants.GRAPH_URL_PREFIX + graphDB + "/" + typeOfNode + "/" + nodeName)
}

func GetEncodedNodeIdAsPerProductConvention(graphDB string, nodeName string, typeOfNode string) string {
	//temp code to retrive orgname
	orgName := strings.TrimSuffix(graphDB, "graph")
	return encodeUTF8ANDLowercaseString(prefix + "graphs/" + orgName + "/" + GetFormatedNodeTypeForNode(typeOfNode) + "/" + strings.ReplaceAll(strings.ToLower(nodeName), " ", "_"))
}

func encodeUTF8ANDLowercaseString(str string) string {
	return url.QueryEscape(strings.ToLower(str))
}

func GetNodeById(graphDB string, nodeId string) model.Node {
	var graph = graphdb.GetGraphDBProvider(graphDB) // we should get client graph specific instance using some factory
	if strings.TrimSpace(nodeId) == "" {
		glog.V(3).Info("GetNodeById : Invalid Node id")
		return model.Node{}
	}
	node, _ := graph.GetNodeById(nodeId) // TODO improve error handling

	return node
}

func CreateNodeWithtype(graphDB string, nodeName string, nodeType string) string {
	var nodeObj = model.Node{}

	nodeObj.Name = nodeName
	nodeObj.Id = GetEncodedLowerCaseNodeId(constants.GH_VIRTUOSO_GRAPHNAME_URL_PREFIX + graphDB + "/" + nodeType + "/" + nodeName)
	var graph = graphdb.GetGraphDBProvider(graphDB) // we should get client graph specific instance using some factory
	nodeId, _ := graph.CreateNode(nodeObj)          // TODO improve error handling

	return nodeId
}

func CrateNode(graphDB string, nodeObj model.Node) string {
	var graph = graphdb.GetGraphDBProvider(graphDB) // we should get client graph specific instance using some factory
	nodeId, _ := graph.CreateNode(nodeObj)          // TODO improve error handling

	return nodeId
}

func DeleteTestGraph(graphDB string) {
	var graph = graphdb.GetGraphDBProvider(graphDB) // we should get client graph specific instance using some factory
	graph.DeleteTestGraph(graphDB)                  // TODO improve error handling
}

func CrateGraphNodeWithType(graphDB string, nodeObj model.Node, nodeTypeId string) string {
	var graph = graphdb.GetGraphDBProvider(graphDB)                // we should get client graph specific instance using some factory
	nodeId, _ := graph.CrateGraphNodeWithType(nodeObj, nodeTypeId) // TODO improve error handling

	return nodeId
}

func UpdateNodeOnlyPropertyProperties(graphDB string, nodeObj model.Node) string {
	var graph = graphdb.GetGraphDBProvider(graphDB)   // we should get client graph specific instance using some factory
	nodeId := graph.UpdateNodeOnlyProperties(nodeObj) // TODO improve error handling

	return nodeId
}

func DeleteNodeById(graphDB, nodeId string) error {
	graph := graphdb.GetGraphDBProvider(graphDB)
	err := graph.DeleteNodeById(nodeId)
	return err
}

func GetEncodedLowerCaseNodeId(id string) string {
	return url.QueryEscape(strings.ToLower(id))
}

func GetDecodedString(encodedStr string) string {
	dec, err := url.QueryUnescape(encodedStr)
	if err != nil {
		log.Fatal(err)
	}
	return dec
}

func containsCheckOnGivenArrayofString(inputArr []string, checkStr string) bool {
	for _, a := range inputArr {
		if a == checkStr {
			return true
		}
	}
	return false
}

/*
The SPARQL 1.1 spec as written has three character escape mechanisms.
	One of them is weird, differs from Turtle, and makes injection attacks very hard to mitigate.
	When an application creates SPARQL query strings by string concatenation, it is potentially vulnerable to SPARQL injection attacks. This attack vector is analogous to SQL injection.
To be safe, an application must apply the appropriate escape sequences to user data before building the query string. How to do this for Turtle is pretty obvious.
		For example, in Javascript (ES6), for use with the triple-quote """...""" and '''...''' string literal forms:
		Reference link : https://github.com/w3c/sparql-12/issues/77
*/
func RefineLiteralValue(value string) string {
	value = strings.ReplaceAll(value, "\\", "\\\\")
	value = strings.ReplaceAll(value, "\"", "\\\"")
	value = strings.ReplaceAll(value, "'", `\'`)
	return value
}

func getIndustryCode(graphDB string) string {
	industryCode := ""
	switch graphDB {
	case "fileoperationgraph":
		industryCode = "fileoperation"
		break
	}
	return industryCode
}
