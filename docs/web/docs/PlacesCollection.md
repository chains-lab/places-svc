# PlacesCollection

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Data** | [**[]PlaceData**](PlaceData.md) |  | 
**Links** | [**PaginationData**](PaginationData.md) |  | 

## Methods

### NewPlacesCollection

`func NewPlacesCollection(data []PlaceData, links PaginationData, ) *PlacesCollection`

NewPlacesCollection instantiates a new PlacesCollection object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewPlacesCollectionWithDefaults

`func NewPlacesCollectionWithDefaults() *PlacesCollection`

NewPlacesCollectionWithDefaults instantiates a new PlacesCollection object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetData

`func (o *PlacesCollection) GetData() []PlaceData`

GetData returns the Data field if non-nil, zero value otherwise.

### GetDataOk

`func (o *PlacesCollection) GetDataOk() (*[]PlaceData, bool)`

GetDataOk returns a tuple with the Data field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetData

`func (o *PlacesCollection) SetData(v []PlaceData)`

SetData sets Data field to given value.


### GetLinks

`func (o *PlacesCollection) GetLinks() PaginationData`

GetLinks returns the Links field if non-nil, zero value otherwise.

### GetLinksOk

`func (o *PlacesCollection) GetLinksOk() (*PaginationData, bool)`

GetLinksOk returns a tuple with the Links field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLinks

`func (o *PlacesCollection) SetLinks(v PaginationData)`

SetLinks sets Links field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


