# PlaceLocalesCollection

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Data** | [**[]PlaceLocalesCollectionDataInner**](PlaceLocalesCollectionDataInner.md) |  | 
**Included** | [**[]PlaceLocaleData**](PlaceLocaleData.md) |  | 

## Methods

### NewPlaceLocalesCollection

`func NewPlaceLocalesCollection(data []PlaceLocalesCollectionDataInner, included []PlaceLocaleData, ) *PlaceLocalesCollection`

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

`func (o *PlaceLocalesCollection) GetData() []PlaceLocalesCollectionDataInner`

GetData returns the Data field if non-nil, zero value otherwise.

### GetDataOk

`func (o *PlaceLocalesCollection) GetDataOk() (*[]PlaceLocalesCollectionDataInner, bool)`

GetDataOk returns a tuple with the Data field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetData

`func (o *PlaceLocalesCollection) SetData(v []PlaceLocalesCollectionDataInner)`

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



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


