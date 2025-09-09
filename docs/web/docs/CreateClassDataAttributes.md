# CreateClassDataAttributes

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Icon** | **string** | class icon | 
**Name** | **string** | class name at english | [default to "New Class"]
**Parent** | Pointer to **string** | parent class code | [optional] 

## Methods

### NewCreateClassDataAttributes

`func NewCreateClassDataAttributes(icon string, name string, ) *CreateClassDataAttributes`

NewCreateClassDataAttributes instantiates a new CreateClassDataAttributes object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCreateClassDataAttributesWithDefaults

`func NewCreateClassDataAttributesWithDefaults() *CreateClassDataAttributes`

NewCreateClassDataAttributesWithDefaults instantiates a new CreateClassDataAttributes object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetIcon

`func (o *CreateClassDataAttributes) GetIcon() string`

GetIcon returns the Icon field if non-nil, zero value otherwise.

### GetIconOk

`func (o *CreateClassDataAttributes) GetIconOk() (*string, bool)`

GetIconOk returns a tuple with the Icon field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIcon

`func (o *CreateClassDataAttributes) SetIcon(v string)`

SetIcon sets Icon field to given value.


### GetName

`func (o *CreateClassDataAttributes) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *CreateClassDataAttributes) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *CreateClassDataAttributes) SetName(v string)`

SetName sets Name field to given value.


### GetParent

`func (o *CreateClassDataAttributes) GetParent() string`

GetParent returns the Parent field if non-nil, zero value otherwise.

### GetParentOk

`func (o *CreateClassDataAttributes) GetParentOk() (*string, bool)`

GetParentOk returns a tuple with the Parent field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetParent

`func (o *CreateClassDataAttributes) SetParent(v string)`

SetParent sets Parent field to given value.

### HasParent

`func (o *CreateClassDataAttributes) HasParent() bool`

HasParent returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


