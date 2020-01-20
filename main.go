package main

import (
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/eclipse/paho.mqtt.golang"
	"github.com/tidwall/gjson"
	"github.com/urfave/cli"
)

type BaseInfo struct {
	NetworkUrl string
	TimeApi string
}
type MapInfo struct {
	Id string
	Url string
	Attributes gjson.Result
}
type CompassInfo struct {
	Services gjson.Result
}
type UosInfo struct {
	path string
	url string
	SimulationCar string
	RealCar string
	Attributes gjson.Result
}
type CertInfo struct {
	Dir string
	Path string
}
type OutputInfo struct {
	Path string
}
type Config struct {
	baseInfo BaseInfo
	mapInfo  MapInfo
	compassInfo CompassInfo
	uosInfo UosInfo
	certInfo CertInfo
	outputInfo OutputInfo
}

func main() {
	app := &cli.App{
		Name:"Compass Detection Tool",
		Version:"v0.0.2",
		Action: func(c *cli.Context) error {
			DetectMain()
			return nil
		},
	}
	err:=app.Run(os.Args)
	if err!=nil{
		log.Fatal(err)
	}
}

//初始化Config
func InitConfig(filePath string) Config {
	configFile:=string(ReadFile(filePath))
	config:=Config{
		BaseInfo{
			NetworkUrl: gjson.Get(configFile,"baseInfo.networkUrl").String(),
			TimeApi:    gjson.Get(configFile,"baseInfo.timeApi").String(),
		},
		MapInfo{
			Id:         gjson.Get(configFile,"mapInfo.id").String(),
			Url:        gjson.Get(configFile,"mapInfo.url").String(),
			Attributes: gjson.Get(configFile,"mapInfo.attributes"),
		},
		CompassInfo{
			Services: gjson.Get(configFile,"compassInfo.services"),
		},
		UosInfo{
			path:          gjson.Get(configFile,"uosInfo.path").String(),
			url:           gjson.Get(configFile,"uosInfo.url").String(),
			SimulationCar: gjson.Get(configFile,"uosInfo.simulationCar").String(),
			RealCar:       gjson.Get(configFile,"uosInfo.realCar").String(),
			Attributes:    gjson.Get(configFile,"uosInfo.attributes"),
		},
		CertInfo{
			Dir:  gjson.Get(configFile,"certInfo.dir").String(),
			Path: gjson.Get(configFile,"certInfo.path").String(),
		},
		OutputInfo{
			Path:gjson.Get(configFile,"outputInfo.path").String(),
		},
	}
	return config
}

func WriteBytes(filePath string, b []byte) (int, error) {
	os.MkdirAll(path.Dir(filePath), os.ModePerm)
	fw, err := os.Create(filePath)
	if err != nil {
		return 0, err
	}
	defer fw.Close()
	return fw.Write(b)
}

//读取文件
func ReadFile(filePath string) []byte{
	config,err:=ioutil.ReadFile(filePath)
	if err!=nil{
		return nil
	}
	return config
}
//判断文件夹是否存在
func Exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}
//更新检测结果
func AddProblems(problems map[int] string,index *int,info string)  {
	problems[*index]=info
	*index+=1
}

//确认UOS移动网络连接正常
func DetectNetworkConnection(url string) (bool,string) {
	info:="网络检测完毕，连接正常！"
	flag:=true
	_,err := net.DialTimeout("tcp",url,2*time.Second)
	if err != nil {
		flag=false
		info="网络连接异常，检查网络连接！"
	}
	return flag,info
}
//确认UOS本地时钟正确
func DetectTime(config Config) (bool,string)  {
	info,TimeApi:="车端时间校正完毕，正常！",config.baseInfo.TimeApi
	flag:=true
	resp,err:=http.Get(TimeApi)
	if err!=nil{
		info=""
		return false,info
	}
	defer resp.Body.Close()
	//解析body
	body,_:=ioutil.ReadAll(resp.Body)
	bj_time_str:=gjson.Get(string(body),"data.t").String()
	bj_time,err:=strconv.Atoi(bj_time_str)
	bj_time=bj_time/1000
	current_time:=int(time.Now().Unix())
	time_err:=math.Abs(float64(current_time-bj_time))
	if time_err>200{
		flag=false
		info="车端时间异常，请校对车端时间设置！"
	}
	return flag,info
}
//检测服务是否启动
func DetectService(service string)(bool,string){
	info:="compass服务检测完毕，所有服务均已正常启动!"
	flag:=true
	cmdStr:=fmt.Sprintf(`ps -ef | grep "%s" | grep -v grep`,service)
	cmd:=exec.Command("sh","-c",cmdStr)
	output,_:=cmd.CombinedOutput()
	if string(output)==""{
		info="compass服务检测异常，"+service+"服务未启动！"
		flag=false
	}
	return flag,info
}
//1.1.1确认compass服务是否正常启动
func DetectCompass(config Config) (bool,string){
	info,CompassServices:="compass启动正常！",config.compassInfo.Services
	flag:=true
	for _,service:=range CompassServices.Array(){
		flag,info=DetectService(service.String())
		if !flag{
			break
		}
	}
	return flag,info
}
//1.1.2查看compass地图是否下载正确
func DetectCompassMap(config Config) (bool,string){
	mapId:=config.mapInfo.Id
	UosConfigPath:=config.uosInfo.path
	//CarMapUrl:=config.mapInfo.Url
	mapAttributes:=config.mapInfo.Attributes.Array()
	info:="车端地图检测完毕，正常！"
	flag:=true
	//cmd := exec.Command("curl","-s",CarMapUrl)
	//output,_:=cmd.CombinedOutput()

	output:=ReadFile("compassMap.json")
	res:=string(output)
	isConnected:=gjson.Get(res,mapAttributes[0].String()).Bool()
	if isConnected{
		flag=false
		info="车端地图检测异常，地图未下载成功！"
	}else {
		localMapId:=gjson.Get(res,mapAttributes[1].String()).String()
		if localMapId==mapId{
			UosConfigPathStr:=string(ReadFile(UosConfigPath))
			uosAttributes:=config.uosInfo.Attributes.Array()
			serverMap:=gjson.Get(UosConfigPathStr,uosAttributes[0].String()).String()
			//isConnected,_:=DetectNetworkConnection(serverMap)
			//if !isConnected{
			//	flag=false
			//	info="无法连接到"+serverMap
			//	return flag,info
			//}
			map_url:=serverMap+"/maps/"+mapId+"/"
			cmd := exec.Command("curl","-s",map_url)
			output,_:=cmd.CombinedOutput()
			versions:=gjson.Get(string(output),mapAttributes[2].String()).Array()
			version_id:=gjson.Get(versions[len(versions)-1].String(),mapAttributes[3].String()).String()
			local_version_id:=gjson.Get(res,mapAttributes[4].String()).String()
			if version_id!=local_version_id{
				flag=false
				info="车端地图检测异常，地图版本与云端不一致，请更新!"
			}
		}else {
			flag=false
			info="车端地图检测异常，车端地图id与云端不一致，请检查！"
		}
	}
	return flag,info
}
//1.1.3查看 Compass 中 UOS 配置正确
func DetectUosPath(config Config) (bool,string) {
	info,UosConfigPath:="UOS配置文件路径检测完毕，正常！",config.uosInfo.path
	flag:=true
	UosConfigPathStr:=string(ReadFile(UosConfigPath))
	uosAttributes:=config.uosInfo.Attributes.Array()
	UosPath:=gjson.Get(UosConfigPathStr,uosAttributes[1].String()).String()
	if !Exists(UosPath){
		flag=false
		info="UOS检测异常，Uos配置"+UosPath+"路径不存在！"
	}
	return flag,info
}
//1.1.4查看 compass 证书是否过期
func DetectCert(config Config)(bool,string){
	//var info string
	info,cert_path:="车端证书有效期检测完毕，证书有效！",config.certInfo.Path
	flag:=true
	certPEMBlock:= ReadFile(cert_path)
	cert,_:= pem.Decode(certPEMBlock)
	if cert != nil {
		x509Cert, err := x509.ParseCertificate(cert.Bytes)
		if err != nil {
			return false,""
		}
		NotBefore,NotAfter:=x509Cert.NotBefore.Format("2006-01-02 15:04"), x509Cert.NotAfter.Format("2006-01-02 15:04")
		CurrentTime:=time.Now().Format("2006-01-02 15:04:05")
		if NotBefore>CurrentTime || CurrentTime>NotAfter{
			flag=false
			info="车端证书有效期检测异常，证书已过期!"
		}
	}
	return flag,info
}
//1.2.1检查 UOS 的配置是否正确

func DetectUosConfig(config Config) (bool,string)  {
	UosConfigPath,SimulationCarPath,RealCarPath,UosUrl:=config.uosInfo.path,config.uosInfo.SimulationCar,config.uosInfo.RealCar,config.uosInfo.url
	info:=""
	mapRoot:="/etc/"
	flag:=true
	var err error
	MOD_uos_config:=string(ReadFile(SimulationCarPath))
	if !Exists(SimulationCarPath){
		isExist,_:=DetectUosPath(config)
		UosConfigPathStr:=string(ReadFile(UosConfigPath))
		uosAttributes:=config.uosInfo.Attributes.Array()
		UosPath:=gjson.Get(UosConfigPathStr,uosAttributes[1].String()).String()
		if isExist{
			Real_car_path:=UosPath+RealCarPath
			Real_car_path_Str:=string(ReadFile(Real_car_path))
			MOD_uos_config=gjson.Get(Real_car_path_Str,uosAttributes[6].String()).String()
			mapRoot=UosPath+"/"
		}
	}
	flag,info=DetectVnameMap(err,info,MOD_uos_config,mapRoot,UosUrl,config)
	return flag,info
}
//1.2.1 检测车辆运行模式、名字、地图
func DetectVnameMap(err error,info string,MOD_uos_config string,mapRoot string,UosUrl string,config Config)(bool,string)  {
	flag:=true
	info="UOS 配置检测完毕，相关配置正确！"
	uosAttributes:=config.uosInfo.Attributes.Array()
	if err!=nil{
		return false,""
	}else {
		run_scene:=gjson.Get(MOD_uos_config,uosAttributes[7].String()).String()
		if run_scene=="real.compass"{
			vehicle_name_config:=gjson.Get(MOD_uos_config,uosAttributes[8].String()).String()

			//cmd := exec.Command("curl","-s",UosUrl)
			//output,_:=cmd.CombinedOutput()

			output:=ReadFile("vehicle.json")
			vehicle_name_true:=gjson.Get(string(output),uosAttributes[9].String()).String()
			mapPath:=mapRoot+gjson.Get(MOD_uos_config,uosAttributes[10].String()).String()
			if vehicle_name_config==vehicle_name_true{
				if !Exists(mapPath){
					flag=false
					info="UOS 配置检测异常，地图文件不存在!"
				}
			}else {
				flag=false
				info="UOS 配置检测异常，车辆名称错误!"
			}
		}else {
			flag=false
			info="UOS 配置检测异常，运行模式错误!"
		}
	}
	return flag,info
}
//1.3.1 检查车端是否连接了正确的云端
func DetectVehicleConnectCloud(config Config)(bool,string)  {
	info,UosConfigPath:="",config.uosInfo.path
	flag:=true
	isExist,_:=DetectUosPath(config)

	if isExist{
		UosConfigPathStr:=string(ReadFile(UosConfigPath))
		uosAttributes:=config.uosInfo.Attributes.Array()
		serverCloud:=gjson.Get(UosConfigPathStr,uosAttributes[2].String())
		mqttBrokerId:=gjson.Get(UosConfigPathStr,uosAttributes[5].String())

		if serverCloud.Exists()&&mqttBrokerId.Exists(){
			serverList:=strings.Split(serverCloud.String(),":")
			isConnected,_:=DetectNetworkConnection(serverList[0]+":80")
			if !isConnected{
				flag=false
				info="车云连接检测异常，车端无法连接到云端！"
			} else{
				isConnected,_=DetectNetworkConnection(serverCloud.String())
				if isConnected {
					mqtt_username,mqtt_password:=gjson.Get(UosConfigPathStr,uosAttributes[3].String()).String(),gjson.Get(UosConfigPathStr,uosAttributes[4].String()).String()
					flag,info=DetectMqtt(serverCloud.String(),mqtt_username,mqtt_password,"test")
				} else{
					flag=false
					info="车云连接检测异常，云端端口错误，无法连接！"
				}
			}
		}else {
			flag=false
			info="车云连接检测异常，车端配置文件中server.cloud 或 mqtt.broker_id 未正确配置！"
		}
	}
	return flag,info
}
//1.3.1 判断 MQTT 能否订阅
func DetectMqtt(server string,uname string,upwd string,topic string)  (bool,string){
	info:="车云连接检测完毕，正常！"
	flag:=true
	//connect mqtt-server and set clientID
	opts := mqtt.NewClientOptions().AddBroker(server).SetClientID("test")
	opts.SetUsername(uname)
	opts.SetPassword(upwd)
	opts.SetKeepAlive(2 * time.Second)
	//create object
	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		flag=false
		info="车云连接检测异常，mqtt连接失败！"
	}
	//subscribe topic
	if token := c.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
		flag=false
		info="车云连接检测异常，mqtt订阅失败！"
	}
	//publish topic
	if token := c.Publish(topic, 0, false,"test"); token.Wait() && token.Error() != nil {
		flag=false
		info="车云连接检测异常，mqtt发布失败！"
	}
	c.Disconnect(250)
	time.Sleep(1 * time.Second)
	return flag,info
}
func DetectMain() error {
	problems:=make(map[int] string)
	index:=0
	config:=InitConfig("Config.json")
	Network_statue,info:=DetectNetworkConnection(config.baseInfo.NetworkUrl)
	AddProblems(problems,&index,info)
	if Network_statue{
		//DetectTime
		_,info=DetectTime(config)
		AddProblems(problems,&index,info)
		//DetectCompass
		_,info=DetectCompass(config)
		AddProblems(problems,&index,info)
		//DetectUosPath
		isExist,info:=DetectUosPath(config)
		AddProblems(problems,&index,info)
		if isExist{
			//DetectCompassMap
			_,info=DetectCompassMap(config)
			AddProblems(problems,&index,info)
		}
		if Exists(config.certInfo.Dir){
			//DetectCert
			_,info=DetectCert(config)
			AddProblems(problems,&index,info)
		}
		//DetectUosConfig
		_,info=DetectUosConfig(config)
		AddProblems(problems,&index,info)
		//DetectVehicleConnectCloud
		_,info=DetectVehicleConnectCloud(config)
		AddProblems(problems,&index,info)
	}
	res, err := json.MarshalIndent(problems, "", "      ")
	if err != nil {
		return err
	}
	//保存到json
	_,err=WriteBytes(config.outputInfo.Path,res)
	if err != nil {
		return err
	}
	var keys []int
	for k :=range problems{
		keys=append(keys,k)
	}
	sort.Ints(keys)
	for _,k:=range keys{
		fmt.Println(k,problems[k])
	}
	fmt.Println("检测结果已保存到"+config.outputInfo.Path)
	return nil
}