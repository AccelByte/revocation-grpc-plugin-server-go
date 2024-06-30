/* groovylint-disable DuplicateStringLiteral, Indentation, NestedBlockDepth */
library(
  identifier: 'jenkins-shared-library@master',
  retriever: modernSCM(
    [
      $class: 'GitSCMSource',
      remote: 'https://github.com/dhanarab/jenkins-pipeline-library.git'
    ]
  )
)

bitbucketHttpsCredentials = 'bitbucket-build-extend-https'
bitbucketCredentialsSsh = 'bitbucket-build-extend-ssh'

bitbucketPayload = null
bitbucketCommitHref = null

pipeline {
  agent none
  stages {
    stage('Prepare') {
      agent {
        label "extend-builder-ci"
      }
      steps {
        script {
          if (env.BITBUCKET_PAYLOAD) {
            bitbucketPayload = readJSON text: env.BITBUCKET_PAYLOAD
            if (bitbucketPayload.pullrequest) {
              bitbucketCommitHref = bitbucketPayload.pullrequest.source.commit.links.self.href
            }
          }
          if (bitbucketCommitHref) {
            bitbucket.setBuildStatus(
              bitbucketHttpsCredentials, bitbucketCommitHref, 'INPROGRESS', env.JOB_NAME,
              "${env.JOB_NAME}-${env.BUILD_NUMBER}", 'Jenkins', "${env.BUILD_URL}console")
          }
        }
      }
    }
    stage('Lint') {
      agent {
        label "extend-builder-ci"
      }
      stages {
        stage('Lint Commits') {
          when {
            expression {
              return env.BITBUCKET_PULL_REQUEST_LATEST_COMMIT_FROM_TARGET_BRANCH
            }
          }
          agent {
            docker {
              image 'commitlint/commitlint:19.3.1'
              args '--entrypoint='
              reuseNode true
            }
          }
          steps {
            sh "git config --add safe.directory '*'"
            sh "commitlint --color false --verbose --from ${env.BITBUCKET_PULL_REQUEST_LATEST_COMMIT_FROM_TARGET_BRANCH}"
          }
        }
        stage('Lint Code') {
          steps {
            sh 'make lint'
          }
        }
      }
    }
    stage('Build') {
      agent {
        label "extend-builder-ci"
      }
      steps {
        sh 'make build'
      }
    }
  }
  post {
    success {
      script {
        if (bitbucketCommitHref) {
          bitbucket.setBuildStatus(
            bitbucketHttpsCredentials, bitbucketCommitHref, 'SUCCESSFUL', env.JOB_NAME,
            "${env.JOB_NAME}-${env.BUILD_NUMBER}", 'Jenkins', "${env.BUILD_URL}console")
        }
      }
    }
    failure {
      script {
        if (bitbucketCommitHref) {
          bitbucket.setBuildStatus(
            bitbucketHttpsCredentials, bitbucketCommitHref, 'FAILED', env.JOB_NAME,
            "${env.JOB_NAME}-${env.BUILD_NUMBER}", 'Jenkins', "${env.BUILD_URL}console")
        }
      }
    }
  }
}
