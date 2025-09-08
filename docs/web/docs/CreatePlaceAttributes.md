# CreatePlaceAttributes

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**CityId** | Pointer to **string** | city id | [optional] 
**DistributorId** | Pointer to **string** | distributor id | [optional] 
**Class** | Pointer to **string** | place class | [optional] 
**Ownership** | Pointer to **string** | place ownership | [optional] 
**Point** | Pointer to [**CreatePlaceAttributesPoint**](CreatePlaceAttributesPoint.md) |  | [optional] 
**Locale** | Pointer to **string** | locale | [optional] 
**Name** | Pointer to **string** | place name | [optional] 
**Address** | Pointer to **string** | place address | [optional] 
**Description** | Pointer to **string** | place description | [optional] 
**Website** | Pointer to **string** | place website | [optional] 
**Phone** | Pointer to **string** | place phone number | [optional] 

## Methods

### NewCreatePlaceAttributes

`func NewCreatePlaceAttributes() *CreatePlaceAttributes`

NewCreatePlaceAttributes instantiates a new CreatePlaceAttributes object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCreatePlaceAttributesWithDefaults

`func NewCreatePlaceAttributesWithDefaults() *CreatePlaceAttributes`

NewCreatePlaceAttributesWithDefaults instantiates a new CreatePlaceAttributes object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetCityId

`func (o *CreatePlaceAttributes) GetCityId() string`

GetCityId returns the CityId field if non-nil, zero value otherwise.

### GetCityIdOk

`func (o *CreatePlaceAttributes) GetCityIdOk() (*string, bool)`

GetCityIdOk returns a tuple with the CityId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCityId

`func (o *CreatePlaceAttributes) SetCityId(v string)`

SetCityId sets CityId field to given value.

### HasCityId

`func (o *CreatePlaceAttributes) HasCityId() bool`

HasCityId returns a boolean if a field has been set.

### GetDistributorId

`func (o *CreatePlaceAttributes) GetDistributorId() string`

GetDistributorId returns the DistributorId field if non-nil, zero value otherwise.

### GetDistributorIdOk

`func (o *CreatePlaceAttributes) GetDistributorIdOk() (*string, bool)`

GetDistributorIdOk returns a tuple with the DistributorId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDistributorId

`func (o *CreatePlaceAttributes) SetDistributorId(v string)`

SetDistributorId sets DistributorId field to given value.

### HasDistributorId

`func (o *CreatePlaceAttributes) HasDistributorId() bool`

HasDistributorId returns a boolean if a field has been set.

### GetClass

`func (o *CreatePlaceAttributes) GetClass() string`

GetClass returns the Class field if non-nil, zero value otherwise.

### GetClassOk

`func (o *CreatePlaceAttributes) GetClassOk() (*string, bool)`

GetClassOk returns a tuple with the Class field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetClass

`func (o *CreatePlaceAttributes) SetClass(v string)`

SetClass sets Class field to given value.

### HasClass

`func (o *CreatePlaceAttributes) HasClass() bool`

HasClass returns a boolean if a field has been set.

### GetOwnership

`func (o *CreatePlaceAttributes) GetOwnership() string`

GetOwnership returns the Ownership field if non-nil, zero value otherwise.

### GetOwnershipOk

`func (o *CreatePlaceAttributes) GetOwnershipOk() (*string, bool)`

GetOwnershipOk returns a tuple with the Ownership field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOwnership

`func (o *CreatePlaceAttributes) SetOwnership(v string)`

SetOwnership sets Ownership field to given value.

### HasOwnership

`func (o *CreatePlaceAttributes) HasOwnership() bool`

HasOwnership returns a boolean if a field has been set.

### GetPoint

`func (o *CreatePlaceAttributes) GetPoint() CreatePlaceAttributesPoint`

GetPoint returns the Point field if non-nil, zero value otherwise.

### GetPointOk

`func (o *CreatePlaceAttributes) GetPointOk() (*CreatePlaceAttributesPoint, bool)`

GetPointOk returns a tuple with the Point field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPoint

`func (o *CreatePlaceAttributes) SetPoint(v CreatePlaceAttributesPoint)`

SetPoint sets Point field to given value.

### HasPoint

`func (o *CreatePlaceAttributes) HasPoint() bool`

HasPoint returns a boolean if a field has been set.

### GetLocale

`func (o *CreatePlaceAttributes) GetLocale() string`

GetLocale returns the Locale field if non-nil, zero value otherwise.

### GetLocaleOk

`func (o *CreatePlaceAttributes) GetLocaleOk() (*string, bool)`

GetLocaleOk returns a tuple with the Locale field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLocale

`func (o *CreatePlaceAttributes) SetLocale(v string)`

SetLocale sets Locale field to given value.

### HasLocale

`func (o *CreatePlaceAttributes) HasLocale() bool`

HasLocale returns a boolean if a field has been set.

### GetName

`func (o *CreatePlaceAttributes) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *CreatePlaceAttributes) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *CreatePlaceAttributes) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *CreatePlaceAttributes) HasName() bool`

HasName returns a boolean if a field has been set.

### GetAddress

`func (o *CreatePlaceAttributes) GetAddress() string`

GetAddress returns the Address field if non-nil, zero value otherwise.

### GetAddressOk

`func (o *CreatePlaceAttributes) GetAddressOk() (*string, bool)`

GetAddressOk returns a tuple with the Address field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAddress

`func (o *CreatePlaceAttributes) SetAddress(v string)`

SetAddress sets Address field to given value.

### HasAddress

`func (o *CreatePlaceAttributes) HasAddress() bool`

HasAddress returns a boolean if a field has been set.

### GetDescription

`func (o *CreatePlaceAttributes) GetDescription() string`

GetDescription returns the Description field if non-nil, zero value otherwise.

### GetDescriptionOk

`func (o *CreatePlaceAttributes) GetDescriptionOk() (*string, bool)`

GetDescriptionOk returns a tuple with the Description field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDescription

`func (o *CreatePlaceAttributes) SetDescription(v string)`

SetDescription sets Description field to given value.

### HasDescription

`func (o *CreatePlaceAttributes) HasDescription() bool`

HasDescription returns a boolean if a field has been set.

### GetWebsite

`func (o *CreatePlaceAttributes) GetWebsite() string`

GetWebsite returns the Website field if non-nil, zero value otherwise.

### GetWebsiteOk

`func (o *CreatePlaceAttributes) GetWebsiteOk() (*string, bool)`

GetWebsiteOk returns a tuple with the Website field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWebsite

`func (o *CreatePlaceAttributes) SetWebsite(v string)`

SetWebsite sets Website field to given value.

### HasWebsite

`func (o *CreatePlaceAttributes) HasWebsite() bool`

HasWebsite returns a boolean if a field has been set.

### GetPhone

`func (o *CreatePlaceAttributes) GetPhone() string`

GetPhone returns the Phone field if non-nil, zero value otherwise.

### GetPhoneOk

`func (o *CreatePlaceAttributes) GetPhoneOk() (*string, bool)`

GetPhoneOk returns a tuple with the Phone field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPhone

`func (o *CreatePlaceAttributes) SetPhone(v string)`

SetPhone sets Phone field to given value.

### HasPhone

`func (o *CreatePlaceAttributes) HasPhone() bool`

HasPhone returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


