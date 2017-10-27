# AR App Service



## 用户注册接口

> 请求: POST /v1/user/reg

> 参数:

| 参数名      | 类型     | 说明    |
| -------- | ------ | ----- |
| username | string | 登录用户名 |
| password | string | 登录密码  |
| nickname | string | 昵称    |
| school | string | 学校    |


> 响应:

| 属性名       | 类型     | 说明                  |
| --------- | ------ | ------------------- |
| resp_code | int    | 响应码 1000 成功,非1000失败 |
| remark    | string | 响应描述                |

```json
{
  "resp_code":1000,
  "remark":"注册成功"
}
```


## 用户登录接口

> 请求: POST /v1/user/sign-in

> 参数:

| 参数名      | 类型     | 说明    |
| -------- | ------ | ----- |
| username | string | 登录用户名 |
| password | string | 登录密码  |
|          |        |       |

> 响应:


| 属性名          | 类型     | 说明                  |
| ------------ | ------ | ------------------- |
| resp_code    | int    | 响应码 1000 成功,非1000失败 |
| remark       | string | 响应描述                |
| username     | string | 用户名                 |
| access_token | string | 用户访问凭证              |
| userid       | string | 用户id                |

```json
{
  "resp_code":1000,
  "remark":"登录成功",
  "username":"张三",
  "access_token":"330058979596026024705875151182809097893",
  "userid":"323558170212446235845673443945363629733"
}
```


## 用户退出登录接口

> 请求: POST /v1/user/sign-out

> 参数:

| 参数名          | 类型     | 说明     |
| ------------ | ------ | ------ |
| userid       | string | 用户id   |
| access_token | string | 用户访问凭证 |
|              |        |        |

> 响应:


| 属性名       | 类型     | 说明                  |
| --------- | ------ | ------------------- |
| resp_code | int    | 响应码 1000 成功,非1000失败 |
| remark    | string | 响应描述                |


```json
{
  "resp_code":1000,
  "remark":"退出成功"
}
```


## 获取进度

> 请求: GET  /v1/schedules

> 参数:

| 参数名          | 类型     | 说明     |
| ------------ | ------ | ------ |
| userid       | string | 用户id   |
| access_token | string | 用户访问凭证 |
|              |        |        |

> 响应:


| 属性名       | 类型     | 说明                  |
| --------- | ------ | ------------------- |
| resp_code | int    | 响应码 1000 成功,非1000失败 |
| remark    | string | 响应描述                |
| tasks| []Object| 进度

>> task Object:

| 属性名      | 类型     | 说明                  |
| -------- | ------ | ------------------- |
| id       | string | 任务id                |
| name     | string | 任务名称                |
| descript | string | 任务说明                |
| status   | string | [任务状态可选值](#任务状态可选值) |


```json
{
  "resp_code":1000,
  "remark":"success",
  "tasks":[
    {
      "id":"t001",
      "name":"任务1",
      "descript":"任务1描述",
      "status":"not-unlock"
    },
    {
      "id":"t002",
      "name":"任务2",
      "descript":"任务2描述",
      "status":"in-process"
    },
    {
      "id":"t003",
      "name":"任务3",
      "descript":"任务3描述",
      "status":"done"
    }
  ]
}
```


## 更新进度


> 请求:  POST  /v1/schedules/update

> 参数:

| 参数名          | 类型     | 说明                         |
| ------------ | ------ | -------------------------- |
| userid       | string | 用户id                       |
| access_token | string | 用户访问凭证                     |
| task_id      | string | 任务id                       |
| task_status  | string | 更新任务状态 [任务状态可选值](#任务状态可选值) |


> 响应:


| 属性名       | 类型     | 说明                  |
| --------- | ------ | ------------------- |
| resp_code | int    | 响应码 1000 成功,非1000失败 |
| remark    | string | 响应描述                |


```json
{
  "resp_code":1000,
  "remark":"更新成功"
}
```




### 任务状态和变化条件
<a name="任务状态可选值"></a>

* not-unlock: 任务未解锁
* in-process:进行中
* done:完成


任务进度更新,新建的任务默认状态是 `not-unlock`:

`not-unlock` 可以更新成 `in-process` 或 `done`

`in-process` 的状态可以更新成 `done`

`done` 不能再做状态更新


## 学校列表接口


> 请求:  GET  /v1/schools

> 参数: 无

> 响应:


| 属性名       | 类型     | 说明                  |
| --------- | ------ | ------------------- |
| resp_code | int    | 响应码 1000 成功,非1000失败 |
| remark    | string | 响应描述                |
|  schools | []Object | 学校列表|


>> school Object:

| 属性名      | 类型     | 说明                  |
| -------- | ------ | ------------------- |
| id       | string | 学校id                |
| name     | string | 学校名称                |



```json
{
  "resp_code":1000,
  "remark":"succes",
  "schools":[
    {
      "id":"学校id",
      "name":"学校名称"
    },
    {
      "id":"学校id",
      "name":"学校名称"
    }
  ]
}
```
