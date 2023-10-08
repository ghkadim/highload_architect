# openapi_client_dialog.DefaultApi

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**dialog_user_id_list_get**](DefaultApi.md#dialog_user_id_list_get) | **GET** /dialog/{user_id}/list | 
[**dialog_user_id_message_message_id_read_put**](DefaultApi.md#dialog_user_id_message_message_id_read_put) | **PUT** /dialog/{user_id}/message/{message_id}/read | 
[**dialog_user_id_send_post**](DefaultApi.md#dialog_user_id_send_post) | **POST** /dialog/{user_id}/send | 


# **dialog_user_id_list_get**
> [DialogMessage] dialog_user_id_list_get(user_id)



### Example

* Bearer Authentication (bearerAuth):

```python
import time
import openapi_client_dialog
from openapi_client_dialog.api import default_api
from openapi_client_dialog.model.dialog_user_id_send_post500_response import DialogUserIdSendPost500Response
from openapi_client_dialog.model.dialog_message import DialogMessage
from pprint import pprint
# Defining the host is optional and defaults to http://localhost
# See configuration.py for a list of all supported configuration parameters.
configuration = openapi_client_dialog.Configuration(
    host = "http://localhost"
)

# The client must configure the authentication and authorization parameters
# in accordance with the API server security policy.
# Examples for each auth method are provided below, use the example that
# satisfies your auth use case.

# Configure Bearer authorization: bearerAuth
configuration = openapi_client_dialog.Configuration(
    access_token = 'YOUR_BEARER_TOKEN'
)

# Enter a context with an instance of the API client
with openapi_client_dialog.ApiClient(configuration) as api_client:
    # Create an instance of the API class
    api_instance = default_api.DefaultApi(api_client)
    user_id = "user_id_example" # str | 

    # example passing only required values which don't have defaults set
    try:
        api_response = api_instance.dialog_user_id_list_get(user_id)
        pprint(api_response)
    except openapi_client_dialog.ApiException as e:
        print("Exception when calling DefaultApi->dialog_user_id_list_get: %s\n" % e)
```


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **user_id** | **str**|  |

### Return type

[**[DialogMessage]**](DialogMessage.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details

| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | Диалог между двумя пользователями |  -  |
**400** | Невалидные данные ввода |  -  |
**401** | Неавторизованный доступ |  -  |
**500** | Ошибка сервера |  * Retry-After - Время, через которое еще раз нужно сделать запрос <br>  |
**503** | Ошибка сервера |  * Retry-After - Время, через которое еще раз нужно сделать запрос <br>  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **dialog_user_id_message_message_id_read_put**
> dialog_user_id_message_message_id_read_put(user_id, message_id)



### Example

* Bearer Authentication (bearerAuth):

```python
import time
import openapi_client_dialog
from openapi_client_dialog.api import default_api
from openapi_client_dialog.model.dialog_user_id_send_post500_response import DialogUserIdSendPost500Response
from pprint import pprint
# Defining the host is optional and defaults to http://localhost
# See configuration.py for a list of all supported configuration parameters.
configuration = openapi_client_dialog.Configuration(
    host = "http://localhost"
)

# The client must configure the authentication and authorization parameters
# in accordance with the API server security policy.
# Examples for each auth method are provided below, use the example that
# satisfies your auth use case.

# Configure Bearer authorization: bearerAuth
configuration = openapi_client_dialog.Configuration(
    access_token = 'YOUR_BEARER_TOKEN'
)

# Enter a context with an instance of the API client
with openapi_client_dialog.ApiClient(configuration) as api_client:
    # Create an instance of the API class
    api_instance = default_api.DefaultApi(api_client)
    user_id = "user_id_example" # str | 
    message_id = "1d535fd6-7521-4cb1-aa6d-031be7123c4d" # str | 

    # example passing only required values which don't have defaults set
    try:
        api_instance.dialog_user_id_message_message_id_read_put(user_id, message_id)
    except openapi_client_dialog.ApiException as e:
        print("Exception when calling DefaultApi->dialog_user_id_message_message_id_read_put: %s\n" % e)
```


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **user_id** | **str**|  |
 **message_id** | **str**|  |

### Return type

void (empty response body)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details

| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | Успешно изменено сообщение |  -  |
**400** | Невалидные данные ввода |  -  |
**401** | Неавторизованный доступ |  -  |
**500** | Ошибка сервера |  * Retry-After - Время, через которое еще раз нужно сделать запрос <br>  |
**503** | Ошибка сервера |  * Retry-After - Время, через которое еще раз нужно сделать запрос <br>  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **dialog_user_id_send_post**
> str dialog_user_id_send_post(user_id)



### Example

* Bearer Authentication (bearerAuth):

```python
import time
import openapi_client_dialog
from openapi_client_dialog.api import default_api
from openapi_client_dialog.model.dialog_user_id_send_post500_response import DialogUserIdSendPost500Response
from openapi_client_dialog.model.dialog_user_id_send_post_request import DialogUserIdSendPostRequest
from pprint import pprint
# Defining the host is optional and defaults to http://localhost
# See configuration.py for a list of all supported configuration parameters.
configuration = openapi_client_dialog.Configuration(
    host = "http://localhost"
)

# The client must configure the authentication and authorization parameters
# in accordance with the API server security policy.
# Examples for each auth method are provided below, use the example that
# satisfies your auth use case.

# Configure Bearer authorization: bearerAuth
configuration = openapi_client_dialog.Configuration(
    access_token = 'YOUR_BEARER_TOKEN'
)

# Enter a context with an instance of the API client
with openapi_client_dialog.ApiClient(configuration) as api_client:
    # Create an instance of the API class
    api_instance = default_api.DefaultApi(api_client)
    user_id = "user_id_example" # str | 
    dialog_user_id_send_post_request = DialogUserIdSendPostRequest(
        text="Привет, как дела?",
    ) # DialogUserIdSendPostRequest |  (optional)

    # example passing only required values which don't have defaults set
    try:
        api_response = api_instance.dialog_user_id_send_post(user_id)
        pprint(api_response)
    except openapi_client_dialog.ApiException as e:
        print("Exception when calling DefaultApi->dialog_user_id_send_post: %s\n" % e)

    # example passing only required values which don't have defaults set
    # and optional values
    try:
        api_response = api_instance.dialog_user_id_send_post(user_id, dialog_user_id_send_post_request=dialog_user_id_send_post_request)
        pprint(api_response)
    except openapi_client_dialog.ApiException as e:
        print("Exception when calling DefaultApi->dialog_user_id_send_post: %s\n" % e)
```


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **user_id** | **str**|  |
 **dialog_user_id_send_post_request** | [**DialogUserIdSendPostRequest**](DialogUserIdSendPostRequest.md)|  | [optional]

### Return type

**str**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details

| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | Успешно отправлено сообщение |  -  |
**400** | Невалидные данные ввода |  -  |
**401** | Неавторизованный доступ |  -  |
**500** | Ошибка сервера |  * Retry-After - Время, через которое еще раз нужно сделать запрос <br>  |
**503** | Ошибка сервера |  * Retry-After - Время, через которое еще раз нужно сделать запрос <br>  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

