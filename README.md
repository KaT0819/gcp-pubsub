# GCR sample

## setup

```
# 以降のコマンドで実行するためにプロジェクトIDを環境変数にセット
$project_id='work-999999'


# クラスタ作成
gcloud container clusters create pubsub-cluster `
--num-nodes=3 `
--zone=asia-northeast1-a `
--enable-stackdriver-kubernetes `
--cluster-version=latest `
--preemptible

# クラスタの確認
kubectl config get-contexts

# サービスアカウントの作成
gcloud projects add-iam-policy-binding $project_id `
--member=serviceAccount:pubsub-publisher@$project_id.iam.gserviceaccount.com `
--role=roles/pubsub.publisher
# サービスアカウントのキー作成
gcloud iam service-accounts keys create credential.json `
--iam-account=pubsub-publisher@$project_id.iam.gserviceaccount.com `
--key-file-type=json 


# トピックの作成
gcloud pubsub topics create sample-topic
# サブスクリプションの作成
gcloud pubsub subscriptions create sample-subscription --topic=sample-topic

# build
docker build -t asia.gcr.io/$project_id/pubsub-publisher .
# push
docker push asia.gcr.io/$project_id/pubsub-publisher

# credential生成設定用のsecret作成
kubectl create secret generic pubsub-credential --from-file=credential.json

# deploy
kubectl apply -f kubernetes/pubsub-publisher.yaml

# 外部IP
kubectl get services pubsub-publisher

# log
kubectl logs -f pubsub-publisher

# サイトにアクセス
curl EXTERNAL_IP

# トピックの取り出し
gcloud pubsub subscriptions pull --auto-ack projects/$project_id/subscriptions/sample-subscription
```

<br>

### 負荷ツール
```
# vegeta install
go get -u github.com/tsenart/vegeta

# 
echo "GET http://$external_ip" | vegeta attack -rate=20 -duration=10s | vegeta encode | vegeta report -type='hist[0,5ms,10ms,15ms,20ms,100ms]'
```


<br>

### Wordpress
```
gcloud sql instances create wordpress `
--availability-type=zonal `
--database-version=MYSQL_5_7 `
--root-password=password `
--zone=asia-northeastl-a

gcloud services enable sqladniin.googleapis.com

# サービスアカウント
gcloud iam service-accounts create wordpress-sql-client `
--display-name="Wordpress SQL Client"
# ロールの割り当て
gcloud projects add-iam-policy-binding $project_id `
--member=serviceAccount:wordpress-sql-client@$project_id.iam.gserviceaccount.com `
--role=roles/cloudsql.client
# 鍵ファイル作成
gcloud iam service-accounts keys create credential.json `
--iam-account=wordpress-sql-client@$project_id.iam.gserviceaccount.com `
--key-file-type=json
# credential生成設定用のsecret作成
kubectl create secret generic wordpress-sql-proxy --from-file=credential.json

pubsub-publisher-secret
gcloud iam service-accounts create pubsub-publisher-secret `
--display-name="PubSub Publisher"
# ロールの割り当て
gcloud projects add-iam-policy-binding $project_id `
--member=serviceAccount:pubsub-publisher-secret@$project_id.iam.gserviceaccount.com `
--role=roles/run.admin
gcloud projects add-iam-policy-binding $project_id `
--member=serviceAccount:pubsub-publisher-secret@$project_id.iam.gserviceaccount.com `
--role=roles/iam.serviceAccountKeyAdmin



```

<br>

### Volume
```
# Dynamic Volume Provisioning
# deploy
kubectl apply -f kubernetes/dynamic-provisioning.yaml

# PVC（Persistent Volume Claim）
kubectl get pvc
# PV（Persistent Volume）
kubectl get pv
# SC（Strage Class）
kubectl get sc

# index.htmlの作成
kubectl exec nginx-with-volume -- touch /var/www/html/index.html
# ファイル確認
kubectl exec nginx-with-volume -- ls /var/www/html

# Pod削除
kubectl delete pods nginx-with-volume
# Pod再作成
kubectl apply -f kubernetes/dynamic-provisioning.yaml
# ファイル確認(変わらず存在)

# Static Volume Provisioning
# 永続ディスク作成
gcloud compute disks create web-volume --zone=asia-northeast1-a --size=10
# dynamic同様ファイルがpodを再作成した後も存在
``` 

<br>

### GKEによるCD
```
# クラスタ作成
gcloud container clusters create gke-cluster `
--num-nodes=3 `
--zone=asia-northeast1-a `
--enable-stackdriver-kubernetes `
--cluster-version=latest `
--preemptible

# クラスタの確認
kubectl config get-contexts

# デプロイ
gcloud builds submit --config .\cloudbuild.yaml

```

# クラスタ作成
gcloud container clusters create cloud-runner `
--addons 'HttpLoadBalancing,CloudRun' `
--enable-stackdriver-kubernetes `
--zone asia-northeast1-a `
--preemptible

# クラスタの確認
kubectl config get-contexts

<br>

### 参考リンク
* クラウド ビルダー  
https://cloud.google.com/build/docs/cloud-builders?hl=ja

* ロールについて  
https://cloud.google.com/iam/docs/understanding-roles#kubernetes-engine-roles
