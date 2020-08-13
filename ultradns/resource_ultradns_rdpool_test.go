package ultradns

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terra-farm/udnssdk"
)

func TestAccUltradnsRdpool(t *testing.T) {
	var record udnssdk.RRSet
	domain, _ := os.LookupEnv("ULTRADNS_DOMAIN")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccRdpoolCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testCfgRdpoolMinimal, domain),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUltradnsRecordExists("ultradns_rdpool.it", &record),
					// Specified
					resource.TestCheckResourceAttr("ultradns_rdpool.it", "zone", domain),
					resource.TestCheckResourceAttr("ultradns_rdpool.it", "name", "test-rdpool-minimal"),
					resource.TestCheckResourceAttr("ultradns_rdpool.it", "ttl", "300"),

					// hashRdatas(): 10.6.0.1 -> 2847814707
					resource.TestCheckResourceAttr("ultradns_rdpool.it", "rdata.2847814707", "10.6.0.1"),
					// Defaults
					resource.TestCheckResourceAttr("ultradns_rdpool.it", "description", "Minimal RD Pool"),
					// Generated
					resource.TestCheckResourceAttr("ultradns_rdpool.it", "id", fmt.Sprintf("test-rdpool-minimal:%s", domain)),
					resource.TestCheckResourceAttr("ultradns_rdpool.it", "hostname", fmt.Sprintf("test-rdpool-minimal.%s.", domain)),
				),
			},

			{
				ResourceName:      "ultradns_rdpool.it",
				ImportState:       true,
				ImportStateVerify: true,
			},


			{
				Config: fmt.Sprintf(testCfgRdpoolMaximal, domain),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUltradnsRecordExists("ultradns_rdpool.it", &record),
					// Specified
					resource.TestCheckResourceAttr("ultradns_rdpool.it", "zone", domain),
					resource.TestCheckResourceAttr("ultradns_rdpool.it", "name", "test-rdpool-maximal"),
					resource.TestCheckResourceAttr("ultradns_rdpool.it", "ttl", "300"),
					resource.TestCheckResourceAttr("ultradns_rdpool.it", "description", "traffic controller pool with all settings tuned"),

					// hashRdatas(): 10.6.1.1 -> 2826722820
					resource.TestCheckResourceAttr("ultradns_rdpool.it", "rdata.2826722820", "10.6.1.1"),

					// hashRdatas(): 10.6.1.2 -> 829755326
					resource.TestCheckResourceAttr("ultradns_rdpool.it", "rdata.829755326", "10.6.1.2"),

					// Generated
					resource.TestCheckResourceAttr("ultradns_rdpool.it", "id", fmt.Sprintf("test-rdpool-maximal:%s", domain)),
					resource.TestCheckResourceAttr("ultradns_rdpool.it", "hostname", fmt.Sprintf("test-rdpool-maximal.%s.", domain)),
				),
			},

			{
				ResourceName:      "ultradns_rdpool.it",
				ImportState:       true,
				ImportStateVerify: true,
			},

		},
	})
}

const testCfgRdpoolMinimal = `
resource "ultradns_rdpool" "it" {
  zone        = "%s"
  name        = "test-rdpool-minimal"
  ttl         = 300
  description = "Minimal RD Pool"
  rdata       = ["10.6.0.1"]
}
`

const testCfgRdpoolMaximal = `
resource "ultradns_rdpool" "it" {
  zone        = "%s"
  name        = "test-rdpool-maximal"
  order       = "ROUND_ROBIN"
  ttl         = 300
  description = "traffic controller pool with all settings tuned"
  rdata       = ["10.6.1.1","10.6.1.2"]
}
`
