# Snips Usage Guide

Snips is a generic HTTP API code generation tool developed by QingCloud using OpenAPI specification (Swagger) v2.0 format to generate various API code for the QingCloud and QingStor SDKs. When it works,it will read data from the API specification file and generate API code files according to the template.

**Note:** API specification files are provide by QingCloud and QingStor, you can learn it or get the latest version files at [qingstor-api-specs](https://github.com/yunify/qingstor-api-specs).

**Note:** About template, you can learn it at [snips_template_for_English.md](https://github.com/yunify/snips/blob/master/docs/snips_template_for_English.md).

**Note:** About how to use Snips and its structure, you can learn it at  [Snips_README.md](https://github.com/yunify/snips/blob/master/README.md)

## Snips Intermediate State

The data from API specification will generate code files according to the template files, and when Snips works, The data will be stored in eight structures at the moment. Only if you understand the eight structures elements in the eight structures, can you product significative template files.

### Snips Intermediate State Transition Diagram

```
+--------------------------------------------------------+
|                                                        |
|                     read                               |
|   API specification ----> Data ---+                    |
|                                   |                    |
|                                   | store              |
|                                   v                    |
|                            Snips Structure             |
|                                   |                    |
|                                   | load               |
|                                   v                    |
|            Template ---------------------> Code Files  |
|                                                        |
+--------------------------------------------------------+
```

### The Structure of Snips Intermediate State

Structure hierarchy diagram:

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

**Note:** The structures of Snips mediacy explained below are sorted from large to small by According to its inclusion and included in the relationship.

### Data

|  Structure 	|      Element     	|      Element Type      	|
|:----------:	|:----------------:	|:----------------------:	|
|    Data    	|      Service     	|        *Service        	|
|            	|    SubServices   	| map[string]*SubService 	|
|            	|  CustomizedTypes 	|  map[string]*Property  	|

`Data` is the largest structure in Snips, and all the API specification are included.

- `Service` stores the data of service.
- `SubServices` stores the data of sub service.
- `CustomizedType` stores all the customized type of SDK.

### Service

|  Structure 	|      Element     	|      Element Type      	|
|:----------:	|:----------------:	|:----------------------:	|
|   Service  	|    APIVersion    	|         string         	|
|            	|       Name       	|         string         	|
|            	|    Description   	|         string         	|
|            	|    Properties    	|        *Property       	|
|            	|    Operations    	|  map[string]*Operation 	|

`Service` stores the data of service.

- `APIVersion` stores the version of service. When you generate code files, you can print this element to the introduction to remind users the version of SDK.
- `Name` stores service's name. Because there is only one service, the name itself is unique. For that, We don't need to add a element called `ID` to identify a service.
- `Description` stores the description to service. We just take it as a introduction.
- `Properties` stores global variables and properties of service.
- `Operations` stores service operations.

### SubService

|  Structure 	|      Element     	|      Element Type      	|
|:----------:	|:----------------:	|:----------------------:	|
| SubService 	|        ID        	|         string         	|
|            	|       Name       	|         string         	|
|            	|    Properties    	|        *Property       	|
|            	|    Operations    	|  map[string]*Operation 	|

`SubService` stores the data of sub service.

- `ID` stores sub service's ID. In general, sub service's ID is the same as its name, but different sub services' ID is unique.
- `Name` stores sub service's name.
- `Properties` stores sub service`s configuration, including global parameters, properties of sub service.
- `Operations` stores sub service's operations.

### Operation

|  Structure 	|      Element     	|      Element Type      	|
|:----------:	|:----------------:	|:----------------------:	|
|  Operation 	|        ID        	|         string         	|
|            	|       Name       	|         string         	|
|            	|    Description   	|         string         	|
|            	| DocumentationURL 	|         string         	|
|            	|      Request     	|        *Request        	|
|            	|     Response     	|    map[int]*Response   	|

`Operation` stores the operations of service and sub services.

- `ID` stores operation's ID. In general, operation's ID is the same as its name, but different operation's ID is unique.
- `Name` stores operation's name.
- `Description` stores description to the operation. In general, this element is only used for introduction.
- `DocumentationURL` stores native path of corresponding operation's document.
- `Request` stores operation's request section. This element indicates the element parameters required for this requestion.
- `Response` stores operation's response section. This element indicates the format, elements and parameters in response section get after this requestion.

### Request

 |  Structure 	|      Element     	|      Element Type      	|
|:----------:	|:----------------:	|:----------------------:	|
|   Request  	|      Method      	|         string         	|
|            	|        URI       	|         string         	|
|            	|      Params      	|        *Property       	|
|            	|      Headers     	|        *Property       	|
|            	|     Elements     	|        *Property       	|
|            	|       Body       	|        *Property       	|

`Request` stores the request section of HTTP(S). This element indicates information which needs in a requestion.

- `Method` stores request method in the request section. For example, `GET`, `POST`, `PUT`.
- `URI` stores native path of target.
- `Params` stores request section's parameters. Correspond to `Parameters` section, these parameters includes required parameters of `Headers` and `Body` which is needed in a requestion, `content-length`, `bucket-name`, `access-key-id` for instance.
- `Headers` stores the headers of request section, but this element only indicates which parameters the header needs, and specific parameters are in `Params`.
- `Elements` stores elements of request section.
- `Body` stores request section's body, but this element only indicates which parameters the body needs, and specific parameters are in `Params`.

### Response

 |  Structure 	|      Element     	|      Element Type      	|
|:----------:	|:----------------:	|:----------------------:	|
|  Response  	|    StatusCodes   	|       *StatusCode      	|
|            	|      Headers     	|        *Property       	|
|            	|     Elements     	|        *Property       	|
|            	|       Body       	|        *Property       	|

`Response` stores the response section of HTTP(S). This element indicates the format, elements and parameters in response section after a requestion.

- `StatusCodes` stores status information of response section.
- `Headers` stores the headers of response section.
- `Elements` stores response section's elements.
- `Body` stores response section's body.

### StatusCode

 |  Structure 	|      Element     	|      Element Type      	|
|:----------:	|:----------------:	|:----------------------:	|
| StatusCode 	|       Code       	|           int          	|
|            	|    Description   	|         string         	|

`StatusCode` stores status information, and it mainly used to store response section of HTTP(S).

- `Code` stores status value, "400", "404" for instance.
- `Description` stores the description to the status. Which is different from the other `Description` is that this element may be useful to users, so that it should be outputed.

### Property

 |  Structure 	|      Element     	|      Element Type      	|
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

`Property` is the mininum structure of Snips intermediate state. It used to store the properties and variables of service, operations,request sections, response sections in API specification.

- `ID` stores the ID of `Property`. In general, ID is the same as its name, but ID is unique among properties.
- `Name` stores the name of `Property`.
- `Description` stores the description to a `Property`. In general, this element is only used for code's note.
- `Type` stores the type of `Property`. For that `Property` used to store properties and variables, `Type` stores the type of properties or variables. you can output the type of `Property` or perform the necessary type conversion according to this element.The types in this element includes `object`, `array`, `map`, `any`. `any` means that the type of the property isn't stored in this element, and you need to load it in `ExtraType`.
- `ExtraType` stores the other unusual parameter types, including `string`, `boolean`, `integer`, `long`, `timestamp`, `binay`. This element also stores the three parameter types of `Type`. That is to say that the type in `Type` may be `any`, but there always exists a definite type in `ExtraType`.
- `Format` stores the format of `Property`. For example, which timing mode `timestamp` uses is specified in `format`.
- `CollectionFormat`stores collection's format. If elements stored in a property can have multiple formats, then these formats will be stored in the collection of this element.
- `Enum` stores variables in the collection of `Property`. if a element has multiple values and each value has the same priority, then these values will be stored in `Enum` disorderly.
- `Default` stores default value of `Property`.
- `IsRequired` indicates if a `Property` is necessary, value is `true` or `false`. Not all of the `Property` must be output, you can according to this element to determine whether the` property` should be output to the code file. 
- `Properties` is itself nested because a `Property` may also contain other multiple sub properties.