# UpdatePlaceVerifyData

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | [**uuid.UUID**](uuid.UUID.md) |  | 
**Type** | **string** |  | 
**Attributes** | [**UpdatePlaceVerifyDataAttributes**](UpdatePlaceVerifyDataAttributes.md) |  | 

## Methods

### NewUpdatePlaceVerifyData

`func NewUpdatePlaceVerifyData(id uuid.UUID, type_ string, attributes UpdatePlaceVerifyDataAttributes, ) *UpdatePlaceVerifyData`

NewUpdatePlaceVerifyData instantiates a new UpdatePlaceVerifyData object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewUpdatePlaceVerifyDataWithDefaults

`func NewUpdatePlaceVerifyDataWithDefaults() *UpdatePlaceVerifyData`

NewUpdatePlaceVerifyDataWithDefaults instantiates a new UpdatePlaceVerifyData object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *UpdatePlaceVerifyData) GetId() uuid.UUID`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *UpdatePlaceVerifyData) GetIdOk() (*uuid.UUID, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *UpdatePlaceVerifyData) SetId(v uuid.UUID)`

SetId sets Id field to given value.


### GetType

`func (o *UpdatePlaceVerifyData) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *UpdatePlaceVerifyData) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *UpdatePlaceVerifyData) SetType(v string)`

SetType sets Type field to given value.


### GetAttributes

`func (o *UpdatePlaceVerifyData) GetAttributes() UpdatePlaceVerifyDataAttributes`

GetAttributes returns the Attributes field if non-nil, zero value otherwise.

### GetAttributesOk

`func (o *UpdatePlaceVerifyData) GetAttributesOk() (*UpdatePlaceVerifyDataAttributes, bool)`

GetAttributesOk returns a tuple with the Attributes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAttributes

`func (o *UpdatePlaceVerifyData) SetAttributes(v UpdatePlaceVerifyDataAttributes)`

SetAttributes sets Attributes field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


