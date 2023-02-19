import pytest

import openapi_client
from openapi_client.api import default_api


def test_unauthorized():
    api = default_api.DefaultApi()
    with pytest.raises(openapi_client.exceptions.UnauthorizedException):
        api.friend_set_user_id_put("0")

    with pytest.raises(openapi_client.exceptions.UnauthorizedException):
        api.friend_delete_user_id_put("0")


def test_add_delete_friend(login, make_user):
    friend = make_user()

    api = default_api.DefaultApi()
    api.friend_set_user_id_put(friend.user_id)
    api.friend_delete_user_id_put(friend.user_id)


def test_add_delete_unknown_friend(login, make_user):
    user_not_exists = "359d95e6-b099-11ed-82fd-0242ac150002"
    friend = make_user()

    api = default_api.DefaultApi()
    with pytest.raises(openapi_client.exceptions.ServiceException):
        api.friend_set_user_id_put(user_not_exists)

    api.friend_delete_user_id_put(user_not_exists)
    api.friend_delete_user_id_put(friend.user_id)
