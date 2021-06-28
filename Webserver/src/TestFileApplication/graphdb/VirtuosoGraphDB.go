package graphdb

import (
	"TestFileApplication/constants"
	"encoding/json"
	"log"
	"TestFileApplication/model"
	"net/url"
	"strings"
	"time"
	"github.com/golang/glog"
	"github.com/knakk/sparql"
	"github.com/pkg/errors"
)

const (
	_SPACE             = " "
	_DOUBLEHASH        = "##"
	_TRIPLEHASH        = "###"
	NODE_NAME_RELATION = "Relation"
	RELATION_TYPE_IS_A = "Is A"
)

type triple map[string]string
type VirtuosoGraphDB struct{}

var databaseName string

func qureyExcecutor(query string) ([]byte, error) {
	glog.V(3).Info("qureyExcecutor [Started]", time.Now().Format("2006-01-02 15:04:05.000000"))
	glog.V(3).Info("qureyExcecutor inputs query=", query)
	log.Println("Executors Final Query :  ", query)
	res, GraphDBErr := executeRequestedQuery("qureyExcecutor", query)

	var returnObject []triple
	log.Println("Result: ")
	if len(res.Solutions()) > 0 {
		log.Println(res.Solutions()[0])

		for _, v := range res.Solutions() {
			m := make(map[string]string)
			if val, ok := v["s"]; ok {
				m["s"] = val.String()
			}

			if val, ok := v["p"]; ok {
				m["p"] = val.String()
			}

			if val, ok := v["o"]; ok {
				m["o"] = val.String()
			}
			returnObject = append(returnObject, m)
		}
	}

	log.Println(returnObject)
	response, err := json.Marshal(returnObject)

	if err != nil {
		errors.Wrap(err, "Query Execution Failed")
		glog.Error(err)
		log.Println(err)
	}

	glog.V(3).Info("qureyExcecutor [Completed]", time.Now().Format("2006-01-02 15:04:05.000000"))
	return response, GraphDBErr
}

func (vdb *VirtuosoGraphDB) SetGraph(graphName string) {
	databaseName = graphName
}

func (vdb *VirtuosoGraphDB) CreateNode(nodeObj model.Node) (string, error) {
	glog.V(3).Info("CreateNode [Started]", time.Now().Format("2006-01-02 15:04:05.000000"))
	glog.V(3).Info("CreateNode nodeObj=", nodeObj)
	nodeId := nodeObj.Id
	log.Println(nodeId)

	query := "INSERT INTO GRAPH " + "<" + constants.GH_VIRTUOSO_GRAPHNAME_URL_PREFIX + databaseName + ">" + " { <" + nodeId + "> rdfs:label \"" + refineLiteralValue(nodeObj.Name) + "\" " + "}"

	qureyExcecutor(query)
	glog.V(3).Info("CreateNode [Completed]", time.Now().Format("2006-01-02 15:04:05.000000"))
	return nodeId, nil
}

func (vdb *VirtuosoGraphDB) DeleteTestGraph(dbname string) {
	glog.V(3).Info("DeleteTestGraph [Started]", time.Now().Format("2006-01-02 15:04:05.000000"))
	glog.V(3).Info("DeleteTestGraph dbanme=", dbname)
	query := "DELETE FROM sys_rdf_schema WHERE RS_NAME LIKE '%pggraph%'"
	qureyExcecutor(query)
	glog.V(3).Info("DeleteTestGraph [Completed]", time.Now().Format("2006-01-02 15:04:05.000000"))
	// return nil
}

/*
this api will add node object with property
	it requerires node object
	and PropertyValueMap
		key : dataPropertyNodeId
		value : string

*/
func (vdb *VirtuosoGraphDB) CrateGraphNodeWithType(nodeObj model.Node, nodeTypeId string) (string, error) {
	glog.V(3).Info("CrateGraphNodeWithType [Started]", time.Now().Format("2006-01-02 15:04:05.000000"))
	glog.V(3).Info("CrateGraphNodeWithType inputs :", " nodeTypeId:", nodeTypeId, " nodeObj:", nodeObj)

	nodeId := nodeObj.Id
	glog.V(3).Info("NodeId: " + nodeId)

	query := "INSERT INTO GRAPH " + "<" + constants.GH_VIRTUOSO_GRAPHNAME_URL_PREFIX + databaseName + ">" + " { <" + nodeId + "> rdf:type <" + nodeTypeId + "> . " + " <" + nodeId + "> rdfs:label \"" + refineLiteralValue(nodeObj.Name) + "\" . "
	if nodeObj.PropertyValueMap != nil && len(nodeObj.PropertyValueMap) > 0 {
		for dataPropertyNodeId, value := range nodeObj.PropertyValueMap {
			query = query + "<" + nodeObj.Id + ">" + "<" + dataPropertyNodeId + ">" + " \"" + refineLiteralValue(value) + "\" . "
		}
	}

	query = query + "}"
	log.Println(query)
	qureyExcecutor(query)
	glog.V(3).Info("CrateGraphNodeWithType [Completed]", time.Now().Format("2006-01-02 15:04:05.000000"))
	return nodeId, nil
}



func (vdb *VirtuosoGraphDB) DeleteNodeById(nodeId string) error {
	glog.V(3).Info("DeleteNodeById [Started]", time.Now().Format("2006-01-02 15:04:05.000000"))
	glog.V(3).Info("DeleteNodeById inputs :", " nodeId:", nodeId)
	log.Println("Processing deleteNodeById with nodeId=" + nodeId)

	query := "DELETE where { GRAPH " + "<" + constants.GH_VIRTUOSO_GRAPHNAME_URL_PREFIX + databaseName + ">" + " { <" + nodeId + "> ?p ?o } }"
	res, err := qureyExcecutor(query)
	log.Println("Delete Response: ", res)

	query = "DELETE where { GRAPH " + "<" + constants.GH_VIRTUOSO_GRAPHNAME_URL_PREFIX + databaseName + ">" + " { ?s ?p  <" + nodeId + "> }}"
	res, err = qureyExcecutor(query)
	log.Println("Delete Response: ", res)
	if err != nil {
		glog.Error("qureyExcecutor failed: ", err)
	}
	glog.V(3).Info("DeleteNodeById [Completed]", time.Now().Format("2006-01-02 15:04:05.000000"))
	return err
}

//Note : Redesign changes apply on GetNodeById api
// 			NOW use rdfs:label instead of <name> AND not applying encde of nodeId and decode/QueryEscape on node name
func (vdb *VirtuosoGraphDB) GetNodeById(nodeId string) (model.Node, error) {
	glog.V(3).Info("GetNodeById [Started]", time.Now().Format("2006-01-02 15:04:05.000000"))
	glog.V(3).Info("Processing getNodeById with nodeId=" + nodeId)
	query := "SELECT * FROM " + "<" + constants.GH_VIRTUOSO_GRAPHNAME_URL_PREFIX + databaseName + ">" + " WHERE { VALUES (?s ?p) { (<" + nodeId + "> rdfs:label) } ?s ?p ?o }"

	glog.V(3).Info("Final query : " + query)
	res, err := qureyExcecutor(query)

	if err != nil {
		errors.Wrap(err, "Query Execution Failed")
		glog.Error(err)
		glog.V(3).Info("GetNodeById [Completed]", time.Now().Format("2006-01-02 15:04:05.000000"))
		return model.Node{}, err
	}

	var response []triple
	jsonErrorResponse := json.Unmarshal(res, &response)
	if jsonErrorResponse != nil {
		log.Println(jsonErrorResponse)
		errors.Wrap(jsonErrorResponse, "Query Execution Failed")
		glog.Error(jsonErrorResponse)
		glog.V(3).Info("GetNodeById [Completed]", time.Now().Format("2006-01-02 15:04:05.000000"))
		return model.Node{}, jsonErrorResponse
	}

	var responseNode model.Node

	if len(response) < 1 {
		glog.Error(errors.New("Current db have 0 instances present for nodeId :" + nodeId))
		glog.V(3).Info("GetNodeById [Completed]", time.Now().Format("2006-01-02 15:04:05.000000"))
		return responseNode, errors.New("Invalid graphDB id")
	} else {
		responseNode.Id = response[0]["s"]
		responseNode.Name = response[0]["o"]
		glog.V(4).Info(getStringObject(responseNode))
		glog.V(3).Info("GetNodeById [Completed]", time.Now().Format("2006-01-02 15:04:05.000000"))
		return responseNode, nil
	}

}

func containsCheckOnGivenArrayofString(inputArr []string, checkStr string) bool {
	for _, a := range inputArr {
		if a == checkStr {
			return true
		}
	}
	return false
}

func containsCheckOnGivenArrayofInt(inputArr []int, checkInt int) bool {
	for _, a := range inputArr {
		if a == checkInt {
			return true
		}
	}
	return false
}

func checkIfNodeTypeIsDataProperty(endNodeTypeIdToIsDataPropertyTypeMap map[string]bool, nodeTypeId string) bool {
	if ifExistInMap, isDataPropertyNodeType := endNodeTypeIdToIsDataPropertyTypeMap[nodeTypeId]; ifExistInMap {
		if isDataPropertyNodeType {
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}

func (vdb *VirtuosoGraphDB) GetNodeIdByName(nodeName string) string {
	glog.V(3).Info("GetNodeIdByName [Started] : " + time.Now().Format("2006-01-02 15:04:05.000000"))
	glog.V(3).Info("GetNodeIdByName inputs: ", " nodeName:", nodeName)
	log.Println("Processing getNodeIdByName with name=" + nodeName)

	query := "SELECT * FROM " + "<" + constants.GH_VIRTUOSO_GRAPHNAME_URL_PREFIX + databaseName + ">" + " WHERE { VALUES (?p ?o) {(rdfs:label \"" + refineLiteralValue(getEnecodedString(nodeName)) + "\")} ?s ?p ?o }"

	results, err := qureyExcecutor(query)

	if err != nil {
		errors.Wrap(err, "Query Execution Failed")
		glog.Error(err)
		log.Println(err)
	}

	log.Println(results)
	var dat []triple
	Jerr := json.Unmarshal(results, &dat)

	if Jerr != nil {
		errors.Wrap(err, "Failed to parse result")
		glog.Error(err)
		log.Println(Jerr)
	}
	glog.V(3).Info("Result : ", dat[0]["s"])
	glog.V(3).Info("GetNodeIdByName [Completed] : " + time.Now().Format("2006-01-02 15:04:05.000000"))
	return dat[0]["s"]
}

func (vdb *VirtuosoGraphDB) UpdateNodeOnlyProperties(nodeObj model.Node) string {
	log.Println("Processing updateNode with node=" + nodeObj.Id)
	glog.V(3).Info("UpdateNodeOnlyProperties [Started]", time.Now().Format("2006-01-02 15:04:05.000000"))
	glog.V(3).Info("inputs nodeObj count=", 1)
	glog.V(4).Info("inputs nodeObj=", nodeObj)

	uid := nodeObj.Id
	if len(uid) <= 0 {
		uid = vdb.GetNodeIdByName(nodeObj.Name)
	}

	if len(uid) <= 0 {
		return ""
	}

	nodeId := uid

	nodeName := nodeObj.Name
	insertPropertiesTriples := ""
	deleteQuery := ""
	propertyValueMap := nodeObj.PropertyValueMap

	if len(propertyValueMap) > 0 {
		for key, _ := range propertyValueMap {
			if key == "name" {
				// Intention here is update properties...we keep one property with name that should not be updated
				continue
			}
			var value = propertyValueMap[key]
			key = strings.ReplaceAll(key, _SPACE, _DOUBLEHASH)
			insertPropertiesTriples = insertPropertiesTriples + " <" + nodeId + "> <" + key + "> \"" + refineLiteralValue(value) + "\" ."

			deleteQuery = "DELETE where { GRAPH <" + constants.GH_VIRTUOSO_GRAPHNAME_URL_PREFIX + databaseName + "> { <" + nodeId + "> <" + key + "> ?o } }"
			qureyExcecutor(deleteQuery)
		}
	}

	query := "DELETE where { GRAPH <" + constants.GH_VIRTUOSO_GRAPHNAME_URL_PREFIX + databaseName + "> { <" + nodeId + "> rdfs:label ?o }  }"

	qureyExcecutor(query)

	query = "INSERT INTO " + "<" + constants.GH_VIRTUOSO_GRAPHNAME_URL_PREFIX + databaseName + ">" + " { <" + nodeId + "> rdfs:label \"" + refineLiteralValue(nodeName) + "\" . " + insertPropertiesTriples + "}"
	qureyExcecutor(query)

	glog.V(3).Info("UpdateNodeOnlyProperties [Completed]", time.Now().Format("2006-01-02 15:04:05.000000"))
	return nodeId
}

func GetNodeIdAsPerProductConvention(graphDB string, nodeName string, typeOfNode string) string {
	orgName := strings.TrimSuffix(graphDB, "graph")
	return "http://www.graphdb.com/" + "graphs/" + orgName + "/" + strings.ReplaceAll(strings.ToLower(typeOfNode), " ", "_") + "/" + strings.ReplaceAll(strings.ToLower(nodeName), " ", "_")
}
func encodeUTF8ANDLowercaseString(str string) string {
	return url.QueryEscape(strings.ToLower(str))
}

func nodeSetQueryExecutor(query string) []model.Node {
	var nodeSet []model.Node
	repo := VirtuosoInstanceProvider()

	log.Println("Final Query :  ", query)
	res, GraphDBErr := repo.Query(query)
	if GraphDBErr != nil {
		errors.Wrap(GraphDBErr, " Failed to execute query ")
		log.Println(GraphDBErr)
	}

	log.Println("Result: ")

	m := make(map[string]string)
	var returnObject []triple

	for _, v := range res.Solutions() {
		var nodeObj model.Node
		if val, ok := v["nodeName"]; ok {
			nodeObj.Name = val.String()
		}
		if val, ok := v["s"]; ok {
			m["s"] = val.String()
			nodeObj.Id = val.String()
		}

		if val, ok := v["p"]; ok {
			m["p"] = val.String()
		}

		log.Println(getDecodedLowerCaseNodeId(nodeObj.Id))
		nodeSet = append(nodeSet, nodeObj)
		returnObject = append(returnObject, m)
	}

	log.Println(nodeSet)
	return nodeSet
}

func getEncodedLowerCaseNodeId(id string) string {
	return url.QueryEscape(strings.ToLower(id))
}

func getEncodedNodeId(id string) string {
	return url.QueryEscape(id)
}

func getDecodedLowerCaseNodeId(id string) string {
	dec, err := url.QueryUnescape(strings.ToLower(id))
	if err != nil {
		log.Println(err)
	}
	return dec
}

func getDecodedString(id string) string {
	dec, err := url.QueryUnescape(id)
	if err != nil {
		log.Println(err)
	}
	return dec
}

func getEnecodedString(id string) string {
	return url.QueryEscape(id)
}

func getStringObject(object interface{}) string {
	objectBytes, _ := json.Marshal(object)
	stringObject := string(objectBytes)
	return stringObject
}

func executeRequestedQuery(funactionName, query string) (*sparql.Results, error) {
	glog.V(3).Info(funactionName + " : executeRequestedQuery [Started] : " + time.Now().Format("2006-01-02 15:04:05.000000"))
	repo := VirtuosoInstanceProvider()
	glog.V(3).Info(funactionName, " Executors Final Query : ", query)
	res, GraphDBErr := repo.Query(query)
	if GraphDBErr != nil {
		errors.Wrap(GraphDBErr, "Query Execution fail:")
		glog.Error(GraphDBErr)
	}
	glog.V(3).Info(funactionName + " executeRequestedQuery [Completed] : " + time.Now().Format("2006-01-02 15:04:05.000000"))
	return res, GraphDBErr
}

/*
The SPARQL 1.1 spec as written has three character escape mechanisms.
	One of them is weird, differs from Turtle, and makes injection attacks very hard to mitigate.
	When an application creates SPARQL query strings by string concatenation, it is potentially vulnerable to SPARQL injection attacks. This attack vector is analogous to SQL injection.
To be safe, an application must apply the appropriate escape sequences to user data before building the query string. How to do this for Turtle is pretty obvious.
		For example, in Javascript (ES6), for use with the triple-quote """...""" and '''...''' string literal forms:
		Reference link : https://github.com/w3c/sparql-12/issues/77
*/
func refineLiteralValue(value string) string {
	value = strings.ReplaceAll(value, "\\", "\\\\")
	value = strings.ReplaceAll(value, "\"", "\\\"")
	value = strings.ReplaceAll(value, "'", `\'`)
	return value
}
