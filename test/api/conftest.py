import pytest
import openapi_client
from openapi_client.api import default_api
from openapi_client.model.user_register_post_request import UserRegisterPostRequest
from openapi_client.model.login_post_request import LoginPostRequest


@pytest.fixture(autouse=True, scope="session")
def configure_app_host():
    conf = openapi_client.Configuration.get_default_copy()
    conf.host = "http://localhost:8080"
    openapi_client.Configuration.set_default(conf)


@pytest.fixture()
def register(configure_app_host):
    api = default_api.DefaultApi()
    resp = api.user_register_post(
        user_register_post_request=UserRegisterPostRequest(
            first_name="first_name",
            second_name="second_name",
            password="password",
        ))

    return resp


@pytest.fixture()
def login(register):
    api = default_api.DefaultApi()
    resp = api.login_post(
        login_post_request=LoginPostRequest(
            id=register.user_id,
            password="password",
        ))

    conf = openapi_client.Configuration.get_default_copy()
    conf.access_token = resp.token
    openapi_client.Configuration.set_default(conf)