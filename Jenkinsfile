pipeline {
    agent any

    environment {
        IMAGE_NAME = 'pudyn/api-testing'
        DOCKER_CREDENTIALS_ID = 'dockerhub'           // ID credentials Docker Hub di Jenkins
        KUBECONFIG_CREDENTIAL_ID = 'kubeconfig'       // ID credentials file kubeconfig di Jenkins
        DEPLOYMENT_NAME = 'api-testing-deployment'    // Nama deployment di K3s
        CONTAINER_NAME = 'api-testing'                // Nama container dalam deployment
        NAMESPACE = 'default'                         // Namespace Kubernetes
    }

    stages {
        stage('Checkout') {
            steps {
                checkout scm
            }
        }

        stage('Build & Tag Docker Image') {
            steps {
                script {
                    // Simpan commit hash ke env biar bisa digunakan antar stage
                    env.COMMIT_HASH = sh(script: "git rev-parse --short HEAD", returnStdout: true).trim()

                    // Build image dengan 2 tag
                    docker.build("${IMAGE_NAME}:latest")
                    docker.build("${IMAGE_NAME}:${COMMIT_HASH}")
                }
            }
        }

        stage('Push to Docker Hub') {
            steps {
                script {
                    docker.withRegistry('https://index.docker.io/v1/', DOCKER_CREDENTIALS_ID) {
                        docker.image("${IMAGE_NAME}:latest").push()
                        docker.image("${IMAGE_NAME}:${COMMIT_HASH}").push()
                    }
                }
            }
        }

        stage('Deploy to K3s') {
            steps {
                withCredentials([file(credentialsId: KUBECONFIG_CREDENTIAL_ID, variable: 'KUBECONFIG')]) {
                    sh '''
                        export KUBECONFIG=$KUBECONFIG
                        kubectl set image deployment/${DEPLOYMENT_NAME} \
                            ${CONTAINER_NAME}=${IMAGE_NAME}:${COMMIT_HASH} \
                            -n ${NAMESPACE}
                    '''
                }
            }
        }
    }

    post {
        failure {
            echo '❌ Build failed!'
        }
        success {
            echo '✅ CI/CD pipeline completed successfully.'
        }
    }
}