//获取Git的ChangeLog             
@NonCPS
def getChanegLogs(){
    def changeLogs = ""
    def currentLogSets = currentBuild.changeSets
    for (int i = 0; i < currentLogSets.size(); i++){
        def entries = currentLogSets[i].items
        for(int j = entries.length - 1; j >= 0; j--){
            def entry = entries[j]
            changeLogs += "[${entry.author}]:${entry.msg}\n"
        }
    }
    if (changeLogs==""){
        changeLogs = "[no Commiter]"
    }
    return changeLogs 
}

def sendNotification() {
   git([url: 'git@gitlab.mobvista.com:QA/cd_tool_adv.git', credentialsId: 'bdb1da6b-e00f-47de-a7a9-e911b9f9169c', branch: 'master'])

   if (env.BRANCH_NAME == 'master') {
      sh "dist/after_test  ${currentBuild.currentResult},${params.JIRA_ISSUE_KEY},${params.username},${env.BUILD_URL}"
    } else {
      sh "dist/before_test  ${currentBuild.currentResult},${params.JIRA_ISSUE_KEY},${params.username},${env.BUILD_URL}"
    }
}

//获取pipeline启动原因
@NonCPS
def getBuildCauser() {
    def causer=""
    UserCause = currentBuild.rawBuild.getCause(Cause.UserIdCause).toString()
    TimerCause = currentBuild.rawBuild.getCause(hudson.triggers.TimerTrigger.TimerTriggerCause).toString()
    if (UserCause != "null") {
        causer=currentBuild.rawBuild.getCause(Cause.UserIdCause).getUserName()
    }else if (TimerCause != "null") {
        causer="TimerTrigger"
    }else{
        causer="Branch Indexing"
    }
    return "[causer]: "+causer+"\n"
}

//每天早上7点定时启动master分支跑CI
String cron_string = BRANCH_NAME == "master" ? "H 7 * * *" : ""
pipeline {
  agent {
    node {
      label 'MSystem-Prebuild'
    }
  }

  triggers { cron(cron_string) }

  parameters {
    string(name: 'JIRA_ISSUE_KEY',  defaultValue: 'JIRA_ISSUE_KEY_defaultValue',description: 'jira issue key ')
    string(name: 'username', defaultValue: 'username_defaultValue', description: 'original user who started job')
  }

  stages {

    stage('CodeAnalysis') {
          parallel {
            stage('SonarAnalysis') {
              steps {
                echo "${env.BRANCH_NAME}"
                build job:'M-Adnet-Sonar', parameters: [ string(name: 'BRANCH_NAME', value: String.valueOf("${env.BRANCH_NAME}")) ]
              }
              post {
                success {
                    script {
                        echo "Sonar analysis success"
                    }
                }
                failure {
                    script {
                        echo "Sonar analysis Failed"
                        stopStage = "Sonar Analysis"
                    }
                }
                unstable {
                    script {
                        echo "Sonar analysis Failed"
                        stopStage = "Sonar Analysis"
                    }
                }
              }
            }
            stage('Gometalinter') {
              steps {
                echo "${env.BRANCH_NAME}"
                build job:'M-Adnet-Gometalinter', parameters: [ string(name: 'BRANCH_NAME', value: String.valueOf("${env.BRANCH_NAME}")) ]
              }
              post {
                success {
                    script {
                        echo "gometalinter analysis success"
                    }
                }
                failure {
                    script {
                        echo "gometalinter analysis Failed"
                        stopStage = "Gometalinter Analysis"
                    }
                }
                unstable {
                    script {
                        echo "gometalinter analysis Failed"
                        stopStage = "Gometalinter Analysis"
                    }
                }
              }
            }
          }
        }

    stage('AdnetPreBuild') {
      steps {
        lock('AdnetPreBuild'){
          echo "${env.BRANCH_NAME}"
          build job:'M-Adnet-PreBuild', parameters: [
            string(name: 'BRANCH_NAME', value: String.valueOf("${env.BRANCH_NAME}")),
            string(name: 'GIT_COMMIT', value: String.valueOf("${env.GIT_COMMIT}"))]
          }
      }
    }

    stage('AdnetWhiteBox') {
      steps {
        lock('AdnetWhiteBox'){
          echo "${env.BRANCH_NAME}"
          build job:'M-Adnet-whiteBox-Test', parameters: [
            string(name: 'BRANCH_NAME', value: String.valueOf("${env.BRANCH_NAME}")),
            string(name: 'GIT_COMMIT', value: String.valueOf("${env.GIT_COMMIT}"))]
          }
      }
    }

    stage('AdnetBlackBox') {
      when {
          branch 'master'
      }
      steps {
        lock('AdnetHbBlackBox'){
          script {
            if (env.BRANCH_NAME == 'master') {
              build job:'M-Adnet-HB-BlackBox-Test-K8s', parameters: [
                string(name: 'BRANCH_NAME', value: String.valueOf("${env.BRANCH_NAME}")),
                string(name: 'GIT_COMMIT', value: String.valueOf("${env.GIT_COMMIT}"))]
            } else {
              //build job:'M-Adnet-BlackBox-Test-Basic-New', parameters: [string(name: 'BRANCH_NAME', value: String.valueOf("${env.BRANCH_NAME}"))],
              // string(name: 'GIT_COMMIT', value: String.valueOf("${env.GIT_COMMIT}"))]
              echo "skip black-box test"
            }
          }
        }
      }
    }
  }

  post {
    //失败时通知
    failure{
      echo 'task failed'
      dingTalk accessToken:'9cc977f872ec2187b7216f623da85069d3dbc0db81000beaac6fe49fbd8e63bf',jenkinsUrl:'http://ci.mobvista.com/bluedingding',
      message:getChanegLogs(),
      notifyPeople:'',
      imageUrl:'http://cdn-adn.rayjump.com/cdn-adn/v2/portal/19/01/09/18/31/5c35cd9fd5789.png'
    }
    unstable{
      echo 'task unstable'
      dingTalk accessToken:'9cc977f872ec2187b7216f623da85069d3dbc0db81000beaac6fe49fbd8e63bf',jenkinsUrl:'http://ci.mobvista.com/bluedingding',
      message:getChanegLogs(),
      notifyPeople:'',
      imageUrl:'http://cdn-adn.rayjump.com/cdn-adn/v2/portal/19/01/09/18/31/5c35cd9fd5789.png'
    }
    //成功时通知
    //test 57e4638d435c5e3552ab663bf62c140a65640c2dc9ae10ad0ba6011ef65da79b
    //online 9cc977f872ec2187b7216f623da85069d3dbc0db81000beaac6fe49fbd8e63bf
    success{
      echo 'task success'
      dingTalk accessToken:'9cc977f872ec2187b7216f623da85069d3dbc0db81000beaac6fe49fbd8e63bf',jenkinsUrl:'http://ci.mobvista.com/bluedingding',
      message:getBuildCauser()+getChanegLogs(),
      notifyPeople:'',
      imageUrl:'http://cdn-adn.rayjump.com/cdn-adn/v2/portal/19/01/09/18/30/5c35cd5065b3f.png'
    }

    always {
      sendNotification()
    }
  }
}
