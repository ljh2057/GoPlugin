node {
    withEnv(["GOPATH=$WORKSPACE"]) {     // 设置stage运行时的环境变量
        env.PATH="${GOPATH}/bin:$PATH"
        stage('Checkout'){
            checkout([$class: 'GitSCM', 
            branches: [[name: '*/master']], 
            doGenerateSubmoduleConfigurations: false, 
            extensions: [], 
            submoduleCfg: [], 
            userRemoteConfigs: [[
                credentialsId: '8b93f470-9b51-48fb-b44b-ed7bbaa963ee', 
                url: 'https://github.com/ljh2057/GoPlugin']]])
        }
        stage('Build') {
                sh 'cd src/${PROJ_DIR}/; go test && CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -v -a -installsuffix cgo -o GoPlugin_arm . && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -a -installsuffix cgo -o GoPlugin_amd .'
        }
        stage('Test'){
                sh 'cd src/${PROJ_DIR}/; ./GoPlugin_amd'
        }
        stage('Deploy'){
                sh 'cd src/${PROJ_DIR}/;tar -zcvf GoPlugin_${VERSION}_arm64.tar.gz config.json GoPlugin_arm && tar -zcvf GoPlugin_${VERSION}_amd64.tar.gz config.json GoPlugin_amd'
        }
        // stage('Get code') {
        //     checkout([                      // git repo
        //         $class: 'GitSCM', 
        //         branches: [[name: '*/master']], 
        //         doGenerateSubmoduleConfigurations: false, 
        //         extensions: [[$class: 'RelativeTargetDirectory', relativeTargetDir: 'src/learningGo']], 
        //         submoduleCfg: [], 
        //         userRemoteConfigs: [[
        //             credentialsId: '3d33494a-ee1a-4540-a5e9-27f068a2b859', 
        //             url: 'git@github.com:peng0208/learningGo.git'
        //         ]]
        //     ])
        // }
        // stage('Build go proejct') {      
        //     sh 'cd ${PROJ_DIR}/daemon/example; go test && go build && go install'
        // }
        // stage('Deploy to test') {           // 部署测试环境
        //     input message: 'deploy to test ?', ok: 'De'
        //     echo 'docker run'
        // }
        // stage('Deploy to qa') {             // 部署预发布环境
        //     input message: 'deploy to qa ?', ok: 'OK!'
        //     echo 'docker run'
        // }
        // stage('Deploy to production') {     // 部署生产环境
        //     input message: 'deploy to production ?', ok: 'OK!'
        //     echo 'docker run'
        // }
    }
}
