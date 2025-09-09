# DeactivateClassData

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **string** | class_code + locale | 
**Type** | **string** |  | 
**Attributes** | [**DeactivateClassDataAttributes**](DeactivateClassDataAttributes.md) |  | 

## Methods

### NewDeactivateClassData

`func NewDeactivateClassData(id string, type_ string, attributes DeactivateClassDataAttributes, ) *DeactivateClassData`

NewDeactivateClassData instantiates a new DeactivateClassData object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewDeactivateClassDataWithDefaults

`func NewDeactivateClassDataWithDefaults() *DeactivateClassData`

NewDeactivateClassDataWithDefaults instantiates a new DeactivateClassData object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *DeactivateClassData) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *DeactivateClassData) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *DeactivateClassData) SetId(v string)`

SetId sets Id field to given value.


### GetType

`func (o *DeactivateClassData) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *DeactivateClassData) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *DeactivateClassData) SetType(v string)`

SetType sets Type field to given value.


### GetAttributes

`func (o *DeactivateClassData) GetAttributes() DeactivateClassDataAttributes`

GetAttributes returns the Attributes field if non-nil, zero value otherwise.

### GetAttributesOk

`func (o *DeactivateClassData) GetAttributesOk() (*DeactivateClassDataAttributes, bool)`

GetAttributesOk returns a tuple with the Attributes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAttributes

`func (o *DeactivateClassData) SetAttributes(v DeactivateClassDataAttributes)`

SetAttributes sets Attributes field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


