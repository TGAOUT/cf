package aliram

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/teamssix/cf/pkg/cloud/alirds"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/teamssix/cf/pkg/cloud/aliecs"

	"github.com/teamssix/cf/pkg/cloud/alioss"

	"github.com/teamssix/cf/pkg/util"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/sts"

	"github.com/teamssix/cf/pkg/cloud"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ram"
	log "github.com/sirupsen/logrus"
)

var header = []string{"序号 (SN)", "策略名称 (PolicyName)", "描述 (Description)"}
var header2 = []string{"序号 (SN)", "可执行的操作 (Available actions)", "描述 (Description)"}
var SN = 1
var data2 [][]string

const (
	osslsAction        = "cf oss ls"
	osslsDescription   = "列出 OSS 资源"
	ecslsAction        = "cf ecs ls"
	ecslsDescription   = "列出 ECS 资源"
	ecsexecAction      = "cf ecs exec"
	ecsexecDescription = "在 ECS 上执行命令"
	rdslsAction        = "cf rds ls"
	rdslsDescription   = "列出 RDS 资源"
	consoleAction      = "cf console"
	consoleDescription = "接管控制台"
)

func ListPermissions() {
	userName := getCallerIdentity()
	log.Infof("当前用户名为 %s (Current username is %s)", userName, userName)
	var data [][]string
	if userName == "root" {
		data = append(data, []string{"1", "AdministratorAccess", "管理所有阿里云资源的权限"})
		var td = cloud.TableData{Header: header, Body: data}
		Caption := "当前凭证具备的权限 (Permissions owned)"
		cloud.PrintTable(td, Caption)
		fmt.Println()
		data2 = appendData(osslsAction, osslsDescription)
		data2 = appendData(ecslsAction, ecslsDescription)
		data2 = appendData(ecsexecAction, ecsexecDescription)
		data2 = appendData(rdslsAction, rdslsDescription)
		data2 = appendData(consoleAction, consoleDescription)
		var td2 = cloud.TableData{Header: header2, Body: data2}
		Caption2 := "当前凭证可以执行的操作 (Available actions)"
		cloud.PrintTable(td2, Caption2)
	} else {
		data, err := listPoliciesForUser(userName)
		if err == nil {
			if len(data) == 0 {
				log.Infoln("当前凭证没有任何权限 (The current access key does not have any permissions)")
			} else {
				var td = cloud.TableData{Header: header, Body: data}
				Caption := "当前凭证具备的权限 (Permissions owned)"
				cloud.PrintTable(td, Caption)
				fmt.Println()
				for _, o := range data {
					switch {
					case strings.Contains(o[1], "AdministratorAccess"):
						data2 = appendData(osslsAction, osslsDescription)
						data2 = appendData(ecslsAction, ecslsDescription)
						data2 = appendData(ecsexecAction, ecsexecDescription)
						data2 = appendData(rdslsAction, rdslsDescription)
						data2 = appendData(consoleAction, consoleDescription)
					case strings.Contains(o[1], "ReadOnlyAccess"):
						data2 = appendData(osslsAction, osslsDescription)
						data2 = appendData(ecslsAction, ecslsDescription)
						data2 = appendData(rdslsAction, rdslsDescription)
					case strings.Contains(o[1], "AliyunECSFullAccess"):
						data2 = appendData(ecslsAction, ecslsDescription)
						data2 = appendData(ecsexecAction, ecsexecDescription)
					case strings.Contains(o[1], "AliyunOSSReadOnlyAccess"):
						data2 = appendData(osslsAction, osslsDescription)
					case strings.Contains(o[1], "AliyunECSReadOnlyAccess"):
						data2 = appendData(ecslsAction, ecslsDescription)
					case strings.Contains(o[1], "AliyunECSAssistantFullAccess"):
						data2 = appendData(ecsexecAction, ecsexecDescription)
					case strings.Contains(o[1], "AliyunRDSReadOnlyAccess"):
						data2 = appendData(rdslsAction, rdslsDescription)
					case strings.Contains(o[1], "AliyunRAMFullAccess"):
						data2 = appendData(consoleAction, consoleDescription)
					}
				}
				if len(data2) == 0 {
					log.Infoln("当前凭证没有可以执行的操作 (Not available actions)")
				} else {
					var td2 = cloud.TableData{Header: header2, Body: data2}
					Caption2 := "当前凭证可以执行的操作 (Available actions)"
					cloud.PrintTable(td2, Caption2)
				}
			}
		} else if strings.Contains(err.Error(), "ErrorCode: NoPermission") {
			log.Debugln("当前凭证不具备 RAM 读权限 (No RAM read permissions)")
			obj1, obj2 := traversalPermissions()
			var data1 = make([][]string, len(obj1))
			var data2 = make([][]string, len(obj2))
			if len(obj1) == 0 {
				log.Infoln("当前凭证没有任何权限 (The current access key does not have any permissions)")
			} else {
				for i, o := range obj1 {
					SN := strconv.Itoa(i + 1)
					data1[i] = []string{SN, o[0], o[1]}
				}
				var td1 = cloud.TableData{Header: header, Body: data1}
				Caption1 := "当前凭证具备的权限 (Permissions owned)"
				cloud.PrintTable(td1, Caption1)
				fmt.Println()
				for j, o := range obj2 {
					SN := strconv.Itoa(j + 1)
					data2[j] = []string{SN, o[0], o[1]}
				}
				var td2 = cloud.TableData{Header: header2, Body: data2}
				Caption2 := "当前凭证可以执行的操作 (Available actions)"
				cloud.PrintTable(td2, Caption2)
			}
		} else {
			log.Debugln(err)
		}
	}
}

func getCallerIdentity() string {
	request := sts.CreateGetCallerIdentityRequest()
	request.Scheme = "https"
	response, err := STSClient().GetCallerIdentity(request)
	util.HandleErr(err)
	accountArn := response.Arn
	var userName string
	if accountArn[len(accountArn)-4:] == "root" {
		userName = "root"
	} else {
		userName = strings.Split(accountArn, "/")[1]
	}
	log.Debugf("获得到当前凭证的用户名为 %s (The user name to get the current credentials is %s)", userName, userName)
	return userName
}

func listPoliciesForUser(userName string) ([][]string, error) {
	request := ram.CreateListPoliciesForUserRequest()
	request.Scheme = "https"
	request.UserName = userName
	response, err := RAMClient().ListPoliciesForUser(request)
	if err == nil {
		log.Debugf("成功获取 crossfire 用户的权限信息 (Successfully obtained permission information for crossfire user)")
	}
	var data [][]string
	for n, i := range response.Policies.Policy {
		SN := strconv.Itoa(n + 1)
		data = append(data, []string{SN, i.PolicyName, i.Description})
	}
	return data, err
}

func traversalPermissions() ([][]string, [][]string) {
	var obj1 [][]string
	var obj2 [][]string
	// 1. cf oss ls
	OSSCollector := &alioss.OSSCollector{}
	_, err1 := OSSCollector.ListBuckets()
	if err1 == nil {
		obj1 = append(obj1, []string{"AliyunOSSReadOnlyAccess", "只读访问对象存储服务(OSS)的权限"})
		obj2 = append(obj2, []string{"cf oss ls", "列出 OSS 资源"})
	} else {
		log.Traceln(err1.Error())
	}
	// 2. cf ecs ls
	request := ecs.CreateDescribeVpcsRequest()
	request.Scheme = "https"
	_, err2 := aliecs.ECSClient("cn-beijing").DescribeVpcs(request)
	if err2 == nil {
		obj1 = append(obj1, []string{"AliyunECSReadOnlyAccess", "只读访问云服务器服务(ECS)的权限"})
		obj2 = append(obj2, []string{"cf ecs ls", "列出 ECS 资源"})
	} else {
		log.Traceln(err2.Error())
	}
	// 3. cf ecs exec
	request3 := ecs.CreateInvokeCommandRequest()
	request3.Scheme = "https"
	request3.CommandId = "abcdefghijklmn"
	request3.InstanceId = &[]string{"abcdefghijklmn"}
	_, err3 := aliecs.ECSClient("cn-beijing").InvokeCommand(request3)
	if !strings.Contains(err3.Error(), "ErrorCode: Forbidden.RAM") {
		obj1 = append(obj1, []string{"AliyunECSAssistantFullAccess", "管理 ECS 云助手服务的权限"})
		obj2 = append(obj2, []string{"cf ecs exec", "在 ECS 上执行命令"})
	} else {
		log.Traceln(err3.Error())
	}
	// 4. cf rds ls
	_, err4 := alirds.DescribeDBInstances("cn-beijing", true, "all", "all")
	if err4 == nil {
		obj1 = append(obj1, []string{"AliyunRDSReadOnlyAccess", "只读访问云数据库服务(RDS)的权限"})
		obj2 = append(obj2, []string{"cf rds ls", "列出 RDS 资源"})
	} else {
		log.Traceln(err4.Error())
	}
	// 5. cf console
	request5 := ram.CreateDetachPolicyFromUserRequest()
	request5.Scheme = "https"
	request5.PolicyType = "System"
	request5.PolicyName = "test"
	request5.UserName = "test"
	_, err5 := RAMClient().DetachPolicyFromUser(request5)
	if !strings.Contains(err5.Error(), "ErrorCode: NoPermission") {
		obj1 = append(obj1, []string{"AliyunRAMFullAccess", "管理访问控制(RAM)的权限，即管理用户以及授权的权限"})
		obj2 = append(obj2, []string{"cf console", "接管控制台"})
	} else {
		log.Traceln(err5.Error())
	}
	return obj1, obj2
}

func appendData(action string, description string) [][]string {
	var actionList []string
	for _, o := range data2 {
		actionList = append(actionList, o[1])
	}
	sort.Strings(actionList)
	index := sort.SearchStrings(actionList, action)
	if index < len(actionList) && actionList[index] == action {
		log.Tracef("当前 data2 中已存在 %s (%s already exists in the current data2 array)", action, action)
	} else {
		data2 = append(data2, []string{strconv.Itoa(SN), action, description})
		SN = SN + 1
	}
	return data2
}
