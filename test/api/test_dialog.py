import itertools
import os

import pytest
import openapi_client
from openapi_client.api import default_api
from openapi_client.model.dialog_user_id_send_post_request import DialogUserIdSendPostRequest
from openapi_client.model.dialog_message import DialogMessage
from user import User


def test_unauthorized():
    api = default_api.DefaultApi()
    with pytest.raises(openapi_client.exceptions.UnauthorizedException):
        api.dialog_user_id_list_get("12345")

    with pytest.raises(openapi_client.exceptions.UnauthorizedException):
        api.dialog_user_id_send_post(
            "12345",
            dialog_user_id_send_post_request=DialogUserIdSendPostRequest(text="hello"))


@pytest.mark.skip("disabled since no check for user existence")
def test_unknown_user(default_user):
    with pytest.raises(Exception):
        default_user.api.dialog_user_id_send_post(
            "12345",
            dialog_user_id_send_post_request=DialogUserIdSendPostRequest(text="hello"))

    messages = default_user.api.dialog_user_id_list_get("12345")
    assert len(messages) == 0


def test_dialog(default_user, make_user):
    friend = make_user()

    default_user.api.dialog_user_id_send_post(
        friend.user_id,
        dialog_user_id_send_post_request=DialogUserIdSendPostRequest(text="hello"))

    messages = default_user.api.dialog_user_id_list_get(friend.user_id)
    assert messages == [DialogMessage(_from=default_user.user_id, to=friend.user_id, text="hello")]

    friend.login()

    friend.api.dialog_user_id_send_post(
        default_user.user_id,
        dialog_user_id_send_post_request=DialogUserIdSendPostRequest(text="hello"))

    messages = default_user.api.dialog_user_id_list_get(friend.user_id)
    assert messages == [
        DialogMessage(_from=friend.user_id, to=default_user.user_id, text="hello"),
        DialogMessage(_from=default_user.user_id, to=friend.user_id, text="hello")
    ]


def make_users(count):
    users = []
    for i in range(count):
        first_name = f'test_dialogs.first_name.{i}'
        second_name = f'test_dialogs.second_name.{i}'
        found_users = User.search(first_name, second_name)
        if len(found_users) != 0:
            u = found_users[0]
        else:
            u = User(first_name, second_name)
        u.login()
        users.append(u)
    return users


def make_dialogs(users, message_count, with_send=True):
    dialogs = {}
    for user_pair in itertools.permutations(users, 2):
        u1, u2 = sorted(user_pair, key=lambda u: u.user_id)
        key = (u1.user_id, u2.user_id)
        if key in dialogs:
            continue
        dialogs[key] = []

        for i in range(message_count):
            for u_from, u_to in ((u1, u2), (u2, u1)):
                text = f'message {i} from {u_from.user_id} to {u_to.user_id}'
                dialogs[key] = [DialogMessage(_from=u_from.user_id, to=u_to.user_id, text=text)] + dialogs[key]
                if with_send:
                    u_from.api.dialog_user_id_send_post(
                        u_to.user_id,
                        dialog_user_id_send_post_request=DialogUserIdSendPostRequest(text=text))
    return dialogs


def check_dialogs(users, dialogs):
    for user_pair in itertools.permutations(users, 2):
        u1, u2 = sorted(user_pair, key=lambda u: u.user_id)
        expected_messages = dialogs[(u1.user_id, u2.user_id)]
        for u_from, u_to in ((u1, u2), (u2, u1)):
            messages = u_from.api.dialog_user_id_list_get(u_to.user_id)
            assert messages == expected_messages


@pytest.mark.skipif("TEST_BEFORE_RESHARDING" not in os.environ, reason="check before resharding")
def test_multiple_dialogs_create():
    users_count = 10
    message_count = 1

    users = make_users(users_count)
    dialogs = make_dialogs(users, message_count, with_send=True)
    check_dialogs(users, dialogs)


@pytest.mark.skipif("TEST_AFTER_RESHARDING" not in os.environ, reason="check after resharding")
def test_multiple_dialogs_read():
    users_count = 10
    message_count = 1

    users = make_users(users_count)
    dialogs = make_dialogs(users, message_count, with_send=False)
    check_dialogs(users, dialogs)
