package cmd

// Table column definitions per resource. These drive the --format table output
// and are kept in one place so column ordering stays consistent.
var (
	colKsAccounts = []string{
		"id", "advertiserId", "accountName", "accountType",
		"authStatus", "deliveryStatus", "activeStatus", "balance",
	}
	colKsCampaigns = []string{
		"id", "campaignId", "campaignName", "advertiserId",
		"putStatus", "status", "campaignType",
	}
	colKsUnits = []string{
		"id", "unitId", "unitName", "campaignId", "advertiserId",
		"putStatus", "status",
	}
	colKsCreatives = []string{
		"id", "creativeId", "creativeName", "unitId", "campaignId",
		"advertiserId", "putStatus", "status",
	}

	colKsAccountReport = []string{
		"id", "advertiserId", "accountName", "statDate", "statHour", "charge",
	}
	colKsCampaignReport = []string{
		"id", "advertiserId", "campaignId", "campaignName", "statDate", "charge",
	}
	colKsUnitReport = []string{
		"id", "advertiserId", "unitId", "unitName", "statDate", "charge",
	}
	colKsCreativeReport = []string{
		"id", "advertiserId", "creativeId", "creativeName", "statDate", "charge",
	}

	// colSummary is shared by ks and oe summary/top replies (same shape).
	colSummary = []string{
		"groupKey", "groupName", "charge", "rowCount",
	}
	colKsMetricMeta = []string{
		"id", "level", "field", "label", "groupName", "agg",
		"valueType", "sortOrder", "enabled", "sortable",
	}

	// Oceanengine (oe) column sets.
	colOeAccounts = []string{
		"id", "advertiserId", "accountName", "accountType",
		"authStatus", "deliveryStatus", "activeStatus", "balance", "budget",
	}
	colOeProjects = []string{
		"id", "projectId", "name", "advertiserId",
		"optStatus", "statusFirst", "deliveryMode", "landingType",
	}
	colOeUnits = []string{
		"id", "promotionId", "name", "projectId", "advertiserId",
		"optStatus", "statusFirst", "learningPhase",
	}

	colOeAccountReport = []string{
		"id", "advertiserId", "statDate", "statHour", "charge",
	}
	colOeProjectReport = []string{
		"id", "advertiserId", "projectId", "projectName", "statDate", "charge",
	}
	colOeUnitReport = []string{
		"id", "advertiserId", "promotionId", "promotionName", "statDate", "charge",
	}

	colOeMetricMeta = []string{
		"id", "level", "field", "label", "groupName", "agg",
		"valueType", "sortOrder", "enabled",
	}
	colOeBudgetVsActual = []string{
		"advertiserId", "accountName", "budgetMode", "budget",
		"totalCharge", "days", "avgDailyCharge", "budgetUsageRate", "balance",
	}

	// Tenant column set.
	colTenants = []string{
		"id", "name", "status", "createdBy", "createdAt", "updatedAt",
	}

	// User column set.
	colUser = []string{
		"id", "username", "name", "status", "currentTenantId", "createdAt", "updatedAt",
	}
)
