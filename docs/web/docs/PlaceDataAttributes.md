# PlaceDataAttributes

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**CityId** | **string** | city id | 
**DistributorId** | Pointer to **string** | distributor id | [optional] 
**Class** | **string** | place class | 
**Status** | **string** | place status | 
**Verified** | **bool** | is place verified | 
**Ownership** | **string** | place ownership | 
**Point** | [**Point**](Point.md) |  | 
**Name** | **string** | place name | 
**Address** | **string** | place address | 
**Description** | **string** | place description | 
**Website** | Pointer to **string** | place website | [optional] 
**Phone** | Pointer to **string** | place phone number | [optional] 
**CreatedAt** | **time.Time** | place creation date | 
**UpdatedAt** | **time.Time** | place last update date | 

## Methods

### NewPlaceDataAttributes

`func NewPlaceDataAttributes(cityId string, class string, status string, verified bool, ownership string, point Point, name string, address string, description string, createdAt time.Time, updatedAt time.Time, ) *PlaceDataAttributes`

NewPlaceDataAttributes instantiates a new PlaceDataAttributes object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewPlaceDataAttributesWithDefaults

`func NewPlaceDataAttributesWithDefaults() *PlaceDataAttributes`

NewPlaceDataAttributesWithDefaults instantiates a new PlaceDataAttributes object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetCityId

`func (o *PlaceDataAttributes) GetCityId() string`

GetCityId returns the CityId field if non-nil, zero value otherwise.

### GetCityIdOk

`func (o *PlaceDataAttributes) GetCityIdOk() (*string, bool)`

GetCityIdOk returns a tuple with the CityId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCityId

`func (o *PlaceDataAttributes) SetCityId(v string)`

SetCityId sets CityId field to given value.


### GetDistributorId

`func (o *PlaceDataAttributes) GetDistributorId() string`

GetDistributorId returns the DistributorId field if non-nil, zero value otherwise.

### GetDistributorIdOk

`func (o *PlaceDataAttributes) GetDistributorIdOk() (*string, bool)`

GetDistributorIdOk returns a tuple with the DistributorId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDistributorId

`func (o *PlaceDataAttributes) SetDistributorId(v string)`

SetDistributorId sets DistributorId field to given value.

### HasDistributorId

`func (o *PlaceDataAttributes) HasDistributorId() bool`

HasDistributorId returns a boolean if a field has been set.

### GetClass

`func (o *PlaceDataAttributes) GetClass() string`

GetClass returns the Class field if non-nil, zero value otherwise.

### GetClassOk

`func (o *PlaceDataAttributes) GetClassOk() (*string, bool)`

GetClassOk returns a tuple with the Class field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetClass

`func (o *PlaceDataAttributes) SetClass(v string)`

SetClass sets Class field to given value.


### GetStatus

`func (o *PlaceDataAttributes) GetStatus() string`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *PlaceDataAttributes) GetStatusOk() (*string, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *PlaceDataAttributes) SetStatus(v string)`

SetStatus sets Status field to given value.


### GetVerified

`func (o *PlaceDataAttributes) GetVerified() bool`

GetVerified returns the Verified field if non-nil, zero value otherwise.

### GetVerifiedOk

`func (o *PlaceDataAttributes) GetVerifiedOk() (*bool, bool)`

GetVerifiedOk returns a tuple with the Verified field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVerified

`func (o *PlaceDataAttributes) SetVerified(v bool)`

SetVerified sets Verified field to given value.


### GetOwnership

`func (o *PlaceDataAttributes) GetOwnership() string`

GetOwnership returns the Ownership field if non-nil, zero value otherwise.

### GetOwnershipOk

`func (o *PlaceDataAttributes) GetOwnershipOk() (*string, bool)`

GetOwnershipOk returns a tuple with the Ownership field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOwnership

`func (o *PlaceDataAttributes) SetOwnership(v string)`

SetOwnership sets Ownership field to given value.


### GetPoint

`func (o *PlaceDataAttributes) GetPoint() Point`

GetPoint returns the Point field if non-nil, zero value otherwise.

### GetPointOk

`func (o *PlaceDataAttributes) GetPointOk() (*Point, bool)`

GetPointOk returns a tuple with the Point field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPoint

`func (o *PlaceDataAttributes) SetPoint(v Point)`

SetPoint sets Point field to given value.


### GetName

`func (o *PlaceDataAttributes) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *PlaceDataAttributes) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *PlaceDataAttributes) SetName(v string)`

SetName sets Name field to given value.


### GetAddress

`func (o *PlaceDataAttributes) GetAddress() string`

GetAddress returns the Address field if non-nil, zero value otherwise.

### GetAddressOk

`func (o *PlaceDataAttributes) GetAddressOk() (*string, bool)`

GetAddressOk returns a tuple with the Address field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAddress

`func (o *PlaceDataAttributes) SetAddress(v string)`

SetAddress sets Address field to given value.


### GetDescription

`func (o *PlaceDataAttributes) GetDescription() string`

GetDescription returns the Description field if non-nil, zero value otherwise.

### GetDescriptionOk

`func (o *PlaceDataAttributes) GetDescriptionOk() (*string, bool)`

GetDescriptionOk returns a tuple with the Description field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDescription

`func (o *PlaceDataAttributes) SetDescription(v string)`

SetDescription sets Description field to given value.


### GetWebsite

`func (o *PlaceDataAttributes) GetWebsite() string`

GetWebsite returns the Website field if non-nil, zero value otherwise.

### GetWebsiteOk

`func (o *PlaceDataAttributes) GetWebsiteOk() (*string, bool)`

GetWebsiteOk returns a tuple with the Website field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWebsite

`func (o *PlaceDataAttributes) SetWebsite(v string)`

SetWebsite sets Website field to given value.

### HasWebsite

`func (o *PlaceDataAttributes) HasWebsite() bool`

HasWebsite returns a boolean if a field has been set.

### GetPhone

`func (o *PlaceDataAttributes) GetPhone() string`

GetPhone returns the Phone field if non-nil, zero value otherwise.

### GetPhoneOk

`func (o *PlaceDataAttributes) GetPhoneOk() (*string, bool)`

GetPhoneOk returns a tuple with the Phone field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPhone

`func (o *PlaceDataAttributes) SetPhone(v string)`

SetPhone sets Phone field to given value.

### HasPhone

`func (o *PlaceDataAttributes) HasPhone() bool`

HasPhone returns a boolean if a field has been set.

### GetCreatedAt

`func (o *PlaceDataAttributes) GetCreatedAt() time.Time`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *PlaceDataAttributes) GetCreatedAtOk() (*time.Time, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *PlaceDataAttributes) SetCreatedAt(v time.Time)`

SetCreatedAt sets CreatedAt field to given value.


### GetUpdatedAt

`func (o *PlaceDataAttributes) GetUpdatedAt() time.Time`

GetUpdatedAt returns the UpdatedAt field if non-nil, zero value otherwise.

### GetUpdatedAtOk

`func (o *PlaceDataAttributes) GetUpdatedAtOk() (*time.Time, bool)`

GetUpdatedAtOk returns a tuple with the UpdatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpdatedAt

`func (o *PlaceDataAttributes) SetUpdatedAt(v time.Time)`

SetUpdatedAt sets UpdatedAt field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


