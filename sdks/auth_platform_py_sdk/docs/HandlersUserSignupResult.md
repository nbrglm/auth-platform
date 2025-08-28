# HandlersUserSignupResult


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**message** | **str** |  | [optional] 
**user_id** | **str** |  | [optional] 

## Example

```python
from auth_platform_py_sdk.models.handlers_user_signup_result import HandlersUserSignupResult

# TODO update the JSON string below
json = "{}"
# create an instance of HandlersUserSignupResult from a JSON string
handlers_user_signup_result_instance = HandlersUserSignupResult.from_json(json)
# print the JSON string representation of the object
print(HandlersUserSignupResult.to_json())

# convert the object into a dict
handlers_user_signup_result_dict = handlers_user_signup_result_instance.to_dict()
# create an instance of HandlersUserSignupResult from a dict
handlers_user_signup_result_from_dict = HandlersUserSignupResult.from_dict(handlers_user_signup_result_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


