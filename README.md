# mrClinic

## Tech Stack

[![Go.Dev reference](https://img.shields.io/badge/-echo%20-05122A?style=flat&logo=go)](https://github.com/labstack/echo)
[![Go.Dev reference](https://img.shields.io/badge/-gorm%20-05122A?style=flat&logo=go)](https://pkg.go.dev/gorm.io/gorm?tab=doc)
[![Go.Dev reference](https://img.shields.io/badge/-aws--s3%20-05122A?style=flat&logo=Amazon%20Aws)](https://github.com/awsdocs/aws-doc-sdk-examples/tree/main/go/example_code/s3)
[![Go.Dev reference](https://img.shields.io/badge/-google--calendar%20-05122A?style=flat&logo=Google%20Calendar)](https://github.com/googleapis/google-api-go-client/blob/main/examples/calendar.go)

## About The Project

<p>
  Lorem ipsum dolor sit amet, consectetur adipiscing elit. Fusce erat ipsum, auctor eu condimentum quis, sagittis in risus. Nulla facilisi. Mauris et malesuada lorem. Nam eget accumsan velit, id aliquet magna. Sed quis ornare nibh. Cras sodales elit ex, ut tempus mi porta id. Donec finibus pulvinar est, sed tempor eros lacinia aliquet. Maecenas et dui sit amet ipsum luctus eleifend sed at nisl. Nam et ligula ac massa feugiat porttitor non eget ipsum. Sed velit mi, pellentesque id velit non, commodo mollis turpis. Praesent volutpat, mi sed ultricies finibus, nisi dolor eleifend augue, congue luctus ipsum diam eget lacus. In eget nulla diam. Curabitur auctor tortor lectus, in dictum turpis aliquet id. Ut sed pellentesque ligula. Aliquam erat volutpat. Pellentesque in nulla ac tellus consequat accumsan.
</p>

<p>
Maecenas vel justo venenatis, rhoncus mauris vel, elementum augue. Aenean volutpat est quis mi pharetra pharetra. Sed accumsan molestie mattis. Vestibulum id erat risus. Nunc placerat et arcu quis euismod. Fusce accumsan nibh nisi, id interdum sapien lacinia quis. Fusce eu enim hendrerit, laoreet lectus sed, egestas est. Aliquam bibendum vestibulum nibh, facilisis maximus magna convallis a.
</p>

[OPEN API](https://app.swaggerhub.com/apis/faliqadlan/mrClinic/1.0.0)

<details>
<summary>ERD</summary>
</details>
<details>
<summary>HLA</summary>
</details>
<details>
<summary>Doctor</summary>

| Feature Doctor | Endpoint        | Query Param | Request Body         | JWT Token | Utility                               |
| -------------- | --------------- | ----------- | -------------------- | --------- | ------------------------------------- |
| POST           | /login          | -           | indentity & password | NO        | login user with indentity & passwords |
| POST           | /doctor         | -           | \_                   | NO        | register doctor                       |
| PUT            | /doctor         | -           | -                    | YES       | update current doctor profile         |
| DELETE         | /doctor         | -           | -                    | YES       | delete current doctor account         |
| GET            | /doctor/profile | \_          | -                    | YES       | get current doctor profile            |
| GET            | /doctor/all     | -           | -                    | YES       | get all doctor                        |

</details>

<details>
<summary>Patient</summary>

| Feature Patient | Endpoint         | Query Param | Request Body         | JWT Token | Utility                               |
| --------------- | ---------------- | ----------- | -------------------- | --------- | ------------------------------------- |
| POST            | /login           | -           | indentity & password | NO        | login user with indentity & passwords |
| POST            | /patient         | -           | \_                   | NO        | register patient                      |
| DELETE          | /patient         | -           | -                    | YES       | delete current patient account        |
| PUT             | /patient         | -           | -                    | YES       | update current patient profile        |
| GET             | /patient/profile | patient_uid | -                    | YES       | get current patient profile           |

</details>

<details>
<summary>Visit</summary>

| Feature Visit | Endpoint          | Query Param                      | Request Body | JWT Token | Utility              |
| ------------- | ----------------- | -------------------------------- | ------------ | --------- | -------------------- |
| POST          | /Visit            | -                                | \_           | NO        | add visit            |
| PUT           | /Visit/:visit_uid | -                                | -            | YES       | update visit detail  |
| DELETE        | /Visit/:visit_uid | -                                | -            | YES       | delete current visit |
| GET           | /Visit            | kind, uid, status, date, grouped | -            | YES       | get visit            |

</details>
<details>
<summary>Testing</summary>
</details>

<p align="right">(<a href="#top">back to top</a>)</p>

# How to Use

### 1. install

```bash
git clone https://github.com/MR-Clinic/be.git
```

### 2.1 create file .env

```bash
touch .env
```

### 2.2 open .env

```
nano .env
```

### 2.3 write .env

```bash
PORT=8080
DB=mysql
DB_Name=<database name>
DB_PORT=3306
DB_HOST=localhot
DB_USERNAME=<username>/root
DB_PASSWORD=<password>
DB_LOC=Local
S3_REGION=<reqion AWS S3>/ap-southeast-1
S3_ID=<KEY ID AWS S3>
S3_SECRET=<SECRETKEY AWS S3>
```
### 3.1 create credential folder

```bash
mkdir credential
```

### 3.2 place credential.json

place credential.json for access google calendar API <br> here the reference link
</br>
[Click Here!](https://developers.google.com/workspace/guides/create-credentials#oauth-client-id)

### 3.3 create token folder

```bash
mkdir credential
```

### 3.4 place token.json

place token.json for access google calendar API <br> here the reference link
</br>
[Click Here!](https://github.com/googleworkspace/go-samples/blob/master/calendar/quickstart/quickstart.go)

### 4. run main.go

```bash
go run main.go
```
