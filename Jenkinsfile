pipeline {
    agent any
    
    environment {
        DOCKER_IMAGE = 'transit-backend'
        DOCKER_TAG = "${BUILD_NUMBER}"
    }
    
    stages {
        stage('Checkout') {
            steps {
                echo 'Checking out code...'
                checkout scm
            }
        }
        
        stage('Test') {
            steps {
                echo 'Running tests...'
                sh '''
                    go mod download
                    go test ./...
                '''
            }
        }
        
        stage('Build') {
            steps {
                echo 'Building application...'
                sh '''
                    go build -o bin/main ./cmd/main.go
                '''
            }
        }
        
        stage('Docker Build') {
            steps {
                echo 'Building Docker image...'
                sh '''
                    cd docker
                    docker build -t ${DOCKER_IMAGE}:${DOCKER_TAG} -t ${DOCKER_IMAGE}:latest -f Dockerfile ..
                '''
            }
        }
        
        stage('Deploy to Dev') {
            when {
                branch 'dev'
            }
            steps {
                echo 'Deploying to development...'
                sh '''
                    cd docker
                    docker-compose -f docker-compose.dev.yaml down || true
                    docker-compose -f docker-compose.dev.yaml up -d
                '''
            }
        }
        
        stage('Deploy to Production') {
            when {
                branch 'main'
            }
            steps {
                input 'Deploy to production?'
                echo 'Deploying to production...'
                sh '''
                    # Start infrastructure services first
                    cd docker/pg && docker-compose up -d
                    cd ../redis && docker-compose up -d
                    cd ../kafka && docker-compose up -d
                    
                    # Wait a bit for services to start
                    sleep 30
                    
                    # Start main application
                    cd .. && docker-compose up -d
                '''
            }
        }
    }
    
    post {
        always {
            echo 'Cleaning up...'
            sh 'docker system prune -f'
        }
        success {
            echo 'Pipeline completed successfully!'
        }
        failure {
            echo 'Pipeline failed!'
        }
    }
}