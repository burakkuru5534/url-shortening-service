# Url Shortening Service
[]: # Author: [Burak Kuru](

## Introduction

In this project, we will build a URL shortening service. The service will take a long URL and return a short URL. When the short URL is visited, the service will redirect the user to the long URL.
Usr must register and login before using the service. There will be two types of usrs: free and premium. Free usrs can shorten up to 1 URLs, while premium usrs can shorten 10  of URLs. 
Usrs can also delete their URLs.
Usrs can also list their active shortened urls.
The service will be implemented using the following technologies:

### Languages and frameworks

Technologies used in this project:

Golang,
postgresql
docker
kubernetes


### Database

Postgresql 

Tables created:

```
table name:usr
columns:
id: serial primary key
password: text not null
email: text not null unique
shortening_url_limit: bigint
account_type: text (could be free or premium)
is_active: bool defult true
zlins_dttm: timestamp default now()
zlupd_dttm: timestamp


table name:url
columns:
id: serial primary key
long_version: text not null
shortened_version: text not null
usr_id: bigint fk to usr.id
zlins_dttm: timestamp default now()
zlupd_dttm: timestamp
constraint: unique (long_version, usr_id)
constraint: unique (shortened_version, usr_id)

table name: log_jwt

id: serial primary key
dttm: timestamp default now()
usr_id: bigint fk to usr.id
jwt: text not null
expires_on: timestamp not null
is_invalid: bool default false


```

## Problem solution

● Create Shortened URL: takes the original URL and returns a shortened URL
● Return Active Shortened URLs: returns url ids, active shortened URLs, original URLs
● Delete Shortened URL: takes URL id and returns generic success/error message
● Redirect: takes a shortened URL and returns the original URL

### Register


Register request url example:

Method: POST

http://localhost:8080/api/register

request Body Example:

 ```json
{
  "email":"testuser@gmail.com",
  "password":"Test123456",
  "IsPremium":true
}
 ```

response example:

for 200:

 ```
{"success":true,"message":"user created successfully","data":"Usr testuser@gmail.com created:"}
 ```

for 400:

```json
{"error": "Bad request"}
```

for 403:
```json
{"error": "User with that email already exists"}
```

for 500:
```json
{"error": "server error"}
```


### Login


Login request url example:

Method: POST

http://localhost:8080/api/login

request Body Example:

 ```json
{
  "email":"testuser@gmail.com",
  "password":"Test123456"
}
 ```

response example:

for 200:

 ```
{"success":true,"message":"user login successfully","data":{"UserID":8,"UserEmail":"testuser@gmail.com","AccessToken":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NjUyMTk5NjksInVzZXJFbWFpbCI6InRlc3R1c2VyQGdtYWlsLmNvbSIsInVzZXJJRCI6OH0.K5qmO9M5Znq-tVIauxHPEHlDZa3_gbQgr7FF4Y6fwNA"}}
 ```

for 400:

```json
{"error": "Bad request"}
```

for 401:
```json
{"error": "wrong password"}
```

for 500:
```json
{"error": "server error"}
```


### Create Shortened Url

Create shortened url request url example:

Method: POST

Note: this api for authenticated users only.

http://localhost:8080/api/url

request Body Example:

 ```json
{
  "LongVersion":"www.facebook.com"

}
 ```

response example:

for 200:

 ```json
{"success":true,"message":"url created successfully","data":"5waxJgD7"}
 ```

for 400:

```json
{"error": "Bad request"}
```
```

for 401:
```
no token found
```

for 403:
```json
{"error": "url with that name already exists"}
```

for 500:
```json
{"error": "server error"}
```

### Return Active Shortened Urls

Return Active Shortened Urls request url example:

Method: GET

Note: this api for authenticated users only.

http://localhost:8080/api/urls



request Body:

response example:

for 200:
 ```json
{"success":true,"message":"url list successfully","data":[{"ID":12,"LongVersion":"www.google.com","ShortenedVersion":"b2Abvk5G","UsrID":5}]}
 ```

for 400:

```json
{"error": "Bad request"}
```

for 401:
```
no token found
```

also 
for 401:
```
token is expired
```

for 500:
```json
{"error": "server error"}
```

### Delete url

Delete url request url example:
Method: DELETE

Note: this api for authenticated users only.

http://localhost:8080/api/url?id=1

id: this id should be one of the url's ids.

response example:

for 200: 

```json
{"success":true,"message":"url deleted successfully","data":12}
```


for 400:

```json
{"error": "Bad request"}
```

for 401:
```
no token found
```

also 
for 401:
```
token is expired
```

for 404:
```json
{"error": "url with that id does not exist"}
```

for 500:
```json
{"error": "server error"}
```

### Redirect 

Redirect  request url example:
Method: GET

http://localhost:8080/api/url?shortenedVersion=2iFdn9qF

response example:

for 200:
 ```json
{"success":true,"message":"get long url successfully","data":"www.facebook.com"}
 ```
for 400:

```json
{"error": "Bad request"}
```

for 401:
```
no token found
```

also
for 401:
```
token is expired
```


for 500:
```json
{"error": "server error"}
```

### Docker & kubernetes

#$ go build
#$ ./url-shortening-service

#create Dockerfile
#$ docker build -t url-shortening-service .
#$ docker tag go-kubernetes burakkuru5534/url-shortening-service:1.0.0
#$ docker login
#$ docker push burakkuru5534/url-shortening-service:1.0.0

#create kubernetes deployment file (.yml)
#$ minikube start
#$ kubectl apply -f k8s-deployment.yml

#$ kubectl get deployments
#$ kubectl get pods
#We can use the kubectl port-forward command to map a local port to a port inside the pod like this:
#$ kubectl port-forward url-shortening-service-69b45499fb-7fh87 8080:8080

#$ kubectl logs -f url-shortening-service-69b45499fb-7fh87

#create kubernetes service
#kubectl apply -f k8s-deployment.yml (we can update this yml or create another yml file)

#$ kubectl get services

#Type the following command to get the URL for the service in the minikube cluster:
#$ minikube service url-shortening-service-service --url

#scale a kubernetes deployment

#$ kubectl scale --replicas=4 deployment/url-shortening-service

#delete a kubernetes deployment

#$ kubectl delete deployment url-shortening-service

#delete a kubernetes service

#$ kubectl delete service url-shortening-service-service

#delete a pod

#$ kubectl delete pod url-shortening-service-69b45499fb-7fh87

## Conclusion

We have successfully implemented the url-shortening-service api service.
Used log library to log errors.
Used jwt library to authenticate users.
Used postgresql as the database language.
Used postman to test the apis.
Used docker to containerize the application.
Used kubernetes to deploy the application.


