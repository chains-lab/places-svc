# CreatePlaceDataAttributes

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**CityId** | Pointer to **string** | city id | [optional] 
**DistributorId** | Pointer to **string** | distributor id | [optional] 
**Class** | Pointer to **string** | place class | [optional] 
**Ownership** | Pointer to **string** | place ownership | [optional] 
**Point** | Pointer to [**CreatePlaceDataAttributesPoint**](CreatePlaceDataAttributesPoint.md) |  | [optional] 
**Locale** | Pointer to **string** | locale | [optional] 
**Name** | Pointer to **string** | place name | [optional] 
**Address** | Pointer to **string** | place address | [optional] 
**Description** | Pointer to **string** | place description | [optional] 
**Website** | Pointer to **string** | place website | [optional] 
**Phone** | Pointer to **string** | place phone number | [optional] 

## Methods

### NewCreatePlaceDataAttributes

`func NewCreatePlaceDataAttributes() *CreatePlaceDataAttributes`

NewCreatePlaceDataAttributes instantiates a new CreatePlaceDataAttributes object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCreatePlaceDataAttributesWithDefaults

`func NewCreatePlaceDataAttributesWithDefaults() *CreatePlaceDataAttributes`

NewCreatePlaceDataAttributesWithDefaults instantiates a new CreatePlaceDataAttributes object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetCityId

`func (o *CreatePlaceDataAttributes) GetCityId() string`

GetCityId returns the CityId field if non-nil, zero value otherwise.

### GetCityIdOk

`func (o *CreatePlaceDataAttributes) GetCityIdOk() (*string, bool)`

GetCityIdOk returns a tuple with the CityId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCityId

`func (o *CreatePlaceDataAttributes) SetCityId(v string)`

SetCityId sets CityId field to given value.

### HasCityId

`func (o *CreatePlaceDataAttributes) HasCityId() bool`

HasCityId returns a boolean if a field has been set.

### GetDistributorId

`func (o *CreatePlaceDataAttributes) GetDistributorId() string`

GetDistributorId returns the DistributorId field if non-nil, zero value otherwise.

### GetDistributorIdOk

`func (o *CreatePlaceDataAttributes) GetDistributorIdOk() (*string, bool)`

GetDistributorIdOk returns a tuple with the DistributorId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDistributorId

`func (o *CreatePlaceDataAttributes) SetDistributorId(v string)`

SetDistributorId sets DistributorId field to given value.

### HasDistributorId

`func (o *CreatePlaceDataAttributes) HasDistributorId() bool`

HasDistributorId returns a boolean if a field has been set.

### GetClass

`func (o *CreatePlaceDataAttributes) GetClass() string`

GetClass returns the Class field if non-nil, zero value otherwise.

### GetClassOk

`func (o *CreatePlaceDataAttributes) GetClassOk() (*string, bool)`

GetClassOk returns a tuple with the Class field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetClass

`func (o *CreatePlaceDataAttributes) SetClass(v string)`

SetClass sets Class field to given value.

### HasClass

`func (o *CreatePlaceDataAttributes) HasClass() bool`

HasClass returns a boolean if a field has been set.

### GetOwnership

`func (o *CreatePlaceDataAttributes) GetOwnership() string`

GetOwnership returns the Ownership field if non-nil, zero value otherwise.

### GetOwnershipOk

`func (o *CreatePlaceDataAttributes) GetOwnershipOk() (*string, bool)`

GetOwnershipOk returns a tuple with the Ownership field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOwnership

`func (o *CreatePlaceDataAttributes) SetOwnership(v string)`

SetOwnership sets Ownership field to given value.

### HasOwnership

`func (o *CreatePlaceDataAttributes) HasOwnership() bool`

HasOwnership returns a boolean if a field has been set.

### GetPoint

`func (o *CreatePlaceDataAttributes) GetPoint() CreatePlaceDataAttributesPoint`

GetPoint returns the Point field if non-nil, zero value otherwise.

### GetPointOk

`func (o *CreatePlaceDataAttributes) GetPointOk() (*CreatePlaceDataAttributesPoint, bool)`

GetPointOk returns a tuple with the Point field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPoint

`func (o *CreatePlaceDataAttributes) SetPoint(v CreatePlaceDataAttributesPoint)`

SetPoint sets Point field to given value.

### HasPoint

`func (o *CreatePlaceDataAttributes) HasPoint() bool`

HasPoint returns a boolean if a field has been set.

### GetLocale

`func (o *CreatePlaceDataAttributes) GetLocale() string`

GetLocale returns the Locale field if non-nil, zero value otherwise.

### GetLocaleOk

`func (o *CreatePlaceDataAttributes) GetLocaleOk() (*string, bool)`

GetLocaleOk returns a tuple with the Locale field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLocale

`func (o *CreatePlaceDataAttributes) SetLocale(v string)`

SetLocale sets Locale field to given value.

### HasLocale

`func (o *CreatePlaceDataAttributes) HasLocale() bool`

HasLocale returns a boolean if a field has been set.

### GetName

`func (o *CreatePlaceDataAttributes) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *CreatePlaceDataAttributes) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *CreatePlaceDataAttributes) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *CreatePlaceDataAttributes) HasName() bool`

HasName returns a boolean if a field has been set.

### GetAddress

`func (o *CreatePlaceDataAttributes) GetAddress() string`

GetAddress returns the Address field if non-nil, zero value otherwise.

### GetAddressOk

`func (o *CreatePlaceDataAttributes) GetAddressOk() (*string, bool)`

GetAddressOk returns a tuple with the Address field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAddress

`func (o *CreatePlaceDataAttributes) SetAddress(v string)`

SetAddress sets Address field to given value.

### HasAddress

`func (o *CreatePlaceDataAttributes) HasAddress() bool`

HasAddress returns a boolean if a field has been set.

### GetDescription

`func (o *CreatePlaceDataAttributes) GetDescription() string`

GetDescription returns the Description field if non-nil, zero value otherwise.

### GetDescriptionOk

`func (o *CreatePlaceDataAttributes) GetDescriptionOk() (*string, bool)`

GetDescriptionOk returns a tuple with the Description field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDescription

`func (o *CreatePlaceDataAttributes) SetDescription(v string)`

SetDescription sets Description field to given value.

### HasDescription

`func (o *CreatePlaceDataAttributes) HasDescription() bool`

HasDescription returns a boolean if a field has been set.

### GetWebsite

`func (o *CreatePlaceDataAttributes) GetWebsite() string`

GetWebsite returns the Website field if non-nil, zero value otherwise.

### GetWebsiteOk

`func (o *CreatePlaceDataAttributes) GetWebsiteOk() (*string, bool)`

GetWebsiteOk returns a tuple with the Website field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWebsite

`func (o *CreatePlaceDataAttributes) SetWebsite(v string)`

SetWebsite sets Website field to given value.

### HasWebsite

`func (o *CreatePlaceDataAttributes) HasWebsite() bool`

HasWebsite returns a boolean if a field has been set.

### GetPhone

`func (o *CreatePlaceDataAttributes) GetPhone() string`

GetPhone returns the Phone field if non-nil, zero value otherwise.

### GetPhoneOk

`func (o *CreatePlaceDataAttributes) GetPhoneOk() (*string, bool)`

GetPhoneOk returns a tuple with the Phone field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPhone

`func (o *CreatePlaceDataAttributes) SetPhone(v string)`

SetPhone sets Phone field to given value.

### HasPhone

`func (o *CreatePlaceDataAttributes) HasPhone() bool`

HasPhone returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


