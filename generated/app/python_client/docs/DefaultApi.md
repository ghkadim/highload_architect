# openapi_client.DefaultApi

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**dialog_user_id_list_get**](DefaultApi.md#dialog_user_id_list_get) | **GET** /dialog/{user_id}/list | 
[**dialog_user_id_send_post**](DefaultApi.md#dialog_user_id_send_post) | **POST** /dialog/{user_id}/send | 
[**friend_delete_user_id_put**](DefaultApi.md#friend_delete_user_id_put) | **PUT** /friend/delete/{user_id} | 
[**friend_set_user_id_put**](DefaultApi.md#friend_set_user_id_put) | **PUT** /friend/set/{user_id} | 
[**login_post**](DefaultApi.md#login_post) | **POST** /login | 
[**post_create_post**](DefaultApi.md#post_create_post) | **POST** /post/create | 
[**post_delete_id_put**](DefaultApi.md#post_delete_id_put) | **PUT** /post/delete/{id} | 
[**post_feed_get**](DefaultApi.md#post_feed_get) | **GET** /post/feed | 
[**post_get_id_get**](DefaultApi.md#post_get_id_get) | **GET** /post/get/{id} | 
[**post_update_put**](DefaultApi.md#post_update_put) | **PUT** /post/update | 
[**user_get_id_get**](DefaultApi.md#user_get_id_get) | **GET** /user/get/{id} | 
[**user_register_post**](DefaultApi.md#user_register_post) | **POST** /user/register | 
[**user_search_get**](DefaultApi.md#user_search_get) | **GET** /user/search | 


# **dialog_user_id_list_get**
> [DialogMessage] dialog_user_id_list_get(user_id)



### Example

* Bearer Authentication (bearerAuth):

```python
import time
import openapi_client
from openapi_client.api import default_api
from openapi_client.model.dialog_message import DialogMessage
from openapi_client.model.login_post500_response import LoginPost500Response
from pprint import pprint
# Defining the host is optional and defaults to http://localhost
# See configuration.py for a list of all supported configuration parameters.
configuration = openapi_client.Configuration(
    host = "http://localhost"
)

# The client must configure the authentication and authorization parameters
# in accordance with the API server security policy.
# Examples for each auth method are provided below, use the example that
# satisfies your auth use case.

# Configure Bearer authorization: bearerAuth
configuration = openapi_client.Configuration(
    access_token = 'YOUR_BEARER_TOKEN'
)

# Enter a context with an instance of the API client
with openapi_client.ApiClient(configuration) as api_client:
    # Create an instance of the API class
    api_instance = default_api.DefaultApi(api_client)
    user_id = "user_id_example" # str | 

    # example passing only required values which don't have defaults set
    try:
        api_response = api_instance.dialog_user_id_list_get(user_id)
        pprint(api_response)
    except openapi_client.ApiException as e:
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

# **dialog_user_id_send_post**
> dialog_user_id_send_post(user_id)



### Example

* Bearer Authentication (bearerAuth):

```python
import time
import openapi_client
from openapi_client.api import default_api
from openapi_client.model.login_post500_response import LoginPost500Response
from openapi_client.model.dialog_user_id_send_post_request import DialogUserIdSendPostRequest
from pprint import pprint
# Defining the host is optional and defaults to http://localhost
# See configuration.py for a list of all supported configuration parameters.
configuration = openapi_client.Configuration(
    host = "http://localhost"
)

# The client must configure the authentication and authorization parameters
# in accordance with the API server security policy.
# Examples for each auth method are provided below, use the example that
# satisfies your auth use case.

# Configure Bearer authorization: bearerAuth
configuration = openapi_client.Configuration(
    access_token = 'YOUR_BEARER_TOKEN'
)

# Enter a context with an instance of the API client
with openapi_client.ApiClient(configuration) as api_client:
    # Create an instance of the API class
    api_instance = default_api.DefaultApi(api_client)
    user_id = "user_id_example" # str | 
    dialog_user_id_send_post_request = DialogUserIdSendPostRequest(
        text="Привет, как дела?",
    ) # DialogUserIdSendPostRequest |  (optional)

    # example passing only required values which don't have defaults set
    try:
        api_instance.dialog_user_id_send_post(user_id)
    except openapi_client.ApiException as e:
        print("Exception when calling DefaultApi->dialog_user_id_send_post: %s\n" % e)

    # example passing only required values which don't have defaults set
    # and optional values
    try:
        api_instance.dialog_user_id_send_post(user_id, dialog_user_id_send_post_request=dialog_user_id_send_post_request)
    except openapi_client.ApiException as e:
        print("Exception when calling DefaultApi->dialog_user_id_send_post: %s\n" % e)
```


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **user_id** | **str**|  |
 **dialog_user_id_send_post_request** | [**DialogUserIdSendPostRequest**](DialogUserIdSendPostRequest.md)|  | [optional]

### Return type

void (empty response body)

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

# **friend_delete_user_id_put**
> friend_delete_user_id_put(user_id)



### Example

* Bearer Authentication (bearerAuth):

```python
import time
import openapi_client
from openapi_client.api import default_api
from openapi_client.model.login_post500_response import LoginPost500Response
from pprint import pprint
# Defining the host is optional and defaults to http://localhost
# See configuration.py for a list of all supported configuration parameters.
configuration = openapi_client.Configuration(
    host = "http://localhost"
)

# The client must configure the authentication and authorization parameters
# in accordance with the API server security policy.
# Examples for each auth method are provided below, use the example that
# satisfies your auth use case.

# Configure Bearer authorization: bearerAuth
configuration = openapi_client.Configuration(
    access_token = 'YOUR_BEARER_TOKEN'
)

# Enter a context with an instance of the API client
with openapi_client.ApiClient(configuration) as api_client:
    # Create an instance of the API class
    api_instance = default_api.DefaultApi(api_client)
    user_id = "user_id_example" # str | 

    # example passing only required values which don't have defaults set
    try:
        api_instance.friend_delete_user_id_put(user_id)
    except openapi_client.ApiException as e:
        print("Exception when calling DefaultApi->friend_delete_user_id_put: %s\n" % e)
```


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **user_id** | **str**|  |

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
**200** | Пользователь успешно удалил из друзей пользователя |  -  |
**400** | Невалидные данные ввода |  -  |
**401** | Неавторизованный доступ |  -  |
**500** | Ошибка сервера |  * Retry-After - Время, через которое еще раз нужно сделать запрос <br>  |
**503** | Ошибка сервера |  * Retry-After - Время, через которое еще раз нужно сделать запрос <br>  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **friend_set_user_id_put**
> friend_set_user_id_put(user_id)



### Example

* Bearer Authentication (bearerAuth):

```python
import time
import openapi_client
from openapi_client.api import default_api
from openapi_client.model.login_post500_response import LoginPost500Response
from pprint import pprint
# Defining the host is optional and defaults to http://localhost
# See configuration.py for a list of all supported configuration parameters.
configuration = openapi_client.Configuration(
    host = "http://localhost"
)

# The client must configure the authentication and authorization parameters
# in accordance with the API server security policy.
# Examples for each auth method are provided below, use the example that
# satisfies your auth use case.

# Configure Bearer authorization: bearerAuth
configuration = openapi_client.Configuration(
    access_token = 'YOUR_BEARER_TOKEN'
)

# Enter a context with an instance of the API client
with openapi_client.ApiClient(configuration) as api_client:
    # Create an instance of the API class
    api_instance = default_api.DefaultApi(api_client)
    user_id = "user_id_example" # str | 

    # example passing only required values which don't have defaults set
    try:
        api_instance.friend_set_user_id_put(user_id)
    except openapi_client.ApiException as e:
        print("Exception when calling DefaultApi->friend_set_user_id_put: %s\n" % e)
```


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **user_id** | **str**|  |

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
**200** | Пользователь успешно указал своего друга |  -  |
**400** | Невалидные данные ввода |  -  |
**401** | Неавторизованный доступ |  -  |
**500** | Ошибка сервера |  * Retry-After - Время, через которое еще раз нужно сделать запрос <br>  |
**503** | Ошибка сервера |  * Retry-After - Время, через которое еще раз нужно сделать запрос <br>  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **login_post**
> LoginPost200Response login_post()



Упрощенный процесс аутентификации путем передачи идентификатор пользователя и получения токена для дальнейшего прохождения авторизации

### Example


```python
import time
import openapi_client
from openapi_client.api import default_api
from openapi_client.model.login_post200_response import LoginPost200Response
from openapi_client.model.login_post500_response import LoginPost500Response
from openapi_client.model.login_post_request import LoginPostRequest
from pprint import pprint
# Defining the host is optional and defaults to http://localhost
# See configuration.py for a list of all supported configuration parameters.
configuration = openapi_client.Configuration(
    host = "http://localhost"
)


# Enter a context with an instance of the API client
with openapi_client.ApiClient() as api_client:
    # Create an instance of the API class
    api_instance = default_api.DefaultApi(api_client)
    login_post_request = LoginPostRequest(
        id="id_example",
        password="Секретная строка",
    ) # LoginPostRequest |  (optional)

    # example passing only required values which don't have defaults set
    # and optional values
    try:
        api_response = api_instance.login_post(login_post_request=login_post_request)
        pprint(api_response)
    except openapi_client.ApiException as e:
        print("Exception when calling DefaultApi->login_post: %s\n" % e)
```


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **login_post_request** | [**LoginPostRequest**](LoginPostRequest.md)|  | [optional]

### Return type

[**LoginPost200Response**](LoginPost200Response.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details

| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | Успешная аутентификация |  -  |
**400** | Невалидные данные |  -  |
**404** | Пользователь не найден |  -  |
**500** | Ошибка сервера |  * Retry-After - Время, через которое еще раз нужно сделать запрос <br>  |
**503** | Ошибка сервера |  * Retry-After - Время, через которое еще раз нужно сделать запрос <br>  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **post_create_post**
> str post_create_post()



### Example

* Bearer Authentication (bearerAuth):

```python
import time
import openapi_client
from openapi_client.api import default_api
from openapi_client.model.post_create_post_request import PostCreatePostRequest
from openapi_client.model.login_post500_response import LoginPost500Response
from pprint import pprint
# Defining the host is optional and defaults to http://localhost
# See configuration.py for a list of all supported configuration parameters.
configuration = openapi_client.Configuration(
    host = "http://localhost"
)

# The client must configure the authentication and authorization parameters
# in accordance with the API server security policy.
# Examples for each auth method are provided below, use the example that
# satisfies your auth use case.

# Configure Bearer authorization: bearerAuth
configuration = openapi_client.Configuration(
    access_token = 'YOUR_BEARER_TOKEN'
)

# Enter a context with an instance of the API client
with openapi_client.ApiClient(configuration) as api_client:
    # Create an instance of the API class
    api_instance = default_api.DefaultApi(api_client)
    post_create_post_request = PostCreatePostRequest(
        text="Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Lectus mauris ultrices eros in cursus turpis massa.",
    ) # PostCreatePostRequest |  (optional)

    # example passing only required values which don't have defaults set
    # and optional values
    try:
        api_response = api_instance.post_create_post(post_create_post_request=post_create_post_request)
        pprint(api_response)
    except openapi_client.ApiException as e:
        print("Exception when calling DefaultApi->post_create_post: %s\n" % e)
```


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **post_create_post_request** | [**PostCreatePostRequest**](PostCreatePostRequest.md)|  | [optional]

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
**200** | Успешно создан пост |  -  |
**400** | Невалидные данные ввода |  -  |
**401** | Неавторизованный доступ |  -  |
**500** | Ошибка сервера |  * Retry-After - Время, через которое еще раз нужно сделать запрос <br>  |
**503** | Ошибка сервера |  * Retry-After - Время, через которое еще раз нужно сделать запрос <br>  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **post_delete_id_put**
> post_delete_id_put(id)



### Example

* Bearer Authentication (bearerAuth):

```python
import time
import openapi_client
from openapi_client.api import default_api
from openapi_client.model.login_post500_response import LoginPost500Response
from pprint import pprint
# Defining the host is optional and defaults to http://localhost
# See configuration.py for a list of all supported configuration parameters.
configuration = openapi_client.Configuration(
    host = "http://localhost"
)

# The client must configure the authentication and authorization parameters
# in accordance with the API server security policy.
# Examples for each auth method are provided below, use the example that
# satisfies your auth use case.

# Configure Bearer authorization: bearerAuth
configuration = openapi_client.Configuration(
    access_token = 'YOUR_BEARER_TOKEN'
)

# Enter a context with an instance of the API client
with openapi_client.ApiClient(configuration) as api_client:
    # Create an instance of the API class
    api_instance = default_api.DefaultApi(api_client)
    id = "1d535fd6-7521-4cb1-aa6d-031be7123c4d" # str | 

    # example passing only required values which don't have defaults set
    try:
        api_instance.post_delete_id_put(id)
    except openapi_client.ApiException as e:
        print("Exception when calling DefaultApi->post_delete_id_put: %s\n" % e)
```


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **str**|  |

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
**200** | Успешно удален пост |  -  |
**400** | Невалидные данные ввода |  -  |
**401** | Неавторизованный доступ |  -  |
**500** | Ошибка сервера |  * Retry-After - Время, через которое еще раз нужно сделать запрос <br>  |
**503** | Ошибка сервера |  * Retry-After - Время, через которое еще раз нужно сделать запрос <br>  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **post_feed_get**
> [Post] post_feed_get()



### Example

* Bearer Authentication (bearerAuth):

```python
import time
import openapi_client
from openapi_client.api import default_api
from openapi_client.model.post import Post
from openapi_client.model.login_post500_response import LoginPost500Response
from pprint import pprint
# Defining the host is optional and defaults to http://localhost
# See configuration.py for a list of all supported configuration parameters.
configuration = openapi_client.Configuration(
    host = "http://localhost"
)

# The client must configure the authentication and authorization parameters
# in accordance with the API server security policy.
# Examples for each auth method are provided below, use the example that
# satisfies your auth use case.

# Configure Bearer authorization: bearerAuth
configuration = openapi_client.Configuration(
    access_token = 'YOUR_BEARER_TOKEN'
)

# Enter a context with an instance of the API client
with openapi_client.ApiClient(configuration) as api_client:
    # Create an instance of the API class
    api_instance = default_api.DefaultApi(api_client)
    offset = 100 # float |  (optional) if omitted the server will use the default value of 0
    limit = 10 # float |  (optional) if omitted the server will use the default value of 10

    # example passing only required values which don't have defaults set
    # and optional values
    try:
        api_response = api_instance.post_feed_get(offset=offset, limit=limit)
        pprint(api_response)
    except openapi_client.ApiException as e:
        print("Exception when calling DefaultApi->post_feed_get: %s\n" % e)
```


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **offset** | **float**|  | [optional] if omitted the server will use the default value of 0
 **limit** | **float**|  | [optional] if omitted the server will use the default value of 10

### Return type

[**[Post]**](Post.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details

| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | Успешно получены посты друзей |  -  |
**400** | Невалидные данные ввода |  -  |
**401** | Неавторизованный доступ |  -  |
**500** | Ошибка сервера |  * Retry-After - Время, через которое еще раз нужно сделать запрос <br>  |
**503** | Ошибка сервера |  * Retry-After - Время, через которое еще раз нужно сделать запрос <br>  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **post_get_id_get**
> Post post_get_id_get(id)



### Example


```python
import time
import openapi_client
from openapi_client.api import default_api
from openapi_client.model.post import Post
from openapi_client.model.login_post500_response import LoginPost500Response
from pprint import pprint
# Defining the host is optional and defaults to http://localhost
# See configuration.py for a list of all supported configuration parameters.
configuration = openapi_client.Configuration(
    host = "http://localhost"
)


# Enter a context with an instance of the API client
with openapi_client.ApiClient() as api_client:
    # Create an instance of the API class
    api_instance = default_api.DefaultApi(api_client)
    id = "1d535fd6-7521-4cb1-aa6d-031be7123c4d" # str | 

    # example passing only required values which don't have defaults set
    try:
        api_response = api_instance.post_get_id_get(id)
        pprint(api_response)
    except openapi_client.ApiException as e:
        print("Exception when calling DefaultApi->post_get_id_get: %s\n" % e)
```


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **str**|  |

### Return type

[**Post**](Post.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details

| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | Успешно получен пост |  -  |
**400** | Невалидные данные ввода |  -  |
**401** | Неавторизованный доступ |  -  |
**500** | Ошибка сервера |  * Retry-After - Время, через которое еще раз нужно сделать запрос <br>  |
**503** | Ошибка сервера |  * Retry-After - Время, через которое еще раз нужно сделать запрос <br>  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **post_update_put**
> post_update_put()



### Example

* Bearer Authentication (bearerAuth):

```python
import time
import openapi_client
from openapi_client.api import default_api
from openapi_client.model.login_post500_response import LoginPost500Response
from openapi_client.model.post_update_put_request import PostUpdatePutRequest
from pprint import pprint
# Defining the host is optional and defaults to http://localhost
# See configuration.py for a list of all supported configuration parameters.
configuration = openapi_client.Configuration(
    host = "http://localhost"
)

# The client must configure the authentication and authorization parameters
# in accordance with the API server security policy.
# Examples for each auth method are provided below, use the example that
# satisfies your auth use case.

# Configure Bearer authorization: bearerAuth
configuration = openapi_client.Configuration(
    access_token = 'YOUR_BEARER_TOKEN'
)

# Enter a context with an instance of the API client
with openapi_client.ApiClient(configuration) as api_client:
    # Create an instance of the API class
    api_instance = default_api.DefaultApi(api_client)
    post_update_put_request = PostUpdatePutRequest(
        id="1d535fd6-7521-4cb1-aa6d-031be7123c4d",
        text="Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Lectus mauris ultrices eros in cursus turpis massa.",
    ) # PostUpdatePutRequest |  (optional)

    # example passing only required values which don't have defaults set
    # and optional values
    try:
        api_instance.post_update_put(post_update_put_request=post_update_put_request)
    except openapi_client.ApiException as e:
        print("Exception when calling DefaultApi->post_update_put: %s\n" % e)
```


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **post_update_put_request** | [**PostUpdatePutRequest**](PostUpdatePutRequest.md)|  | [optional]

### Return type

void (empty response body)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details

| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | Успешно изменен пост |  -  |
**400** | Невалидные данные ввода |  -  |
**401** | Неавторизованный доступ |  -  |
**500** | Ошибка сервера |  * Retry-After - Время, через которое еще раз нужно сделать запрос <br>  |
**503** | Ошибка сервера |  * Retry-After - Время, через которое еще раз нужно сделать запрос <br>  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **user_get_id_get**
> User user_get_id_get(id)



Получение анкеты пользователя

### Example


```python
import time
import openapi_client
from openapi_client.api import default_api
from openapi_client.model.user import User
from openapi_client.model.login_post500_response import LoginPost500Response
from pprint import pprint
# Defining the host is optional and defaults to http://localhost
# See configuration.py for a list of all supported configuration parameters.
configuration = openapi_client.Configuration(
    host = "http://localhost"
)


# Enter a context with an instance of the API client
with openapi_client.ApiClient() as api_client:
    # Create an instance of the API class
    api_instance = default_api.DefaultApi(api_client)
    id = "id_example" # str | Идентификатор пользователя

    # example passing only required values which don't have defaults set
    try:
        api_response = api_instance.user_get_id_get(id)
        pprint(api_response)
    except openapi_client.ApiException as e:
        print("Exception when calling DefaultApi->user_get_id_get: %s\n" % e)
```


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **str**| Идентификатор пользователя |

### Return type

[**User**](User.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details

| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | Успешное получение анкеты пользователя |  -  |
**400** | Невалидные данные |  -  |
**404** | Анкета не найдена |  -  |
**500** | Ошибка сервера |  * Retry-After - Время, через которое еще раз нужно сделать запрос <br>  |
**503** | Ошибка сервера |  * Retry-After - Время, через которое еще раз нужно сделать запрос <br>  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **user_register_post**
> UserRegisterPost200Response user_register_post()



Регистрация нового пользователя

### Example


```python
import time
import openapi_client
from openapi_client.api import default_api
from openapi_client.model.user_register_post200_response import UserRegisterPost200Response
from openapi_client.model.login_post500_response import LoginPost500Response
from openapi_client.model.user_register_post_request import UserRegisterPostRequest
from pprint import pprint
# Defining the host is optional and defaults to http://localhost
# See configuration.py for a list of all supported configuration parameters.
configuration = openapi_client.Configuration(
    host = "http://localhost"
)


# Enter a context with an instance of the API client
with openapi_client.ApiClient() as api_client:
    # Create an instance of the API class
    api_instance = default_api.DefaultApi(api_client)
    user_register_post_request = UserRegisterPostRequest(
        first_name="Имя",
        second_name="Фамилия",
        age=18,
        birthdate=dateutil_parser('Wed Feb 01 03:00:00 MSK 2017').date(),
        biography="Хобби, интересы и т.п.",
        city="Москва",
        password="Секретная строка",
    ) # UserRegisterPostRequest |  (optional)

    # example passing only required values which don't have defaults set
    # and optional values
    try:
        api_response = api_instance.user_register_post(user_register_post_request=user_register_post_request)
        pprint(api_response)
    except openapi_client.ApiException as e:
        print("Exception when calling DefaultApi->user_register_post: %s\n" % e)
```


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **user_register_post_request** | [**UserRegisterPostRequest**](UserRegisterPostRequest.md)|  | [optional]

### Return type

[**UserRegisterPost200Response**](UserRegisterPost200Response.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details

| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | Успешная регистрация |  -  |
**400** | Невалидные данные |  -  |
**500** | Ошибка сервера |  * Retry-After - Время, через которое еще раз нужно сделать запрос <br>  |
**503** | Ошибка сервера |  * Retry-After - Время, через которое еще раз нужно сделать запрос <br>  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **user_search_get**
> [User] user_search_get(first_name, last_name)



Поиск анкет

### Example


```python
import time
import openapi_client
from openapi_client.api import default_api
from openapi_client.model.user import User
from openapi_client.model.login_post500_response import LoginPost500Response
from pprint import pprint
# Defining the host is optional and defaults to http://localhost
# See configuration.py for a list of all supported configuration parameters.
configuration = openapi_client.Configuration(
    host = "http://localhost"
)


# Enter a context with an instance of the API client
with openapi_client.ApiClient() as api_client:
    # Create an instance of the API class
    api_instance = default_api.DefaultApi(api_client)
    first_name = "Конст" # str | Условие поиска по имени
    last_name = "Оси" # str | Условие поиска по фамилии

    # example passing only required values which don't have defaults set
    try:
        api_response = api_instance.user_search_get(first_name, last_name)
        pprint(api_response)
    except openapi_client.ApiException as e:
        print("Exception when calling DefaultApi->user_search_get: %s\n" % e)
```


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **first_name** | **str**| Условие поиска по имени |
 **last_name** | **str**| Условие поиска по фамилии |

### Return type

[**[User]**](User.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details

| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | Успешные поиск пользователя |  -  |
**400** | Невалидные данные |  -  |
**500** | Ошибка сервера |  * Retry-After - Время, через которое еще раз нужно сделать запрос <br>  |
**503** | Ошибка сервера |  * Retry-After - Время, через которое еще раз нужно сделать запрос <br>  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

