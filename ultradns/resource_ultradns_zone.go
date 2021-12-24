package ultradns

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	log "github.com/sirupsen/logrus"
	udnssdk "github.com/ultradns/ultradns-sdk-go"
)

type zoneProperties struct {
	Name                 string `json:"name"`
	AccountName          string `json:"accountName"`
	Type                 string `json:"type"`
	DnssecStatus         string `json:"dnssecStatus"`
	Status               string `json:"status"`
	Owner                string `json:"owner"`
	ResourceRecordCount  int    `json:"resourceRecordCount"`
	LastModifiedDateTime string `json:"lastModifiedDateTime"`
}

type zoneTsigDTO struct {
	TsigKeyName   string `json:"tsigKeyName,omitempty"`
	TsigKeyValue  string `json:"tsigKeyValue,omitempty"`
	Description   string `json:"description,omitempty"`
	TsigAlgorithm string `json:"tsigAlgorithm,omitempty"`
}

type zonePrimaryCreateInfo struct {
	ForceImport bool   `json:"forceImport,omitempty"`
	CreateType  string `json:"createType,omitempty"`
	NameServer  struct {
		Ip            string `json:"ip,omitempty"`
		TsigKey       string `json:"tsigKey,omitempty"`
		TsigKeyValue  string `json:"tsigKeyValue,omitempty"`
		TsigAlgorithm string `json:"tsigAlgorithm,omitempty"`
	} `json:"nameServer,omitempty"`
	OriginalZoneName string       `json:"originalZoneName"`
	RestrictIPList   []string     `json:"restrictIPList"`
	Tsig             *zoneTsigDTO `json:"tsig,omitempty"`
	NotifyAddresses  []string     `json:"notifyAddresses"`
	Inherit          string       `json:"inherit,omitempty"`
}

type zoneSecondaryCreateInfo struct {
	PrimaryNameServers       udnssdk.NameServerLists `json:"primaryNameServers"`
	NotificationEmailAddress string                  `json:"notificationEmailAddress"`
}

type zoneAliasCreateDTO struct {
	OriginalZoneName string `json:"originalZoneName,omitempty"`
}

type zoneCreateDTO struct {
	Properties          zoneProperties          `json:"properties"`
	PrimaryCreateInfo   zonePrimaryCreateInfo   `json:"primaryCreateInfo,omitempty"`
	SecondaryCreateInfo zoneSecondaryCreateInfo `json:"secondaryCreateInfo,omitempty"`
	AliasCreateInfo     zoneAliasCreateDTO      `json:"aliasCreateInfo,omitempty"`
	ChangeComment       string                  `json:"changeComment"`
}

func newZoneResource(d *schema.ResourceData) (zoneCreateDTO, error) {
	log.Infof("Schema = %+v", d)
	z := zoneCreateDTO{}

	if attr, ok := d.GetOk("type"); ok {
		switch attr.(string) {
		case "PRIMARY":
			if attr, ok := d.GetOk("create_type"); ok {
				z.PrimaryCreateInfo.CreateType = attr.(string)
				if attr.(string) == "COPY" {
					if attr, ok := d.GetOk("original_zone_name"); ok {
						z.PrimaryCreateInfo.OriginalZoneName = attr.(string)
					}
				}
			}
		case "SECONDARY":
			// z.SecondaryCreateInfo.PrimaryNameServers.NameServerList["nameServerIP1"] = map[interface{}]interface{}{"ip": "8.8.8.8"}
			ipList1 := make(map[string]interface{})
			ipList1["ip"] = "8.8.8.8"
			outerList := make(map[string]interface{})
			outerList["nameServerIp1"] = ipList1
			z.SecondaryCreateInfo.PrimaryNameServers.NameServerList = outerList
		case "ALIAS":
			if attr, ok := d.GetOk("alias_target"); ok {
				z.AliasCreateInfo.OriginalZoneName = attr.(string)
			}
		default:
			return zoneCreateDTO{}, errors.New("type must be one of PRIMARY, SECONDARY, or ALIAS")
		}
		z.Properties.Type = attr.(string)
	}

	if attr, ok := d.GetOk("name"); ok {
		z.Properties.Name = attr.(string)
		d.SetId(attr.(string))
	}

	if attr, ok := d.GetOk("account"); ok {
		z.Properties.AccountName = attr.(string)
	}

	z.Properties.Status = "ACTIVE"

	return z, nil
}

func validateNewZone(zoneType interface{}, other string) ([]string, []error) {
	switch zoneType {
	case "PRIMARY":
		break
	case "SECONDARY":
		break
	case "ALIAS":
		break
	default:
		return nil, []error{errors.New("invalid zone type: must be one of PRIMARY, SECONDARY, ALIAS")}
	}
	return nil, nil
}

func validateFqdn(zoneName interface{}, other string) ([]string, []error) {
	z := zoneName.(string)
	pattern := regexp.MustCompile(`\.$`)
	if !pattern.MatchString(z) {
		return nil, []error{errors.New("name must be an FQDN ending in a period")}
	}
	return nil, nil
}

func resourceUltradnsZone() *schema.Resource {
	return &schema.Resource{
		Create: resourceUltraDNSZoneCreate,
		Read:   resourceUltraDNSZoneRead,
		Update: resourceUltraDNSZoneUpdate,
		Delete: resourceUltraDNSZoneDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateFqdn,
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateNewZone,
			},
			"account": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"create_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "NEW",
			},
			"original_zone_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"alias_target": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},

		Importer: &schema.ResourceImporter{
			State: resourceUltraDNSZoneImport,
		},
	}
}

// CRUD Operations

func resourceUltraDNSZoneCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*udnssdk.Client)

	z, err := newZoneResource(d)
	if err != nil {
		return err
	}

	log.Printf("[INFO] ultradns_zone create: %+v", z)

	// UltraDNS SDK doesn't implement this; we emulate the SDK's API requests here
	var ignored interface{}

	_, err = client.Do("post", "/zones", z, ignored)

	if err != nil {
		return fmt.Errorf("create failed: %v", err)
	}
	d.SetId(z.Properties.Name)
	log.Printf("[INFO] ultradns_record.id: %v", d.Id())

	return resourceUltraDNSZoneRead(d, meta)
}

func resourceUltraDNSZoneRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*udnssdk.Client)

	z, err := newZoneResource(d)
	if err != nil {
		return err
	}

	uri := fmt.Sprintf("/zones/%s", z.Properties.Name)
	var zone udnssdk.Zone

	zoneJson, err := client.Do("get", uri, new(interface{}), &zone)
	log.Infof("Zone Info: %v", zoneJson)

	if err != nil {
		uderr, ok := err.(*udnssdk.ErrorResponseList)
		if ok {
			for _, r := range uderr.Responses {
				if r.ErrorCode == 70002 {
					d.SetId("")
					return nil
				}
				return fmt.Errorf("not found: %v", err)
			}
		}
		return fmt.Errorf("not found: %v", err)
	}

	return populateResourceDataFromZone(zone, d)
}

func resourceUltraDNSZoneUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*udnssdk.Client)

	z, err := newZoneResource(d)
	if err != nil {
		return err
	}

	// UltraDNS SDK doesn't implement this; we emulate the SDK's API requests here
	var ignored interface{}
	uri := fmt.Sprintf("/zones/%s", z.Properties.Name)
	_, err = client.Do("put", uri, z, ignored)

	if err != nil {
		return fmt.Errorf("update failed: %v", err)
	}
	log.Printf("[INFO] ultradns_zone updated: %+v", z)

	return nil
}

func resourceUltraDNSZoneDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*udnssdk.Client)

	z, err := newZoneResource(d)
	if err != nil {
		return err
	}

	// UltraDNS SDK doesn't implement this; we emulate the SDK's API requests here
	var ignored interface{}
	uri := fmt.Sprintf("/zones/%s", z.Properties.Name)
	_, err = client.Do("delete", uri, ignored, ignored)

	if err != nil {
		return fmt.Errorf("delete failed: %v", err)
	}
	log.Printf("[INFO] ultradns_zone deleted: %+v", z)

	return nil
}

func resourceUltraDNSZoneImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	return setResourceAndParseId(d)
}

func populateResourceDataFromZone(z udnssdk.Zone, d *schema.ResourceData) error {
	d.Set("name", z.Properties.Name)
	d.Set("type", z.Properties.Type)
	d.Set("account", z.Properties.AccountName)
	d.Set("status", z.Properties.Status)
	d.Set("resource_record_count", z.Properties.ResourceRecordCount)
	d.Set("last_modified", z.Properties.LastModifiedDateTime)

	return nil

}
