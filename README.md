### 本项目为 Compass 端自检工具

#### 已实现功能：

参照 《WIP: 云脑日常运维操作手册》

1. 确认 UOS 本地时钟正确

2. 确认 UOS 移动网络连接正常

3. 车端 Compass 检测

     - 确认 compass 服务是否正常启动
     - 查看 compass 地图是否下载正确
     - 查看 compass 中 UOS 配置正确
     - 查看 compass 证书是否过期

4. UOS 检测

     - 检查 uos 的配置是否正确
     - 检测运行模式是否正确
     - 检测车辆名字是否正确
     - 检测地图是否存在

5. 检测 Compass 车云通讯通道的可用性

     - 检查车端是否连接了正确的云端

     - 判断 mosquitto 能否订阅

#### 使用方法：

在车端终端执行以下命令后会自动创建 CompassChecker 文件夹并进行自检。

```bash
curl -sL https://github.com/ljh2057/GoPlugin/releases/download/v0.0.2/deploy.sh | sh
```

CompassChecker 文件夹目录如下：

```
CompassChecker/
├── config.json
├── GoPlugin
├── README.md
└── result.json
```

其中 config.json 文件包含相关配置信息，GoPlugin 为检测工具，result.json 保存工具检测结果。
config.json 文件参考如下，可根据需要修改该文件：如
	修改 outputInfo.path 可以指定输出结果文件保存位置。 
	通过添加 compassInfo.services 可添加多个服务进行检测。
	attributes 表示工具中使用到的相关 json 文件中的属性，注：如属性中含有"." (例如 "server.map": "http://10.0.165.2:9090" )  需要使用   "server\\.map" 对 "." 进行解析。

```json
{
  "baseInfo": {
    "networkUrl": "www.baidu.com:80",
    "timeApi": "http://api.m.taobao.com/rest/api3.do?api=mtop.common.getTimestamp"
  },
  "mapInfo": {
    "id": "20181005",
    "url" : "localhost:8000/apis/v2/maps/current",
    "attributes": [
      "error_code",
      "data.id",
      "versions",
      "version",
      "data.version"
    ]
  },
  "compassInfo" :{
    "services": [
      "/usr/sbin/sshd -D"
    ]
  },
  "uosInfo": {
    "path" : "/opt/compass/vehicle_backend/web_service/config.json",
    "simulationCar" : "/etc/uos_config.json",
    "realCar" : "/uos_common.json",
    "url" : "localhost:8000/apis/v2/vehicle",
    "attributes": [
      "server\\.map",
      "uos\\.path",
      "server\\.cloud",
      "mqtt\\.username",
      "mqtt\\.password",
      "mqtt\\.broker_id",
      "_MOD_uos_config",
      "run_scene",
      "vehicle_name",
      "data.vin",
      "roadmap_fname"
    ]
  },
  "certInfo": {
    "dir" : "/opt/certs/",
    "path" : "/opt/certs/device_cert.crt"
  },
  "outputInfo":{
    "path": "result.json"
  }
}
```

