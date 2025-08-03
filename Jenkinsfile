pipeline {
    agent any

    environment {
        DOCKER_IMAGE = 'transit-backend'
        DOCKER_TAG   = "${BUILD_NUMBER}"
    }

    stages {
        stage('Checkout') {
            steps {
                echo 'Checking out source…'
                checkout scm
            }
        }

        stage('Test') {
            steps {
                script {
                    if (isUnix()) {
                        sh '''
                            set -e
                            go mod download
                            go test ./...
                        '''
                    } else {
                        bat '''
                            go mod download
                            go test ./...
                        '''
                    }
                }
            }
        }

        stage('Build') {
            steps {
                script {
                    if (isUnix()) {
                        sh '''
                            set -e
                            mkdir -p bin
                            go build -o bin/main ./cmd/main.go
                        '''
                    } else {
                        bat '''
                            if not exist bin mkdir bin
                            go build -o bin\\main.exe ./cmd/main.go
                        '''
                    }
                }
            }
        }

        stage('Docker Build') {
            steps {
                script {
                    if (isUnix()) {
                        sh '''
                            cd docker
                            docker build \
                                -t ${DOCKER_IMAGE}:${DOCKER_TAG} \
                                -t ${DOCKER_IMAGE}:latest \
                                -f Dockerfile \
                                ..
                        '''
                    } else {
                        bat '''
                            cd docker
                            docker build ^
                                -t %DOCKER_IMAGE%:%DOCKER_TAG% ^
                                -t %DOCKER_IMAGE%:latest ^
                                -f Dockerfile ^
                                ..
                        '''
                    }
                }
            }
        }

        stage('Deploy to Dev') {
            when { branch 'dev' }
            steps {
                script {
                    if (isUnix()) {
                        sh '''
                            cd docker
                            docker-compose -f docker-compose.dev.yaml down || true
                            docker-compose -f docker-compose.dev.yaml up -d
                        '''
                    } else {
                        bat '''
                            cd docker
                            docker-compose -f docker-compose.dev.yaml down || echo "No existing containers"
                            docker-compose -f docker-compose.dev.yaml up -d
                        '''
                    }
                }
            }
        }

        stage('Deploy to Production') {
            when { branch 'main' }
            steps {
                input message: 'Deploy to production?'
                script {
                    if (isUnix()) {
                        sh '''
                            # Start infra stacks
                            cd docker/pg     && docker-compose up -d
                            cd ../redis      && docker-compose up -d
                            cd ../kafka      && docker-compose up -d
                            sleep 30
                            # Start main stack
                            cd ..
                            docker-compose up -d
                        '''
                    } else {
                        bat '''
                            cd docker\\pg     && docker-compose up -d
                            cd ..\\redis      && docker-compose up -d
                            cd ..\\kafka      && docker-compose up -d
                            timeout /t 30
                            cd ..
                            docker-compose up -d
                        '''
                    }
                }
            }
        }
    }

    post {
        always {
            echo 'Cleaning up Docker artefacts…'
            script {
                if (isUnix()) {
                    sh 'docker system prune -f || true'
                } else {
                    bat 'docker system prune -f || echo "No dangling objects"'
                }
            }
        }
        success { echo 'Pipeline completed successfully!' }
        failure { echo 'Pipeline failed!' }
    }
}
