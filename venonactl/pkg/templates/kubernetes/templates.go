// Code generated by go generate; DO NOT EDIT.
// using data from templates/kubernetes
package kubernetes

func TemplatesMap() map[string]string {
	templatesMap := make(map[string]string)

	templatesMap["cluster-role-binding.dind-volume-provisioner.vp.yaml"] = `---
kind: ClusterRoleBinding
apiVersion: rbac.authorization.k8s.io/v1beta1
metadata:
  name: volume-provisioner-{{ .AppName }}-{{ .Namespace }}
  labels:
    app: dind-volume-provisioner-{{ .AppName }}
subjects:
  - kind: ServiceAccount
    name: volume-provisioner-{{ .AppName }}
    namespace: {{ .Namespace }}
roleRef:
  kind: ClusterRole
  name: volume-provisioner-{{ .AppName }}-{{ .Namespace }}
  apiGroup: rbac.authorization.k8s.io
`

	templatesMap["cluster-role-binding.engine.yaml"] = `kind: ClusterRoleBinding
apiVersion: rbac.authorization.k8s.io/v1beta1
metadata:
  name: {{ .AppName }}-{{ .Namespace }}-engine
subjects:
- kind: ServiceAccount
  name: engine
  namespace: {{ .Namespace }}
roleRef:
  kind: ClusterRole
  name: cluster-admin
  apiGroup: rbac.authorization.k8s.io`

	templatesMap["cluster-role-binding.venona.yaml"] = `kind: ClusterRoleBinding
apiVersion: rbac.authorization.k8s.io/v1beta1
metadata:
  name: {{ .AppName }}-{{ .Namespace }}
subjects:
- kind: ServiceAccount
  name: {{ .AppName }}
  namespace: {{ .Namespace }}
roleRef:
  kind: ClusterRole
  name: system:discovery
  apiGroup: rbac.authorization.k8s.io`

	templatesMap["cluster-role.dind-volume-provisioner.vp.yaml"] = `kind: ClusterRole
apiVersion: rbac.authorization.k8s.io/v1beta1
metadata:
  name: volume-provisioner-{{ .AppName }}-{{ .Namespace }}
  labels:
    app: dind-volume-provisioner
rules:
  - apiGroups: [""]
    resources: ["persistentvolumes"]
    verbs: ["get", "list", "watch", "create", "delete"]
  - apiGroups: [""]
    resources: ["persistentvolumeclaims"]
    verbs: ["get", "list", "watch", "update"]
  - apiGroups: ["storage.k8s.io"]
    resources: ["storageclasses"]
    verbs: ["get", "list", "watch"]
  - apiGroups: [""]
    resources: ["events"]
    verbs: ["list", "watch", "create", "update", "patch"]
  - apiGroups: [""]
    resources: ["secrets"]
    verbs: ["get", "list"]
  - apiGroups: [""]
    resources: ["nodes"]
    verbs: ["get", "list", "watch"]
  - apiGroups: [""]
    resources: ["pods"]
    verbs: ["get", "list", "watch", "create", "delete", "patch"]
  - apiGroups: [""]
    resources: ["endpoints"]
    verbs: ["get", "list", "watch", "create", "update", "delete"]
`

	templatesMap["cluster-role.engine.yaml"] = `kind: ClusterRole
apiVersion: rbac.authorization.k8s.io/v1beta1
metadata:
  name: {{ .AppName }}-{{ .Namespace }}-engine
rules:
  - apiGroups: [""]
    resources: ["pods", "persistentvolumeclaims"]
    verbs: ["get", "create", "delete", "list"]
`

	templatesMap["codefresh-certs-server-secret.re.yaml"] = `apiVersion: v1
type: Opaque
kind: Secret
metadata:
  labels:
    app: venona
  name: codefresh-certs-server
  namespace: {{ .Namespace }}
data:
  server-cert.pem: {{ .ServerCert.Cert }}
  server-key.pem: {{ .ServerCert.Key }}
  ca.pem: {{ .ServerCert.Ca }}

`

	templatesMap["daemonset.dind-lv-monitor.vp.yaml"] = `apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: dind-lv-monitor-{{ .AppName }}
  namespace: {{ .Namespace }}
  labels:
    app: dind-lv-monitor
spec:
  selector:
    matchLabels:
      app: dind-lv-monitor
  template:
    metadata:
      labels:
        app: dind-lv-monitor
      annotations:
        prometheus_port: "9100"
        prometheus_scrape: "true"
    spec:
      serviceAccountName: volume-provisioner-{{ .AppName }}
      # Debug:
      # hostNetwork: true
      # nodeSelector:
      #   kubernetes.io/role: "node"
      {{ if ne .NodeSelector "" }}
      nodeSelector:
        {{ .NodeSelector }}
      {{ end }}
      tolerations:
        - key: 'codefresh/dind'
          operator: 'Exists'
          effect: 'NoSchedule'
      {{ if ne .Tolerations "" }}
        {{ .Tolerations | indent 8 }}
      {{ end }}
      containers:
        - image: codefresh/dind-volume-utils:v5
          name: lv-cleaner
          imagePullPolicy: Always
          command:
          - /bin/local-volumes-agent
          env:
            - name: NODE_NAME
              valueFrom:
                fieldRef:
                  fieldPath: spec.nodeName
            - name: VOLUME_PARENT_DIR
              value: /var/lib/codefresh/dind-volumes
#              Debug:
#            - name: DRY_RUN
#              value: "1"
#            - name: DEBUG
#              value: "1"
#            - name: SLEEP_INTERVAL
#              value: "3"
#            - name: LOG_DF_EVERY
#              value: "60"
#            - name: KB_USAGE_THRESHOLD
#              value: "20"

          volumeMounts:
          - mountPath: /var/lib/codefresh/dind-volumes
            readOnly: false
            name: dind-volume-dir
      volumes:
      - name: dind-volume-dir
        hostPath:
          path: /var/lib/codefresh/dind-volumes
`

	templatesMap["deployment.dind-volume-provisioner.vp.yaml"] = `apiVersion: apps/v1
kind: Deployment
metadata:
  name: dind-volume-provisioner-{{ .AppName }}
  namespace: {{ .Namespace }}
  labels:
    app: dind-volume-provisioner
spec:
  selector:
    matchLabels:
      app: dind-volume-provisioner
  replicas: 1
  strategy:
    type: Recreate
  template:
    metadata:
      labels:
        app: dind-volume-provisioner
    spec:
      {{ if ne .NodeSelector "" }}
      nodeSelector:
        {{ .NodeSelector }}
      {{ end }}
      serviceAccount: volume-provisioner-{{ .AppName }}
      tolerations:
      - effect: NoSchedule
        key: node-role.kubernetes.io/master
        operator: "Exists"
      {{ if ne .Tolerations "" }}
        {{ .Tolerations | indent 8 }}
      {{ end }}
      containers:
      - name: dind-volume-provisioner
        image: {{ .VolumeProvisionerImage.Name }}:{{ .VolumeProvisionerImage.Tag }}
        imagePullPolicy: Always
        resources:
          requests:
            cpu: "300m"
            memory: "400Mi"
          limits:
            cpu: "1000m"
            memory: "6000Mi"
        command:
          - /usr/local/bin/dind-volume-provisioner
          - -v=4
          - --resync-period=50s
        env:
        - name: PROVISIONER_NAME
          value: codefresh.io/dind-volume-provisioner-{{ .AppName }}-{{ .Namespace }}
`

	templatesMap["deployment.venona.yaml"] = `apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: {{ .AppName }}
    version: {{ .Version }}
  name: {{ .AppName }}
  namespace: {{ .Namespace }}
spec:
  selector:
    matchLabels:
      app: {{ .AppName }}
      version: {{ .Version }}
  replicas: 1
  revisionHistoryLimit: 5
  strategy:
    rollingUpdate:
      maxSurge: 50%
      maxUnavailable: 50%
    type: RollingUpdate
  template:
    metadata:
      labels:
        app: {{ .AppName }}
        version: {{ .Version }}
    spec:
      serviceAccountName: {{ .AppName }}
      {{ if ne .NodeSelector "" }}
      nodeSelector:
        {{ .NodeSelector }}
      {{ end }}
      {{ if ne .Tolerations "" }}
      tolerations:
        {{ .Tolerations | indent 8 }}
      {{ end }}
      containers:
      - env:
        - name: SELF_DEPLOYMENT_NAME
          valueFrom:
            fieldRef:
              fieldPath: metadata.name
        - name: CODEFRESH_TOKEN
          valueFrom:
            secretKeyRef:
              name: {{ .AppName }}
              key: codefresh.token
        - name: CODEFRESH_HOST
          value: {{ .CodefreshHost }}
        - name: AGENT_MODE
          value: {{ .Mode }}
        - name: AGENT_NAME
          value: {{ .AppName }}
        image: {{ .Image.Name }}:{{ .Image.Tag }}
        imagePullPolicy: Always
        name: {{ .AppName }}
      restartPolicy: Always
`

	templatesMap["dind-daemon-conf.re.yaml"] = `---
apiVersion: v1
kind: ConfigMap
metadata:
  name: codefresh-dind-config
  namespace: {{ .Namespace }}
data:
  daemon.json: |
    {
      "hosts": [ "unix:///var/run/docker.sock",
                 "tcp://0.0.0.0:1300"],
      "storage-driver": "overlay2",
      "tlsverify": true,  
      "tls": true,
      "tlscacert": "/etc/ssl/cf-client/ca.pem",
      "tlscert": "/etc/ssl/cf/server-cert.pem",
      "tlskey": "/etc/ssl/cf/server-key.pem",
      "insecure-registries" : ["192.168.99.100:5000"],
      "metrics-addr" : "0.0.0.0:9323",
      "experimental" : true
    }
`

	templatesMap["dind-headless-service.re.yaml"] = `---
apiVersion: v1
kind: Service
metadata:
  labels:
    app: dind
  name: dind
  namespace: {{ .Namespace }}
spec:
  ports:
  - name: "dind-port"
    port: 1300
    protocol: TCP

  # This is a headless service, Kubernetes won't assign a VIP for it.
  # *.dind.default.svc.cluster.local
  clusterIP: None
  selector:
    app: dind

`

	templatesMap["role-binding.venona.yaml"] = `kind: RoleBinding
apiVersion: rbac.authorization.k8s.io/v1beta1
metadata:
  name: {{ .AppName }}
  namespace: {{ .Namespace }}
subjects:
- kind: ServiceAccount
  name: {{ .AppName }}
roleRef:
  kind: Role
  name: {{ .AppName }}
  apiGroup: rbac.authorization.k8s.io`

	templatesMap["role.venona.yaml"] = `kind: Role
apiVersion: rbac.authorization.k8s.io/v1beta1
metadata:
  name: {{ .AppName }}
  namespace: {{ .Namespace }}
rules:
- apiGroups: [""]
  resources: ["pods", "persistentvolumeclaims"]
  verbs: ["get", "create", "delete"]
`

	templatesMap["secret.venona.yaml"] = `apiVersion: v1
kind: Secret
type: Opaque
metadata:
  name: {{ .AppName }}
  namespace: {{ .Namespace }}
data:
  codefresh.token: {{ .AgentToken }}`

	templatesMap["service-account.dind-volume-provisioner.vp.yaml"] = `---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: volume-provisioner-{{ .AppName }}
  namespace: {{ .Namespace }}
  labels:
    app: dind-volume-provisioner
`

	templatesMap["service-account.engine.yaml"] = `apiVersion: v1
kind: ServiceAccount
metadata:
  name: engine
  namespace: {{ .Namespace }}`

	templatesMap["service-account.venona.yaml"] = `apiVersion: v1
kind: ServiceAccount
metadata:
  name: {{ .AppName }}
  namespace: {{ .Namespace }}`

	templatesMap["storageclass.dind-local-volume-provisioner.vp.yaml"] = `---
kind: StorageClass
apiVersion: storage.k8s.io/v1
metadata:
  name: dind-local-volumes-{{ .AppName }}-{{ .Namespace }}
  labels:
    app: dind-volume-provisioner
provisioner: codefresh.io/dind-volume-provisioner-{{ .AppName }}-{{ .Namespace }}
parameters:
  volumeBackend: local
`

	return templatesMap
}
