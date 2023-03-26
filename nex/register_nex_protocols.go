package nex

import (
	"github.com/PretendoNetwork/mario-vs-dk-tipping-stars-secure/globals"
	nex_datastore "github.com/PretendoNetwork/mario-vs-dk-tipping-stars-secure/nex/datastore"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

func registerNEXProtocols() {
	datastoreProtocol := nexproto.NewDataStoreProtocol(globals.NEXServer)

	datastoreProtocol.DeleteObject(nex_datastore.DeleteObject)
	datastoreProtocol.GetMeta(nex_datastore.GetMeta)
	datastoreProtocol.PostMetaBinary(nex_datastore.PostMetaBinary)
	datastoreProtocol.CompletePostObjects(nex_datastore.CompletePostObjects)
	datastoreProtocol.ChangeMeta(nex_datastore.ChangeMeta)
}
