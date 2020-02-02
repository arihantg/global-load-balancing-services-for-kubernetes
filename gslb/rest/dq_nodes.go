package rest

import (
	"errors"
	"os"
	"sync"

	avimodels "github.com/avinetworks/sdk/go/models"
	"github.com/davecgh/go-spew/spew"
	"gitlab.eng.vmware.com/orion/container-lib/utils"
	avicache "gitlab.eng.vmware.com/orion/mcc/gslb/cache"
	"gitlab.eng.vmware.com/orion/mcc/gslb/gslbutils"
	"gitlab.eng.vmware.com/orion/mcc/gslb/nodes"
)

var restLayer *RestOperations
var restOnce sync.Once

type RestOperations struct {
	cache             *avicache.AviCache
	aviRestPoolClient *utils.AviRestClientPool
}

func NewRestOperations(cache *avicache.AviCache, aviRestPoolClient *utils.AviRestClientPool) *RestOperations {
	restOnce.Do(func() {
		restLayer = &RestOperations{cache: cache, aviRestPoolClient: aviRestPoolClient}
	})
	return restLayer
}

func (restOp *RestOperations) DqNodes(key string) {
	gslbutils.Logf("key: %s, msg: starting rest layer sync", key)
	// got the key from graph layer, let's fetch the model
	ok, aviModelIntf := nodes.SharedAviGSGraphLister().Get(key)
	if !ok {
		gslbutils.Logf("key: %s, msg: %s", key, "no model found for the key")
		return
	}

	tenant, gsName := utils.ExtractNamespaceObjectName(key)
	gsKey := avicache.TenantName{Tenant: tenant, Name: gsName}
	gsCacheObj := restOp.getGSCacheObj(gsKey, key)
	gslbutils.Logf("GS cache object: %v", gsCacheObj)
	if aviModelIntf == nil {
		if gsCacheObj != nil {
			gslbutils.Logf("key: %s, msg: %s", key, "no model found, this is a GS deletion case")
			// Delete case can have two sub-cases:
			// 1. Just delete the IP if present
			// if len(gsCacheObj.Members) > 1 {
			// 	for _, member := gsCacheObj.Members {
			// 	}
			// }
			// 2. Delete the GS object only if this is the last IP present.
			restOp.deleteGSOper(gsCacheObj, tenant, key)
		}
		return
	}
	if ok && aviModelIntf != nil {
		aviModel := aviModelIntf.(*nodes.AviGSObjectGraph)
		gslbutils.Logf("key: %s, msg: GS create/update", key)
		if aviModel.GSNode == nil {
			gslbutils.Warnf("key: %s, msg: %s", key, "no gslbservice in the model")
			return
		}
		restOp.RestOperation(gsName, tenant, aviModel.GSNode, gsCacheObj, key)
	}
}

func (restOp *RestOperations) RestOperation(gsName, tenant string, aviGSNode *nodes.AviGSNode,
	gsCacheObj *avicache.AviGSCache, key string) {
	gsKey := avicache.TenantName{Tenant: tenant, Name: gsName}
	var operation *utils.RestOp
	if gsCacheObj != nil {
		var restOps []*utils.RestOp
		var cksum uint32
		cksum = aviGSNode.GetChecksum()
		if gsCacheObj.CloudConfigCksum == cksum {
			gslbutils.Logf("key: %s, GSLBService: %s, msg: %s", key, gsName,
				"the checksums are same for the GSLB service, ignoring")
			return
		}
		gslbutils.Logf("key: %s, GSLBService: %s, oldCksum: %s, newCksum: %s, msg: %s", key, gsName,
			gsCacheObj.CloudConfigCksum, cksum, "checksums are different for the GSLB Service")
		// it should be a PUT call
		operation = restOp.AviGSBuild(aviGSNode, utils.RestPut, gsCacheObj, key)
		gslbutils.Logf("gsKey: %s, restOps: %v, operation: %v", gsKey, restOps, operation)
	} else {
		// its a post operation
		gslbutils.Logf("key: %s, operation: POST, msg: GS not found in cache", key)
		operation = restOp.AviGSBuild(aviGSNode, utils.RestPost, nil, key)
	}
	restOp.ExecuteRestAndPopulateCache(operation, gsKey, key)
}

func (restOp *RestOperations) ExecuteRestAndPopulateCache(operation *utils.RestOp, gsKey avicache.TenantName, key string) {
	// Choose a AVI client based on the model name hash. This would ensure that the same worker queue processes updates for a
	// given GS everytime.
	bkt := utils.Bkt(key, utils.NumWorkersGraph)
	gslbutils.Logf("key: %s, queue: %d, msg: processing in rest queue", key, bkt)
	if len(restOp.aviRestPoolClient.AviClient) > 0 {
		aviClient := restOp.aviRestPoolClient.AviClient[bkt]
		err := restOp.aviRestPoolClient.AviRestOperate(aviClient, []*utils.RestOp{operation})
		if err != nil {
			gslbutils.Errf("key: %s, msg: rest operation error: %s", key, err)
			return
		}
		// rest call executed successfully
		gslbutils.Logf("key: %s, msg: rest call executed successfully, will update cache", key)
		if operation.Err == nil && (operation.Method == utils.RestPost || operation.Method == utils.RestPut) {
			if operation.Model == "GSLBService" {
				restOp.AviGSCacheAdd(operation, key)
			} else {
				gslbutils.Errf("key: %s, method: %s, model: %s, msg: invalid model", key, operation.Method,
					operation.Model)
			}
		} else {
			if operation.Model == "GSLBService" {
				// delete call executed
				restOp.AviGSCacheDel(restOp.cache, operation, key)
			}
		}
	}
}

func (restOp *RestOperations) AviGSBuild(gsMeta *nodes.AviGSNode, restMethod utils.RestMethod,
	cacheObj *avicache.AviGSCache, key string) *utils.RestOp {
	gslbutils.Logf("key: %s, msg: creating rest operation", key)
	// description field needs references
	var gslbPoolMembers []*avimodels.GslbPoolMember
	var gslbSvcGroups []*avimodels.GslbPool
	for id, ip := range gsMeta.Members {
		if ip == "" {
			continue
		}
		enabled := true
		// ipAddr := ip
		ipVersion := "V4"
		ratio := int32(1)

		gslbPoolMember := avimodels.GslbPoolMember{
			Enabled: &enabled,
			IP:      &avimodels.IPAddr{Addr: &gsMeta.Members[id], Type: &ipVersion},
			Ratio:   &ratio,
		}
		gslbPoolMembers = append(gslbPoolMembers, &gslbPoolMember)
	}
	// Now, build a GSLB pool
	algorithm := "GSLB_ALGORITHM_ROUND_ROBIN"
	poolEnabled := true
	poolName := gsMeta.Name + "-10"
	priority := int32(10)
	gslbPool := avimodels.GslbPool{
		Algorithm: &algorithm,
		Enabled:   &poolEnabled,
		Members:   gslbPoolMembers,
		Name:      &poolName,
		Priority:  &priority,
	}
	gslbSvcGroups = append(gslbSvcGroups, &gslbPool)

	// Now, build the GSLB service
	ctrlHealthStatusEnabled := true
	createdBy := "mcc-gslb"
	// TODO: description to be appropriately filled
	gsEnabled := true
	healthMonitorScope := "GSLB_SERVICE_HEALTH_MONITOR_ALL_MEMBERS"
	isFederated := true
	minMembers := int32(0)
	gsName := gsMeta.Name
	poolAlgorithm := "GSLB_SERVICE_ALGORITHM_PRIORITY"
	resolveCname := false
	sitePersistenceEnabled := false
	tenantRef := "https://" + os.Getenv("CTRL_IPADDRESS") + "/api/tenant/" + utils.ADMIN_NS
	useEdnsClientSubnet := true
	wildcardMatch := false

	aviGslbSvc := avimodels.GslbService{
		ControllerHealthStatusEnabled: &ctrlHealthStatusEnabled,
		CreatedBy:                     &createdBy,
		DomainNames:                   gsMeta.DomainNames,
		Enabled:                       &gsEnabled,
		Groups:                        gslbSvcGroups,
		HealthMonitorScope:            &healthMonitorScope,
		IsFederated:                   &isFederated,
		MinMembers:                    &minMembers,
		Name:                          &gsName,
		PoolAlgorithm:                 &poolAlgorithm,
		ResolveCname:                  &resolveCname,
		SitePersistenceEnabled:        &sitePersistenceEnabled,
		UseEdnsClientSubnet:           &useEdnsClientSubnet,
		WildcardMatch:                 &wildcardMatch,
		TenantRef:                     &tenantRef,
	}

	path := "/api/gslbservice/"
	operation := utils.RestOp{Path: path, Obj: aviGslbSvc, Tenant: gsMeta.Tenant, Model: "GSLBService", Version: AviVersion}

	if restMethod == utils.RestPost {
		operation.Method = utils.RestPost
		return &operation
	}
	// Else, its a PUT call
	operation.Path = path + cacheObj.Uuid
	operation.Method = utils.RestPut

	gslbutils.Logf(spew.Sprintf("key: %s, gsMeta: %s, msg: GS rest operation %v\n", key, *gsMeta, utils.Stringify(operation)))
	return &operation
}

func (restOp *RestOperations) getGSCacheObj(gsKey avicache.TenantName, key string) *avicache.AviGSCache {
	gsCache, found := restOp.cache.AviCacheGet(gsKey)
	if found {
		gsCacheObj, ok := gsCache.(*avicache.AviGSCache)
		if !ok {
			gslbutils.Warnf("key: %s, msg: %s", key, "invalid GS object found, ignoring...")
			return nil
		}
		return gsCacheObj
	}
	gslbutils.Logf("key: %s, gsKey: %v, msg: GS cache object not found", key, gsKey)
	return nil
}

func (restOp *RestOperations) AviGSDel(uuid string, tenant string, key string) *utils.RestOp {
	path := "/api/gslbservice/" + uuid
	operation := utils.RestOp{Path: path, Method: "DELETE", Tenant: tenant, Model: "GSLBService",
		Version: utils.CtrlVersion}
	gslbutils.Logf(spew.Sprintf("GSLB Service DELETE Restop %v\n", utils.Stringify(operation)))
	return &operation
}

func (restOp *RestOperations) deleteGSOper(gsCacheObj *avicache.AviGSCache, tenant string, key string) bool {
	var restOps *utils.RestOp
	bkt := utils.Bkt(key, utils.NumWorkersGraph)
	aviclient := restOp.aviRestPoolClient.AviClient[bkt]
	if gsCacheObj != nil {
		operation := restOp.AviGSDel(gsCacheObj.Uuid, tenant, key)
		restOps = operation
		err := restOp.aviRestPoolClient.AviRestOperate(aviclient, []*utils.RestOp{restOps})
		if err != nil {
			// TODO: Just log it for now, will add a retry logic later
			gslbutils.Warnf("key: %s, GSLBService: %s, msg: %s", key, gsCacheObj.Uuid,
				"failed to delete GSLB Service")
			return false
		}
		// Clear all the cache objects which were deleted
		restOp.AviGSCacheDel(restOp.cache, operation, key)
		return true
	}
	return false
}

func (restOp *RestOperations) AviGSCacheDel(gsCache *avicache.AviCache, op *utils.RestOp, key string) {
	gsKey := avicache.TenantName{Tenant: op.Tenant, Name: op.ObjName}
	gslbutils.Logf("key: %s, msg: deleting gs cache", key, gsKey)
	gsCache.AviCacheDelete(gsKey)
}

func (restOp *RestOperations) AviGSCacheAdd(operation *utils.RestOp, key string) error {
	if (operation.Err != nil) || (operation.Response == nil) {
		gslbutils.Warnf("key: %s, response: %s, msg: rest operation has err or no response for VS: %s", key,
			operation.Response, operation.Err)
		return errors.New("rest operation errored")
	}

	respElem, err := RestRespArrToObjByType(operation, "gslbservice", key)
	if err != nil || respElem == nil {
		gslbutils.Warnf("key: %s, resp: %s, msg: unable to find GS object in resp", key, operation.Response)
		return errors.New("GS not found")
	}
	name, ok := respElem["name"].(string)
	if !ok {
		gslbutils.Warnf("key: %s, resp: %s, msg: name not present in response", key, respElem)
		return errors.New("name not present in response")
	}
	uuid, ok := respElem["uuid"].(string)
	if !ok {
		gslbutils.Warnf("key: %s, resp: %s, msg: uuid not present in response", key, respElem)
		return errors.New("uuid not present in response")
	}

	cksum, members, err := GetDetailsFromAviGSLB(respElem)
	if err != nil {
		gslbutils.Errf("key: %s, resp: %v, msg: error in getting checksum for gslb svc: %s", key, respElem, err)
	}
	gslbutils.Logf("key: %s, resp: %s, cksum: %d, msg: GS information", key, utils.Stringify(respElem), cksum)
	k := avicache.TenantName{Tenant: operation.Tenant, Name: name}
	gsCache, ok := restOp.cache.AviCacheGet(k)
	if ok {
		gsCacheObj, found := gsCache.(*avicache.AviGSCache)
		if found {
			gsCacheObj.Uuid = uuid
			gsCacheObj.CloudConfigCksum = cksum
			gsCacheObj.Members = members
			gslbutils.Logf(spew.Sprintf("key: %s, cacheKey: %v, value: %v, msg: updated GS cache\n", key, k,
				utils.Stringify(gsCacheObj)))
		}
	} else {
		// New cache object
		gsCacheObj := avicache.AviGSCache{
			Name:             name,
			Tenant:           operation.Tenant,
			Uuid:             uuid,
			Members:          members,
			CloudConfigCksum: cksum,
		}
		restOp.cache.AviCacheAdd(k, &gsCacheObj)
		gslbutils.Logf(spew.Sprintf("key: %s, cacheKey: %v, value: %v, msg: added GS to the cache", key, k,
			utils.Stringify(gsCacheObj)))
	}

	return nil
}

func SyncFromNodesLayer(key string) error {
	cache := avicache.NewAviCache()
	gslbutils.Logf("cache values: %v", cache)
	aviclient := avicache.SharedAviClients()
	restLayerF := NewRestOperations(cache, aviclient)
	restLayerF.DqNodes(key)
	return nil
}