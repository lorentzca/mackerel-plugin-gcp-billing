package main

import (
	"flag"
	"fmt"
	"golang.org/x/net/context"
	"log"
	"os"
	"strings"

	"cloud.google.com/go/bigquery"
	"github.com/jessevdk/go-flags"
	mp "github.com/mackerelio/go-mackerel-plugin-helper"
	"google.golang.org/api/iterator"
)

var opts struct {
	ProjectId *string `short:"p" long:"projectid" required:"true" description:"Project ID"`
	Dataset   *string `short:"d" long:"dataset" required:"true" description:"Dataset"`
	Table     *string `short:"t" long:"table" required:"true" description:"Table"`
}

type GcpBillingPlugin struct {
	Prefix string
}

func (g GcpBillingPlugin) GraphDefinition() map[string](mp.Graphs) {
	labelPrefix := strings.Title(g.Prefix)
	return map[string](mp.Graphs){
		g.Prefix: mp.Graphs{
			Label: labelPrefix,
			Unit:  "float",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "cost", Label: "Cost"},
			},
		},
	}
}

func (g GcpBillingPlugin) FetchMetrics() (map[string]interface{}, error) {
	bl, err := billing()
	if err != nil {
		return nil, fmt.Errorf("Faild to fetch billing metrics: %s", err)
	}
	return map[string]interface{}{"cost": bl}, nil
}

func billing() (float64, error) {
	ctx := context.Background()

	_, err := flags.Parse(&opts)
	if err != nil {
		os.Exit(1)
	}

	projectID := *opts.ProjectId
	dataset := *opts.Dataset
	table := *opts.Table

	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	q := client.Query("SELECT SUM(cost) AS sum_cost, LEFT (FORMAT_UTC_USEC( UTC_USEC_TO_MONTH(TIMESTAMP_TO_USEC(start_time))), 7) AS month FROM " + dataset + "." + table + " GROUP BY month")

	it, err := q.Read(ctx)
	if err != nil {
		log.Fatalf("Failed to querying: %v", err)
	}

	var values []bigquery.Value
	for {
		err := it.Next(&values)
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterator: %v", err)
		}
	}
	return values[0].(float64), err
}

func main() {
	optPrefix := flag.String("metric-key-prefix", "billing", "Metric key prefix")
	optTempfile := flag.String("tempfile", "", "Temp file name")

	g := GcpBillingPlugin{
		Prefix: *optPrefix,
	}
	helper := mp.NewMackerelPlugin(g)
	helper.Tempfile = *optTempfile
	if helper.Tempfile == "" {
		helper.Tempfile = fmt.Sprintf("/tmp/mackerel-plugin-%s", *optPrefix)
	}
	helper.Run()
}
