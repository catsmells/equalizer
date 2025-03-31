package main
import (
	"fmt"
	"io/ioutil"
	"net/http"
	"github.com/spf13/cobra"
)
func fetchCompanyData(realmID, companyID string) {
	url := fmt.Sprintf("https://api.simcotools.com/v1/realms/%s/companies/%s", realmID, companyID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	req.Header.Set("Accept", "*/*")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}
	fmt.Println("Response Status:", resp.Status)
	fmt.Println("Response Body:", string(body))
}
func main() {
	var rootCmd = &cobra.Command{
		Use:   "fetch",
		Short: "Fetch company data from the SimCompanies API",
		Run: func(cmd *cobra.Command, args []string) {
			realmID, _ := cmd.Flags().GetString("realm")
			companyID, _ := cmd.Flags().GetString("company")
			if realmID == "" || companyID == "" {
				fmt.Println("Both --realm and --company flags are required.")
				return
			}
			fetchCompanyData(realmID, companyID)
		},
	}
	rootCmd.Flags().StringP("realm", "r", "", "Realm ID (required)")
	rootCmd.Flags().StringP("company", "c", "", "Company ID (required)")
	rootCmd.MarkFlagRequired("realm")
	rootCmd.MarkFlagRequired("company")
	rootCmd.Execute()
}
