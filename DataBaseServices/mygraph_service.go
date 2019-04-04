// file: services/movie_service.go
/*eslint-disable */
package DataBaseServices

import (
	"context"
	"fmt"
	"log"

	"github.com/dgraph-io/dgo"
	"github.com/dgraph-io/dgo/protos/api"

	grpc "google.golang.org/grpc"
)

type MygraphService interface {
	Mutate(string) string
	Query([]byte) []byte
	VQuery([]byte, map[string]string) []byte
}

func NewMygraphService() *mygraphService {
	return &mygraphService{
		key: "dgraph",
		//dg:  newClient(),
		txn: newClient().NewTxn(),
	}
}

//newClient :To create a client, dial a connection to Dgraph’s external gRPC port (typically 9080). The following code snippet shows just one connection.
func newClient() *dgo.Dgraph {
	// Dial a gRPC connection. The address to dial to can be configured when
	// setting up the dgraph cluster.
	d, err := grpc.Dial("localhost:9080", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	return dgo.NewDgraphClient(
		api.NewDgraphClient(d),
	)
}

type mygraphService struct {
	key string
	//dg  *dgo.Dgraph
	txn *dgo.Txn
}

/*
all mutation api type
	SetJson             []byte   `protobuf:"bytes,1,opt,name=set_json,json=setJson,proto3" json:"set_json,omitempty"`
	DeleteJson          []byte   `protobuf:"bytes,2,opt,name=delete_json,json=deleteJson,proto3" json:"delete_json,omitempty"`
	SetNquads           []byte   `protobuf:"bytes,3,opt,name=set_nquads,json=setNquads,proto3" json:"set_nquads,omitempty"`
	DelNquads           []byte   `protobuf:"bytes,4,opt,name=del_nquads,json=delNquads,proto3" json:"del_nquads,omitempty"`
*/
//MutateRDF : mutate database with  JSON and RDF N-Quad.
//allowed type here is
// Ttype : transaction type
//	* "set" : insert or update data
//	* "delete" : delete data
func (s *mygraphService) MutateRDF(qry []byte, Ttype string) interface{} {
	defer s.txn.Discard(context.Background())
	switch Ttype {
	case "set":
		res, err := s.txn.Mutate(context.Background(), &api.Mutation{
			CommitNow: true,
			SetNquads: qry,
		})
		if err != nil {
			log.Fatal(err)
		}
		return res
	case "delete":
		res, err := s.txn.Mutate(context.Background(), &api.Mutation{
			CommitNow: true,
			DelNquads: qry,
		})
		if err != nil {
			log.Fatal(err)
		}
		return res
	}
	return "null"
}

//Mutate : mutate the data base with json qry
func (s *mygraphService) Mutate(qry []byte) (FULL string, UID string) {
	defer s.txn.Discard(context.Background())
	res, err := s.txn.Mutate(context.Background(), &api.Mutation{
		CommitNow: true,
		SetJson:   qry,
	})

	if err != nil {
		log.Fatal(err)
	}
	return fmt.Sprint("%+v", res), res.Uids["blank-0"]
}

func (s *mygraphService) Query(qry string) []byte {
	res, err := s.txn.Query(context.Background(), qry)
	if err != nil {
		log.Fatal(err)
	}
	return res.Json
}
func (s *mygraphService) VQuery(qry string, data map[string]string) []byte {
	res, err := s.txn.QueryWithVars(context.Background(), qry, data)
	if err != nil {
		log.Fatal(err)
	}
	return res.Json
}
