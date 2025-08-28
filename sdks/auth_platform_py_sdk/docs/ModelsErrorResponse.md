# ModelsErrorResponse


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**code** | **int** | Code is an error code that can be used for programmatic handling of errors. | [optional] 
**debug** | **str** | DebugMessage is a technical message that can be used for debugging. | [optional] 
**message** | **str** | Message is a user-friendly message that can be displayed to the end user. | [optional] 

## Example

```python
from auth_platform_py_sdk.models.models_error_response import ModelsErrorResponse

# TODO update the JSON string below
json = "{}"
# create an instance of ModelsErrorResponse from a JSON string
models_error_response_instance = ModelsErrorResponse.from_json(json)
# print the JSON string representation of the object
print(ModelsErrorResponse.to_json())

# convert the object into a dict
models_error_response_dict = models_error_response_instance.to_dict()
# create an instance of ModelsErrorResponse from a dict
models_error_response_from_dict = ModelsErrorResponse.from_dict(models_error_response_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


