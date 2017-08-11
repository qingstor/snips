# Snips 使用指南

Snips 是青云开发的通用的 HTTP API 代码生成工具，使用 OpenAPI 规范（Swagger）v2.0 格式的 API 规范为 QingCloud 和 QingStor SDK 生成各种 API 代码。它的运行方式为从 API 规范文件中读取数据，按照模板利用读取的数据生成 API 代码文件。

**附注1：**API 规范由 QingCloud 和 QingStor 提供，你可以在 [qingstor-api-specs](https://github.com/yunify/qingstor-api-specs) 了解它并获取 qingstor SDK API 规范的最新版本。

**附注2：**关于模板，你可以在[snips_template.md](snips_template.md) 了解它。

**附注3：** 关于 Snips 的数据流图及其使用，你可以在 Snips 的 [README.md](https://github.com/yunify/snips/blob/master/README.md) 了解它。

## Snips 中间状态

Snips 从 API 规范中读取的数据会根据模板文件的规定的格式生成为代码文件，而在 Snips 运行过程中，从 API 规范中读取的数据会根据各自的特点暂存在若八个结构体中。你需要理解这八个结构体的意义及其中各元素存储的内容才能写出有意义的模板文件。

### Snips 中间状态结构

Snips 结构体层次图：

```
Snips
 └── Data
      ├── Service
             ├── Property
             ├── Operation
                    ├── Request
                           ├── StatusCode
                           ├── Property
                    ├── Response
                           ├── Property
      ├── SubService
             ├── Property
             ├── Operation
                    ├── Request
                           ├── StatusCode
                           ├── Property
                    ├── Response
                           ├── Property
      ├── Property
```

**附注：**以下讲解的中间状态是根据其包含与被包含的关系由大到小排列的。

### Data

|  结构体名  	|       元素       	|        元素类型        	|
|:----------:	|:----------------:	|:----------------------:	|
|    Data    	|      Service     	|        *Service        	|
|            	|    SubServices   	| map[string]*SubService 	|
|            	|  CustomizedTypes 	|  map[string]*Property  	|

`Data` 是 Snips 中最大的数据结构，所有的 API 规范都在该字段下。

- `Service` 保存着主服务数据。
- `SubService` 保存着子服务数据。
- `CustomizedTypes` 保存着所有的 SDK 的自定义类型。

### Service

|  结构体名  	|       元素       	|        元素类型        	|
|:----------:	|:----------------:	|:----------------------:	|
|   Service  	|    APIVersion    	|         string         	|
|            	|       Name       	|         string         	|
|            	|    Description   	|         string         	|
|            	|    Properties    	|        *Property       	|
|            	|    Operations    	|  map[string]*Operation 	|

`Service` 保存着主服务数据。

- `APIVersion` 保存着主服务的版本，你可以在生成代码文件时，将这个字段输出到说明中以提醒使用者 SDK 的版本。
- `Name` 保存着主服务的名字，因为主服务只有一个，所以其名字本身就具有唯一性，所以不需要 ID 字段。
- `Description` 保存着对主服务的描述，这个字段只作为说明使用。
- `Properties` 保存着主服务内的全局变量以及属性等。
- `Operations` 保存着主服务的操作。

### SubService

|  结构体名  	|       元素       	|        元素类型        	|
|:----------:	|:----------------:	|:----------------------:	|
| SubService 	|        ID        	|         string         	|
|            	|       Name       	|         string         	|
|            	|    Properties    	|        *Property       	|
|            	|    Operations    	|  map[string]*Operation 	|

`SubService` 保存着子服务数据。

- `ID` 保存着子服务 ID，一般来说子服务 ID 与其名字相同，但是不同子服务的 ID 具有唯一性。
- `Name` 保存着子服务的名字。
- `Properties` 保存着该子服务的配置，包括子服务内的全局参数、属性等。
- `Operation` 保存着子服务的操作。

### Operation

|  结构体名  	|       元素       	|        元素类型        	|
|:----------:	|:----------------:	|:----------------------:	|
|  Operation 	|        ID        	|         string         	|
|            	|       Name       	|         string         	|
|            	|    Description   	|         string         	|
|            	| DocumentationURL 	|         string         	|
|            	|      Request     	|        *Request        	|
|            	|     Response     	|    map[int]*Response   	|

`Operation` 保存着服务以及子服务的操作。

- `ID` 保存着操作的 ID。一般来说，操作 ID 与其名字相同，但是不同操作的 ID 具有唯一性。
- `Name` 保存着操作的名字。
- `Description` 保存着操作的描述，一般来说，这个字段只作为说明使用。
- `DocumentationURL` 保存着对应操作的文档的相对位置。
- `Request` 保存着对应操作的 request 字段，该元素内标明了本次操作请求需要的元素参数等。
- `Response` 保存着对应操作的 response 字段，该元素内著名了本次操作请求后的回复的格式、元素以及参数。

### Request

|  结构体名  	|       元素       	|        元素类型        	|
|:----------:	|:----------------:	|:----------------------:	|
|   Request  	|      Method      	|         string         	|
|            	|        URI       	|         string         	|
|            	|      Params      	|        *Property       	|
|            	|      Headers     	|        *Property       	|
|            	|     Elements     	|        *Property       	|
|            	|       Body       	|        *Property       	|

`Request` 存储着 HTTP(S) 的 request 字段。该字段内指明了某个 request 请求需要的所有必须、必要信息。

- `Method` 存储着 request 字段的请求方法，如 `GET`, `POST`, `PUT` 等。
- `URI` 存储着目标地址的相对位置。
- `Params` 存储着 request 字段的参数。对应 API 规范中的 `Parameters` 字段，这些参数包含某个 request 字段所必须的 `Headers` 和 `Body` 内的参数，如 `content-length`, `bucket-name`, `access-key-id` 等。
- `Headers` 存储着 request 字段的头，但是在头内一般只指明该头需要什么参数，具体参数在 `Params` 中。
- `Elements` 存储着 request 字段的元素。
- `Body` 存储着 request 字段的 Body，但是在 Body 内一般只指明该 Body 需要什么参数，具体参数在 `Params` 中。

### Response

|  结构体名  	|       元素       	|        元素类型        	|
|:----------:	|:----------------:	|:----------------------:	|
|  Response  	|    StatusCodes   	|       *StatusCode      	|
|            	|      Headers     	|        *Property       	|
|            	|     Elements     	|        *Property       	|
|            	|       Body       	|        *Property       	|

`Response` 存储着 HTTP(S) 的 response 字段。该字段指明了在某个 request 请求之后获得的 response 字段的各个参数及其格式。

- `StatusCodes` 存储着 reponse 字段的状态信息。
- `Headers` 存储着 reponse 字段的头，头内指明了各个返回的参数。
- `Elements` 存储着 reponse 字段的元素。
- `Body` 存储着 reponse 字段的 Body。

### StatusCode

|  结构体名  	|       元素       	|        元素类型        	|
|:----------:	|:----------------:	|:----------------------:	|
| StatusCode 	|       Code       	|           int          	|
|            	|    Description   	|         string         	|

`StatusCode` 用于存储状态信息，一般用于存储 HTTP(S) 的 response 字段。

- `Code` 存储着状态值，如 "404", "400" 等。
- `Description` 存储着该状态的状态信息描述，与其他的 `Description` 字段不同，该字段或许对用户有用，因而它应该被输出。

### Property

|  结构体名  	|       元素       	|        元素类型        	|
|:----------:	|:----------------:	|:----------------------:	|
|  Property  	|        ID        	|         string         	|
|            	|       Name       	|         string         	|
|            	|    Description   	|         string         	|
|            	|       Type       	|         string         	|
|            	|     ExtraType    	|         string         	|
|            	|      Format      	|         string         	|
|            	| CollectionFormat 	|         string         	|
|            	|       Enum       	|        []string        	|
|            	|      Default     	|         string         	|
|            	|    IsRequired    	|          bool          	|
|            	|    Properties    	|  map[string]*Property  	|

`Property` 是 Snips 中间状态中最小的结构体。它用于存储 API 规范中服务、操作、request 段或 response 段的属性值以及变量值。

- `ID` 用于存储 `Property` 的 ID。一般来说，ID 即指该 `Property` 的名字，但是 ID 较之 `Name` 具有唯一性。
- `Name` 用于存储 `Property` 的名字。
- `Description` 用于存储对于某个 `Property` 的描述。一般来说，所有的 `Description` 字段都只作为生成的代码文件的注释出现，以下出现的该字段不做更多说明。
- `Type` 用于保存对应 `Property` 的类型。因为 `Property` 用于存储属性值和变量值，`Type` 则存储了该属性值或变量值的类型。你可以根据这个字段输出 `Property` 的类型或进行必要的类型转换。该字段中的类型包括 `object`, `array`, `map`, `any`。其中 `any` 指该 `Property` 的类型没有保存在该字段中，需要在 `ExtraType` 字段中读取。
- `ExtraType` 用于保存其它不常见的变量类型，包括 `string`, `boolean`, `integer`, `long`, `timestamp`, `binary`。同时它其中也会保存 `Type` 的三种变量类型。也就是说，`Type` 内的类型可能是 `ang` ，单 `ExtraType` 内永远都有一个类型。
- `Format` 存储着该 `Property` 的格式，例如 `timestamp` 类型采用哪种计时方式。
- `CollectionFormat` 存储着集合的格式。如果某个 `Property` 存储的元素可以有多个格式，那么就会以集合的方式存储在该字段中。
- `Enum` 存储着 `Property` 中可选存在的集合。如果某个元素有多个值，而且每个值都具有相同的优先级，那么它们就会无序地存储在该字段中。
- `Default` 存储着 `Property` 的默认值。
- `IsRequired` 指明了该 `Property` 是否是必须的，取值为 `true` 或 `false` 。并不是所有的 `Property` 都是必须输出的，你可以根据这个字段判断是否应该将这个 `Property` 输出到代码文件中。
- `Properties` 是本身的嵌套，因为一个属性也有可能包含其它多个子属性。