package check

import (
	"errors"
	"fmt"

	"github.com/vmware/govmomi/vim25/mo"
	v1 "k8s.io/api/core/v1"
	"k8s.io/klog/v2"
)

// CheckResourcePoolPermissions confirms that resources associated with the node maintain required privileges.
type CheckResourcePoolPermissions struct {
	resourcePools map[string]*mo.ResourcePool
}

var _ NodeCheck = &CheckResourcePoolPermissions{}

func (c *CheckResourcePoolPermissions) Name() string {
	return "CheckResourcePoolPermissions"
}

func (c *CheckResourcePoolPermissions) StartCheck() error {
	c.resourcePools = make(map[string]*mo.ResourcePool)
	return nil
}

func (c *CheckResourcePoolPermissions) checkResourcePoolPrivileges(ctx *CheckContext, vm *mo.VirtualMachine) *CheckError {
	resourcePool, resourcePoolPath, err := getResourcePool(ctx, vm.Reference())
	if err != nil {
		klog.Info("resource pool could not be obtained for %v", vm.Reference())
		return nil
	}

	if _, ok := c.resourcePools[resourcePoolPath]; ok {
		klog.Infof("privileges for resource pool %s have already been checked", resourcePoolPath)
	}
	c.resourcePools[resourcePoolPath] = resourcePool

	if _, ok := ctx.VMConfig.VirtualCenter[ctx.VMConfig.Workspace.VCenterIP]; !ok {
		return NewCheckError(VcenterNotFound, errors.New("vcenter instance not found in the virtual center map"))
	}

	if err := comparePrivileges(ctx.Context, ctx.Username, resourcePool.Reference(), ctx.AuthManager, permissions[permissionCluster]); err != nil {
		return NewCheckError(ResourcePoolMissingPermissions, fmt.Errorf("missing privileges for resource pool %s: %s", resourcePoolPath, err.Error()))
	}

	return nil
}

func (c *CheckResourcePoolPermissions) CheckNode(ctx *CheckContext, node *v1.Node, vm *mo.VirtualMachine) *CheckError {

	// Skip permission checks if pre-existing resource pool was not defined
	if ctx.VMConfig.Workspace.ResourcePoolPath == "" {
		return nil
	}

	errCheck := c.checkResourcePoolPrivileges(ctx, vm)
	if errCheck != nil {
		return errCheck
	}
	return nil
}

func (c *CheckResourcePoolPermissions) FinishCheck(ctx *CheckContext) {
}
