pipeline {
    agent any  //在任何jenkins节点上都可运行
    tools{
        go 'go-1.14'
    }
    environment {
        APP_NAME = 'GoPlugin'
    }
    stages {
        stage('Init'){
            steps{
                script{
                    sh "echo ${env.WORKSPACE}"
                    sh "ls ${env.WORKSPACE}"
                    //def root = tool name: 'go-1.14', type: 'go'
                    //withEnv(["GOPATH=${env.WORKSPACE}/go", "GOROOT=${root}", "GOBIN=${root}/bin", "PATH+GO=${root}/bin"]) {
                    //    sh "mkdir -p ${env.WORKSPACE}/go/src"
                    //}
                }
            }
        }
      
        stage('Build') {    // buid 阶段
            steps {        //build 步骤
                sh "pwd"
                sh "go env -w GO111MODULE=on;go env -w GOPROXY=https://goproxy.cn,direct"
                
                sh 'go get -v && go test && CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -v -a -installsuffix cgo -o GoPlugin_arm . && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -a -installsuffix cgo -o GoPlugin_amd .'
                sh 'ls -l'
            }
        }
        // stage('Test') {     // test 阶段
        //     steps {
        //         // 
        //     }
        // }
        // stage('Deploy') {   // 部署 阶段
        //     steps {
        //         // 
        //     }
        // }
    }
}
