# Point

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Lon** | **float64** | Longitude | 
**Lat** | **float64** | Latitude | 

## Methods

### NewPoint

`func NewPoint(lon float64, lat float64, ) *Point`

NewPoint instantiates a new Point object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewPointWithDefaults

`func NewPointWithDefaults() *Point`

NewPointWithDefaults instantiates a new Point object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetLon

`func (o *Point) GetLon() float64`

GetLon returns the Lon field if non-nil, zero value otherwise.

### GetLonOk

`func (o *Point) GetLonOk() (*float64, bool)`

GetLonOk returns a tuple with the Lon field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLon

`func (o *Point) SetLon(v float64)`

SetLon sets Lon field to given value.


### GetLat

`func (o *Point) GetLat() float64`

GetLat returns the Lat field if non-nil, zero value otherwise.

### GetLatOk

`func (o *Point) GetLatOk() (*float64, bool)`

GetLatOk returns a tuple with the Lat field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLat

`func (o *Point) SetLat(v float64)`

SetLat sets Lat field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


