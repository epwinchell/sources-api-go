package dao

import (
	m "github.com/RedHatInsights/sources-api-go/model"
	"github.com/RedHatInsights/sources-api-go/util"
	"github.com/redhatinsights/platform-go-middlewares/identity"
)

// GetTenantDao is a function definition that can be replaced in runtime in case some other DAO provider is
// needed.
var GetTenantDao func() TenantDao

// getDefaultRhcConnectionDao gets the default DAO implementation which will have the given tenant ID.
func getDefaultTenantDao() TenantDao {
	return &tenantDaoImpl{}
}

// init sets the default DAO implementation so that other packages can request it easily.
func init() {
	GetTenantDao = getDefaultTenantDao
}

type tenantDaoImpl struct{}

func (t *tenantDaoImpl) GetOrCreateTenantID(identity *identity.Identity) (int64, error) {
	// Start setting up the query.
	query := DB.
		Debug().
		Model(&m.Tenant{})

	// Query by OrgId or EBS account number.
	var tenant m.Tenant
	if identity.OrgID != "" {
		tenant.OrgID = identity.OrgID

		query.Where("org_id = ?", tenant.OrgID)
	} else {
		tenant.ExternalTenant = identity.AccountNumber

		query.Where("external_tenant = ?", tenant.ExternalTenant)
	}

	// Find the tenant, scanning into the struct above
	result := query.
		First(&tenant)

	// Looks like we didn't find it, create it and return the ID.
	if result.Error != nil {
		result = DB.Create(&tenant)
	}

	return tenant.Id, result.Error
}

func (t *tenantDaoImpl) TenantByIdentity(id *identity.Identity) (*m.Tenant, error) {
	// Start setting up the query.
	query := DB.
		Debug().
		Model(&m.Tenant{})

	// Query by OrgId or EBS account number.
	var tenant m.Tenant
	if id.OrgID != "" {
		tenant.OrgID = id.OrgID

		query.Where("org_id = ?", tenant.OrgID)
	} else {
		tenant.ExternalTenant = id.AccountNumber

		query.Where("external_tenant = ?", tenant.ExternalTenant)
	}

	// Find the tenant, scanning into the struct above
	err := query.
		First(&tenant).
		Error

	if err != nil {
		return nil, util.NewErrNotFound("tenant")
	}

	return &tenant, nil
}
