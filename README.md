# GCR sample

## setup

```
# クラスタ作成
gcloud container clusters create pubsub --num-nodes=3 --zone=asia-northeast1-a --enable-stackdriver-kubernetes --cluster-version=latest --preemptible

# クラスタの確認
kubectl config get-contexts

# サービスアカウントの作成
gcloud projects add-iam-policy-binding work-999999 --member=serviceAccount:pubsub-publisher@work-999999.iam.gserviceaccount.com --role=roles/pubsub.publisher
# サービスアカウントのキー作成
gcloud iam service-accounts keys create credential.json --iam-account=pubsub-publisher@work-999999.iam.gserviceaccount.com --key-file-type=json 


# トピックの作成
gcloud pubsub topics create sample-topic
# サブスクリプションの作成
gcloud pubsub subscriptions create sample-subscription --topic=sample-topic

# build
docker build -t asia.gcr.io/work-999999/pubsub-publisher .
# push
docker push asia.gcr.io/work-999999/pubsub-publisher

# credential生成設定用のsecret作成
kubectl create secret generic pubsub-credential --from-file=credential.json

# deploy
kubectl apply -f kubernetes/pubsub-publisher.yaml

# 外部IP
kubectl get services pubsub-publisher

# log
kubectl logs -f pubsub-publisher-no-credential

# サイトにアクセス
curl EXTERNAL_IP

# トピックの取り出し
gcloud pubsub subscriptions pull --auto-ack projects/work-999999/subscriptions/sample-subscription
```

<br>

### kubectx, kubens
コマンド
```
# install（Win scoop）
scoop install kubectx
# コンテキスト一覧
kubectx
# コンテキスト切り替え
kubectx gke-cluster
# コンテキスト名の変更
kubectx 変更前⁼変更後

# Namespace一覧
kubens
# Namespace切り替え
kubens default
``` 

<br>


<br>

### 参考リンク
* クラウド ビルダー  
https://cloud.google.com/build/docs/cloud-builders?hl=ja

* ロールについて  
https://cloud.google.com/iam/docs/understanding-roles#kubernetes-engine-roles
