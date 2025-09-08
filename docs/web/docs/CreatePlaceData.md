# CreatePlaceData

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **string** | place id | 
**Type** | **string** |  | 
**Attributes** | [**CreatePlaceDataAttributes**](CreatePlaceDataAttributes.md) |  | 

## Methods

### NewCreatePlaceData

`func NewCreatePlaceData(id string, type_ string, attributes CreatePlaceDataAttributes, ) *CreatePlaceData`

NewCreatePlaceData instantiates a new CreatePlaceData object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCreatePlaceDataWithDefaults

`func NewCreatePlaceDataWithDefaults() *CreatePlaceData`

NewCreatePlaceDataWithDefaults instantiates a new CreatePlaceData object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *CreatePlaceData) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *CreatePlaceData) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *CreatePlaceData) SetId(v string)`

SetId sets Id field to given value.


### GetType

`func (o *CreatePlaceData) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *CreatePlaceData) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *CreatePlaceData) SetType(v string)`

SetType sets Type field to given value.


### GetAttributes

`func (o *CreatePlaceData) GetAttributes() CreatePlaceDataAttributes`

GetAttributes returns the Attributes field if non-nil, zero value otherwise.

### GetAttributesOk

`func (o *CreatePlaceData) GetAttributesOk() (*CreatePlaceDataAttributes, bool)`

GetAttributesOk returns a tuple with the Attributes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAttributes

`func (o *CreatePlaceData) SetAttributes(v CreatePlaceDataAttributes)`

SetAttributes sets Attributes field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


