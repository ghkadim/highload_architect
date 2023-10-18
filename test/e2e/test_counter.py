import pytest
import openapi_client_counter
from openapi_client_counter.api import default_api
from openapi_client_counter.model.counter_counter_id_get200_response import CounterCounterIdGet200Response
from user import User


def test_unauthorized():
    api = default_api.DefaultApi()
    with pytest.raises(openapi_client_counter.exceptions.UnauthorizedException):
        api.counter_counter_id_get("unread_messages")


def test_zero_counter(default_user: User):
    res = default_user.counter_api.counter_counter_id_get("unread_messages")
    assert res.value == 0
