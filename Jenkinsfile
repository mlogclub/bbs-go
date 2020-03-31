podTemplate(cloud:'kubernetes',containers: [
                            containerTemplate(name: 'ci-base',privileged:true,
                            image: '10.0.3.200:32382/ci-base', ttyEnabled: true,
                            command: 'cat'),
                            // containerTemplate(name: 'docker', image: 'docker',privileged:true, ttyEnabled: true, command: 'cat')
],volumes: [
      hostPathVolume(mountPath: '/var/run/docker.sock', hostPath: '/var/run/docker.sock'),
      hostPathVolume(mountPath: '/root/.m2', hostPath: '/data')      
  ])
{
    node(POD_LABEL) {
       checkout scm
       container('ci-base') {
                       def sdf = System.currentTimeMillis()  
                       def reg='10.0.3.200:32382/'
                       stage('Build image') {    
                           sh 'kubectl get nodes'
                           sh 'docker build server/ -t '+reg+'bbs-server:'+sdf
                           sh 'docker build site/ -t '+reg+'bbs-site:'+sdf
                         //  sh 'docker build admin/ -t '+reg+'bbs-admin:'+sdf
                       }
                        stage('push image') {    
                           sh 'docker push '+reg+'bbs-server:'+sdf
                           sh 'docker push '+reg+'bbs-site:'+sdf
                         //  sh 'docker push '+reg+'bbs-admin:'+sdf
                       }
                       stage('update deployment') { 
                           sh 'kubectl set image deployment bbs-server bbs-server='+reg+'bbs-server:'+sdf+' -n bbs'
                           sh 'kubectl set image deployment bbs-site bbs-site='+reg+'bbs-site:'+sdf+' -n bbs'
                          // sh 'kubectl set image deployment bbs-admin bbs-admin='+reg+'bbs-admin:'+sdf+' -n bbs'
                       }
       }
      
    }
}