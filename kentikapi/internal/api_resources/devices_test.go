package api_resources_test

import (
	"testing"
	"time"

	"github.com/kentik/community_sdk_golang/kentikapi/internal/api_connection"
	"github.com/kentik/community_sdk_golang/kentikapi/internal/api_resources"
	"github.com/kentik/community_sdk_golang/kentikapi/internal/utils"
	"github.com/kentik/community_sdk_golang/kentikapi/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateDeviceRouter(t *testing.T) {
	// arrange
	createResponsePayload := `
    {
        "device": {
            "id": "42",
            "company_id": "74333",
            "device_name": "testapi_router_router_full",
            "device_type": "router",
            "device_status": "V",
            "device_description": "testapi router with full config",
            "site": {
                "id": 8483,
                "site_name": null,
                "lat": null,
                "lon": null,
                "company_id": null
            },
            "plan": {
                "active": null,
                "bgp_enabled": null,
                "cdate": null,
                "company_id": null,
                "description": null,
                "deviceTypes": [],
                "devices": [],
                "edate": null,
                "fast_retention": null,
                "full_retention": null,
                "id": 11466,
                "max_bigdata_fps": null,
                "max_devices": null,
                "max_fps": null,
                "name": null,
                "metadata": null
            },
            "labels": [],
            "all_interfaces": [],
            "device_flow_type": "auto",
            "device_sample_rate": "1",
            "sending_ips": [
                "128.0.0.10"
            ],
            "device_snmp_ip": "127.0.0.1",
            "device_snmp_community": "",
            "minimize_snmp": false,
            "device_bgp_type": "device",
            "device_bgp_neighbor_ip": "127.0.0.2",
            "device_bgp_neighbor_ip6": null,
            "device_bgp_neighbor_asn": "77",
            "device_bgp_flowspec": true,
            "device_bgp_password": "******************ord",
            "use_bgp_device_id": null,
            "custom_columns": "",
            "custom_column_data": [],
            "device_chf_client_port": null,
            "device_chf_client_protocol": null,
            "device_chf_interface": null,
            "device_agent_type": null,
            "max_flow_rate": null,
            "max_big_flow_rate": null,
            "device_proxy_bgp": "",
            "device_proxy_bgp6": "",
            "created_date": "2021-01-08T08:17:07.338Z",
            "updated_date": "2021-01-08T08:17:07.338Z",
            "device_snmp_v3_conf": {
                "UserName": "John",
                "AuthenticationProtocol": "MD5",
                "AuthenticationPassphrase": "Auth_Pass",
                "PrivacyProtocol": "DES",
                "PrivacyPassphrase": "******ass"
            },
            "bgpPeerIP4": "208.76.14.223",
            "bgpPeerIP6": "2620:129:1:2::1",
            "snmp_last_updated": null,
            "device_subtype": "router"
        }
    }`
	transport := &api_connection.StubTransport{ResponseBody: createResponsePayload}
	devicesAPI := api_resources.NewDevicesAPI(transport)

	// act
	snmpv3conf := models.NewSNMPv3Conf("John")
	snmpv3conf = snmpv3conf.WithAuthentication(models.AuthenticationProtocolMD5, "Auth_Pass")
	snmpv3conf = snmpv3conf.WithPrivacy(models.PrivacyProtocolDES, "Priv_Pass")
	router := models.NewDeviceRouter(
		"testapi_router_router_full",
		models.DeviceSubtypeRouter,
		1,
		models.ID(11466),
		[]string{"128.0.0.10"},
		false,
	).WithBGPTypeDevice("77")
	models.SetOptional(&router.DeviceDescription, "testapi router with full config")
	models.SetOptional(&router.DeviceSNMNPIP, "127.0.0.1")
	models.SetOptional(&router.DeviceSNMPv3Conf, *snmpv3conf)
	models.SetOptional(&router.DeviceBGPNeighborIP, "127.0.0.2")
	models.SetOptional(&router.DeviceBGPPassword, "bgp-optional-password")
	models.SetOptional(&router.SiteID, 8483)
	models.SetOptional(&router.DeviceBGPFlowSpec, true)
	device, err := devicesAPI.Create(nil, *router)

	// assert request properly formed
	assert := assert.New(t)
	require := require.New(t)

	require.NoError(err)
	payload := utils.NewJSONPayloadInspector(t, transport.RequestBody)
	require.NotNil(payload.Get("device"))
	assert.Equal("testapi_router_router_full", payload.String("device/device_name"))
	assert.Equal("router", payload.String("device/device_type"))
	assert.Equal("router", payload.String("device/device_subtype"))
	assert.Equal(1, payload.Count("device/sending_ips"))
	assert.Equal("128.0.0.10", payload.String("device/sending_ips/*[1]")) // xpath [1] means first element
	assert.Equal("1", payload.String("device/device_sample_rate"))
	assert.Equal("testapi router with full config", payload.String("device/device_description"))
	assert.Equal("127.0.0.1", payload.String("device/device_snmp_ip"))
	assert.Equal(11466, payload.Int("device/plan_id"))
	assert.Equal(8483, payload.Int("device/site_id"))
	assert.False(payload.Bool("device/minimize_snmp"))
	assert.NotNil(payload.Get("device/device_snmp_v3_conf"))
	assert.Equal("John", payload.String("device/device_snmp_v3_conf/UserName"))
	assert.Equal("MD5", payload.String("device/device_snmp_v3_conf/AuthenticationProtocol"))
	assert.Equal("Auth_Pass", payload.String("device/device_snmp_v3_conf/AuthenticationPassphrase"))
	assert.Equal("DES", payload.String("device/device_snmp_v3_conf/PrivacyProtocol"))
	assert.Equal("Priv_Pass", payload.String("device/device_snmp_v3_conf/PrivacyPassphrase"))
	assert.Equal("device", payload.String("device/device_bgp_type"))
	assert.Equal("77", payload.String("device/device_bgp_neighbor_asn"))
	assert.Equal("127.0.0.2", payload.String("device/device_bgp_neighbor_ip"))
	assert.Equal("bgp-optional-password", payload.String("device/device_bgp_password"))
	assert.True(payload.Bool("device/device_bgp_flowspec"))

	// and response properly parsed
	assert.Equal(models.ID(42), device.ID)
	assert.Equal(models.ID(74333), device.CompanyID)
	assert.Equal("testapi_router_router_full", device.DeviceName)
	assert.Equal(models.DeviceTypeRouter, device.DeviceType)
	assert.Equal("testapi router with full config", *device.DeviceDescription)
	require.NotNil(device.Site)
	assert.Equal(models.ID(8483), *device.Site.ID)
	assert.Nil(device.Site.SiteName)
	assert.Nil(device.Site.Latitude)
	assert.Nil(device.Site.Longitude)
	assert.Nil(device.Site.CompanyID)
	require.NotNil(device.Plan)
	assert.Nil(device.Plan.Active)
	assert.Nil(device.Plan.BGPEnabled)
	assert.Nil(device.Plan.CreatedDate)
	assert.Nil(device.Plan.CompanyID)
	assert.Nil(device.Plan.Description)
	assert.Equal(0, len(device.Plan.DeviceTypes))
	assert.Equal(0, len(device.Plan.Devices))
	assert.Nil(device.Plan.UpdatedDate)
	assert.Nil(device.Plan.FastRetention)
	assert.Nil(device.Plan.FullRetention)
	assert.Equal(models.ID(11466), *device.Plan.ID)
	assert.Nil(device.Plan.MaxBigdataFPS)
	assert.Nil(device.Plan.MaxDevices)
	assert.Nil(device.Plan.MaxFPS)
	assert.Nil(device.Plan.Name)
	assert.Equal(0, len(device.Labels))
	assert.Equal(0, len(device.AllInterfaces))
	assert.Equal("auto", *device.DeviceFlowType)
	assert.Equal(1, device.DeviceSampleRate)
	assert.Equal(1, len(device.SendingIPS))
	assert.Equal("128.0.0.10", device.SendingIPS[0])
	assert.Equal("127.0.0.1", *device.DeviceSNMNPIP)
	assert.Equal("", *device.DeviceSNMPCommunity)
	assert.False(*device.MinimizeSNMP)
	assert.Equal(models.DeviceBGPTypeDevice, *device.DeviceBGPType)
	assert.Equal("127.0.0.2", *device.DeviceBGPNeighborIP)
	assert.Nil(device.DeviceBGPNeighborIPv6)
	assert.Equal("77", *device.DeviceBGPNeighborASN)
	assert.True(*device.DeviceBGPFlowSpec)
	assert.Equal("******************ord", *device.DeviceBGPPassword)
	assert.Nil(device.UseBGPDeviceID)
	assert.Equal(time.Date(2021, 1, 8, 8, 17, 7, 338*1000000, time.UTC), device.CreatedDate)
	assert.Equal(time.Date(2021, 1, 8, 8, 17, 7, 338*1000000, time.UTC), device.UpdatedDate)
	require.NotNil(device.DeviceSNMPv3Conf)
	assert.Equal("John", device.DeviceSNMPv3Conf.UserName)
	assert.Equal(models.AuthenticationProtocolMD5, *device.DeviceSNMPv3Conf.AuthenticationProtocol)
	assert.Equal("Auth_Pass", *device.DeviceSNMPv3Conf.AuthenticationPassphrase)
	assert.Equal(models.PrivacyProtocolDES, *device.DeviceSNMPv3Conf.PrivacyProtocol)
	assert.Equal("******ass", *device.DeviceSNMPv3Conf.PrivacyPassphrase)
	assert.Equal("208.76.14.223", *device.BGPPeerIP4)
	assert.Equal("2620:129:1:2::1", *device.BGPPeerIP6)
	assert.Nil(device.SNMPLastUpdated)
	assert.Equal(models.DeviceSubtypeRouter, device.DeviceSubType)
}

func TestCreateDeviceDNS(t *testing.T) {
	// arrange
	createResponsePayload := `
    {
        "device": {
            "id": "43",
            "company_id": "74333",
            "device_name": "testapi_dns_aws_subnet_bgp_other_device",
            "device_type": "host-nprobe-dns-www",
            "device_status": "V",
            "device_description": "testapi dns with minimal config",
            "site": {
                "id": 8483,
                "site_name": null,
                "lat": null,
                "lon": null,
                "company_id": null
            },
            "plan": {
                "active": null,
                "bgp_enabled": null,
                "cdate": null,
                "company_id": null,
                "description": null,
                "deviceTypes": [],
                "devices": [],
                "edate": null,
                "fast_retention": null,
                "full_retention": null,
                "id": 11466,
                "max_bigdata_fps": null,
                "max_devices": null,
                "max_fps": null,
                "name": null,
                "metadata": null
            },
            "labels": [],
            "all_interfaces": [],
            "device_flow_type": "auto",
            "device_sample_rate": "1",
            "sending_ips": [],
            "device_snmp_ip": null,
            "device_snmp_community": "",
            "minimize_snmp": false,
            "device_bgp_type": "other_device",
            "use_bgp_device_id": 42,
            "device_bgp_flowspec": true,
            "custom_columns": "",
            "custom_column_data": [],
            "device_chf_client_port": null,
            "device_chf_client_protocol": null,
            "device_chf_interface": null,
            "device_agent_type": null,
            "max_flow_rate": null,
            "max_big_flow_rate": null,
            "device_proxy_bgp": "",
            "device_proxy_bgp6": "",
            "created_date": "2021-01-08T11:10:33.465Z",
            "updated_date": "2021-01-08T11:10:33.465Z",
            "device_snmp_v3_conf": null,
            "cdn_attr": "Y",
            "bgpPeerIP4": "208.76.14.223",
            "bgpPeerIP6": "2620:129:1:2::1",
            "snmp_last_updated": null,
            "device_subtype": "aws_subnet"
        }
    }`
	transport := &api_connection.StubTransport{ResponseBody: createResponsePayload}
	devicesAPI := api_resources.NewDevicesAPI(transport)

	// act
	dns := models.NewDeviceDNS(
		"testapi_dns-aws_subnet_bgp_other_device",
		models.DeviceSubtypeAwsSubnet,
		1,
		models.ID(11466),
		models.CDNAttributeYes,
	).WithBGPTypeOtherDevice(models.ID(42))
	models.SetOptional(&dns.DeviceDescription, "testapi dns with minimal config")
	models.SetOptional(&dns.SiteID, 8483)
	models.SetOptional(&dns.DeviceBGPFlowSpec, true)
	device, err := devicesAPI.Create(nil, *dns)

	// assert request properly formed
	assert := assert.New(t)
	require := require.New(t)

	require.NoError(err)
	payload := utils.NewJSONPayloadInspector(t, transport.RequestBody)
	require.NotNil(payload.Get("device"))
	assert.Equal("testapi_dns-aws_subnet_bgp_other_device", payload.String("device/device_name"))
	assert.Equal("host-nprobe-dns-www", payload.String("device/device_type"))
	assert.Equal("aws_subnet", payload.String("device/device_subtype"))
	assert.Equal("Y", payload.String("device/cdn_attr"))
	assert.Equal("1", payload.String("device/device_sample_rate"))
	assert.Equal("testapi dns with minimal config", payload.String("device/device_description"))
	assert.Equal(11466, payload.Int("device/plan_id"))
	assert.Equal(8483, payload.Int("device/site_id"))
	assert.Equal("other_device", payload.String("device/device_bgp_type"))
	assert.Equal(42, payload.Int("device/use_bgp_device_id"))
	assert.True(payload.Bool("device/device_bgp_flowspec"))
	assert.False(payload.Exists("device/sending_ips")) // empty array should not be included in payload

	// # and response properly parsed
	assert.Equal(models.ID(43), device.ID)
	assert.Equal("testapi_dns_aws_subnet_bgp_other_device", device.DeviceName)
	assert.Equal(models.DeviceTypeHostNProbeDNSWWW, device.DeviceType)
	assert.Equal("testapi dns with minimal config", *device.DeviceDescription)
	assert.NotNil(device.Site)
	assert.Equal(models.ID(8483), *device.Site.ID)
	assert.Nil(device.Site.SiteName)
	assert.Nil(device.Site.Latitude)
	assert.Nil(device.Site.Longitude)
	assert.Nil(device.Site.CompanyID)
	assert.Nil(device.Plan.Active)
	assert.Nil(device.Plan.BGPEnabled)
	assert.Nil(device.Plan.CreatedDate)
	assert.Nil(device.Plan.CompanyID)
	assert.Nil(device.Plan.Description)
	assert.Equal(0, len(device.Plan.DeviceTypes))
	assert.Equal(0, len(device.Plan.Devices))
	assert.Nil(device.Plan.UpdatedDate)
	assert.Nil(device.Plan.FastRetention)
	assert.Nil(device.Plan.FullRetention)
	assert.Equal(models.ID(11466), *device.Plan.ID)
	assert.Nil(device.Plan.MaxBigdataFPS)
	assert.Nil(device.Plan.MaxDevices)
	assert.Nil(device.Plan.MaxFPS)
	assert.Nil(device.Plan.Name)
	assert.Equal(0, len(device.Labels))
	assert.Equal(0, len(device.AllInterfaces))
	assert.Equal("auto", *device.DeviceFlowType)
	assert.Equal(1, device.DeviceSampleRate)
	assert.Equal(0, len(device.SendingIPS))
	assert.Nil(device.DeviceSNMNPIP)
	assert.Equal("", *device.DeviceSNMPCommunity)
	assert.False(*device.MinimizeSNMP)
	assert.Equal(models.DeviceBGPTypeOtherDevice, *device.DeviceBGPType)
	assert.True(*device.DeviceBGPFlowSpec)
	assert.Equal(models.ID(42), *device.UseBGPDeviceID)
	assert.Equal(time.Date(2021, 1, 8, 11, 10, 33, 465*1000000, time.UTC), device.CreatedDate)
	assert.Equal(time.Date(2021, 1, 8, 11, 10, 33, 465*1000000, time.UTC), device.UpdatedDate)
	assert.Nil(device.DeviceSNMPv3Conf)
	assert.Equal("208.76.14.223", *device.BGPPeerIP4)
	assert.Equal("2620:129:1:2::1", *device.BGPPeerIP6)
	assert.Nil(device.SNMPLastUpdated)
	assert.Equal(models.DeviceSubtypeAwsSubnet, device.DeviceSubType)
}

func TestUpdatetDeviceRouter(t *testing.T) {
	// arrange
	updateResponsePayload := `
    {
        "device": {
            "id": "42",
            "company_id": "74333",
            "device_name": "testapi_router_paloalto_minimal",
            "device_type": "router",
            "device_status": "V",
            "device_description": "updated description",
            "site": {
                "id": 8483,
                "site_name": null,
                "lat": null,
                "lon": null,
                "company_id": null
            },
            "plan": {
                "active": null,
                "bgp_enabled": null,
                "cdate": null,
                "company_id": null,
                "description": null,
                "deviceTypes": [],
                "devices": [],
                "edate": null,
                "fast_retention": null,
                "full_retention": null,
                "id": 11466,
                "max_bigdata_fps": null,
                "max_devices": null,
                "max_fps": null,
                "name": null,
                "metadata": null
            },
            "labels": [],
            "all_interfaces": [],
            "device_flow_type": "auto",
            "device_sample_rate": "10",
            "sending_ips": [
                "128.0.0.10",
                "128.0.0.11"
            ],
            "device_snmp_ip": "127.0.0.10",
            "device_snmp_community": "",
            "minimize_snmp": true,
            "device_bgp_type": "device",
            "device_bgp_neighbor_ip": null,
            "device_bgp_neighbor_ip6": "2001:db8:85a3:8d3:1319:8a2e:370:7348",
            "device_bgp_neighbor_asn": "77",
            "device_bgp_flowspec": true,
            "device_bgp_password": "******************ord",
            "use_bgp_device_id": null,
            "custom_columns": "",
            "custom_column_data": [],
            "device_chf_client_port": null,
            "device_chf_client_protocol": null,
            "device_chf_interface": null,
            "device_agent_type": null,
            "max_flow_rate": null,
            "max_big_flow_rate": null,
            "device_proxy_bgp": "",
            "device_proxy_bgp6": "",
            "created_date": "2021-01-08T13:02:45.733Z",
            "updated_date": "2021-01-08T13:11:57.795Z",
            "device_snmp_v3_conf": {
                "UserName": "John",
                "AuthenticationProtocol": "SHA",
                "AuthenticationPassphrase": "Auth_Pass",
                "PrivacyProtocol": "AES",
                "PrivacyPassphrase": "******ass"
            },
            "bgpPeerIP4": "208.76.14.223",
            "bgpPeerIP6": "2620:129:1:2::1",
            "snmp_last_updated": null,
            "device_subtype": "paloalto"
        }
    }`
	transport := &api_connection.StubTransport{ResponseBody: updateResponsePayload}
	devicesAPI := api_resources.NewDevicesAPI(transport)

	// act
	snmpv3conf := models.NewSNMPv3Conf("John")
	snmpv3conf = snmpv3conf.WithAuthentication(models.AuthenticationProtocolSHA, "Auth_Pass")
	snmpv3conf = snmpv3conf.WithPrivacy(models.PrivacyProtocolAES, "Priv_Pass")
	deviceID := models.ID(42)

	router := models.Device{
		ID:               deviceID,
		SendingIPS:       []string{"128.0.0.10", "128.0.0.11"},
		DeviceSampleRate: 10,
		DeviceSNMPv3Conf: snmpv3conf,
	}
	models.SetOptional(&router.DeviceDescription, "updated description")
	models.SetOptional(&router.DeviceSNMNPIP, "127.0.0.10")
	models.SetOptional(&router.MinimizeSNMP, true)
	models.SetOptional(&router.PlanID, models.ID(11466))
	models.SetOptional(&router.SiteID, models.ID(8483))
	models.SetOptional(&router.DeviceBGPType, models.DeviceBGPTypeDevice)
	models.SetOptional(&router.DeviceBGPNeighborASN, "77")
	models.SetOptional(&router.DeviceBGPNeighborIPv6, "2001:db8:85a3:8d3:1319:8a2e:370:7348")
	models.SetOptional(&router.DeviceBGPPassword, "bgp-optional-password")
	models.SetOptional(&router.DeviceBGPFlowSpec, true)
	device, err := devicesAPI.Update(nil, router)

	// assert request properly formed
	assert := assert.New(t)
	require := require.New(t)

	require.NoError(err)
	payload := utils.NewJSONPayloadInspector(t, transport.RequestBody)
	require.NotNil(payload.Get("device"))
	assert.Equal(2, payload.Count("device/sending_ips/*"))
	assert.Equal("128.0.0.10", payload.String("device/sending_ips/*[1]")) // xpath [1] means first element
	assert.Equal("128.0.0.11", payload.String("device/sending_ips/*[2]")) // xpath [2] means second element
	assert.Equal("10", payload.String("device/device_sample_rate"))
	assert.Equal("updated description", payload.String("device/device_description"))
	assert.Equal("127.0.0.10", payload.String("device/device_snmp_ip"))
	assert.Equal(11466, payload.Int("device/plan_id"))
	assert.Equal(8483, payload.Int("device/site_id"))
	assert.True(payload.Bool("device/minimize_snmp"))
	assert.NotNil(payload.Get("device/device_snmp_v3_conf"))
	assert.Equal("John", payload.String("device/device_snmp_v3_conf/UserName"))
	assert.Equal("SHA", payload.String("device/device_snmp_v3_conf/AuthenticationProtocol"))
	assert.Equal("Auth_Pass", payload.String("device/device_snmp_v3_conf/AuthenticationPassphrase"))
	assert.Equal("AES", payload.String("device/device_snmp_v3_conf/PrivacyProtocol"))
	assert.Equal("Priv_Pass", payload.String("device/device_snmp_v3_conf/PrivacyPassphrase"))
	assert.Equal("device", payload.String("device/device_bgp_type"))
	assert.Equal("77", payload.String("device/device_bgp_neighbor_asn"))
	assert.Equal("2001:db8:85a3:8d3:1319:8a2e:370:7348", payload.String("device/device_bgp_neighbor_ip6"))
	assert.Equal("bgp-optional-password", payload.String("device/device_bgp_password"))
	assert.True(payload.Bool("device/device_bgp_flowspec"))

	// # and response properly parsed
	assert.Equal(models.ID(42), device.ID)
	assert.Equal(models.ID(74333), device.CompanyID)
	assert.Equal("testapi_router_paloalto_minimal", device.DeviceName)
	assert.Equal(models.DeviceTypeRouter, device.DeviceType)
	assert.Equal("updated description", *device.DeviceDescription)
	require.NotNil(device.Site)
	assert.Equal(models.ID(8483), *device.Site.ID)
	assert.Nil(device.Site.SiteName)
	assert.Nil(device.Site.Latitude)
	assert.Nil(device.Site.Longitude)
	assert.Nil(device.Site.CompanyID)
	require.NotNil(device.Plan)
	assert.Nil(device.Plan.Active)
	assert.Nil(device.Plan.BGPEnabled)
	assert.Nil(device.Plan.CreatedDate)
	assert.Nil(device.Plan.CompanyID)
	assert.Nil(device.Plan.Description)
	assert.Equal(0, len(device.Plan.DeviceTypes))
	assert.Equal(0, len(device.Plan.Devices))
	assert.Nil(device.Plan.UpdatedDate)
	assert.Nil(device.Plan.FastRetention)
	assert.Nil(device.Plan.FullRetention)
	assert.Equal(models.ID(11466), *device.Plan.ID)
	assert.Nil(device.Plan.MaxBigdataFPS)
	assert.Nil(device.Plan.MaxDevices)
	assert.Nil(device.Plan.MaxFPS)
	assert.Nil(device.Plan.Name)
	assert.Equal(0, len(device.Labels))
	assert.Equal(0, len(device.AllInterfaces))
	assert.Equal("auto", *device.DeviceFlowType)
	assert.Equal(10, *&device.DeviceSampleRate)
	assert.Equal(2, len(device.SendingIPS))
	assert.Equal("128.0.0.10", device.SendingIPS[0])
	assert.Equal("128.0.0.11", device.SendingIPS[1])
	assert.Equal("127.0.0.10", *device.DeviceSNMNPIP)
	assert.Equal("", *device.DeviceSNMPCommunity)
	assert.True(*device.MinimizeSNMP)
	assert.Equal(models.DeviceBGPTypeDevice, *device.DeviceBGPType)
	assert.Equal("2001:db8:85a3:8d3:1319:8a2e:370:7348", *device.DeviceBGPNeighborIPv6)
	assert.Nil(device.DeviceBGPNeighborIP)
	assert.Equal("77", *device.DeviceBGPNeighborASN)
	assert.True(*device.DeviceBGPFlowSpec)
	assert.Equal("******************ord", *device.DeviceBGPPassword)
	assert.Nil(device.UseBGPDeviceID)
	assert.Equal(time.Date(2021, 1, 8, 13, 2, 45, 733*1000000, time.UTC), device.CreatedDate)
	assert.Equal(time.Date(2021, 1, 8, 13, 11, 57, 795*1000000, time.UTC), device.UpdatedDate)
	require.NotNil(device.DeviceSNMPv3Conf)
	assert.Equal("John", device.DeviceSNMPv3Conf.UserName)
	assert.Equal(models.AuthenticationProtocolSHA, *device.DeviceSNMPv3Conf.AuthenticationProtocol)
	assert.Equal("Auth_Pass", *device.DeviceSNMPv3Conf.AuthenticationPassphrase)
	assert.Equal(models.PrivacyProtocolAES, *device.DeviceSNMPv3Conf.PrivacyProtocol)
	assert.Equal("******ass", *device.DeviceSNMPv3Conf.PrivacyPassphrase)
	assert.Equal("208.76.14.223", *device.BGPPeerIP4)
	assert.Equal("2620:129:1:2::1", *device.BGPPeerIP6)
	assert.Nil(device.SNMPLastUpdated)
	assert.Equal(models.DeviceSubtypePaloalto, device.DeviceSubType)
}

func TestGetDeviceRouter(t *testing.T) {
	// arrange
	getResponsePayload := `
    {
        "device": {
            "id": "42",
            "company_id": "74333",
            "device_name": "testapi_router_full_1",
            "device_type": "router",
            "device_status": "V",
            "device_description": "testapi router with full config",
            "site": {
                "id": 8483,
                "site_name": "marina gdańsk",
                "lat": 54.348972,
                "lon": 18.659791,
                "company_id": 74333
            },
            "plan": {
                "active": true,
                "bgp_enabled": true,
                "cdate": "2020-09-03T08:41:57.489Z",
                "company_id": 74333,
                "description": "Your Free Trial includes 6 devices (...)",
                "deviceTypes": [],
                "devices": [],
                "edate": "2020-09-03T08:41:57.489Z",
                "fast_retention": 30,
                "full_retention": 30,
                "id": 11466,
                "max_bigdata_fps": 30,
                "max_devices": 6,
                "max_fps": 1000,
                "name": "Free Trial Plan",
                "metadata": {}
            },
            "labels": [
                        {
                            "id": 2590,
                            "name": "AWS: terraform-demo-aws",
                            "description": null,
                            "edate": "2020-10-05T15:28:00.276Z",
                            "cdate": "2020-10-05T15:28:00.276Z",
                            "user_id": "133210",
                            "company_id": "74333",
                            "color": "#5340A5",
                            "order": null,
                            "_pivot_device_id": "77715",
                            "_pivot_label_id": "2590"
                        },
                        {
                            "id": 2751,
                            "name": "GCP: traffic-generator-gcp",
                            "description": null,
                            "edate": "2020-11-20T12:54:49.575Z",
                            "cdate": "2020-11-20T12:54:49.575Z",
                            "user_id": "136885",
                            "company_id": "74333",
                            "color": "#5289D9",
                            "order": null,
                            "_pivot_device_id": "77373",
                            "_pivot_label_id": "2751"
                        }
                    ],
            "all_interfaces": [
                {
                    "interface_description": "testapi-interface-1",
                    "initial_snmp_speed": null,
                    "device_id": "42",
                    "snmp_speed": "75"
                },
                {
                    "interface_description": "testapi-interface-2",
                    "initial_snmp_speed": "7",
                    "device_id": "42",
                    "snmp_speed": "7"
                }
            ],
            "device_flow_type": "auto",
            "device_sample_rate": "1001",
            "sending_ips": [
                "128.0.0.11",
                "128.0.0.12"
            ],
            "device_snmp_ip": "129.0.0.1",
            "device_snmp_community": "",
            "minimize_snmp": false,
            "device_bgp_type": "device",
            "device_bgp_neighbor_ip": "127.0.0.1",
            "device_bgp_neighbor_ip6": null,
            "device_bgp_neighbor_asn": "11",
            "device_bgp_flowspec": true,
            "device_bgp_password": "*********ass",
            "use_bgp_device_id": null,
            "custom_columns": "",
            "custom_column_data": [],
            "device_chf_client_port": null,
            "device_chf_client_protocol": null,
            "device_chf_interface": null,
            "device_agent_type": null,
            "max_flow_rate": 1000,
            "max_big_flow_rate": 30,
            "device_proxy_bgp": "",
            "device_proxy_bgp6": "",
            "created_date": "2020-12-17T08:24:45.074Z",
            "updated_date": "2020-12-17T08:24:45.074Z",
            "device_snmp_v3_conf": {
                "UserName": "John",
                "AuthenticationProtocol": "MD5",
                "AuthenticationPassphrase": "john_md5_pass",
                "PrivacyProtocol": "DES",
                "PrivacyPassphrase": "**********ass"
            },
            "bgpPeerIP4": "208.76.14.223",
            "bgpPeerIP6": "2620:129:1:2::1",
            "snmp_last_updated": null,
            "device_subtype": "router"
        }
    }`
	transport := &api_connection.StubTransport{ResponseBody: getResponsePayload}
	devicesAPI := api_resources.NewDevicesAPI(transport)
	deviceID := models.ID(42)

	// act
	device, err := devicesAPI.Get(nil, deviceID)

	// assert request properly formed
	assert := assert.New(t)
	require := require.New(t)

	require.NoError(err)
	assert.Zero(transport.RequestBody)

	// and response properly parsed
	assert.Equal(models.ID(42), device.ID)
	assert.Equal(models.ID(74333), device.CompanyID)
	assert.Equal("testapi_router_full_1", device.DeviceName)
	assert.Equal(models.DeviceTypeRouter, device.DeviceType)
	assert.Equal("testapi router with full config", *device.DeviceDescription)
	require.NotNil(device.Site)
	assert.Equal(models.ID(8483), *device.Site.ID)
	assert.Equal("marina gdańsk", *device.Site.SiteName)
	assert.Equal(54.348972, *device.Site.Latitude)
	assert.Equal(18.659791, *device.Site.Longitude)
	assert.Equal(models.ID(74333), *device.Site.CompanyID)
	require.NotNil(device.Plan)
	assert.True(*device.Plan.Active)
	assert.True(*device.Plan.BGPEnabled)
	assert.Equal(time.Date(2020, 9, 3, 8, 41, 57, 489*1000000, time.UTC), *device.Plan.CreatedDate)
	assert.Equal(models.ID(74333), *device.Plan.CompanyID)
	assert.Equal("Your Free Trial includes 6 devices (...)", *device.Plan.Description)
	assert.Equal(0, len(device.Plan.DeviceTypes))
	assert.Equal(0, len(device.Plan.Devices))
	assert.Equal(time.Date(2020, 9, 3, 8, 41, 57, 489*1000000, time.UTC), *device.Plan.UpdatedDate)
	assert.Equal(30, *device.Plan.FastRetention)
	assert.Equal(30, *device.Plan.FullRetention)
	assert.Equal(models.ID(11466), *device.Plan.ID)
	assert.Equal(30, *device.Plan.MaxBigdataFPS)
	assert.Equal(6, *device.Plan.MaxDevices)
	assert.Equal(1000, *device.Plan.MaxFPS)
	assert.Equal("Free Trial Plan", *device.Plan.Name)
	assert.Equal(2, len(device.Labels))
	assert.Equal(models.ID(2590), device.Labels[0].ID)
	assert.Equal("AWS: terraform-demo-aws", device.Labels[0].Name)
	assert.Equal(time.Date(2020, 10, 5, 15, 28, 00, 276*1000000, time.UTC), device.Labels[0].UpdatedDate)
	assert.Equal(time.Date(2020, 10, 5, 15, 28, 00, 276*1000000, time.UTC), device.Labels[0].CreatedDate)
	assert.Equal(models.ID(133210), *device.Labels[0].UserID)
	assert.Equal(models.ID(74333), device.Labels[0].CompanyID)
	assert.Equal("#5340A5", device.Labels[0].Color)
	assert.Equal(models.ID(2751), device.Labels[1].ID)
	assert.Equal("GCP: traffic-generator-gcp", device.Labels[1].Name)
	assert.Equal(time.Date(2020, 11, 20, 12, 54, 49, 575*1000000, time.UTC), device.Labels[1].UpdatedDate)
	assert.Equal(time.Date(2020, 11, 20, 12, 54, 49, 575*1000000, time.UTC), device.Labels[1].CreatedDate)
	assert.Equal(models.ID(136885), *device.Labels[1].UserID)
	assert.Equal(models.ID(74333), device.Labels[1].CompanyID)
	assert.Equal("#5289D9", device.Labels[1].Color)
	assert.Equal(2, len(device.AllInterfaces))
	assert.Equal("testapi-interface-1", device.AllInterfaces[0].InterfaceDescription)
	assert.Nil(device.AllInterfaces[0].InitialSNMPSpeed)
	assert.Equal(models.ID(42), device.AllInterfaces[0].DeviceID)
	assert.Equal(75.0, device.AllInterfaces[0].SNMPSpeed)
	assert.Equal("testapi-interface-2", device.AllInterfaces[1].InterfaceDescription)
	assert.Equal(7.0, *device.AllInterfaces[1].InitialSNMPSpeed)
	assert.Equal(models.ID(42), device.AllInterfaces[1].DeviceID)
	assert.Equal(7.0, device.AllInterfaces[1].SNMPSpeed)
	assert.Equal("auto", *device.DeviceFlowType)
	assert.Equal(1001, *&device.DeviceSampleRate)
	assert.Equal(2, len(device.SendingIPS))
	assert.Equal("128.0.0.11", device.SendingIPS[0])
	assert.Equal("128.0.0.12", device.SendingIPS[1])
	assert.Equal("129.0.0.1", *device.DeviceSNMNPIP)
	assert.Equal("", *device.DeviceSNMPCommunity)
	assert.False(*device.MinimizeSNMP)
	assert.Equal(models.DeviceBGPTypeDevice, *device.DeviceBGPType)
	assert.Equal("127.0.0.1", *device.DeviceBGPNeighborIP)
	assert.Nil(device.DeviceBGPNeighborIPv6)
	assert.Equal("11", *device.DeviceBGPNeighborASN)
	assert.True(*device.DeviceBGPFlowSpec)
	assert.Equal("*********ass", *device.DeviceBGPPassword)
	assert.Nil(device.UseBGPDeviceID)
	assert.Equal(time.Date(2020, 12, 17, 8, 24, 45, 74*1000000, time.UTC), device.CreatedDate)
	assert.Equal(time.Date(2020, 12, 17, 8, 24, 45, 74*1000000, time.UTC), device.UpdatedDate)
	require.NotNil(device.DeviceSNMPv3Conf)
	assert.Equal("John", device.DeviceSNMPv3Conf.UserName)
	assert.Equal(models.AuthenticationProtocolMD5, *device.DeviceSNMPv3Conf.AuthenticationProtocol)
	assert.Equal("john_md5_pass", *device.DeviceSNMPv3Conf.AuthenticationPassphrase)
	assert.Equal(models.PrivacyProtocolDES, *device.DeviceSNMPv3Conf.PrivacyProtocol)
	assert.Equal("**********ass", *device.DeviceSNMPv3Conf.PrivacyPassphrase)
	assert.Equal("208.76.14.223", *device.BGPPeerIP4)
	assert.Equal("2620:129:1:2::1", *device.BGPPeerIP6)
	assert.Nil(device.SNMPLastUpdated)
	assert.Equal(models.DeviceSubtypeRouter, device.DeviceSubType)
}

func TestGetDeviceDNS(t *testing.T) {
	// arrange
	getResponsePayload := `
    {
        "device": {
            "id": "43",
            "company_id": "74333",
            "device_name": "testapi_dns_minimal_1",
            "device_type": "host-nprobe-dns-www",
            "device_status": "V",
            "device_description": "testapi dns with minimal config",
            "site": {},
            "plan": {
                "active": true,
                "bgp_enabled": true,
                "cdate": "2020-09-03T08:41:57.489Z",
                "company_id": 74333,
                "description": "Your Free Trial includes 6 devices (...)",
                "deviceTypes": [],
                "devices": [],
                "edate": "2020-09-03T08:41:57.489Z",
                "fast_retention": 30,
                "full_retention": 30,
                "id": 11466,
                "max_bigdata_fps": 30,
                "max_devices": 6,
                "max_fps": 1000,
                "name": "Free Trial Plan",
                "metadata": {}
            },
            "labels": [],
            "all_interfaces": [],
            "device_flow_type": "auto",
            "device_sample_rate": "1",
            "sending_ips": [],
            "device_snmp_ip": null,
            "device_snmp_community": "",
            "minimize_snmp": false,
            "device_bgp_type": "none",
            "device_bgp_neighbor_ip": null,
            "device_bgp_neighbor_ip6": null,
            "device_bgp_neighbor_asn": null,
            "device_bgp_flowspec": false,
            "device_bgp_password": null,
            "use_bgp_device_id": null,
            "custom_columns": "",
            "custom_column_data": [],
            "device_chf_client_port": null,
            "device_chf_client_protocol": null,
            "device_chf_interface": null,
            "device_agent_type": null,
            "max_flow_rate": 1000,
            "max_big_flow_rate": 30,
            "device_proxy_bgp": "",
            "device_proxy_bgp6": "",
            "created_date": "2020-12-17T12:53:01.025Z",
            "updated_date": "2020-12-17T12:53:01.025Z",
            "device_snmp_v3_conf": null,
            "cdn_attr": "Y",
            "snmp_last_updated": null,
            "device_subtype": "aws_subnet"
        }
    }`
	transport := &api_connection.StubTransport{ResponseBody: getResponsePayload}
	devicesAPI := api_resources.NewDevicesAPI(transport)
	deviceID := models.ID(43)

	// act
	device, err := devicesAPI.Get(nil, deviceID)

	// assert request properly formed
	assert := assert.New(t)
	require := require.New(t)

	require.NoError(err)
	assert.Zero(transport.RequestBody)

	// and response properly parsed
	assert.Equal(models.ID(43), device.ID)
	assert.Equal(models.ID(74333), device.CompanyID)
	assert.Equal("testapi_dns_minimal_1", device.DeviceName)
	assert.Equal(models.DeviceTypeHostNProbeDNSWWW, device.DeviceType)
	assert.Equal("testapi dns with minimal config", *device.DeviceDescription)
	assert.Zero(*device.Site) // empty site in response body
	assert.True(*device.Plan.Active)
	assert.True(*device.Plan.BGPEnabled)
	assert.Equal(time.Date(2020, 9, 3, 8, 41, 57, 489*1000000, time.UTC), *device.Plan.CreatedDate)
	assert.Equal(models.ID(74333), *device.Plan.CompanyID)
	assert.Equal("Your Free Trial includes 6 devices (...)", *device.Plan.Description)
	assert.Equal(0, len(device.Plan.DeviceTypes))
	assert.Equal(0, len(device.Plan.Devices))
	assert.Equal(time.Date(2020, 9, 3, 8, 41, 57, 489*1000000, time.UTC), *device.Plan.UpdatedDate)
	assert.Equal(30, *device.Plan.FastRetention)
	assert.Equal(30, *device.Plan.FullRetention)
	assert.Equal(models.ID(11466), *device.Plan.ID)
	assert.Equal(30, *device.Plan.MaxBigdataFPS)
	assert.Equal(6, *device.Plan.MaxDevices)
	assert.Equal(1000, *device.Plan.MaxFPS)
	assert.Equal("Free Trial Plan", *device.Plan.Name)
	assert.Equal(0, len(device.Labels))
	assert.Equal(0, len(device.AllInterfaces))
	assert.Equal("auto", *device.DeviceFlowType)
	assert.Equal(1, device.DeviceSampleRate)
	assert.Equal(0, len(device.SendingIPS))
	assert.Nil(device.DeviceSNMNPIP)
	assert.Equal("", *device.DeviceSNMPCommunity)
	assert.False(*device.MinimizeSNMP)
	assert.Equal(models.DeviceBGPTypeNone, *device.DeviceBGPType)
	assert.Nil(device.DeviceBGPNeighborIP)
	assert.Nil(device.DeviceBGPNeighborIPv6)
	assert.Nil(device.DeviceBGPNeighborASN)
	assert.False(*device.DeviceBGPFlowSpec)
	assert.Nil(device.UseBGPDeviceID)
	assert.Equal(time.Date(2020, 12, 17, 12, 53, 1, 25*1000000, time.UTC), device.CreatedDate)
	assert.Equal(time.Date(2020, 12, 17, 12, 53, 1, 25*1000000, time.UTC), device.UpdatedDate)
	assert.Nil(device.DeviceSNMPv3Conf)
	assert.Equal(models.CDNAttributeYes, *device.CDNAttr)
	assert.Nil(device.BGPPeerIP4)
	assert.Nil(device.BGPPeerIP6)
	assert.Nil(device.SNMPLastUpdated)
	assert.Equal(models.DeviceSubtypeAwsSubnet, device.DeviceSubType)
}

func TestGetAllDevices(t *testing.T) {
	// arrange
	getResponsePayload := `
    {
        "devices": [
            {
                "id": "42",
                "company_id": "74333",
                "device_name": "testapi_router_full_1",
                "device_type": "router",
                "device_status": "V",
                "device_description": "testapi router with full config",
                "site": {
                    "id": 8483,
                    "site_name": "marina gdańsk",
                    "lat": 54.348972,
                    "lon": 18.659791,
                    "company_id": 74333
                },
                "plan": {
                    "active": true,
                    "bgp_enabled": true,
                    "cdate": "2020-09-03T08:41:57.489Z",
                    "company_id": 74333,
                    "description": "Your Free Trial includes 6 devices (...)",
                    "deviceTypes": [],
                    "devices": [],
                    "edate": "2020-09-03T08:41:57.489Z",
                    "fast_retention": 30,
                    "full_retention": 30,
                    "id": 11466,
                    "max_bigdata_fps": 30,
                    "max_devices": 6,
                    "max_fps": 1000,
                    "name": "Free Trial Plan",
                    "metadata": {}
                },
                "labels": [
                            {
                                "id": 2590,
                                "name": "AWS: terraform-demo-aws",
                                "description": null,
                                "edate": "2020-10-05T15:28:00.276Z",
                                "cdate": "2020-10-05T15:28:00.276Z",
                                "user_id": "133210",
                                "company_id": "74333",
                                "color": "#5340A5",
                                "order": null,
                                "_pivot_device_id": "77715",
                                "_pivot_label_id": "2590"
                            },
                            {
                                "id": 2751,
                                "name": "GCP: traffic-generator-gcp",
                                "description": null,
                                "edate": "2020-11-20T12:54:49.575Z",
                                "cdate": "2020-11-20T12:54:49.575Z",
                                "user_id": null,
                                "company_id": "74333",
                                "color": "#5289D9",
                                "order": null,
                                "_pivot_device_id": "77373",
                                "_pivot_label_id": "2751"
                            }
                        ],
                "all_interfaces": [],
                "device_flow_type": "auto",
                "device_sample_rate": "1001",
                "sending_ips": [
                    "128.0.0.11",
                    "128.0.0.12"
                ],
                "device_snmp_ip": "129.0.0.1",
                "device_snmp_community": "",
                "minimize_snmp": false,
                "device_bgp_type": "device",
                "device_bgp_neighbor_ip": "127.0.0.1",
                "device_bgp_neighbor_ip6": null,
                "device_bgp_neighbor_asn": "11",
                "device_bgp_flowspec": true,
                "device_bgp_password": "*********ass",
                "use_bgp_device_id": null,
                "custom_columns": "",
                "custom_column_data": [],
                "device_chf_client_port": null,
                "device_chf_client_protocol": null,
                "device_chf_interface": null,
                "device_agent_type": null,
                "max_flow_rate": 1000,
                "max_big_flow_rate": 30,
                "device_proxy_bgp": "",
                "device_proxy_bgp6": "",
                "created_date": "2020-12-17T08:24:45.074Z",
                "updated_date": "2020-12-17T08:24:45.074Z",
                "device_snmp_v3_conf": {
                    "UserName": "John",
                    "AuthenticationProtocol": "MD5",
                    "AuthenticationPassphrase": "john_md5_pass",
                    "PrivacyProtocol": "DES",
                    "PrivacyPassphrase": "**********ass"
                },
                "bgpPeerIP4": "208.76.14.223",
                "bgpPeerIP6": "2620:129:1:2::1",
                "snmp_last_updated": null,
                "device_subtype": "router"
            },
            {
                "id": "43",
                "company_id": "74333",
                "device_name": "testapi_dns_minimal_1",
                "device_type": "host-nprobe-dns-www",
                "device_status": "V",
                "device_description": "testapi dns with minimal config",
                "site": {
                    "id": null,
                    "site_name": null,
                    "lat": null,
                    "lon": null,
                    "company_id": null
                },
                "plan": {
                    "active": true,
                    "bgp_enabled": true,
                    "cdate": "2020-09-03T08:41:57.489Z",
                    "company_id": 74333,
                    "description": "Your Free Trial includes 6 devices (...)",
                    "deviceTypes": [],
                    "devices": [],
                    "edate": "2020-09-03T08:41:57.489Z",
                    "fast_retention": 30,
                    "full_retention": 30,
                    "id": 11466,
                    "max_bigdata_fps": 30,
                    "max_devices": 6,
                    "max_fps": 1000,
                    "name": "Free Trial Plan",
                    "metadata": {}
                },
                "labels": [],
                "all_interfaces": [],
                "device_flow_type": "auto",
                "device_sample_rate": "1",
                "sending_ips": [],
                "device_snmp_ip": null,
                "device_snmp_community": "",
                "minimize_snmp": false,
                "device_bgp_type": "none",
                "device_bgp_neighbor_ip": null,
                "device_bgp_neighbor_ip6": null,
                "device_bgp_neighbor_asn": null,
                "device_bgp_flowspec": false,
                "device_bgp_password": null,
                "use_bgp_device_id": null,
                "custom_columns": "",
                "custom_column_data": [],
                "device_chf_client_port": null,
                "device_chf_client_protocol": null,
                "device_chf_interface": null,
                "device_agent_type": null,
                "max_flow_rate": 1000,
                "max_big_flow_rate": 30,
                "device_proxy_bgp": "",
                "device_proxy_bgp6": "",
                "created_date": "2020-12-17T12:53:01.025Z",
                "updated_date": "2020-12-17T12:53:01.025Z",
                "device_snmp_v3_conf": null,
                "cdn_attr": "Y",
                "snmp_last_updated": null,
                "device_subtype": "aws_subnet"
            }
        ]
    }`
	transport := &api_connection.StubTransport{ResponseBody: getResponsePayload}
	devicesAPI := api_resources.NewDevicesAPI(transport)

	// act
	devices, err := devicesAPI.GetAll(nil)

	// assert request properly formed
	assert := assert.New(t)
	require := require.New(t)

	require.NoError(err)
	assert.Zero(transport.RequestBody)

	// and response properly parsed
	require.Equal(2, len(devices))
	// device 0
	device := devices[0]
	assert.Equal(models.ID(42), device.ID)
	assert.Equal(models.ID(74333), device.CompanyID)
	assert.Equal("testapi_router_full_1", device.DeviceName)
	assert.Equal(models.DeviceTypeRouter, device.DeviceType)
	assert.Equal("testapi router with full config", *device.DeviceDescription)
	require.NotNil(device.Site)
	assert.Equal(models.ID(8483), *device.Site.ID)
	assert.Equal("marina gdańsk", *device.Site.SiteName)
	assert.Equal(54.348972, *device.Site.Latitude)
	assert.Equal(18.659791, *device.Site.Longitude)
	assert.Equal(models.ID(74333), *device.Site.CompanyID)
	require.NotNil(device.Plan)
	assert.True(*device.Plan.Active)
	assert.True(*device.Plan.BGPEnabled)
	assert.Equal(time.Date(2020, 9, 3, 8, 41, 57, 489*1000000, time.UTC), *device.Plan.CreatedDate)
	assert.Equal(models.ID(74333), *device.Plan.CompanyID)
	assert.Equal("Your Free Trial includes 6 devices (...)", *device.Plan.Description)
	assert.Equal(0, len(device.Plan.DeviceTypes))
	assert.Equal(0, len(device.Plan.Devices))
	assert.Equal(time.Date(2020, 9, 3, 8, 41, 57, 489*1000000, time.UTC), *device.Plan.UpdatedDate)
	assert.Equal(30, *device.Plan.FastRetention)
	assert.Equal(30, *device.Plan.FullRetention)
	assert.Equal(models.ID(11466), *device.Plan.ID)
	assert.Equal(30, *device.Plan.MaxBigdataFPS)
	assert.Equal(6, *device.Plan.MaxDevices)
	assert.Equal(1000, *device.Plan.MaxFPS)
	assert.Equal("Free Trial Plan", *device.Plan.Name)
	assert.Equal(2, len(device.Labels))
	assert.Equal(models.ID(2590), device.Labels[0].ID)
	assert.Equal("AWS: terraform-demo-aws", device.Labels[0].Name)
	assert.Equal(time.Date(2020, 10, 5, 15, 28, 00, 276*1000000, time.UTC), device.Labels[0].UpdatedDate)
	assert.Equal(time.Date(2020, 10, 5, 15, 28, 00, 276*1000000, time.UTC), device.Labels[0].CreatedDate)
	assert.Equal(models.ID(133210), *device.Labels[0].UserID)
	assert.Equal(models.ID(74333), device.Labels[0].CompanyID)
	assert.Equal("#5340A5", device.Labels[0].Color)
	assert.Equal(models.ID(2751), device.Labels[1].ID)
	assert.Equal("GCP: traffic-generator-gcp", device.Labels[1].Name)
	assert.Equal(time.Date(2020, 11, 20, 12, 54, 49, 575*1000000, time.UTC), device.Labels[1].UpdatedDate)
	assert.Equal(time.Date(2020, 11, 20, 12, 54, 49, 575*1000000, time.UTC), device.Labels[1].CreatedDate)
	assert.Nil(device.Labels[1].UserID)
	assert.Equal(models.ID(74333), device.Labels[1].CompanyID)
	assert.Equal("#5289D9", device.Labels[1].Color)
	require.Equal(0, len(device.AllInterfaces))
	assert.Equal("auto", *device.DeviceFlowType)
	assert.Equal(1001, *&device.DeviceSampleRate)
	assert.Equal(2, len(device.SendingIPS))
	assert.Equal("128.0.0.11", device.SendingIPS[0])
	assert.Equal("128.0.0.12", device.SendingIPS[1])
	assert.Equal("129.0.0.1", *device.DeviceSNMNPIP)
	assert.Equal("", *device.DeviceSNMPCommunity)
	assert.False(*device.MinimizeSNMP)
	assert.Equal(models.DeviceBGPTypeDevice, *device.DeviceBGPType)
	assert.Equal("127.0.0.1", *device.DeviceBGPNeighborIP)
	assert.Nil(device.DeviceBGPNeighborIPv6)
	assert.Equal("11", *device.DeviceBGPNeighborASN)
	assert.True(*device.DeviceBGPFlowSpec)
	assert.Equal("*********ass", *device.DeviceBGPPassword)
	assert.Nil(device.UseBGPDeviceID)
	assert.Equal(time.Date(2020, 12, 17, 8, 24, 45, 74*1000000, time.UTC), device.CreatedDate)
	assert.Equal(time.Date(2020, 12, 17, 8, 24, 45, 74*1000000, time.UTC), device.UpdatedDate)
	require.NotNil(device.DeviceSNMPv3Conf)
	assert.Equal("John", device.DeviceSNMPv3Conf.UserName)
	assert.Equal(models.AuthenticationProtocolMD5, *device.DeviceSNMPv3Conf.AuthenticationProtocol)
	assert.Equal("john_md5_pass", *device.DeviceSNMPv3Conf.AuthenticationPassphrase)
	assert.Equal(models.PrivacyProtocolDES, *device.DeviceSNMPv3Conf.PrivacyProtocol)
	assert.Equal("**********ass", *device.DeviceSNMPv3Conf.PrivacyPassphrase)
	assert.Equal("208.76.14.223", *device.BGPPeerIP4)
	assert.Equal("2620:129:1:2::1", *device.BGPPeerIP6)
	assert.Nil(device.SNMPLastUpdated)
	assert.Equal(models.DeviceSubtypeRouter, device.DeviceSubType)

	// device 1
	device = devices[1]
	assert.Equal(models.ID(43), device.ID)
	assert.Equal(models.DeviceTypeHostNProbeDNSWWW, device.DeviceType)
	assert.Equal(models.DeviceSubtypeAwsSubnet, device.DeviceSubType)
	assert.Equal(models.DeviceBGPTypeNone, *device.DeviceBGPType)
}

func TestDeleteDevice(t *testing.T) {
	// arrange
	deleteResponsePayload := "" // deleting device responds with empty body
	transport := &api_connection.StubTransport{ResponseBody: deleteResponsePayload}
	devicesAPI := api_resources.NewDevicesAPI(transport)

	// act
	deviceID := models.ID(42)
	err := devicesAPI.Delete(nil, deviceID)

	// assert
	assert := assert.New(t)
	require := require.New(t)
	require.NoError(err)
	assert.Zero(transport.RequestBody)
}

func TestApplyLabels(t *testing.T) {
	// arrange
	applyLabelsResponsePayload := `
    {
        "id": "42",
        "device_name": "test_router",
        "labels": [
            {
                "id": 3011,
                "name": "apitest-label-red",
                "description": null,
                "edate": "2021-01-11T08:38:08.678Z",
                "cdate": "2021-01-11T08:38:08.678Z",
                "user_id": "144319",
                "company_id": "74333",
                "color": "#FF0000",
                "order": null,
                "_pivot_device_id": "79175",
                "_pivot_label_id": "3011"
            },
            {
                "id": 3012,
                "name": "apitest-label-blue",
                "description": null,
                "edate": "2021-01-11T08:38:42.627Z",
                "cdate": "2021-01-11T08:38:42.627Z",
                "user_id": "144319",
                "company_id": "74333",
                "color": "#0000FF",
                "order": null,
                "_pivot_device_id": "79175",
                "_pivot_label_id": "3012"
            }
        ]
    }`
	transport := &api_connection.StubTransport{ResponseBody: applyLabelsResponsePayload}
	devicesAPI := api_resources.NewDevicesAPI(transport)

	// act
	deviceID := models.ID(42)
	labels := []models.ID{models.ID(3011), models.ID(3012)}
	result, err := devicesAPI.ApplyLabels(nil, deviceID, labels)

	// assert request properly formed
	assert := assert.New(t)
	require := require.New(t)
	payload := utils.NewJSONPayloadInspector(t, transport.RequestBody)

	require.NoError(err)
	require.NotNil(payload.Get("labels"))
	assert.Equal(2, payload.Count("labels/*"))
	assert.Equal(models.ID(3011), payload.Int("labels/*[1]/id"))
	assert.Equal(models.ID(3012), payload.Int("labels/*[2]/id"))

	// and response properly parsed
	assert.Equal(models.ID(42), result.ID)
	assert.Equal("test_router", result.DeviceName)
	assert.Equal(2, len(result.Labels))
	assert.Equal(models.ID(3011), result.Labels[0].ID)
	assert.Equal("apitest-label-red", result.Labels[0].Name)
	assert.Equal(time.Date(2021, 1, 11, 8, 38, 8, 678*1000000, time.UTC), result.Labels[0].CreatedDate)
	assert.Equal(time.Date(2021, 1, 11, 8, 38, 8, 678*1000000, time.UTC), result.Labels[0].UpdatedDate)
	assert.Equal(models.ID(144319), *result.Labels[0].UserID)
	assert.Equal(models.ID(74333), result.Labels[0].CompanyID)
	assert.Equal("#FF0000", result.Labels[0].Color)
	assert.Equal(models.ID(3012), result.Labels[1].ID)
	assert.Equal("apitest-label-blue", result.Labels[1].Name)
	assert.Equal(time.Date(2021, 1, 11, 8, 38, 42, 627*1000000, time.UTC), result.Labels[1].CreatedDate)
	assert.Equal(time.Date(2021, 1, 11, 8, 38, 42, 627*1000000, time.UTC), result.Labels[1].UpdatedDate)
	assert.Equal(models.ID(144319), *result.Labels[1].UserID)
	assert.Equal(models.ID(74333), result.Labels[1].CompanyID)
	assert.Equal("#0000FF", result.Labels[1].Color)
}

func TestGetInterfaceMinimal(t *testing.T) {
	// arrange
	getResponsePayload := `
    {
        "interface": {
            "id": "43",
            "company_id": "74333",
            "device_id": "42",
            "snmp_id": "1",
            "snmp_speed": "15",
            "snmp_type": null,
            "snmp_alias": null,
            "interface_ip": null,
            "interface_description": null,
            "cdate": "2021-01-13T08:50:37.068Z",
            "edate": "2021-01-13T08:55:59.403Z",
            "initial_snmp_id": null,
            "initial_snmp_alias": null,
            "initial_interface_description": null,
            "initial_snmp_speed": null,
            "interface_ip_netmask": null,
            "top_nexthop_asns": null,
            "provider": null,
            "vrf_id": null,
            "vrf": null,
            "secondary_ips": null
        }
    }`
	transport := &api_connection.StubTransport{ResponseBody: getResponsePayload}
	devicesAPI := api_resources.NewDevicesAPI(transport)

	// act
	deviceID := models.ID(42)
	interfaceID := models.ID(43)
	intf, err := devicesAPI.Interfaces.Get(nil, deviceID, interfaceID)

	// assert request properly formed
	assert := assert.New(t)
	require := require.New(t)

	require.NoError(err)
	assert.Zero(transport.RequestBody)

	// and response properly parsed
	assert.Equal(models.ID(43), intf.ID)
	assert.Equal(models.ID(74333), intf.CompanyID)
	assert.Equal(models.ID(42), intf.DeviceID)
	assert.Equal(models.ID(1), intf.SNMPID)
	assert.Equal(15, intf.SNMPSpeed)
	assert.Nil(intf.SNMPAlias)
	assert.Nil(intf.InterfaceIP)
	assert.Nil(intf.InterfaceDescription)
	assert.Equal(time.Date(2021, 1, 13, 8, 50, 37, 68*1000000, time.UTC), intf.CreatedDate)
	assert.Equal(time.Date(2021, 1, 13, 8, 55, 59, 403*1000000, time.UTC), intf.UpdatedDate)
	assert.Nil(intf.InitialSNMPID)
	assert.Nil(intf.InitialSNMPAlias)
	assert.Nil(intf.InitialInterfaceDescription)
	assert.Nil(intf.InitialSNMPSpeed)
	assert.Nil(intf.InterfaceIPNetmask)
	assert.Equal(0, len(intf.TopNextHopASNs))
	assert.Nil(intf.VRFID)
	assert.Nil(intf.VRF)
	assert.Equal(0, len(intf.SecondaryIPS))
}

func TestGetInterfaceFull(t *testing.T) {
	// arrange
	getResponsePayload := `
    {
        "interface": {
            "id": "43",
            "company_id": "74333",
            "device_id": "42",
            "snmp_id": "1",
            "snmp_speed": "15",
            "snmp_type": null,
            "snmp_alias": "interace-description-1",
            "interface_ip": "127.0.0.1",
            "interface_description": "testapi-interface-1",
            "interface_kvs": "",
            "interface_tags": "",
            "interface_status": "V",
            "extra_info": {},
            "cdate": "2021-01-13T08:50:37.068Z",
            "edate": "2021-01-13T08:55:59.403Z",
            "initial_snmp_id": "150",
            "initial_snmp_alias": "initial-interace-description-1",
            "initial_interface_description": "initial-testapi-interface-1",
            "initial_snmp_speed": "7",
            "interface_ip_netmask": "255.255.255.0",
            "connectivity_type": "",
            "network_boundary": "",
            "initial_connectivity_type": "",
            "initial_network_boundary": "",
            "top_nexthop_asns": [
                {
                    "ASN": 20,
                    "packets":30100
                },
                {
                    "ASN": 21,
                    "fala": "hala",
                    "packets":30101
                }
            ],
            "provider": "",
            "initial_provider": "",
            "vrf_id": "39902",
            "vrf": {
                "id": 39902,
                "company_id": "74333",
                "description": "vrf-description",
                "device_id": "79175",
                "name": "vrf-name",
                "route_distinguisher": "11.121.111.13:3254",
                "route_target": "101:100"
            },
            "secondary_ips": [
                {
                "address": "198.186.193.51",
                "netmask": "255.255.255.240"
                },
                {
                "address": "198.186.193.63",
                "netmask": "255.255.255.225"
                }
            ]
        }
    }`
	transport := &api_connection.StubTransport{ResponseBody: getResponsePayload}
	devicesAPI := api_resources.NewDevicesAPI(transport)

	// act
	deviceID := models.ID(42)
	interfaceID := models.ID(43)
	intf, err := devicesAPI.Interfaces.Get(nil, deviceID, interfaceID)

	// assert request properly formed
	assert := assert.New(t)
	require := require.New(t)

	require.NoError(err)
	assert.Zero(transport.RequestBody)

	// and response properly parsed
	assert.Equal(models.ID(43), intf.ID)
	assert.Equal(models.ID(74333), intf.CompanyID)
	assert.Equal(models.ID(42), intf.DeviceID)
	assert.Equal(models.ID(1), intf.SNMPID)
	assert.Equal(15, intf.SNMPSpeed)
	assert.Equal("interace-description-1", *intf.SNMPAlias)
	assert.Equal("127.0.0.1", *intf.InterfaceIP)
	assert.Equal("testapi-interface-1", *intf.InterfaceDescription)
	assert.Equal(time.Date(2021, 1, 13, 8, 50, 37, 68*1000000, time.UTC), intf.CreatedDate)
	assert.Equal(time.Date(2021, 1, 13, 8, 55, 59, 403*1000000, time.UTC), intf.UpdatedDate)
	assert.Equal("150", *intf.InitialSNMPID)
	assert.Equal("initial-interace-description-1", *intf.InitialSNMPAlias)
	assert.Equal("initial-testapi-interface-1", *intf.InitialInterfaceDescription)
	assert.Equal(7, *intf.InitialSNMPSpeed)
	assert.Equal("255.255.255.0", *intf.InterfaceIPNetmask)
	assert.Equal(2, len(intf.TopNextHopASNs))
	assert.Equal(20, intf.TopNextHopASNs[0].ASN)
	assert.Equal(30100, intf.TopNextHopASNs[0].Packets)
	assert.Equal(21, intf.TopNextHopASNs[1].ASN)
	assert.Equal(30101, intf.TopNextHopASNs[1].Packets)
	assert.Equal(models.ID(39902), *intf.VRFID)
	require.NotNil(intf.VRF)
	assert.Equal(models.ID(74333), intf.VRF.CompanyID)
	assert.Equal("vrf-description", *intf.VRF.Description)
	assert.Equal(models.ID(79175), intf.VRF.DeviceID)
	assert.Equal("vrf-name", intf.VRF.Name)
	assert.Equal("11.121.111.13:3254", intf.VRF.RouteDistinguisher)
	assert.Equal("101:100", intf.VRF.RouteTarget)
	assert.Equal(2, len(intf.SecondaryIPS))
	assert.Equal("198.186.193.51", intf.SecondaryIPS[0].Address)
	assert.Equal("255.255.255.240", intf.SecondaryIPS[0].Netmask)
	assert.Equal("198.186.193.63", intf.SecondaryIPS[1].Address)
	assert.Equal("255.255.255.225", intf.SecondaryIPS[1].Netmask)
}

func TestGetAllInterfaces(t *testing.T) {
	// arrange
	getResponsePayload := `
    [
        {
            "id": "43",
            "company_id": "74333",
            "device_id": "42",
            "snmp_id": "1",
            "snmp_speed": "15",
            "snmp_type": null,
            "snmp_alias": "interace-description-1",
            "interface_ip": "127.0.0.1",
            "interface_description": "testapi-interface-1",
            "interface_kvs": "",
            "interface_tags": "",
            "interface_status": "V",
            "extra_info": {},
            "cdate": "2021-01-13T08:50:37.068Z",
            "edate": "2021-01-13T08:55:59.403Z",
            "initial_snmp_id": "150",
            "initial_snmp_alias": "initial-interace-description-1",
            "initial_interface_description": "initial-testapi-interface-1",
            "initial_snmp_speed": "7",
            "interface_ip_netmask": "255.255.255.0",
            "connectivity_type": "",
            "network_boundary": "",
            "initial_connectivity_type": "",
            "initial_network_boundary": "",
            "top_nexthop_asns": [
                {
                    "ASN": 20,
                    "packets":30100
                },
                {
                    "ASN": 21,
                    "packets":30101
                }
            ],
            "provider": "",
            "initial_provider": "",
            "vrf_id": "39902",
            "vrf": {
                "id": 39902,
                "company_id": "74333",
                "description": "vrf-description",
                "device_id": "79175",
                "name": "vrf-name",
                "route_distinguisher": "11.121.111.13:3254",
                "route_target": "101:100"
            },
            "secondary_ips": [
                {
                "address": "198.186.193.51",
                "netmask": "255.255.255.240"
                },
                {
                "address": "198.186.193.63",
                "netmask": "255.255.255.225"
                }
            ]
        },
        {
            "id": "44",
            "company_id": "74333",
            "device_id": "42",
            "snmp_id": "1",
            "snmp_speed": "15",
            "snmp_type": null,
            "snmp_alias": "interace-description-1",
            "interface_ip": "127.0.0.1",
            "interface_description": "testapi-interface-1",
            "interface_kvs": "",
            "interface_tags": "",
            "interface_status": "V",
            "extra_info": {},
            "cdate": "2021-01-13T08:50:37.068Z",
            "edate": "2021-01-13T08:50:37.074Z",
            "initial_snmp_id": "",
            "initial_snmp_alias": null,
            "initial_interface_description": null,
            "initial_snmp_speed": null,
            "interface_ip_netmask": "255.255.255.0",
            "secondary_ips": null,
            "connectivity_type": "",
            "network_boundary": "",
            "initial_connectivity_type": "",
            "initial_network_boundary": "",
            "top_nexthop_asns": null,
            "provider": "",
            "initial_provider": "",
            "vrf_id": "39902",
            "vrf": {
                "id": 39902,
                "company_id": "74333",
                "description": "vrf-description",
                "device_id": "42",
                "name": "vrf-name",
                "route_distinguisher": "11.121.111.13:3254",
                "route_target": "101:100"
            }
        },
        {
            "id": "45",
            "company_id": "74333",
            "device_id": "42",
            "snmp_id": "1",
            "snmp_speed": "15",
            "snmp_type": null,
            "snmp_alias": "interace-description-1",
            "interface_ip": "127.0.0.1",
            "interface_description": "testapi-interface-1",
            "interface_kvs": "",
            "interface_tags": "",
            "interface_status": "V",
            "extra_info": {},
            "cdate": "2021-01-13T08:50:37.068Z",
            "edate": "2021-01-13T08:50:37.074Z",
            "initial_snmp_id": "",
            "initial_snmp_alias": null,
            "initial_interface_description": null,
            "initial_snmp_speed": null,
            "interface_ip_netmask": "255.255.255.0",
            "secondary_ips": null,
            "connectivity_type": "",
            "network_boundary": "",
            "initial_connectivity_type": "",
            "initial_network_boundary": "",
            "top_nexthop_asns": null,
            "provider": "",
            "initial_provider": "",
            "vrf_id": "39902",
            "vrf": {}
        }
    ]`
	transport := &api_connection.StubTransport{ResponseBody: getResponsePayload}
	devicesAPI := api_resources.NewDevicesAPI(transport)

	// act
	deviceID := models.ID(42)
	interfaces, err := devicesAPI.Interfaces.GetAll(nil, deviceID)

	// assert request properly formed
	assert := assert.New(t)
	require := require.New(t)

	require.NoError(err)
	assert.Zero(transport.RequestBody)

	// and response properly parsed
	assert.Equal(3, len(interfaces))
	intf := interfaces[0]
	assert.Equal(models.ID(43), intf.ID)
	assert.Equal(models.ID(74333), intf.CompanyID)
	assert.Equal(models.ID(42), intf.DeviceID)
	assert.Equal(models.ID(1), intf.SNMPID)
	assert.Equal(15, intf.SNMPSpeed)
	assert.Equal("interace-description-1", *intf.SNMPAlias)
	assert.Equal("127.0.0.1", *intf.InterfaceIP)
	assert.Equal("testapi-interface-1", *intf.InterfaceDescription)
	assert.Equal(time.Date(2021, 1, 13, 8, 50, 37, 68*1000000, time.UTC), intf.CreatedDate)
	assert.Equal(time.Date(2021, 1, 13, 8, 55, 59, 403*1000000, time.UTC), intf.UpdatedDate)
	assert.Equal("150", *intf.InitialSNMPID)
	assert.Equal("initial-interace-description-1", *intf.InitialSNMPAlias)
	assert.Equal("initial-testapi-interface-1", *intf.InitialInterfaceDescription)
	assert.Equal(7, *intf.InitialSNMPSpeed)
	assert.Equal("255.255.255.0", *intf.InterfaceIPNetmask)
	assert.Equal(2, len(intf.TopNextHopASNs))
	assert.Equal(20, intf.TopNextHopASNs[0].ASN)
	assert.Equal(30100, intf.TopNextHopASNs[0].Packets)
	assert.Equal(21, intf.TopNextHopASNs[1].ASN)
	assert.Equal(30101, intf.TopNextHopASNs[1].Packets)
	assert.Equal(models.ID(39902), *intf.VRFID)
	require.NotNil(intf.VRF)
	assert.Equal(models.ID(74333), intf.VRF.CompanyID)
	assert.Equal("vrf-description", *intf.VRF.Description)
	assert.Equal(models.ID(79175), intf.VRF.DeviceID)
	assert.Equal("vrf-name", intf.VRF.Name)
	assert.Equal("11.121.111.13:3254", intf.VRF.RouteDistinguisher)
	assert.Equal("101:100", intf.VRF.RouteTarget)
	assert.Equal(2, len(intf.SecondaryIPS))
	assert.Equal("198.186.193.51", intf.SecondaryIPS[0].Address)
	assert.Equal("255.255.255.240", intf.SecondaryIPS[0].Netmask)
	assert.Equal("198.186.193.63", intf.SecondaryIPS[1].Address)
	assert.Equal("255.255.255.225", intf.SecondaryIPS[1].Netmask)
}

func TestCreateInterfaceMinimal(t *testing.T) {
	// arrange
	getResponsePayload := `
    {
        "snmp_id": "2",
        "snmp_speed": 8,
        "interface_description": "testapi-interface-2",
        "interface_kvs": "",
        "company_id": "74333",
        "device_id": "42",
        "edate": "2021-01-13T08:41:16.191Z",
        "cdate": "2021-01-13T08:41:16.191Z",
        "id": "43"
    }`
	transport := &api_connection.StubTransport{ResponseBody: getResponsePayload}
	devicesAPI := api_resources.NewDevicesAPI(transport)

	// act
	deviceID := models.ID(42)
	intf := models.NewInterface(
		deviceID,
		models.ID(2),
		8,
		"testapi-interface-2",
	)
	created, err := devicesAPI.Interfaces.Create(nil, *intf)

	// assert request properly formed
	assert := assert.New(t)
	require := require.New(t)

	require.NoError(err)
	payload := utils.NewJSONPayloadInspector(t, transport.RequestBody)
	assert.Equal(2, payload.Int("snmp_id"))
	assert.Equal(8, payload.Int("snmp_speed"))
	assert.Equal("testapi-interface-2", payload.String("interface_description"))

	// and response properly parsed
	assert.Equal(deviceID, created.DeviceID)
	assert.Equal(models.ID(2), created.SNMPID)
	assert.Equal(models.ID(74333), created.CompanyID)
	assert.Equal(8, created.SNMPSpeed)
	assert.Equal("testapi-interface-2", *created.InterfaceDescription)
	assert.Equal(time.Date(2021, 1, 13, 8, 41, 16, 191*1000000, time.UTC), created.CreatedDate)
	assert.Equal(time.Date(2021, 1, 13, 8, 41, 16, 191*1000000, time.UTC), created.UpdatedDate)
	assert.Equal(0, len(created.SecondaryIPS))
	assert.Nil(created.SNMPAlias)
	assert.Nil(created.InterfaceIP)
	assert.Nil(created.InterfaceIPNetmask)
	assert.Nil(created.VRFID)
	assert.Nil(created.VRF)
}

func TestCreateInterfaceFull(t *testing.T) {
	// arrange
	getResponsePayload := `
    {
        "snmp_id": "243205880",
        "snmp_alias": "interace-description-1",
        "snmp_speed": 8,
        "interface_description": "testapi-interface-1",
        "interface_ip": "127.0.0.1",
        "interface_ip_netmask": "255.255.255.0",
        "interface_kvs": "",
        "company_id": "74333",
        "device_id": "42",
        "edate": "2021-01-13T08:31:40.629Z",
        "cdate": "2021-01-13T08:31:40.619Z",
        "id": "43",
        "vrf_id": 39903,
        "secondary_ips": [
            {
                "address": "198.186.193.51",
                "netmask": "255.255.255.240"
            },
            {
                "address": "198.186.193.63",
                "netmask": "255.255.255.225"
            }
        ]
    }`
	transport := &api_connection.StubTransport{ResponseBody: getResponsePayload}
	devicesAPI := api_resources.NewDevicesAPI(transport)

	// act
	vrf := models.NewVRFAttributes(
		"vrf-name",
		"101:100",
		"11.121.111.13:3254",
	)
	models.SetOptional(&vrf.Description, "vrf-description")
	models.SetOptional(&vrf.ExtRouteDistinguisher, "15")
	secondaryIP1 := models.SecondaryIP{Address: "127.0.0.2", Netmask: "255.255.255.0"}
	secondaryIP2 := models.SecondaryIP{Address: "127.0.0.3", Netmask: "255.255.255.0"}
	deviceID := models.ID(42)
	intf := models.NewInterface(
		deviceID,
		models.ID(2),
		8,
		"testapi-interface-2",
	)
	models.SetOptional(&intf.SNMPAlias, "interace-description-1")
	models.SetOptional(&intf.InterfaceIP, "127.0.0.1")
	models.SetOptional(&intf.InterfaceIPNetmask, "255.255.255.0")
	intf.SecondaryIPS = []models.SecondaryIP{secondaryIP1, secondaryIP2}
	intf.VRF = vrf
	created, err := devicesAPI.Interfaces.Create(nil, *intf)

	// assert request properly formed
	assert := assert.New(t)
	require := require.New(t)

	require.NoError(err)
	payload := utils.NewJSONPayloadInspector(t, transport.RequestBody)
	assert.Equal(2, payload.Int("snmp_id"))
	assert.Equal(8, payload.Int("snmp_speed"))
	assert.Equal("testapi-interface-2", payload.String("interface_description"))

	assert.Equal("interace-description-1", payload.String("snmp_alias"))
	assert.Equal("127.0.0.1", payload.String("interface_ip"))
	assert.Equal("255.255.255.0", payload.String("interface_ip_netmask"))
	assert.Equal("vrf-name", payload.String("vrf/name"))
	assert.Equal("vrf-description", payload.String("vrf/description"))
	assert.Equal("101:100", payload.String("vrf/route_target"))
	assert.Equal("11.121.111.13:3254", payload.String("vrf/route_distinguisher"))
	assert.Equal(15, payload.Int("vrf/ext_route_distinguisher"))
	assert.Equal("127.0.0.2", payload.String("secondary_ips/*[1]/address"))
	assert.Equal("255.255.255.0", payload.String("secondary_ips/*[1]/netmask"))
	assert.Equal("127.0.0.3", payload.String("secondary_ips/*[2]/address"))
	assert.Equal("255.255.255.0", payload.String("secondary_ips/*[2]/netmask"))

	// and response properly parsed
	assert.Equal(deviceID, created.DeviceID)
	assert.Equal(models.ID(243205880), created.SNMPID)
	assert.Equal(models.ID(74333), created.CompanyID)
	assert.Equal(8, created.SNMPSpeed)
	assert.Equal("testapi-interface-1", *created.InterfaceDescription)
	assert.Equal(time.Date(2021, 1, 13, 8, 31, 40, 619*1000000, time.UTC), created.CreatedDate)
	assert.Equal(time.Date(2021, 1, 13, 8, 31, 40, 629*1000000, time.UTC), created.UpdatedDate)
	assert.Equal(2, len(created.SecondaryIPS))
	assert.Equal("198.186.193.51", created.SecondaryIPS[0].Address)
	assert.Equal("255.255.255.240", created.SecondaryIPS[0].Netmask)
	assert.Equal("198.186.193.63", created.SecondaryIPS[1].Address)
	assert.Equal("255.255.255.225", created.SecondaryIPS[1].Netmask)
	assert.Equal("interace-description-1", *created.SNMPAlias)
	assert.Equal("127.0.0.1", *created.InterfaceIP)
	assert.Equal("255.255.255.0", *created.InterfaceIPNetmask)
	assert.Equal(models.ID(39903), *created.VRFID)
	assert.Nil(created.VRF)
}

func TestDeleteInterface(t *testing.T) {
	// arrange
	deleteResponsePayload := "{}"
	transport := &api_connection.StubTransport{ResponseBody: deleteResponsePayload}
	devicesAPI := api_resources.NewDevicesAPI(transport)

	// act
	deviceID := models.ID(42)
	interfaceID := models.ID(43)
	err := devicesAPI.Interfaces.Delete(nil, deviceID, interfaceID)

	// assert
	assert := assert.New(t)
	require := require.New(t)
	require.NoError(err)
	assert.Zero(transport.RequestBody)
}
