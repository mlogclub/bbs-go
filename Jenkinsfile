podTemplate(cloud:'kubernetes',containers: [
                            containerTemplate(name: 'maven',privileged:true,
                            image: '10.0.3.200:32382/maven:3.0', ttyEnabled: true,
                            command: 'cat'),
                            // containerTemplate(name: 'docker', image: 'docker',privileged:true, ttyEnabled: true, command: 'cat')
],volumes: [
      hostPathVolume(mountPath: '/var/run/docker.sock', hostPath: '/var/run/docker.sock'),
      hostPathVolume(mountPath: '/root/.m2', hostPath: '/data')      
  ])
{
    node(POD_LABEL) {
       checkout scm
       container('maven') {
                       stage('Build a Maven project xbcapi') {
                           def sdf = System.currentTimeMillis()  
                           sh 'cd xbcapi && ls && mvn  clean && mvn package'
                           sh 'kubectl get nodes'
                           sh 'docker build xbcapi/ -t 10.0.3.200:32382/xbcapi:'+sdf
                           sh 'docker push 10.0.3.200:32382/xbcapi:'+sdf
                           sh 'docker build mediaupload/ -t 10.0.3.200:32382/mediaupload:'+sdf
                           sh 'docker push 10.0.3.200:32382/mediaupload:'+sdf
                           sh 'kubectl set image deployment xbcapi xbcapi=10.0.3.200:32382/xbcapi:'+sdf+' -n xbc'
                           sh 'kubectl set image deployment mediaupload mediaupload=10.0.3.200:32382/mediaupload:'+sdf+' -n xbc'

                       }
       }
      
    }
}