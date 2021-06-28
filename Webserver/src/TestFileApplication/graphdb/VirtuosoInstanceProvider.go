package graphdb

import (
	"TestFileApplication/constants"
	"log"
	"time"

	"github.com/golang/glog"
	"github.com/knakk/sparql"
	"github.com/pkg/errors"
)

func VirtuosoInstanceProvider() *sparql.Repo {
	repo, GraphDBErr := sparql.NewRepo(constants.GH_VIRTUOSO_CONNECTION_URL,
		sparql.DigestAuth(constants.GH_VIRTUOSO_USER, constants.GH_VIRTUOSO_PASSWORD),
		sparql.Timeout(time.Millisecond*150000),
		// was getting lakhs rows..that query was timing out with default 35 seconds...changed it to 2.5 minutes for now...
	)
	if GraphDBErr != nil {
		errors.Wrap(GraphDBErr, "sparql NewRepo failed")
		glog.Error("VirtuosoInstanceProvider: sparql NewRepo failed ", GraphDBErr)
		log.Fatal(GraphDBErr)
	}
	return repo
}
