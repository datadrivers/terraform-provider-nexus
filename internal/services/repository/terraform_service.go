package repository

import (
	"reflect"
	"unsafe"

	nexus "github.com/datadrivers/go-nexus-client/nexus3"
	nexusclient "github.com/datadrivers/go-nexus-client/nexus3/pkg/client"
	"github.com/datadrivers/go-nexus-client/nexus3/pkg/repository/common"
)

const terraformAPIEndpoint = common.RepositoryAPIEndpoint + "/terraform"

func terraformHostedService(n *nexus.NexusClient) *common.RepositoryService[TerraformHostedRepository] {
	return common.NewRepositoryService[TerraformHostedRepository](terraformAPIEndpoint+"/hosted", extractNexusClient(n))
}

func terraformProxyService(n *nexus.NexusClient) *common.RepositoryService[TerraformProxyRepository] {
	return common.NewRepositoryService[TerraformProxyRepository](terraformAPIEndpoint+"/proxy", extractNexusClient(n))
}

func terraformGroupService(n *nexus.NexusClient) *common.RepositoryService[TerraformGroupRepository] {
	return common.NewRepositoryService[TerraformGroupRepository](terraformAPIEndpoint+"/group", extractNexusClient(n))
}

// extractNexusClient accesses the unexported raw HTTP client from NexusClient via reflection.
// This bridge is required because go-nexus-client v1.20.0 does not yet include Terraform
// repository support; the private field name and type are stable across that release.
func extractNexusClient(n *nexus.NexusClient) *nexusclient.Client {
	v := reflect.ValueOf(n).Elem().FieldByName("client")
	return reflect.NewAt(v.Type(), unsafe.Pointer(v.UnsafeAddr())).Elem().Interface().(*nexusclient.Client)
}
