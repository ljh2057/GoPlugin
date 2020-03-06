pipeline {
    // Lets Jenkins use Docker for us later.
    agent any    

    // If anything fails, the whole Pipeline stops.
    stages {
        stage('Build & Test') {   
            // Use golang.
            agent { docker { image 'golang' } }

            steps {                                           
                // Create our project directory.
                sh 'cd ${GOPATH}/src'
                sh 'mkdir -p ${GOPATH}/src/GoPlugin'

                // Copy all files in our Jenkins workspace to our project directory.                
                sh 'cp -r ${WORKSPACE}/* ${GOPATH}/src/GoPlugin'

                // Copy all files in our "vendor" folder to our "src" folder.
                sh 'cp -r ${WORKSPACE}/vendor/* ${GOPATH}/src'

                // Build the app.
                sh 'CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -v -a -installsuffix cgo -o GoPlugin_arm . && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -a -installsuffix cgo -o GoPlugin_amd .'               
            }            
        }

        stage('Test') {
            // Use golang.
            agent { docker { image 'golang' } }

            steps {                 
                // Create our project directory.
                sh 'cd ${GOPATH}/src'
                sh 'mkdir -p ${GOPATH}/src/GoPlugin'

                // Copy all files in our Jenkins workspace to our project directory.                
                sh 'cp -r ${WORKSPACE}/* ${GOPATH}/src/GoPlugin'

                // Copy all files in our "vendor" folder to our "src" folder.
                sh 'cp -r ${WORKSPACE}/vendor/* ${GOPATH}/src'

                // Remove cached test results.
                sh 'go clean -cache'

                // Run Unit Tests.
                sh 'go test ./... -v -short'            
            }
        }      

    }

    post {
        always {
            // Clean up our workspace.
            deleteDir()
        }
    }
}  
