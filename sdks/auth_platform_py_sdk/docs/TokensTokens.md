# TokensTokens


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**refresh_token** | **str** | RefreshToken is the generated refresh token.  This is base64.RawURLEncoding encoded. DO NOT DECODE IT WHILE RETRIEVING IT FROM THE COOKIES/CLIENT. | [optional] 
**refresh_token_expiry** | **str** | RefreshTokenExpiry is the expiration time of the refresh token. | [optional] 
**session_id** | **str** | SessionId is the unique identifier for the session. | [optional] 
**session_token** | **str** | SessionToken is the generated session token.  This is a jwt which is base64.RawURLEncoding encoded. YOU NEED TO DECODE IT WHILE RETRIEVING IT FROM THE COOKIES/CLIENT. DO NOT USE IT AS IS. VALIDATION WILL FAIL WITHOUT DECODING. ONLY WHEN DECODED, YOU SHOULD PASS IT TO THE THINGS THAT REQUIRE THE SESSION TOKEN. | [optional] 
**session_token_expiry** | **str** | SessionTokenExpiry is the expiration time of the session token. | [optional] 

## Example

```python
from auth_platform_py_sdk.models.tokens_tokens import TokensTokens

# TODO update the JSON string below
json = "{}"
# create an instance of TokensTokens from a JSON string
tokens_tokens_instance = TokensTokens.from_json(json)
# print the JSON string representation of the object
print(TokensTokens.to_json())

# convert the object into a dict
tokens_tokens_dict = tokens_tokens_instance.to_dict()
# create an instance of TokensTokens from a dict
tokens_tokens_from_dict = TokensTokens.from_dict(tokens_tokens_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


