# PlaceData

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | [**uuid.UUID**](uuid.UUID.md) | place id | 
**Type** | **string** |  | 
**Attributes** | [**PlaceDataAttributes**](PlaceDataAttributes.md) |  | 
**Relationships** | [**PlaceRelationships**](PlaceRelationships.md) |  | 

## Methods

### NewPlaceData

`func NewPlaceData(id uuid.UUID, type_ string, attributes PlaceDataAttributes, relationships PlaceRelationships, ) *PlaceData`

NewPlaceData instantiates a new PlaceData object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewPlaceDataWithDefaults

`func NewPlaceDataWithDefaults() *PlaceData`

NewPlaceDataWithDefaults instantiates a new PlaceData object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *PlaceData) GetId() uuid.UUID`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *PlaceData) GetIdOk() (*uuid.UUID, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *PlaceData) SetId(v uuid.UUID)`

SetId sets Id field to given value.


### GetType

`func (o *PlaceData) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *PlaceData) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *PlaceData) SetType(v string)`

SetType sets Type field to given value.


### GetAttributes

`func (o *PlaceData) GetAttributes() PlaceDataAttributes`

GetAttributes returns the Attributes field if non-nil, zero value otherwise.

### GetAttributesOk

`func (o *PlaceData) GetAttributesOk() (*PlaceDataAttributes, bool)`

GetAttributesOk returns a tuple with the Attributes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAttributes

`func (o *PlaceData) SetAttributes(v PlaceDataAttributes)`

SetAttributes sets Attributes field to given value.


### GetRelationships

`func (o *PlaceData) GetRelationships() PlaceRelationships`

GetRelationships returns the Relationships field if non-nil, zero value otherwise.

### GetRelationshipsOk

`func (o *PlaceData) GetRelationshipsOk() (*PlaceRelationships, bool)`

GetRelationshipsOk returns a tuple with the Relationships field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRelationships

`func (o *PlaceData) SetRelationships(v PlaceRelationships)`

SetRelationships sets Relationships field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


