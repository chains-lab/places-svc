# CreatePlaceDataAttributes

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**CityId** | [**uuid.UUID**](uuid.UUID.md) | city id | 
**DistributorId** | Pointer to [**uuid.UUID**](uuid.UUID.md) | distributor id | [optional] 
**Class** | **string** | place class | 
**Point** | [**Point**](Point.md) |  | 
**Locale** | **string** | locale | 
**Name** | **string** | place name | 
**Description** | **string** | place description | 
**Website** | Pointer to **string** | place website | [optional] 
**Phone** | Pointer to **string** | place phone number | [optional] 

## Methods

### NewCreatePlaceDataAttributes

`func NewCreatePlaceDataAttributes(cityId uuid.UUID, class string, point Point, locale string, name string, description string, ) *CreatePlaceDataAttributes`

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

`func (o *CreatePlaceDataAttributes) GetCityId() uuid.UUID`

GetCityId returns the CityId field if non-nil, zero value otherwise.

### GetCityIdOk

`func (o *CreatePlaceDataAttributes) GetCityIdOk() (*uuid.UUID, bool)`

GetCityIdOk returns a tuple with the CityId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCityId

`func (o *CreatePlaceDataAttributes) SetCityId(v uuid.UUID)`

SetCityId sets CityId field to given value.


### GetDistributorId

`func (o *CreatePlaceDataAttributes) GetDistributorId() uuid.UUID`

GetDistributorId returns the DistributorId field if non-nil, zero value otherwise.

### GetDistributorIdOk

`func (o *CreatePlaceDataAttributes) GetDistributorIdOk() (*uuid.UUID, bool)`

GetDistributorIdOk returns a tuple with the DistributorId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDistributorId

`func (o *CreatePlaceDataAttributes) SetDistributorId(v uuid.UUID)`

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


### GetPoint

`func (o *CreatePlaceDataAttributes) GetPoint() Point`

GetPoint returns the Point field if non-nil, zero value otherwise.

### GetPointOk

`func (o *CreatePlaceDataAttributes) GetPointOk() (*Point, bool)`

GetPointOk returns a tuple with the Point field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPoint

`func (o *CreatePlaceDataAttributes) SetPoint(v Point)`

SetPoint sets Point field to given value.


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


