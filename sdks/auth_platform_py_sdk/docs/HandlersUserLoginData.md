# HandlersUserLoginData


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**email** | **str** |  | 
**flow_return_to** | **str** | Optional field to store in the flow data which can be fetched by the client after login This can be used to redirect the user to a specific page after login or to maintain the state of the application. It is recommended to validate this field on the client side to prevent open redirect vulnerabilities. | [optional] 
**password** | **str** |  | 

## Example

```python
from auth_platform_py_sdk.models.handlers_user_login_data import HandlersUserLoginData

# TODO update the JSON string below
json = "{}"
# create an instance of HandlersUserLoginData from a JSON string
handlers_user_login_data_instance = HandlersUserLoginData.from_json(json)
# print the JSON string representation of the object
print(HandlersUserLoginData.to_json())

# convert the object into a dict
handlers_user_login_data_dict = handlers_user_login_data_instance.to_dict()
# create an instance of HandlersUserLoginData from a dict
handlers_user_login_data_from_dict = HandlersUserLoginData.from_dict(handlers_user_login_data_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


