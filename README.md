# mariadb-operator
k8s operator for mariadb

## Prerequisites
You'll need [operator-sdk from CoreOS](https://github.com/operator-framework/operator-sdk).
Build and install according to the instructions

## Building the operator
```sh
$ cd mariadb-operator
$ operator-sdk build cflondonservices/mysql-operator
$ docker push cflondonservices/mysql-operator
```

## Using the operator
See [example mariadb helm chart](https://github.com/kieron-pivotal/mariadb-helm), in particular `templates/operator_*.yaml`, 
for an example of installing the operator.

The operator will then listen for new `mysqlbinding` objects. See `deploy/cr.yaml` for an example.

```sh
$ kubectl create -f deploy/cr.yaml
```

Run the above to create a new mysqlbinding object, and the operator will create a new database, user and grant all to that user
on the DB. Details will be stored into the CR object:

```
$ kubectl get mysqlbinding example -o yaml
apiVersion: binding.mysql.pivotal.io/v1alpha1
kind: MysqlBinding
metadata:
  clusterName: ""
  creationTimestamp: 2018-06-28T15:11:31Z
  generation: 1
  name: example
  namespace: default
  resourceVersion: "936765"
  selfLink: /apis/binding.mysql.pivotal.io/v1alpha1/namespaces/default/mysqlbindings/example
  uid: 8a019f56-7ae5-11e8-bb87-42010a9a0003
spec:
  database: QDIWWMNPPQ2
  hostname: wrapping-stingray-mariadb
  password: LMODYLUMGU
  username: MWJGSAPGAT
status:
  status: ""
```
