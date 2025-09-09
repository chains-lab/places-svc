# ClassesCollection

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Data** | [**[]ClassData**](ClassData.md) |  | 
**Links** | [**PaginationData**](PaginationData.md) |  | 

## Methods

### NewClassesCollection

`func NewClassesCollection(data []ClassData, links PaginationData, ) *ClassesCollection`

NewClassesCollection instantiates a new ClassesCollection object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewClassesCollectionWithDefaults

`func NewClassesCollectionWithDefaults() *ClassesCollection`

NewClassesCollectionWithDefaults instantiates a new ClassesCollection object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetData

`func (o *ClassesCollection) GetData() []ClassData`

GetData returns the Data field if non-nil, zero value otherwise.

### GetDataOk

`func (o *ClassesCollection) GetDataOk() (*[]ClassData, bool)`

GetDataOk returns a tuple with the Data field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetData

`func (o *ClassesCollection) SetData(v []ClassData)`

SetData sets Data field to given value.


### GetLinks

`func (o *ClassesCollection) GetLinks() PaginationData`

GetLinks returns the Links field if non-nil, zero value otherwise.

### GetLinksOk

`func (o *ClassesCollection) GetLinksOk() (*PaginationData, bool)`

GetLinksOk returns a tuple with the Links field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLinks

`func (o *ClassesCollection) SetLinks(v PaginationData)`

SetLinks sets Links field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


