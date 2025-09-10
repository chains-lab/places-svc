# TimetableInterval

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**From** | [**TimeMoment**](TimeMoment.md) |  | 
**To** | [**TimeMoment**](TimeMoment.md) |  | 

## Methods

### NewTimetableInterval

`func NewTimetableInterval(from TimeMoment, to TimeMoment, ) *TimetableInterval`

NewTimetableInterval instantiates a new TimetableInterval object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewTimetableIntervalWithDefaults

`func NewTimetableIntervalWithDefaults() *TimetableInterval`

NewTimetableIntervalWithDefaults instantiates a new TimetableInterval object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetFrom

`func (o *TimetableInterval) GetFrom() TimeMoment`

GetFrom returns the From field if non-nil, zero value otherwise.

### GetFromOk

`func (o *TimetableInterval) GetFromOk() (*TimeMoment, bool)`

GetFromOk returns a tuple with the From field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFrom

`func (o *TimetableInterval) SetFrom(v TimeMoment)`

SetFrom sets From field to given value.


### GetTo

`func (o *TimetableInterval) GetTo() TimeMoment`

GetTo returns the To field if non-nil, zero value otherwise.

### GetToOk

`func (o *TimetableInterval) GetToOk() (*TimeMoment, bool)`

GetToOk returns a tuple with the To field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTo

`func (o *TimetableInterval) SetTo(v TimeMoment)`

SetTo sets To field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


