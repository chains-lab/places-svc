# TimetableData

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | [**uuid.UUID**](uuid.UUID.md) | place id | 
**Type** | **string** |  | 
**Attributes** | [**TimetableDataAttributes**](TimetableDataAttributes.md) |  | 

## Methods

### NewTimetableData

`func NewTimetableData(id uuid.UUID, type_ string, attributes TimetableDataAttributes, ) *TimetableData`

NewTimetableData instantiates a new TimetableData object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewTimetableDataWithDefaults

`func NewTimetableDataWithDefaults() *TimetableData`

NewTimetableDataWithDefaults instantiates a new TimetableData object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *TimetableData) GetId() uuid.UUID`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *TimetableData) GetIdOk() (*uuid.UUID, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *TimetableData) SetId(v uuid.UUID)`

SetId sets Id field to given value.


### GetType

`func (o *TimetableData) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *TimetableData) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *TimetableData) SetType(v string)`

SetType sets Type field to given value.


### GetAttributes

`func (o *TimetableData) GetAttributes() TimetableDataAttributes`

GetAttributes returns the Attributes field if non-nil, zero value otherwise.

### GetAttributesOk

`func (o *TimetableData) GetAttributesOk() (*TimetableDataAttributes, bool)`

GetAttributesOk returns a tuple with the Attributes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAttributes

`func (o *TimetableData) SetAttributes(v TimetableDataAttributes)`

SetAttributes sets Attributes field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


