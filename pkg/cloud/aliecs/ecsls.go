package aliecs

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	log "github.com/sirupsen/logrus"
	"github.com/teamssix/cf/pkg/cloud"
	"github.com/teamssix/cf/pkg/util"
	"github.com/teamssix/cf/pkg/util/cmdutil"
)

type Instances struct {
	InstanceId       string
	InstanceName     string
	OSName           string
	OSType           string
	Status           string
	PrivateIpAddress string
	PublicIpAddress  string
	RegionId         string
}

var (
	ECSCacheFilePath = cmdutil.ReturnECSCacheFile()
	header           = []string{"序号 (SN)", "实例 ID (Instance ID)", "实例名称 (Instance Name)", "系统名称 (OS Name)", "系统类型 (OS Type)", "状态 (Status)", "私有 IP (Private Ip Address)", "公网 IP (Public Ip Address)", "区域 ID (Region ID)"}
)

func DescribeInstances(region string, running bool, SpecifiedInstanceID string) []Instances {
	var out []Instances
	var response *ecs.DescribeInstancesResponse
	request := ecs.CreateDescribeInstancesRequest()
	request.Scheme = "https"
	if running == true {
		request.Status = "Running"
	}
	if SpecifiedInstanceID != "all" {
		request.InstanceIds = fmt.Sprintf("[\"%s\"]", SpecifiedInstanceID)
	}
	response, err := ECSClient(region).DescribeInstances(request)
	util.HandleErr(err)
	InstancesList := response.Instances.Instance
	log.Tracef("正在 %s 区域中查找实例 (Looking for instances in the %s region)", region, region)
	if len(InstancesList) != 0 {
		log.Debugf("在 %s 区域下找到 %d 个实例 (Found %d instances in %s region)", region, len(InstancesList), len(InstancesList), region)
		for _, i := range InstancesList {
			// When the instance has multiple IPs, it is presented in a different format.
			var PrivateIpAddressList []string
			var PublicIpAddressList []string
			var PrivateIpAddress string
			var PublicIpAddress string
			for _, m := range i.NetworkInterfaces.NetworkInterface {
				for _, n := range m.PrivateIpSets.PrivateIpSet {
					PrivateIpAddressList = append(PrivateIpAddressList, n.PrivateIpAddress)
				}
			}
			a, _ := json.Marshal(PrivateIpAddressList)

			if len(PrivateIpAddressList) == 1 {
				PrivateIpAddress = PrivateIpAddressList[0]
			} else {
				PrivateIpAddress = string(a)
			}

			PublicIpAddressList = i.PublicIpAddress.IpAddress
			b, _ := json.Marshal(PublicIpAddressList)
			if len(PublicIpAddressList) == 1 {
				PublicIpAddress = i.PublicIpAddress.IpAddress[0]
			} else {
				PublicIpAddress = string(b)
			}
			obj := Instances{
				InstanceId:       i.InstanceId,
				InstanceName:     i.InstanceName,
				OSName:           i.OSName,
				OSType:           i.OSType,
				Status:           i.Status,
				PrivateIpAddress: PrivateIpAddress,
				PublicIpAddress:  PublicIpAddress,
				RegionId:         i.RegionId,
			}
			out = append(out, obj)
		}
	}
	return out
}

func ReturnInstancesList(region string, running bool, specifiedInstanceID string) []Instances {
	var InstancesList []Instances
	var Instance []Instances
	if region == "all" {
		for _, j := range GetECSRegions() {
			region := j.RegionId
			Instance = DescribeInstances(region, running, specifiedInstanceID)
			for _, i := range Instance {
				InstancesList = append(InstancesList, i)
			}
		}
	} else {
		InstancesList = DescribeInstances(region, running, specifiedInstanceID)
	}
	return InstancesList
}

func PrintInstancesListRealTime(region string, running bool, specifiedInstanceID string) {
	InstancesList := ReturnInstancesList(region, running, specifiedInstanceID)
	var data = make([][]string, len(InstancesList))
	for i, o := range InstancesList {
		SN := strconv.Itoa(i + 1)
		data[i] = []string{SN, o.InstanceId, o.InstanceName, o.OSName, o.OSType, o.Status, o.PrivateIpAddress, o.PublicIpAddress, o.RegionId}
	}
	var td = cloud.TableData{Header: header, Body: data}
	if len(data) == 0 {
		log.Info("未发现 ECS，可能是因为当前访问凭证权限不够 (No ECS found, Probably because the current Access Key do not have enough permissions)")
		cmdutil.WriteCacheFile(td, ECSCacheFilePath)
	} else {
		Caption := "ECS 资源 (ECS resources)"
		cloud.PrintTable(td, Caption)
		cmdutil.WriteCacheFile(td, ECSCacheFilePath)
	}
	util.WriteTimeStamp(util.ReturnECSTimeStampFile())
}

func PrintInstancesListHistory(region string, running bool, specifiedInstanceID string) {
	if cmdutil.FileExists(ECSCacheFilePath) {
		cmdutil.PrintECSCacheFile(ECSCacheFilePath, header, region, specifiedInstanceID)
	} else {
		PrintInstancesListRealTime(region, running, specifiedInstanceID)
	}
}

func PrintInstancesList(region string, running bool, specifiedInstanceID string, ecsFlushCache bool) {
	if ecsFlushCache {
		PrintInstancesListRealTime(region, running, specifiedInstanceID)
	} else {
		oldTimeStamp := util.ReadTimeStamp(util.ReturnECSTimeStampFile())
		if oldTimeStamp == 0 {
			PrintInstancesListRealTime(region, running, specifiedInstanceID)
		} else if util.IsFlushCache(oldTimeStamp) {
			PrintInstancesListRealTime(region, running, specifiedInstanceID)
		} else {
			util.TimeDifference(oldTimeStamp)
			PrintInstancesListHistory(region, running, specifiedInstanceID)
		}
	}
}
