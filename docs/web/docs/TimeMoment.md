# TimeMoment

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Weekday** | **string** | Day of the week. | 
**Time** | **string** | Time of the day in 24-hour format (HH:MM). | 

## Methods

### NewTimeMoment

`func NewTimeMoment(weekday string, time string, ) *TimeMoment`

NewTimeMoment instantiates a new TimeMoment object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewTimeMomentWithDefaults

`func NewTimeMomentWithDefaults() *TimeMoment`

NewTimeMomentWithDefaults instantiates a new TimeMoment object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetWeekday

`func (o *TimeMoment) GetWeekday() string`

GetWeekday returns the Weekday field if non-nil, zero value otherwise.

### GetWeekdayOk

`func (o *TimeMoment) GetWeekdayOk() (*string, bool)`

GetWeekdayOk returns a tuple with the Weekday field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWeekday

`func (o *TimeMoment) SetWeekday(v string)`

SetWeekday sets Weekday field to given value.


### GetTime

`func (o *TimeMoment) GetTime() string`

GetTime returns the Time field if non-nil, zero value otherwise.

### GetTimeOk

`func (o *TimeMoment) GetTimeOk() (*string, bool)`

GetTimeOk returns a tuple with the Time field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTime

`func (o *TimeMoment) SetTime(v string)`

SetTime sets Time field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


