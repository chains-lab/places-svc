# PaginationData

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**PageNumber** | **int64** | The current page number. | 
**PageSize** | **int64** | The number of items per page. | 
**TotalItems** | **int64** | The total number of items available. | 

## Methods

### NewPaginationData

`func NewPaginationData(pageNumber int64, pageSize int64, totalItems int64, ) *PaginationData`

NewPaginationData instantiates a new PaginationData object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewPaginationDataWithDefaults

`func NewPaginationDataWithDefaults() *PaginationData`

NewPaginationDataWithDefaults instantiates a new PaginationData object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetPageNumber

`func (o *PaginationData) GetPageNumber() int64`

GetPageNumber returns the PageNumber field if non-nil, zero value otherwise.

### GetPageNumberOk

`func (o *PaginationData) GetPageNumberOk() (*int64, bool)`

GetPageNumberOk returns a tuple with the PageNumber field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPageNumber

`func (o *PaginationData) SetPageNumber(v int64)`

SetPageNumber sets PageNumber field to given value.


### GetPageSize

`func (o *PaginationData) GetPageSize() int64`

GetPageSize returns the PageSize field if non-nil, zero value otherwise.

### GetPageSizeOk

`func (o *PaginationData) GetPageSizeOk() (*int64, bool)`

GetPageSizeOk returns a tuple with the PageSize field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPageSize

`func (o *PaginationData) SetPageSize(v int64)`

SetPageSize sets PageSize field to given value.


### GetTotalItems

`func (o *PaginationData) GetTotalItems() int64`

GetTotalItems returns the TotalItems field if non-nil, zero value otherwise.

### GetTotalItemsOk

`func (o *PaginationData) GetTotalItemsOk() (*int64, bool)`

GetTotalItemsOk returns a tuple with the TotalItems field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTotalItems

`func (o *PaginationData) SetTotalItems(v int64)`

SetTotalItems sets TotalItems field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


