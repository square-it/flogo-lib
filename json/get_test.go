package json

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

var jsonData = `{
    "City": [
        {
            "Array": [
                {
                    "id": "11111"
                },
                {
                    "id": "2222"
                }
            ],
            "InUS": true,
            "Name": "Sugar Land",
            "Park": {
                "Emails": null,
                "Location": "location",
                "Maps": {
                    "bb": "bb",
                    "cc": "cc",
                    "dd": "dd"
                },
                "Name": "Name"
            }
        }
    ],
    "Emails": [
        "123@123.com",
        "456@456.com"
    ],
    "Id": 1234,
    "In": "string222",
    "Maps": {
        "bb": "bb",
        "cc": "cc",
        "dd": "dd"
    },
    "State": "Taxes",
    "Streat": "311 wind st",
    "ZipCode": "77477",
    "hello world":"CHINA",
    "tag  world":"CHINA"
}`

func TestRootArray(t *testing.T) {
	value, err := GetFieldValue(jsonData, "City[0]")
	assert.Nil(t, err)
	assert.NotNil(t, value)
	logger.Debug("Value is:", value)
}

func TestRoot(t *testing.T) {
	value, err := GetFieldValue(jsonData, "City")
	assert.Nil(t, err)
	assert.NotNil(t, value)
	logger.Debug("Value is:", value)
}

func TestGetFieldWithSpaces(t *testing.T) {
	value, err := GetFieldValue(jsonData, "hello world")
	assert.Nil(t, err)
	assert.NotNil(t, value)
	logger.Debug("Value is:", value)
}

func TestGetFieldWithTag(t *testing.T) {
	value, err := GetFieldValue(jsonData, "tag  world")
	assert.Nil(t, err)
	assert.NotNil(t, value)
	logger.Info("Value is:", value)
}

func TestGetArray(t *testing.T) {
	value, err := GetFieldValue(jsonData, "Emails[0]")
	assert.Nil(t, err)
	assert.NotNil(t, value)
	logger.Debug("Value is:", value)

}

func TestMultipleLevelArray(t *testing.T) {
	value, err := GetFieldValue(jsonData, "City[0].Array[1].id")
	assert.Nil(t, err)
	assert.NotNil(t, value)
	logger.Debug("Value:", value)
}

func TestMultipleLevelObject(t *testing.T) {
	value, err := GetFieldValue(jsonData, "City[0].Park.Maps.bb")
	assert.Nil(t, err)
	assert.NotNil(t, value)
	logger.Debug("Value:", value)
}

func TestID(t *testing.T) {
	value, err := GetFieldValue(jsonData, "Id")
	assert.Nil(t, err)
	assert.NotNil(t, value)
	logger.Debug("Value:", value)
}

func TestGetFieldValue(t *testing.T) {
	account := `{
    "Account": {
        "records": [
            {
                "AccountNumber": "12356",
                "AccountSource": "Test Source",
                "Active__c": "Yes",
                "AnnualRevenue": "324556",
                "BillingCity": "Palo Alto",
                "BillingCountry": "USA",
                "BillingGeocodeAccuracy": null,
                "BillingLatitude": null,
                "BillingLongitude": null,
                "BillingPostalCode": "94207",
                "BillingState": "California",
                "BillingStreet": "3330 hillview ave",
                "CleanStatus": "Pending",
                "CustomerPriority__c": "High",
                "Description": "Sample Description for the account",
                "DunsNumber": "32653",
                "Fax": "345272",
                "Industry": "Engineering",
                "Jigsaw": "Test",
                "NaicsCode": "34583",
                "NaicsDesc": "Test Description",
                "Name": "may24_a",
                "Ownership": "Private",
                "ParentId": null,
                "Phone": "1234567890",
                "Rating": "Warm",
                "SLAExpirationDate__c": "2017-08-27",
                "SLA__c": "23453",
                "ShippingCity": "San Francisco",
                "ShippingCountry": "USA",
                "ShippingGeocodeAccuracy": null,
                "ShippingLatitude": null,
                "ShippingLongitude": null,
                "ShippingPostalCode": 45692,
                "ShippingState": "California",
                "ShippingStreet": "1234 Hillview Ave",
                "Sic": "Gold",
                "SicDesc": null,
                "Site": "www.example2.com",
                "TickerSymbol": null,
                "Tradestyle": "Regular",
                "Type": "Custumer-Direct",
                "UpsellOpportunity__c": "Yes",
                "Website": "www.example.com",
                "YearStarted": "2015"
            }
        ]
    }
}
`

	value, err := GetFieldValue(account, "Account.records[0].Name")
	logger.Infof("Value:%s", value)

	assert.Nil(t, err)
	assert.NotNil(t, value)
}

func TestConcurrentGetk(t *testing.T) {
	w := sync.WaitGroup{}
	var recovered interface{}
	//Create factory

	for r := 0; r < 100000; r++ {
		w.Add(1)
		go func(i int) {
			defer w.Done()
			defer func() {
				if r := recover(); r != nil {
					recovered = r
				}
			}()
			value, err := GetFieldValue(jsonData, "City[0].Park.Maps.bb")
			assert.Nil(t, err)
			assert.NotNil(t, value)
		}(r)

	}
	w.Wait()
	assert.Nil(t, recovered)
}
