package main
import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"time"

	"github.com/guptarohit/asciigraph"
	"github.com/spf13/cobra"
)
type HistoryEntry struct {
	Date  string  `json:"date"`
	Value float64 `json:"value"`
}
type HistoryResponse struct {
	History []HistoryEntry `json:"history"`
}
type CompanyInfo struct {
	ID                    int     `json:"id"`
	Realm                 int     `json:"realm"`
	Name                  string  `json:"name"`
	Logo                  string  `json:"logo"`
	Level                 int     `json:"level"`
	Tier                  int     `json:"tier"`
	Year                  int     `json:"year"`
	DateJoined            string  `json:"dateJoined"`
	DateReset             string  `json:"dateReset"`
	Country               string  `json:"country"`
	Rank                  int     `json:"rank"`
	Rating                string  `json:"rating"`
	Value                 float64 `json:"value"`
	BuildingValue         float64 `json:"buildingValue"`
	TotalBuildings        int     `json:"totalBuildings"`
	Workers               int     `json:"workers"`
	AdministrationOverhead float64 `json:"administrationOverhead"`
	PatentsValue          float64 `json:"patentsValue"`
	BondsSold             float64 `json:"bondsSold"`
	Buildings             map[string]int `json:"buildings"`
	Tags                  map[string]string `json:"tags"`
	UpdatedAt             string  `json:"updatedAt"`
}
type CompanyResponse struct {
	Company CompanyInfo `json:"company"`
}
func fetchCompanyHistory(realmID, companyID string) {
	url := fmt.Sprintf("https://api.simcotools.com/v1/realms/%s/companies/%s/history", realmID, companyID)
	resp, err := http.Get(url)
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
	var apiResponse HistoryResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}
	if len(apiResponse.History) == 0 {
		fmt.Println("No historical data available.")
		return
	}
	sort.Slice(apiResponse.History, func(i, j int) bool {
		ti, _ := time.Parse(time.RFC3339, apiResponse.History[i].Date)
		tj, _ := time.Parse(time.RFC3339, apiResponse.History[j].Date)
		return ti.Before(tj)
	})
	values := []float64{}
	dates := []string{}
	for _, entry := range apiResponse.History {
		values = append(values, entry.Value)
		parsedTime, _ := time.Parse(time.RFC3339, entry.Date)
		dates = append(dates, parsedTime.Format("Jan 02"))
	}
	graph := asciigraph.Plot(values, asciigraph.Height(10), asciigraph.Caption("Company Value Over Time"))
	fmt.Println(graph)
	fmt.Println("Dates:")
	for i, date := range dates {
		if i%(len(dates)/5) == 0 || i == len(dates)-1 { 
			fmt.Printf("[%d] %s  ", i, date)
		}
	}
	fmt.Println()
}
func fetchCompanyInfo(realmID, companyID string) {
	url := fmt.Sprintf("https://api.simcotools.com/v1/realms/%s/companies/%s/", realmID, companyID)
	resp, err := http.Get(url)
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
	var apiResponse CompanyResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}
	company := apiResponse.Company
	fmt.Printf("Company: %s\n", company.Name)
	fmt.Printf("Rating: %s\n", company.Rating)
	fmt.Printf("Value: %.2f\n", company.Value)
	fmt.Printf("Rank: %d\n\n", company.Rank)
}
func main() {
	var showValueGraph bool
	var rootCmd = &cobra.Command{
		Use:   "fetch",
		Short: "Fetch company data from the API",
		Run: func(cmd *cobra.Command, args []string) {
			realmID, _ := cmd.Flags().GetString("realm")
			companyID, _ := cmd.Flags().GetString("company")

			if realmID == "" || companyID == "" {
				fmt.Println("Both --realm and --company flags are required.")
				return
			}
			fetchCompanyInfo(realmID, companyID)
			if showValueGraph {
				fetchCompanyHistory(realmID, companyID)
			}
		},
	}
	rootCmd.Flags().StringP("realm", "r", "", "Realm ID (required)")
	rootCmd.Flags().StringP("company", "c", "", "Company ID (required)")
	rootCmd.Flags().BoolVarP(&showValueGraph, "value", "v", false, "Show company value history as a graph")
	rootCmd.MarkFlagRequired("realm")
	rootCmd.MarkFlagRequired("company")
	rootCmd.Execute()
}
