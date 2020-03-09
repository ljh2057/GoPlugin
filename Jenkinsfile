pipeline {
    agent any  //在任何jenkins节点上都可运行
    tools{
        go 'go-1.14'
    }
    environment {
        APP_NAME = 'GoPlugin'
        VERSION = "v0.1.7"
    }
    stages {
        stage('Init'){
            steps{
                script{
                    sh "echo ${env.WORKSPACE}"
                    sh "ls ${env.WORKSPACE}"
                }
            }
        }
      
        stage('Build') {    // buid 阶段
            steps {        //build 步骤
                script{
                    //sh "go env -w GO111MODULE=on;go env -w GOPROXY=https://goproxy.cn,direct"
                    sh "go get -u github.com/gpmgo/gopm"
                    sh "gopm get -v"
                    sh 'go test && CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -v -a -installsuffix cgo -o GoPlugin_arm . && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -a -installsuffix cgo -o GoPlugin_amd .'
                    sh 'ls -l'
                }
            }
        }
        
        stage('Deploy'){
            steps{
                script{
                    sh 'tar -zcvf GoPlugin_$VERSION_arm64.tar.gz config.json GoPlugin_arm && tar -zcvf GoPlugin_$VERSION_amd64.tar.gz config.json GoPlugin_amd'
                    sh 'ls -l'
                }
            }
        }
    }
}
