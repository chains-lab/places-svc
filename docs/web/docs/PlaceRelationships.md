# PlaceRelationships

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Class** | [**ClassRelationshipsParent**](ClassRelationshipsParent.md) |  | 
**City** | [**ClassRelationshipsParent**](ClassRelationshipsParent.md) |  | 
**Distributor** | Pointer to [**ClassRelationshipsParent**](ClassRelationshipsParent.md) |  | [optional] 

## Methods

### NewPlaceRelationships

`func NewPlaceRelationships(class ClassRelationshipsParent, city ClassRelationshipsParent, ) *PlaceRelationships`

NewPlaceRelationships instantiates a new PlaceRelationships object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewPlaceRelationshipsWithDefaults

`func NewPlaceRelationshipsWithDefaults() *PlaceRelationships`

NewPlaceRelationshipsWithDefaults instantiates a new PlaceRelationships object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetClass

`func (o *PlaceRelationships) GetClass() ClassRelationshipsParent`

GetClass returns the Class field if non-nil, zero value otherwise.

### GetClassOk

`func (o *PlaceRelationships) GetClassOk() (*ClassRelationshipsParent, bool)`

GetClassOk returns a tuple with the Class field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetClass

`func (o *PlaceRelationships) SetClass(v ClassRelationshipsParent)`

SetClass sets Class field to given value.


### GetCity

`func (o *PlaceRelationships) GetCity() ClassRelationshipsParent`

GetCity returns the City field if non-nil, zero value otherwise.

### GetCityOk

`func (o *PlaceRelationships) GetCityOk() (*ClassRelationshipsParent, bool)`

GetCityOk returns a tuple with the City field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCity

`func (o *PlaceRelationships) SetCity(v ClassRelationshipsParent)`

SetCity sets City field to given value.


### GetDistributor

`func (o *PlaceRelationships) GetDistributor() ClassRelationshipsParent`

GetDistributor returns the Distributor field if non-nil, zero value otherwise.

### GetDistributorOk

`func (o *PlaceRelationships) GetDistributorOk() (*ClassRelationshipsParent, bool)`

GetDistributorOk returns a tuple with the Distributor field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDistributor

`func (o *PlaceRelationships) SetDistributor(v ClassRelationshipsParent)`

SetDistributor sets Distributor field to given value.

### HasDistributor

`func (o *PlaceRelationships) HasDistributor() bool`

HasDistributor returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


