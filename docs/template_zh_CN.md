# 模板

模板文件由使用者创建，模板文件采用 Go 语言的 text/template 包要求的格式。你需要按照你的编程语言规范和格式利用 API 规范中的数据编写你的模板。

## 模板目录结构

该模板的内的文件名都必须是固定的，你需要按照下面的目录结构安排你的模板目录：

``` sh
└──template
      └──manifest.json 或 manifest.yaml
      └──service.tmpl
      └──shared.tmpl
      └──sub_service.tmpl
      └──types.tmpl
```
## 模板文件及模板配置文件

Snips模板文件包括 `service.tmpl`, `shared.tmpl`, `sub_service.tmpl`, `types.tmpl` 以及一个模板配置文件 `manifest`。

### 模板配置文件

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

模板文件用于读取 API 规范文件并将其组织为对应编程语言的代码文件。模板文件应使用 Go 语言的 text/template 包格式要求编写，后缀名为 ".tmpl"。Snips 模板文件包括 `service.tmpl`, `shared.tmpl`, `sub_service.tmpl`, `types.tmpl`。

#### service.tmpl

`service.tmpl` 是用于生成主服务代码的模板文件，Snips 会根据这个模板文件单独生成一个以此主服务命名的代码文件。该代码文件需要包含的功能有：

- 主服务的初始化函数
- 与主服务操作对应的结构体或类
- 各结构体或类的输入输出函数

**例子：** 以 QingStor Go SDK ([qingstor-sdk-go](https://github.com/yunify/qingstor-sdk-go)) 为例，该模板文件用于生成 [qingstor.go](https://github.com/yunify/qingstor-sdk-go/tree/master/service/qingstor.go) 代码文件。在该 SDK 的模板文件 [service.tmpl](https://github.com/yunify/qingstor-sdk-go/blob/master/template/service.tmpl) 中，包含了 qingstor service 的结构体，初始化函数以及 QingStor 服务操作 `ListBucket` 的输入输出函数。

#### shared.tmpl

`shared.tmpl` 用于定义数据组织规范和格式的函数，是其它三个模板文件要调用的函数的集合。该模板函数需要大量地引用 Snips 中用于保存API规范数据的结构体。

#### sub_service.tmpl

`sub_service.tmpl` 是用于生成对应服务的子服务的模板文件，Snips 会根据该模板文件生成一个以该子服务命名的代码文件。该代码文件需要包含的代码有：

- 子服务的初始化函数
- 与子服务操作对应的结构体或类
- 各结构体或类的输入输出函数

**例子：** 以 QingStor Go SDK ([qingstor-sdk-go](https://github.com/yunify/qingstor-sdk-go)) 为例，该模板文件用于生成 [bucket.go](https://github.com/yunify/qingstor-sdk-go/tree/master/service/bucket.go) 代码文件。在该 SDK 的模板文件 [sub_service.tmpl](https://github.com/yunify/qingstor-sdk-go/blob/master/template/sub_service.tmpl) 中，包含了 QingStor bucket 的初始化函数以及 bucket 子服务操作的结构体及其输入输出函数。

#### types.tmpl

`types.tmpl` 是用于生成服务自定义类型代码文件的模板文件。这些自定义类型的变量用于存储服务的各个数据集合。Snips 会根据该模板文件生成以types命名的代码文件。

**例子：** 以 QingStor Go SDK ([qingstor-sdk-go](https://github.com/yunify/qingstor-sdk-go)) 为例，该模板文件用于生成 [types.go](https://github.com/yunify/qingstor-sdk-go/tree/master/service/types.go) 代码文件。在该 SDK 的模板文件 [types.tmpl](https://github.com/yunify/qingstor-sdk-go/blob/master/template/types.tmpl) 中，包含了 QingStor 主服务以及 bucket 子服务的自定义类型。
