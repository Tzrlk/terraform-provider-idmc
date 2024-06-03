// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccExampleDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccExampleDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.scaffolding_example.test", "id", "example-id"),
				),
			},
		},
	})
}

const testAccExampleDataSourceConfig = `
data "scaffolding_example" "test" {
  configurable_attribute = "example"
}
`
// TODO: Mock this JSON response:
/*
[
    {
        "@type": "activityLogEntry",
        "id": "000001C100000000000C",
        "type": "DRS",
        "objectName": "drstask1",
        "runId": 1,
        "runtimeEnvironmentId": "00000C25000000000002",
	       "startTime": "2012-07-30T13:30:30.000Z",
        "endTime": "2012-07-30T13:32:30.000Z",
        "state": 1,
        "failedSourceRows": 0,
        "successSourceRows": 39,
        "failedTargetRows": 0,
        "successTargetRows": 39,
        "errorMsg": null,
        "entries": [
            {
                "@type": "activityLogEntry",
                "id": "128964732",
                "type": "DRS",
                "objectName": "Contact",
                "runId": 0,
                "runtimeEnvironmentId": "00000C25000000000002",
                "agentId: "01000008000000000006",
                "startTime": "2012-07-30T13:32:31.000Z",
                "endTime": "2012-07-30T13:35:31.000Z",
                "state": 1,
                "isStopped": FALSE,
                "failedSourceRows": 0,
                "successSourceRows": 39,
                "failedTargetRows": 0,
                "successTargetRows": 39,
                "errorMsg": "No errors encountered.",
                "entries": []
            },
        ]
    },
    {
        "@type": "activityLogEntry",
        "id": "010000C1000000000PGP",
        "type": "MTT_TEST",
        "objectId": "0100000Z00000000001N",
        "objectName": "Mapping-MultiSource",
        "runId": 12,
        "startTime": "2020-03-27T08:05:56.000Z",
        "endTime": "2020-03-27T08:06:07.000Z",
        "startTimeUtc": "2020-03-27T12:05:56.000Z",
        "endTimeUtc": "2020-03-27T12:06:07.000Z",
        "state": 2,
        "failedSourceRows": 0,
        "successSourceRows": 800,
        "failedTargetRows": 200,
        "successTargetRows": 600,
        "startedBy": "di@infa.com",
        "runContextType": "ICS_UI",
        "entries": [
            {
                "@type": "activityLogEntry",
                "id": "118964723",
                "type": "MTT_TEST",
                "objectName": "",
                "runId": 12,
                "agentId": "01000008000000000004",
                "runtimeEnvironmentId": "01000025000000000004",
                "startTime": "2020-03-27T08:05:56.000Z",
                "endTime": "2020-03-27T08:06:07.000Z",
                "startTimeUtc": "2020-03-27T12:05:56.000Z",
                "endTimeUtc": "2020-03-27T12:06:07.000Z",
                "state": 2,
                "failedSourceRows": 0,
                "successSourceRows": 800,
                "failedTargetRows": 200,
                "successTargetRows": 600,
                "errorMsg": null,
                "startedBy": "di@infa.com",
                "runContextType": "ICS_UI",
                "entries": [],
                "subTaskEntries": [],
                "logEntryItemAttrs": {
                    "CONSUMED_COMPUTE_UNITS": "0.0",
                    "ERROR_CODE": "0",
                    "IS_SERVER_LESS": "false",
                    "REQUESTED_COMPUTE_UNITS": "0.0",
                    "Session Log File Name": "s_mtt_0Sr7LdcbAG2ldG33Lp8koQ_2.log"
                },
                "totalSuccessRows": 0,
                "totalFailedRows": 0,
                "stopOnError": false,
                "hasStopOnErrorRecord": false,
                "contextExternalId": "0100000Z00000000001N",
                "transformationEntries": [
                    {
                        "@type": "transformationLogEntry",
                        "id": "141332309",
                        "txName": "FFSource2",
                        "txType": "SOURCE",
                        "successRows": 600,
                        "failedRows": 0
                    },
                    {
                        "@type": "transformationLogEntry",
                        "id": "141332310",
                        "txName": "FFSource1",
                        "txType": "SOURCE",
                        "successRows": 200,
                        "failedRows": 0
                    },
                    {
                        "@type": "transformationLogEntry",
                        "id": "141332311",
                        "txName": "FFTarget.csv",
                        "txType": "TARGET",
                        "successRows": 600,
                        "affectedRows": 600,
                        "failedRows": 0
                    },
                    {
                        "@type": "transformationLogEntry",
                        "id": "141332312",
                        "txName": "MYSQLTarget",
                        "txType": "TARGET",
                        "successRows": 0,
                        "affectedRows": 0,
                        "failedRows": 200
                    }
                ]
            }
        ]
    }
]
 */
