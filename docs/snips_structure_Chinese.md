# Snips 数据结构

Snips 是青云的 SDK 代码生成器，使用 OpenAPI 规范（Swagger）v2.0 格式的API规范为 QingCloud 和 QingStor SDK 生成各种 API 代码。它的运行方式为从 API 规范文件中读取数据，按照模板利用读取的数据生成 API 代码文件。

## API规范

API规范由 QingCloud 和 QingStor 提供，API 规范中指定了 API 的版本，规格格式，规格以及 API 数据。Snips 的 API 规范采用 json 格式保存。你可以在 [qingstor-api-specs](https://github.com/yunify/qingstor-api-specs/tree/master) 获取 qingstor SDK API 规范的最新版本。

### API 规范目录结构

API规范目录的结构有以下级别：

- API 版本
- 规格格式
- 规格

例如：

``` sh
└── 2016-01-06
    └── swagger
        └── api_v2.0.json
        ├── bucket.json
        ├── bucket_acl.json
        ...
```
## 模板

模板文件由使用者创建，模板文件的格式为 go 模板语言。你需要按照你的编程语言规范和格式利用 API 规范中的数据编写你的模板。

### 模板目录结构

该模板的内的文件名都必须是固定的，你需要按照下面的目录结构安排你的模板目录：

``` sh
└──template
      └──manifest.json 或 manifest.yaml
      └──service.tmpl
      └──shared.tmpl
      └──sub_service.tmpl
      └──types.tmpl
```
### 模板配置文件及模板配置文件

Snips模板文件包括 `service.tmpl`, `shared.tmpl`, `sub_service.tmpl`, `types.tmpl` 以及一个模板配置文件 `manifest`。

#### 模板配置文件

模板配置文件支持 json 格式和 yaml 格式，可命名为 `manifest.json` 或 `manifest.yaml`。该文件用于指定生成文件的路径，代码的命名风格，后缀名以及各模板文件的文件路径以及对应模板文件生成的代码文件的附加前缀和后缀。

配置文件格式：

- `output` 集合用于配置所有输出文件的共同属性。其中的 `file_naming` 集合用于配置输出文件的路径以及代码的命名规范。`file_naming` 集合下的`style` 属性规定了输出代码文件使用的命名法，可选驼峰命名法（camel_case) 或匈牙利命名法(snake_case)，`extension` 属性规定了输出文件的后缀名。`output` 集合有且只有一个。
- `template_files` 集合中的子集合规定了各个输出文件的特殊属性。每个子集合需要分别以模板文件名命名。子集合中的 `file_path` 属性指定模板文件的路径，`output_file_naming` 中的 `prefix` 属性规定了输出文件名附加前缀，`suffix`属性规定了输出文件名附加后缀。

例如：

json 格式:

``` json
{
  "output": {
    "file_naming": {
      "style": "camel_case",
      "extension": ".go"
    }
  },
  "template_files": {
    "service": {
      "file_path": "service.tmpl",
      "output_file_naming": {
        "prefix": "",
        "suffix": ""
      }
    },
    "sub_service": {
      "file_path": "sub_services.tmpl",
      "output_file_naming": {
        "prefix": "",
        "suffix": ""
      }
    },
    "types": {
      "file_path": "types.tmpl",
      "output_file_naming": {
        "prefix": "",
        "suffix": ""
      }
    }
  }
}
```

yaml 格式：

``` yaml
template:
  # Template files format
  # Validate formats: Go
  # Default: Go
  format: Go

# Output
output:
  file_naming:
    # Naming style to use in the output file.
    # Available styles: snake_case, camel_case
    # Default: snake_case
    # Example: bucket_acl (snake_case), BucketACL (camel_case)
    style: snake_case
    extension: .go

# Template files to read and execute.
#   Currently, there's are three of them.
template_files:
  # Shared template file.
  # This file will be concatenated with each other template to provide shared
  # nested template definitions.
  shared:
    # Relative file path to load the shared template file.
    # Default: shared.tmpl
    file_path: shared.tmpl

  # Service template file.
  # In this case, a file named "qs_qingstor_service.go" will be generated.
  service:
    # Relative file path to load service template file.
    # Default: service.tmpl
    file_path: service.tmpl
    # Naming options for output file.
    output_file_naming:
      prefix: qs_
      suffix: _service

  # Service template file.
  # In this case, multiple files named like "qs_bucket_sub_service.go"
  # will be generated.
  sub_service:
    # Relative file path to load sub service template file.
    # Default: sub_service.tmpl
    file_path: sub_service.tmpl
    output_file_naming:
      prefix: qs_
      suffix: _sub_service
  # Types template file.
  # In this case, a file named "qs_types.go" will be generated.
  types:
    # Relative file path to load types template file.
    # Default: types.tmpl
    file_path: types.tmpl
    output_file_naming:
      prefix: qs_
      suffix:

# Supporting files to copy directly.
supporting_files:
  - utils.go
  - utils_test.go
  - README.md
	
```

### 模板文件

模板文件用于读取 API 规范文件并将其组织为对应编程语言的代码文件。模板文件应使用 go 模板语言编写，后缀名为 ".tmpl"。Snips 模板文件包括 `service.tmpl`, `shared.tmpl`, `sub_service.tmpl`, `types.tmpl`。

#### service.tmpl

`service.tmpl` 是用于生成主服务代码的模板文件，Snips会根据这个模板文件单独生成一个以此主服务命名的代码文件。该代码文件需要包含的功能有：

- 主服务的初始化函数
- 与主服务操作对应的结构体或类
- 各结构体或类的输入输出函数

**例子：**以 go 语言 qingstor SDK ([qingstor-sdk-go](https://github.com/yunify/qingstor-sdk-go)) 为例，该模板文件用于生成 [qingstor.go](https://github.com/yunify/qingstor-sdk-go/tree/master/service/qingstor.go) 代码文件。在该 SDK 的模板文件 [service.tmpl](https://github.com/yunify/qingstor-sdk-go/blob/master/template/service.tmpl) 中，包含了 qingstor service 的结构体，初始化函数以及 qingstor 服务操作 `ListBucket` 的输入输出函数。

#### shared.tmpl

`shared.tmpl` 用于定义数据组织规范和格式的函数，是其它三个模板文件要调用的函数的集合。该模板函数需要大量地引用 Snips 中用于保存API规范数据的结构体，关于Snips中的结构体，在文档的后半部分会着重讲解。

#### sub_service.tmpl

`sub_service.tmpl` 是用于生成对应服务的子服务的模板文件，Snips 会根据该模板文件生成一个以该子服务命名的代码文件。该代码文件需要包含的代码有：

- 子服务的初始化函数
- 与子服务操作对应的结构体或类
- 各结构体或类的输入输出函数

**例子：**以 go 语言 qingstor SDK ([qingstor-sdk-go](https://github.com/yunify/qingstor-sdk-go)) 为例，该模板文件用于生成 [bucket.go](https://github.com/yunify/qingstor-sdk-go/tree/master/service/bucket.go) 代码文件。在该 SDK 的模板文件 [sub_service.tmpl](https://github.com/yunify/qingstor-sdk-go/blob/master/template/sub_service.tmpl) 中，包含了 qingstor bucket 的初始化函数以及 bucket 子服务操作的结构体及其输入输出函数。

#### types.tmpl

`types.tmpl` 是用于生成服务自定义类型代码文件的模板文件。这些自定义类型的变量用于存储服务的各个数据集合。Snips 会根据该模板文件生成以types命名的代码文件。

**例子：**以 go 语言 qingstor SDK ([qingstor-sdk-go](https://github.com/yunify/qingstor-sdk-go)) 为例，该模板文件用于生成 [types.go](https://github.com/yunify/qingstor-sdk-go/tree/master/service/types.go) 代码文件。在该 SDK 的模板文件 [types.tmpl](https://github.com/yunify/qingstor-sdk-go/blob/master/template/types.tmpl) 中，包含了 qingstor 主服务以及 bucket 子服务的自定义类型。

### Snips 中间状态

Snips 从 API 规范中读取的数据会根据模板文件的规定的格式生成为代码文件，而在 Snips 运行过程中，从 API 规范中读取的数据会根据各自的特点暂存在若八个结构体中。你需要理解这八个结构体的意义及其中各元素存储的内容才能写出有意义的模板文件。

**附注：**以下讲解的中间状态是根据其包含与被包含的关系由小到大排列的。

结构体构成：

|  结构体名  	|       元素       	|      元素类型(go)      	|
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
| StatusCode 	|       Code       	|           int          	|
|            	|    Description   	|         string         	|
|  Response  	|    StatusCodes   	|       *StatusCode      	|
|            	|      Headers     	|        *Property       	|
|            	|     Elements     	|        *Property       	|
|            	|       Body       	|        *Property       	|
|   Request  	|      Method      	|         string         	|
|            	|        URI       	|         string         	|
|            	|      Params      	|        *Property       	|
|            	|      Headers     	|        *Property       	|
|            	|     Elements     	|        *Property       	|
|            	|       Body       	|        *Property       	|
|  Operation 	|        ID        	|         string         	|
|            	|       Name       	|         string         	|
|            	|    Description   	|         string         	|
|            	| DocumentationURL 	|         string         	|
|            	|      Request     	|        *Request        	|
|            	|     Response     	|    map[int]*Response   	|
| SubService 	|        ID        	|         string         	|
|            	|       Name       	|         string         	|
|            	|    Properties    	|        *Property       	|
|            	|    Operations    	|  map[string]*Operation 	|
|   Service  	|    APIVersion    	|         string         	|
|            	|       Name       	|         string         	|
|            	|    Description   	|         string         	|
|            	|    Properties    	|        *Property       	|
|            	|    Operations    	|  map[string]*Operation 	|
|    Data    	|      Service     	|        *Service        	|
|            	|    SubServices   	| map[string]*SubService 	|
|            	|  CustomizedTypes 	|  map[string]*Property  	|

#### Property

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

#### StatusCode

`StatusCode` 用于存储状态信息，一般用于存储 HTTP(S) 的 response 字段。

- `Code` 存储着状态值，如 "404", "400" 等。
- `Description` 存储着该状态的状态信息描述，与其他的 `Description` 字段不同，该字段或许对用户有用，因而它应该被输出。

#### Request

`Request` 存储着 HTTP(S) 的 request 字段。该字段内指明了某个 request 请求需要的所有必须、必要信息。

- `Method` 存储着 request 字段的请求方法，如 `GET`, `POST`, `PUT` 等。
- `URI` 存储着目标地址的相对位置。
- `Params` 存储着 request 字段的参数。对应 API 规范中的 `Parameters` 字段，这些参数包含某个 request 字段所必须的 `Headers` 和 `Body` 内的参数，如 `content-length`, `bucket-name`, `access-key-id` 等。
- `Headers` 存储着 request 字段的头，但是在头内一般只指明该头需要什么参数，具体参数在 `Params` 中。
- `Elements` 存储着 request 字段的元素。
- `Body` 存储着 request 字段的 Body，但是在 Body 内一般只指明该 Body 需要什么参数，具体参数在 `Params` 中。

#### Response

`Response` 存储着 HTTP(S) 的 response 字段。该字段指明了在某个 request 请求之后获得的 response 字段的各个参数及其格式。

- `StatusCodes` 存储着 reponse 字段的状态信息。
- `Headers` 存储着 reponse 字段的头，头内指明了各个返回的参数。
- `Elements` 存储着 reponse 字段的元素。
- `Body` 存储着 reponse 字段的 Body。

#### Operation

`Operation` 保存着服务以及子服务的操作。

- `ID` 保存着操作的 ID。一般来说，操作 ID 与其名字相同，但是不同操作的 ID 具有唯一性。
- `Name` 保存着操作的名字。
- `Description` 保存着操作的描述，一般来说，这个字段只作为说明使用。
- `DocumentationURL` 保存着对应操作的文档的相对位置。
- `Request` 保存着对应操作的 request 字段，该元素内标明了本次操作请求需要的元素参数等。
- `Response` 保存着对应操作的 response 字段，该元素内著名了本次操作请求后的回复的格式、元素以及参数。

#### SubService

`SubService` 保存着子服务数据。

- `ID` 保存着子服务 ID，一般来说子服务 ID 与其名字相同，但是不同子服务的 ID 具有唯一性。
- `Name` 保存着子服务的名字。
- `Properties` 保存着该子服务的配置，包括子服务内的全局参数、属性等。
- `Operation` 保存着子服务的操作。

#### Service

`Service` 保存着主服务数据。

- `APIVersion` 保存着主服务的版本，你可以在生成代码文件时，将这个字段输出到说明中以提醒使用者 SDK 的版本。
- `Name` 保存着主服务的名字，因为主服务只有一个，所以其名字本身就具有唯一性，所以不需要 ID 字段。
- `Description` 保存着对主服务的描述，这个字段只作为说明使用。
- `Properties` 保存着主服务内的全局变量以及属性等。
- `Operations` 保存着主服务的操作。

#### Data

`Data` 是 Snips 中最大的数据结构，所有的 API 规范都在该字段下。

- `Service` 保存着主服务数据。
- `SubService` 保存着子服务数据。
- `CustomizedTypes` 保存着所有的 SDK 的自定义类型。