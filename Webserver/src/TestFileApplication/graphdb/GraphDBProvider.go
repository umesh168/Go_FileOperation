package graphdb

import (
	"time"

	"github.com/golang/glog"
)

func GetGraphDBProvider(databaseName string) *VirtuosoGraphDB {
	glog.V(3).Info("GetGraphDBProvider [Started]", time.Now().Format("2006-01-02 15:04:05.000000"))
	glog.V(3).Info("GetGraphDBProvider inputs databaseName=" + databaseName)
	var graph = new(VirtuosoGraphDB)
	graph.SetGraph(databaseName)
	glog.V(4).Info("Result graph=", graph)
	glog.V(3).Info("GetGraphDBProvider [Completed]", time.Now().Format("2006-01-02 15:04:05.000000"))
	return graph
}
