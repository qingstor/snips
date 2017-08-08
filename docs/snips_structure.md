#Snips Structure

Snips is a code generation for QingStor SDK.When Snips works,Snips loads data from spec files and generate  code files depending on template files.

##Template Content Structure

When you use Snips,you should specify template content and its content structure is fixed as well.

You should manage template content as follow: 

``` sh
└──template
      └──manifest.json
      └──service.tmpl
      └──shared.tmpl
      └──sub_service.tmpl
      └──types.tmpl
```
##Introduction of Template Files' Functions

###manifest.json

`manifest.json` is used to specify output files' naming style,extension and tmpl files' name and file path.
 
 You May can manage `manifest.json` as follow:
 
``` Json
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
###service.tmpl

`service.tmpl` is used to generate the code file  towards to service operations.

As an example, `service.tmpl` generated `qingstor` code file of QingStor SDK.(e.g. [qingstor.go](https://github.com/yunify/qingstor-sdk-go/blob/master/service/qingstor.go))

In this template,you should specify the template to generate follow codes:

- Initialization of service.
- Functions pointing to service.

###shared.tmpl

`shared.tmpl` should define shared functions which other template files need.

###sub_service.tmpl

`sub_service.tmpl`is used to generate the code file  towards to sub service operations.

As an example, `sub_service.tmpl` generated `bucket` code file of QingStor SDK.([bucket.go](https://github.com/yunify/qingstor-sdk-go/blob/master/service/bucket.go))

In this template,you should specify the template to generate follow codes:

- Initalization of sub service.
- Functions pointing to sub service.

###types.tmpl

`types.tmpl` is used to generate types code files.

As an example,`types.tmpl` generated `types` code file of QingStor SDK.([types.go](https://github.com/yunify/qingstor-sdk-go/tree/master/service/types.go)).

In this template,you should specify the template to generate follow codes:

- All types functions.

##Introduction of Snips Intermediate Variables

The intermediate variables are used to store data loading from spec files,and then they will be used to generate code files according to the template files.

###Intermediate Variables Structure

``` golang
type Data struct {
	Service         *Service
	SubServices     map[string]*SubService
	CustomizedTypes map[string]*Property
}

type Service struct {
	APIVersion  string
	Name        string
	Description string
	Properties  *Property
	Operations  map[string]*Operation
}

type SubService struct {
	ID         string
	Name       string
	Properties *Property
	Operations map[string]*Operation
}

type Operation struct {
	ID               string
	Name             string
	Description      string
	DocumentationURL string
	Request          *Request
	Response         map[int]*Response
}

type Request struct {
	Method   string
	URI      string
	Params   *Property
	Headers  *Property
	Elements *Property
	Body     *Property
}

type Response struct {
	StatusCodes *StatusCode
	Headers     *Property
	Elements    *Property
	Body        *Property
}

type StatusCode struct {
	Code int
	Description string
}

type Property struct {
	ID               string
	Name             string
	Description      string
	Type             string
	ExtraType        string
	Format           string
	CollectionFormat string
	Enum             []string
	Default          string
	IsRequired       bool
	Properties       map[string]*Property
}

```
###Introduction of Parameters In Structures

####Property

``` golang
type Property struct {
	ID               string
	Name             string
	Description      string
	Type             string
	ExtraType        string
	Format           string
	CollectionFormat string
	Enum             []string
	Default          string
	IsRequired       bool
	Properties       map[string]*Property
}
```
`Property` is the minimum structure.It is used to store data of elements from larger structures.It is usually from Json file's `properties` section.
- `ID` stores the property's ID. 
- `Name` stores the property's name.
- `Description` stores the property's description.
- `Type` stores the property's type.
- `ExtraType` stores the property's extra type.
- `Format` stores the property's format.
- `CollectionFormat` stores the property's collection format.
- `Enum` limits the range of values or types for `Property`.
- `Default` stores the property's default value.
- `IsRequired` specifies if the `Property` is necessary.`True` for necessary and `False` for not.
- `Properties` stores the  son properties of `Property`.

####StatusCode

``` golang
type StatusCode struct {
	Code int
	Description string
}
```
`StatusCode` is used to store the data of status code.

- `Code` stores the code number of status.
- `Description` stores the description of status.

####Response

``` golang
type Response struct {
	StatusCodes *StatusCode
	Headers     *Property
	Elements    *Property
	Body        *Property
}
```
` Response` stores the data of response section.

- `StatusCodes` stores the information of status in a response section.
- `Headers` stores the header's data of a response section.
- `Elements` stores elements of a response section.
- `Body` stores the body of a response section.

####request

``` golang
type Request struct {
	Method   string
	URI      string
	Params   *Property
	Headers  *Property
	Elements *Property
	Body     *Property
}
```
`Request` stores the data of request section.

- `Method` stores a request method of a request section.
- `URI` stores a native target path of  a request section.
- `Params` stores parameters of a request section.
- `Headers` stores headers of a request section.
- `Elements` stores elements of a request section.
- `Body` stores the body of a request section.

####Operation

``` golang
type Operation struct {
	ID               string
	Name             string
	Description      string
	DocumentationURL string
	Request          *Request
	Response         map[int]*Response
}
```
`Operation` stores the data of operations towards to services or sub services.

- `ID` stores the operation's id.
- `Name` stores the operation's name.
- `Description` stores the operation's description.
- `DocumentationURL` stores the docunment's URL which used to explain the operation.
- `Request` stores a operation's request sections.
- `Response` stores a operation's response sections.

####SubService

``` golang
type SubService struct {
	ID         string
	Name       string
	Properties *Property
	Operations map[string]*Operation
}
```
`SubService` stores the data of sub service.
- `ID` stores the sub service's ID.
- `Name` stores the sub service's name.
- `Properties` stores the sub service's properties,such as variables.
- `Operations` stores the sub service's operations.

####Service

``` golang
type Service struct {
	APIVersion  string
	Name        string
	Description string
	Properties  *Property
	Operations  map[string]*Operation
}
```
`Service` stores the data of service.
- `APIVersion` stores the version of the API.
- `Name` stores the service's name.
- `Description` stores the service's description.
- `Properties` stores the service's properties.
- `Operations` stores the service's operations.

####Data
``` golang
type Data struct {
	Service         *Service
	SubServices     map[string]*SubService
	CustomizedTypes map[string]*Property
}
```
`Data` is the miximum data struction.And all data from spec files is in this structure.

- `Service` stores the service.
- `SubServices` stores all the sub services.
- `CustomizedTypes` stores all the customized types.
