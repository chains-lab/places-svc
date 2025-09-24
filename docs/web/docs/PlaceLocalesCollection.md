# PlaceLocalesCollection

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Data** | [**[]RelationshipDataObject**](RelationshipDataObject.md) |  | 
**Included** | [**[]PlaceLocaleData**](PlaceLocaleData.md) |  | 
**Links** | [**PaginationData**](PaginationData.md) |  | 

## Methods

### NewPlaceLocalesCollection

`func NewPlaceLocalesCollection(data []RelationshipDataObject, included []PlaceLocaleData, links PaginationData, ) *PlaceLocalesCollection`

NewPlaceLocalesCollection instantiates a new PlaceLocalesCollection object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewPlaceLocalesCollectionWithDefaults

`func NewPlaceLocalesCollectionWithDefaults() *PlaceLocalesCollection`

NewPlaceLocalesCollectionWithDefaults instantiates a new PlaceLocalesCollection object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetData

`func (o *PlaceLocalesCollection) GetData() []RelationshipDataObject`

GetData returns the Data field if non-nil, zero value otherwise.

### GetDataOk

`func (o *PlaceLocalesCollection) GetDataOk() (*[]RelationshipDataObject, bool)`

GetDataOk returns a tuple with the Data field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetData

`func (o *PlaceLocalesCollection) SetData(v []RelationshipDataObject)`

SetData sets Data field to given value.


### GetIncluded

`func (o *PlaceLocalesCollection) GetIncluded() []PlaceLocaleData`

GetIncluded returns the Included field if non-nil, zero value otherwise.

### GetIncludedOk

`func (o *PlaceLocalesCollection) GetIncludedOk() (*[]PlaceLocaleData, bool)`

GetIncludedOk returns a tuple with the Included field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIncluded

`func (o *PlaceLocalesCollection) SetIncluded(v []PlaceLocaleData)`

SetIncluded sets Included field to given value.


### GetLinks

`func (o *PlaceLocalesCollection) GetLinks() PaginationData`

GetLinks returns the Links field if non-nil, zero value otherwise.

### GetLinksOk

`func (o *PlaceLocalesCollection) GetLinksOk() (*PaginationData, bool)`

GetLinksOk returns a tuple with the Links field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLinks

`func (o *PlaceLocalesCollection) SetLinks(v PaginationData)`

SetLinks sets Links field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


