def VERSION='v0.1.7'
pipeline {
    agent any
    
    stages {
        stage('Checkout'){
            steps{
            checkout([$class: 'GitSCM', 
            branches: [[name: '*/master']], 
            doGenerateSubmoduleConfigurations: false, 
            extensions: [], 
            submoduleCfg: [], 
            userRemoteConfigs: [[
                credentialsId: '8b93f470-9b51-48fb-b44b-ed7bbaa963ee', 
                url: 'https://github.com/ljh2057/GoPlugin']]])
            }
        }
        stage('Build') {
            steps{
                sh 'cd src/GoPlugin/; go test && CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -v -a -installsuffix cgo -o GoPlugin_arm . && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -a -installsuffix cgo -o GoPlugin_amd .'
            }
        }
        stage('Test'){
            steps{
                sh 'cd src/GoPlugin/; ./GoPlugin_amd'
            }
        }
        stage('Deploy'){
            steps{
                sh 'cd src/GoPlugin/;tar -zcvf GoPlugin_${VERSION}_arm64.tar.gz config.json GoPlugin_arm && tar -zcvf GoPlugin_${VERSION}_amd64.tar.gz config.json GoPlugin_amd'
            }
        }
        // stage('Build') {    // buid 阶段
        //     steps {        //build 步骤          

        //     }
        // }
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
