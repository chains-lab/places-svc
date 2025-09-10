# Place

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Data** | [**PlaceData**](PlaceData.md) |  | 
**Included** | [**[]TimetableData**](TimetableData.md) | Array of related resources | 

## Methods

### NewPlace

`func NewPlace(data PlaceData, included []TimetableData, ) *Place`

NewPlace instantiates a new Place object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewPlaceWithDefaults

`func NewPlaceWithDefaults() *Place`

NewPlaceWithDefaults instantiates a new Place object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetData

`func (o *Place) GetData() PlaceData`

GetData returns the Data field if non-nil, zero value otherwise.

### GetDataOk

`func (o *Place) GetDataOk() (*PlaceData, bool)`

GetDataOk returns a tuple with the Data field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetData

`func (o *Place) SetData(v PlaceData)`

SetData sets Data field to given value.


### GetIncluded

`func (o *Place) GetIncluded() []TimetableData`

GetIncluded returns the Included field if non-nil, zero value otherwise.

### GetIncludedOk

`func (o *Place) GetIncludedOk() (*[]TimetableData, bool)`

GetIncludedOk returns a tuple with the Included field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIncluded

`func (o *Place) SetIncluded(v []TimetableData)`

SetIncluded sets Included field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


