package cmd

import (
	"github.com/spf13/cobra"
)

// newKsAccountsCmd creates "adex ks accounts" (ListKsAdAccounts).
func newKsAccountsCmd(f *Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "accounts",
		Short: "List Kuaishou ad accounts",
		Long: `List Kuaishou ad accounts (GET /v1/ks/ad-accounts).

Examples:
  adex ks accounts --tenant 6 --page-size 20
  adex ks accounts --tenant 6 --order-by balance --order-desc --format table
  adex ks accounts --tenant 6 --page-all --jq '.items[].advertiserId'`,
		RunE: func(cmd *cobra.Command, args []string) error {
			tenant, err := f.resolveTenant(cmd)
			if err != nil {
				return err
			}
			params := map[string]interface{}{"tenant_id": tenant}

			advertiserID, _ := cmd.Flags().GetString("advertiser")
			accountName, _ := cmd.Flags().GetString("account-name")
			accountType, _ := cmd.Flags().GetString("account-type")
			authStatus, _ := cmd.Flags().GetString("auth-status")
			deliveryStatus, _ := cmd.Flags().GetString("delivery-status")
			activeStatus, _ := cmd.Flags().GetString("active-status")
			ownerUserID, _ := cmd.Flags().GetInt("owner-user")

			setString(params, "advertiser_id", advertiserID)
			setString(params, "account_name", accountName)
			setString(params, "account_type", accountType)
			setString(params, "auth_status", authStatus)
			setString(params, "delivery_status", deliveryStatus)
			setString(params, "active_status", activeStatus)
			setInt(params, "owner_user_id", ownerUserID, 0)

			applyOrder(cmd, params)
			applyPaging(cmd, params)

			return f.runList(cmd, "/v1/ks/ad-accounts", params, colKsAccounts)
		},
	}

	addTenantFlag(cmd)
	cmd.Flags().String("advertiser", "", "advertiser ID filter")
	cmd.Flags().String("account-name", "", "account name fuzzy match")
	cmd.Flags().String("account-type", "", "account type filter")
	cmd.Flags().String("auth-status", "", "auth status filter")
	cmd.Flags().String("delivery-status", "", "delivery status filter")
	cmd.Flags().String("active-status", "", "active status filter")
	cmd.Flags().Int("owner-user", 0, "owner user ID filter")
	addOrderFlags(cmd, "id")
	addPagingFlags(cmd)
	addJQFlag(cmd)

	return cmd
}
