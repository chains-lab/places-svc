# ClassLocaleData

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **string** | class_code + locale | 
**Type** | **string** |  | 
**Attributes** | [**ClassLocaleDataAttributes**](ClassLocaleDataAttributes.md) |  | 

## Methods

### NewClassLocaleData

`func NewClassLocaleData(id string, type_ string, attributes ClassLocaleDataAttributes, ) *ClassLocaleData`

NewClassLocaleData instantiates a new ClassLocaleData object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewClassLocaleDataWithDefaults

`func NewClassLocaleDataWithDefaults() *ClassLocaleData`

NewClassLocaleDataWithDefaults instantiates a new ClassLocaleData object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *ClassLocaleData) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *ClassLocaleData) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *ClassLocaleData) SetId(v string)`

SetId sets Id field to given value.


### GetType

`func (o *ClassLocaleData) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *ClassLocaleData) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *ClassLocaleData) SetType(v string)`

SetType sets Type field to given value.


### GetAttributes

`func (o *ClassLocaleData) GetAttributes() ClassLocaleDataAttributes`

GetAttributes returns the Attributes field if non-nil, zero value otherwise.

### GetAttributesOk

`func (o *ClassLocaleData) GetAttributesOk() (*ClassLocaleDataAttributes, bool)`

GetAttributesOk returns a tuple with the Attributes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAttributes

`func (o *ClassLocaleData) SetAttributes(v ClassLocaleDataAttributes)`

SetAttributes sets Attributes field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


