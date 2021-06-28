package model

type Node struct {
	Id   string
	Name string
	RawValue string
	NType            string            `json:"NType,omitempty"`
	PropertyValueMap map[string]string `json:",omitempty"`
}

func NewNode(id string, name string) Node {
	var node = new(Node)
	node.Id = id
	node.Name = name
	return *node
}

func NewNodeWithEmptyPropertyMap(id string, name string) Node {
	var node = new(Node)
	node.Id = id
	node.Name = name
	node.PropertyValueMap = map[string]string{}
	return *node
}

func (node *Node) SetId(id string) {
	node.Id = id
}

func (node Node) GetId() string {
	return node.Id
}

func (node *Node) SetName(name string) {
	node.Name = name
}

func (node Node) GetName() string {
	return node.Name
}

func (node *Node) SetPropertyValueMap(propertyValueMap map[string]string) {
	node.PropertyValueMap = propertyValueMap
}

func (node Node) GetPropertyValueMap() map[string]string {
	return node.PropertyValueMap
}

func (node Node) GetPropertyValueMapValueByKey(key string) string {
	return node.PropertyValueMap[key]
}
