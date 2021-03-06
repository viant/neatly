# Neatly - a neat format for representing nested structured data.

[![GoReportCard](https://goreportcard.com/badge/github.com/viant/neatly)](https://goreportcard.com/report/github.com/viant/neatly)
[![GoDoc](https://godoc.org/github.com/viant/neatly?status.svg)](https://godoc.org/github.com/viant/neatly)

This library is compatible with Go 1.8+

Please refer to [`CHANGELOG.md`](CHANGELOG.md) if you encounter breaking changes.

- [Motivation](#Motivation)
- [Getting started](https://github.com/adrianwit/neatly-introduction)
- [Neatly](#Neatly)
- [User Defined Function](#udf)
- [Usage](#Usage)
- [License](#License)
- [Credits and Acknowledgements](#Credits-and-Acknowledgements)



<a name="Motivation"></a>
# Motivation 

Most of the data-driven application use extensively nested structured data to instruct how the application should behave.
There is various data format that comes handy like vanilla CSV, XML, JSON, YAML which are more or less human-friendly. 
Once the data grows, organizing it with these formats becomes the mundane task.
As a matter of fact, in the absence of UI or/and data store (either RDBMS or NoSQL) data organization is challenging.
Neatly has been designing in mind to address all these concerns promoting flexibility, reusability, data cohesion, delegation 
and data organization


<a name="Neatly"></a>
# Neatly

Neatly is a neat format for representing nested structured data, with a simple tabular approach.

Neatly document uses tabular format thus can be easily stored in CSV or other delimited formats,

The first column in a row represents an **object tag**,  followed by columns defining the object fields.
Next row/s would define object tag values, in this case, the first column would be left empty.
You can think of an object tag as it is an object definition, whereas object values would create one or more object instances.
An object tag can represent a single instance or an instance element within an array.
In the latter case, object tag would be created with **[]** prefix.

The very first object tag in the document represents a **root object**,  it can be mapped to either to user-defined struct type or a simple map.  


The field defines a path to the leaf object of an object tag.

Field name should start with an uppercase letter to map to the underlying data structure.

It can be prefixed with 
   1) **square bracket '[]'** to denote that field is an array, all rows below will be elements for the array, unless there is empty line
   2) **slash '/'** to denote that field belongs to root object rather then preceding tag object
   3) **colon ':'**  or expression stating with **lowercase** character denotes a virtual object field that can be used as data substitution source referenced with dollar($) sign expression.
   
   
On top of that field definition can use dot (.) to define nested object of object or array type. 
For instance 
   1) field1 defines 'field1' field.
   2) []field2.attr2 defines 'field2' filed  of an array of an object type, which has 'attr2' field.
   3) field3.[]a.attr3 defines 'field3' field of an object type, which has another 'a' field  of an array of object type with 'attr3' field.
   4) /data.[]tables defines in the root object 'data' field of an object type which has another 'tables' field of an array type.
   5) :v.attr5 defines a virtual object which has 'v' field  of an object type which has another 'attr5' field .


**This** keyword used in field expression expands corresponding key/value directly to the current object tag.


 


## Neatly capabilities

Use cases 
### Basic data structure with repeated fields.
 
 Take as example the following data structure:
 
```json

{
  "UseCase": "case 1",
  "Requests": [
    {
      "URL": "http://127.0.0.1/test1",
      "Method": "GET",
      "Cookies": {
        "Cookie1": "value1",
        "Cookie2": "value2"
      }
    },
    {
      "URL": "http://127.0.0.1/test2",
      "Method": "GET",
      "Cookies": {
        "Cookie1": "value1",
        "Cookie2": "value2"
      }
    }
  ],
  "Expect": [
    {
      "StatusCode": 200,
      "Body": "test1Content"
    },
    {
      "StatusCode": 404
    }
  ]
}
 
```

Neatly tabular representation.

 
| Root | UseCase | []Requests.Method  | []Requests.URL| []Requests.Cookies | []Expect.StatusCode | []Expect.Body | Comments |
| --- | --- | --- | --- | --- |--- | --- | --- |
| | case 1 |  GET | http://127.0.0.1/test1 | {"Cookie1":"value1", "Cookie1":"value2"}   |  200  | testContent | {!123!} |
| | | GET | http://127.0.0.1/test2 | {"Cookie1":"value1", "Cookie1":"value2"}   |  404  |  | |

  
In this case we have only one object tag called root which has simple and repeated fields.  
Note that {} or [] prefix, sufix in an object value converts value to an object or array respectively.
You can escape **{** ... **}**  with **{!** ... **!}** quote 
or  **[** ... **]** with **[!** ... **!]** to represent value as text instead.
 
###  One to many with forward object tag reference use case
 
   Take as example the following data structure:
  
   ```json
{
  "CreateTime": "2017-10-23 10:00",
  "Orders": [
    {
      "Id": 1,
      "Name": "Order 1",
      "SubTotal": 100.0,
      "LineItems": [
        {
          "Product": "Magic Mouse",
          "Quantity": 5,
          "Price": 10.0
        },
        {
          "Position": 1,
          "Product": "TrackPad",
          "Quantity": 5,
          "Price": 10.0
        }
      ]
    },
    {
      "Id": 2,
      "Name": "Order 2",
      "SubTotal": 150.0,
      "LineItems": [
        {
          "Product": "Keyboard",
          "Quantity": 10,
          "Price": 10.0
        },
        {
          "Product": "TrackPad",
          "Quantity": 5,
          "Price": 10.0
        }
      ]
    }
  ]
}

```

Neatly tabular representation.


| Root | CreateTime | Orders | | |
| --- | --- | ---| --- | --- |
| | 2017-10-23 10:00 | %Orders |
|**[]Orders**| **Id** | **Name** | **LineItems** | **SubTotal** |
| |1 | Order 1 | %LineItems1 | 100 |
|**[]LineItems1**| **Product** | **Quantity** | **Price** | |
| |Magic Mouse| 5 | 10.0 ||
| |TrackPad| 5 | 10.0 ||
|**[]Orders**| **Id** | **Name** | **LineItems** | **SubTotal** |
| |1 | Order 2 | %LineItems2 | 150 |
|**[]LineItems2**| **Product** | **Quantity** | **Price** | **Note** |
| |Keyboard| 10 | 10.0 | Requested extended keyboard |
| |TrackPad| 5 | 10.0 ||

In this case we have 5 tags each defining its own objects, 
Note that percentage (%) prefixed object's value will be substitute with the object matching tag's value.
**Percentage (%)** denotes forward reference, which means that referencing tag definition takes places in the following rows.

You can escape **%** ... **%%**.

### Root object tag field with data cohesion use case.

   ```json
{
  "Bonus": [
    {
      "EmpNo": 1,
      "Name": "Smith",
      "Amount": 10000
    },
    {
      "EmpNo": 2,
      "Name": "Kowalczyk",
      "Amount": 8000
    },
    {
      "EmpNo": 3,
      "Name": "Schmidt",
      "Amount": 4000
    }
  ],
  "Merits": [
    {
      "EmpNo": 1,
      "Description": "Increase sales by 400%"
    },
    {
      "EmpNo": 2,
      "Description": "Reduced cost by 30%"
    },
    {
      "EmpNo": 3,
      "Description": "Improve resource reusibility by 40%"
    }
  ]
}

```

Neatly tabular representation.


| SomeName | Merits | Created | |
| --- | --- | --- | --- |
| |	%Merits| 2017-10-10 | |
|**[]Merits**|**Empno**| **Description** | **/[]Bonus** |
| |	1 |	Increased sales by 400% |	{"EmpNo":1, "Name":"Smith", "Amount":10000} |
| |	2 |	Reduced cost by 30% |	{"EmpNo":2, "Name":"Kowalczyk", "Amount":8000} |
| |	3 |	Improved resource reusibility by 40% |	{"EmpNo":3, "Name":"Schmidt", "Amount":4000} |

In this case "Bonus filed" on Merits tag is actual root object field reference. 
Cohesion has been achieved by placing  data related the same employee in the same row.
Note that root object uses **slash(/)** in the field name

<a name="vobject"></a>

###  Virtual Objects for data sharing and delegation use case.

The previous example address cohesion somehow, however we can see that empNo is repeated twice in the same row, 
Op top of that bonus value uses json notation, which may not be too elegant.


| Root|	Merits | Created | | | | |
| --- | --- | --- | --- | --- | --- | --- |
| |	%Merits| 2017-10-10 | | |  |  | 
|**[]Merits**|**Empno**| **Description** | **/[]Bonus** | **emp.EmpNo** | **emp.Name** | **:emp.Amount** |
| |	$emp.EmpNo |	Increased sales by 400% | $emp |	 1 | Smith | 10000  |
| |	$emp.EmpNo |	Reduced cost by 30% | $emp | 2 |  Kowalczyk | 8000 |
| |	$emp.EmpNo |	Improved resource reusibility by 40% | $emp | 3 | Schmidt |4000 |

In this case the virtual object emp was defined by 3 fields either lowercase or prefixed with colon (:) sing. 
In order to reference virtual  object dollar sign expression is being used.


<a name="tagiterator"></a>

### Loading repeated data with tags template.   


   Take as example the following data structure:
  
   ```json

{
  "Repeated": [
    {
      "Id": 1,
      "Name": "Name 01"
    },
    {
      "Id": 2,
      "Name": "Name 02"
    },
    {
      "Id": 3,
      "Name": "Name 03"
    },
    {
      "Id": 4,
      "Name": "Name 04"
    },
    {
      "Id": 5,
      "Name": "Name 05"
    }
  ]
}

```

Neatly tabular representation.

| Root|	Repeated |  |
| --- | --- | --- |
| |	%Repeated|  |
|**[]Repeated{1 .. 05}**|**Id**| **Name** | 
| |	 $index | Name $index  |



In this case Repeted tag uses **Tag Range** expression  **{  }**, which repeats all value rows as (max - min) from the range. 

#### Variable substitution


#### Variable substitution rules
1)  '$' path expression will be only substituted if path is present in the current context
2)  UDF wil be substituted only if it returns no error and in case it takes $ expression as parameter, expression path is present in context.


**$index** / ${index} is a special/reserved keyword in this context it would expand to 01 in the firs iteration followed by 02, 03, 04 and 05.
If there are more than one value row, all can be expanded to the same index value within the same iteration.


Note that number of digits in the upper bound range value creates a padding template.
1 .. 010 -> would '0' left pad an index to 3 digits
1 .. 00100 -> would '0' left pad an index to 5 digits.

**$tag**  expands to the current object tag.

**$tagId** expands to current root object Name if specified, followed by object tag, object tag index, and subpath if specified.

**$subPath** expands to current subPath if specified for the tag

**$path** expands to current path if sub path is specified for the tag

**$pathMatch** expands to resolved wildcard subPath fragment  


### Data delegation and loading external resources use case


   Take as example the following data structure:
  
   ```json
{
  "Students": [
    {
      "Id": 1,
      "Name": "Smith",
      "Scores": [
        {
          "Subject":"Math",
          "Score": 3.2
        },
        {
          "Subject":"English",
          "Score":3.5
        }
      ]
    },
    {
      "Id": 2,
      "Name": "Kowalczyk",
      "Scores": [
        {
          "Subject":"Math",
          "Score": 3.7
        },
        {
          "Subject":"English",
          "Score": 3.2
        }
      ]
    }
  ]
}
```

Neatly tabular representation.

| Root | Students | | |
| --- | ---| --- | --- |
| | %Students | |  |
|**[]Students**| **Id** | **Name** | **Scores** |
| | 1 | Smith | \@scores1.json | 
| | 2 | Kowalczyk | \@scores2.json | 
 

In this case scores are loaded from local json file.

External resource starts with  **(\@)**  or \# sing, resource URL can be relative, absolute path or a valid URL to any content. 
In case of json or yaml files, the content is treated as data structure.
**To escape '@' use '@@', or '#' use '##' .**

Where 
\@scores1.json

```json
  [
        {
           "Subject":"Math",
           "Score": 3.2
        },
        {
          "Subject":"English",
           "Score":3.5
        }
  ]

```

\@scores2.json

```json
  [
        {
           "Subject":"Math",
           "Score": 3.7
        },
        {
           "Subject":"English",
           "Score":3.2
        }
  ]

```

### Data delegation and loading external resources with subpath use case.

   Take as example the following data structure:
  
   ```json
{
  "Setup": {
    "MyDb": {
      "Customer": [
        {
          "ID": 1,
          "NAME": "Smith",
          "DAILY_CAP": 100,
          "OVERALL_CAP": 1000
        },
        {
          "ID": 2,
          "NAME": "Kowalczyk",
          "DAILY_CAP": 400,
          "OVERALL_CAP": 8000
        }
      ]
    }
  },
  "UseCases": [
    {
      "Id": "1",
      "Description": "use case 1"
    },
    {
      "Id": "2",
      "Description": "use case 2"
    }
  ]
}
```

Neatly tabular representation.

| Root | UseCases | | | | 
| --- | --- | --- | --- | --- |
|  |%UseCases | | | |
| **[]UseCases{1..002}**  | **Subpath** | **Id** | **Description** | **/Setup.MyDb.[]Customer** |
| | usecase7/${index} / | $index | \@use_case.txt |  \@customer.json |


Where

\@usecase7/001/use_case.txt

```text
use case 1
```

\@usecase7/001/customer.json
```json
{
  "ID": 1,
  "NAME": "Smith",
  "DAILY_CAP": 100,
  "OVERALL_CAP": 1000
}
```


\@usecase7/002/use_case.txt

```text
use case 2
```

\@usecase7/002/customer.json
```json
{
  "ID": 2,
  "NAME": "Kowalczyk",
  "DAILY_CAP": 400,
  "OVERALL_CAP": 8000
}
```


Note that external lookup first check if resource exists in subpath if not it will check if it exists in the 
same directory as a master file (source of the tabular data).


<a name="external_resource_substitution></a>
### External resources loading with content substitution and user defined function (udf) use case.


   Take as example the following data structure:
  
   ```json
{
  "Setup": {
    "MyDb": {
      "Customer": [
        {
          "ID": 1,
          "NAME": "Smith",
          "DAILY_CAP": 100,
          "OVERALL_CAP": 1000
        },
        {
          "ID": 2,
          "NAME": "Kowalczyk",
          "DAILY_CAP": 100,
          "OVERALL_CAP": 1000
        }
      ]
    }
  },
  "UseCases": [
    {
      "Id": "1",
      "Description": "use case 1"
    },
    {
      "Id": "2",
      "Description": "use case 2"
    }
  ]
}
```


Neatly tabular representation.


| Root | UseCases | | | | 
| --- | --- | --- | --- | --- |
|  |%UseCases | | | |
| **[]UseCases{1..002}**  | **Subpath** | **Id** | **Description** | **/Setup.MyDb.[]Customer** |
| | usecase8/${index} / | $index | \@use_case.txt|  \@customer.json\| {"dailyCap":100, "overallCap":1000} |


Where

\@usecase8/001/use_case.txt

```text
use case 1
```

\@usecase8/001/customer.json
```json
{
   "ID": 1,
   "NAME": "Smith",
   "DAILY_CAP": "$AsInt($dailyCap)",
   "OVERALL_CAP": "$AsInt($overallCap)"
 
}
```


\@usecase8/002/use_case.txt

```text
use case 2
```

\@usecase8/002/customer.json
```json
{
  "ID": 2,
  "NAME": "Kowalczyk",
   "DAILY_CAP": "$AsInt($dailyCap)",
   "OVERALL_CAP": "$AsInt($overallCap)"
}
```

In this case content of the @customer.json were substituted with dailyCap and overallCap values.
Since customer and cap data represents valid data structure, dailyCap and overallCap are substituted.


The following construct in a cell
 
    @customer.json | {"dailyCap":100, "overallCap":1000}

can be replaced with two or more external JSON or YAML resources (with space or pipe) providing desired data structure.

    @customer.json @cap.json

    
where 
@cap.json

```json
 {"dailyCap":100, "overallCap":1000}
```     


Note that xxx_CAP values were expanded within double quote thus they became of text type.
In order to convert then to int data type, AsInt user defined function is being used.



### External resources arguments inlining

In case where the whole content of subsequent arguments needs to be inlined, you can use  $argsX reference.

Neatly tabular alternative representation.


| Root | UseCases | | | | 
| --- | --- | --- | --- | --- |
|  |%UseCases | | | |
| **[]UseCases{1..2}**  | **Subpath** | **Id** | **Description** | **/Setup.MyDb.[]Customer** |
| | usecase9/${index} / | $index | \@use_case.txt |  \@customer.json\| {"DAILY_CAP":100, "OVERALL_CAP":1000} |


Where


\@usecase9/001/customer.json
```json
{
   "ID": 1,
   "NAME": "Smith",
   $args0
}
```



\@usecase9/002/customer.json
```json
{
  "ID": 2,
  "NAME": "Kowalczyk",
   $args0
}
```


The following special variables are available for substitution:
    
  1) $args{index} - piping content stripped from first and last characters. 
  2) $arg{index} - full piping content.
     
  Where  index corresponds to piping number starting with 0 


### Dynamic subpath discovery 

Subpath can use '*' as suffix to dynamically discover actual supath directory 


Take as example the following

Neatly tabular alternative representation.

| Root | UseCases | | | | 
| --- | --- | --- | --- | --- |
|  |%UseCases | | | |
| **[]UseCases{1..2}**  | **Subpath** | **Id** | **Description** | **/Setup.MyDb.[]Customer** |
| | usecase10/${index}* / | $index | \@use_case.txt |  \@customer.json\| {"DAILY_CAP":100, "OVERALL_CAP":1000} |




Where


\@usecase10/001_description1/customer.json
```json
{
   "ID": 1,
   "NAME": "Smith",
   $args0
}
```



\@usecase10/002_description2/customer.json
```json
{
  "ID": 2,
  "NAME": "Kowalczyk",
   $args0
}
```


<a name="udf"></a>
### User defined functions (udf)

The user defined system allowed to dynamically convert value from one form to another.
To invoke udf value of data structure has to start with **$FUNC_NAME(ARGS)** 
for instance $AsFloat("123"), $AsFloat($key1)

In order to define udf please use the following function signature:

```go
    type Udf func(interface{}, Map) (interface{}, error)

```

Build-in udf'

-  HasResource returns true if external resource exists
-  Md5 generates md5 for provided parameter
-  WorkingDirectory returns working directory joined with supplied sub path,  '../' is supported.
-  LoadNeatly loads neatly document as data structure.
-  Zip, takes []byte or string to compress it.
-  Unzip, takes []byte to uncompress it into []byte.
-  UnzipText takes []byte to uncompress it into string.
-  Markdown generate HTML for suppied markdown
-  Cat returns content of supplied filename

UDF Defined in [toolbox/data/udf](https://github.com/viant/toolbox/tree/master/data/udf)


### External resources loading with virtual object value substitution use case.


   Take as example the following data structure:
  
   ```json
{
  "Setup": {
    "MyDb": {
      "Customer": [
        {
          "ID": 1,
          "NAME": "Smith",
          "DAILY_CAP": "200",
          "OVERALL_CAP": "3000"
        },
        {
          "ID": 2,
          "NAME": "Kowalczyk",
          "DAILY_CAP": "100",
          "OVERALL_CAP": "1000"
        }
      ]
    }
  },
  "UseCases": [
    {
      "Id": "1",
      "Description": "use case 1"
    },
    {
      "Id": "2",
      "Description": "use case 2"
    }
  ]
}
```


Neatly tabular representation.


| Root | UseCases | | | | | |
| --- | --- | --- | --- | --- | ---| --- |
|  |%UseCases | | | | | |
| **[]UseCases**  | **Subpath** | **Id** | **Description** | **/Setup.MyDb.[]Customer** | **:data.DAILY_CAP** | **:data.OVERALL_CAP** |
| | usecase10/001 / | 1 | \@useCase.json |  \@customer.json\| $data | 200 | 3000 |
| | usecase10/002 / | 2 | \@useCase.json |  \@customer.json\| $data | 100 | 1000 |



### Loading external resources lookups

1) For valid URL, new resource if returned with owner resource credential
2) For asset starting  with /, a new file resource if returned with owner resource credential
3) For asset is a relative path then the following lookup are being used with the first successful to be used.
	a) owner resource path with subpath if provided and asset name
	b) owner resource path  without subpath and asset name
	c) Local/remoteResourceRepo and asset name


## Inline array use case


   Take as example the following data structure:
  
   ```json
{
  "Actions": [
    {
      "Send": {
        "Request": [
          {
            "Method": "GET",
            "URL": "http://127.0.0.1/path1"
          },
          {
            "Method": "GET",
            "URL": "http://127.0.0.1/path2"
          }
        ],
        "Udf": "MyUdf"
      },
      "Expect": [
        {
          "Code": 200
        },
        {
          "Code": 404
        }
      ]
    },
    {
      "Send": {
        "Request": [
          {
            "Method": "GET",
            "URL": "http://127.0.0.1/path3"
          },
          {
            "Method": "GET",
            "URL": "http://127.0.0.1/path4"
          }
        ],
        "Udf": "MyUdf"
      },
      "Expect": [
        {
          "Code": 404
        },
        {
          "Code": 200
        }
      ]
    },
    {
          "Send": {
            "Request": [
              {
                "Method": "GET",
                "URL": "http://127.0.0.1/path5"
              },
              {
                "Method": "GET",
                "URL": "http://127.0.0.1/path6"
              }
            ],
            "Udf": "MyUdf"
          },
          "Expect": [
            {
              "Code": 200
            },
            {
              "Code": 200
            }
          ]
        }
  ]
}
```


|Document |	Actions | | | |		
| --- | --- | --- | --- | --- |	
| |	%Actions | | | |			
|[]Actions|Send.Udf|Send.[]Requests.Method|Send.[]Requests.URL|[]Expect.Code|
| |	MyUdf|	GET	|http://127.0.0.1/path1 |200 | 
|  | |	GET |	http://127.0.0.1/path2 | 404 |
| | | | | |				
|[]Actions|Send.Udf|Send.[]Requests.Method|Send.[]Requests.URL|[]Expect.Code|
| |	MyUdf|	GET|	http://127.0.0.1/path3|	404|
| |	|	GET|	http://127.0.0.1/path4|	200|
|-| MyUdf|GET|http://127.0.0.1/path5|200|
|	| |	GET|	http://127.0.0.1/path6|	200|


If tag object field uses **[]** neatly will scan all rows down to add values to the array unless encounter next tag object or hyphen **'-'** in the first column.



### Accessing meta data 

For every object the following attributes will be set as object values, 
thus they should be treated as reserved keyword,  unless object needs to expose them.

 1) **Tag** name of currently processing tag.
 1) **TagId** tag id of currently processing tag (concatenation of Tag name, Group if specified, Index if within Tag Range, Subpath is specified).
 2) **TagIndex** index value if within tag range.
 3) **Subpath**  defines subpath.
 4) **PathMatch**  defines matched subpath.


### Comments

To prevent line loading into document , **//** can be used at the beginning of the line, optionally followed by some comments.  


<a name="Usage"></a>



## Usage 

[Download](https://github.com/viant/neatly/releases) or build latest neatly.



```bash
neatly -i=test/use_case1.csv -f=yaml
```


```text
 $ neatly -h
Usage of neatly:
  -f string
    	<output format> json or yaml (default "json")
  -i string
    	<neatly document path> 
  -v	neatly version
  -f string
    	<output format> json or yaml (default "json")
  -m	include meta data -m=true

```


To convert neatly document into go data structure.

```go

    import (
    		"github.com/viant/neatly"
        	"github.com/viant/toolbox/data"
        	"github.com/viant/toolbox/url"
    )


    var localAssetRepo, remoteAssetRepo string
	dao := neatly.NewDao(localAssetRepo, remoteAssetRepo, "yyyy-MM-dd h:mm:ss", nil)
	
	//user defined function
	var nilUdf = fun(source interface{}, state data.Map) (interface{}, error) {
		return nil, nil
	}
	
	var context = data.NewMap() //data is toolbox/data package
	//optionally register UDF
	context.Put("Nil", nilUdf)
    
	var targetObject = &MyStruct{} // or map[string]interface{}
	err := dao.Load(context, url.NewResource("mystruct.csv"), targetObject)
    if err != nil {
    	log.Fatal(err)
    }


```

	
<a name="License"></a>
## License

The source code is made available under the terms of the Apache License, Version 2, as stated in the file `LICENSE`.

Individual files may be made available under their own specific license,
all compatible with Apache License, Version 2. Please see individual files for details.


<a name="Credits-and-Acknowledgements"></a>

##  Credits and Acknowledgements

**Library Author:** Adrian Witas

