# HandlersUserSignupData


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**confirm_password** | **str** |  | 
**email** | **str** |  | 
**first_name** | **str** |  | 
**invite_token** | **str** | Optional invite token for signup | [optional] 
**last_name** | **str** |  | 
**password** | **str** |  | 

## Example

```python
from auth_platform_py_sdk.models.handlers_user_signup_data import HandlersUserSignupData

# TODO update the JSON string below
json = "{}"
# create an instance of HandlersUserSignupData from a JSON string
handlers_user_signup_data_instance = HandlersUserSignupData.from_json(json)
# print the JSON string representation of the object
print(HandlersUserSignupData.to_json())

# convert the object into a dict
handlers_user_signup_data_dict = handlers_user_signup_data_instance.to_dict()
# create an instance of HandlersUserSignupData from a dict
handlers_user_signup_data_from_dict = HandlersUserSignupData.from_dict(handlers_user_signup_data_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


