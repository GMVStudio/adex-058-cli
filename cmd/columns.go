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

	colKsSummary = []string{
		"groupKey", "groupName", "charge", "rowCount",
	}
	colKsMetricMeta = []string{
		"id", "level", "field", "label", "groupName", "agg",
		"valueType", "sortOrder", "enabled", "sortable",
	}
)
