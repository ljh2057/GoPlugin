def VERSION='v0.1.7'
pipeline {
    agent any  //在任何jenkins节点上都可运行
    stages {
        stage('Checkout'){
            steps{
            checkout([$class: 'GitSCM', branches: [[name: '*/master']], doGenerateSubmoduleConfigurations: false, extensions: [], submoduleCfg: [], userRemoteConfigs: [[credentialsId: '8b93f470-9b51-48fb-b44b-ed7bbaa963ee', url: 'https://github.com/ljh2057/GoPlugin']]])
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
