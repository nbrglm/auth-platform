# HandlersLogoutResult


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**message** | **str** |  | [optional] 
**success** | **bool** |  | [optional] 

## Example

```python
from auth_platform_py_sdk.models.handlers_logout_result import HandlersLogoutResult

# TODO update the JSON string below
json = "{}"
# create an instance of HandlersLogoutResult from a JSON string
handlers_logout_result_instance = HandlersLogoutResult.from_json(json)
# print the JSON string representation of the object
print(HandlersLogoutResult.to_json())

# convert the object into a dict
handlers_logout_result_dict = handlers_logout_result_instance.to_dict()
# create an instance of HandlersLogoutResult from a dict
handlers_logout_result_from_dict = HandlersLogoutResult.from_dict(handlers_logout_result_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


