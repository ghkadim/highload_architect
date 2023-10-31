import itertools
import os

import pytest
import openapi_client
import openapi_client_dialog
from openapi_client.api import default_api
from openapi_client.model.dialog_user_id_send_post_request import DialogUserIdSendPostRequest
from openapi_client_dialog.model.dialog_user_id_send_post_request import DialogUserIdSendPostRequest
from openapi_client.model.dialog_message import DialogMessage
from openapi_client_dialog.model.dialog_message import DialogMessage

from user import User


@pytest.mark.parametrize(
    "client_lib", [openapi_client, openapi_client_dialog])
def test_unauthorized(client_lib):
    api = client_lib.api.default_api.DefaultApi()
    with pytest.raises(client_lib.exceptions.UnauthorizedException):
        api.dialog_user_id_list_get("12345")

    with pytest.raises(client_lib.exceptions.UnauthorizedException):
        api.dialog_user_id_send_post(
            "12345",
            dialog_user_id_send_post_request=
                client_lib.model.dialog_user_id_send_post_request.DialogUserIdSendPostRequest(text="hello"))


@pytest.mark.skip("disabled since no check for user existence")
def test_unknown_user(default_user):
    with pytest.raises(Exception):
        default_user.api.dialog_user_id_send_post(
            "12345",
            dialog_user_id_send_post_request=DialogUserIdSendPostRequest(text="hello"))

    messages = default_user.api.dialog_user_id_list_get("12345")
    assert len(messages) == 0


class LegacyDialogApi:
    DialogMessage = openapi_client.model.dialog_message.DialogMessage
    DialogUserIdSendPostRequest = \
        openapi_client.model.dialog_user_id_send_post_request.DialogUserIdSendPostRequest

    def send(self, from_user, to_user, text):
        from_user.api.dialog_user_id_send_post(
            to_user.user_id,
            dialog_user_id_send_post_request=self.DialogUserIdSendPostRequest(text=text))

    def list(self, from_user, to_user):
        return from_user.api.dialog_user_id_list_get(to_user.user_id)


class NewDialogApi:
    DialogMessage = openapi_client_dialog.model.dialog_message.DialogMessage
    DialogUserIdSendPostRequest = \
        openapi_client_dialog.model.dialog_user_id_send_post_request.DialogUserIdSendPostRequest

    def send(self, from_user, to_user, text):
        from_user.dialog_api.dialog_user_id_send_post(
            to_user.user_id,
            dialog_user_id_send_post_request=self.DialogUserIdSendPostRequest(text=text))

    def list(self, from_user, to_user):
        return from_user.dialog_api.dialog_user_id_list_get(to_user.user_id)


def cmp_messages(expected, got):
    for expect_msg, got_msg in zip(expected, got):
        assert expect_msg._from == got_msg._from
        assert expect_msg.to == got_msg.to
        assert expect_msg.text == got_msg.text

@pytest.mark.parametrize(
    "dialog_api", [LegacyDialogApi(), NewDialogApi()])
def test_dialog(default_user, make_user, dialog_api):
    friend = make_user()
    friend.login()

    dialog_api.send(default_user, friend, "hello")
    messages = dialog_api.list(default_user, friend)
    cmp_messages(messages, [dialog_api.DialogMessage(_from=default_user.user_id, to=friend.user_id, text="hello")])

    dialog_api.send(friend, default_user, "hello")
    messages = dialog_api.list(friend, default_user)
    cmp_messages(messages, [
        dialog_api.DialogMessage(_from=friend.user_id, to=default_user.user_id, text="hello"),
        dialog_api.DialogMessage(_from=default_user.user_id, to=friend.user_id, text="hello")
    ])

    messages_on_friend = dialog_api.list(default_user, friend)
    cmp_messages(messages, messages_on_friend)


def test_backward_compatibility(default_user, make_user):
    legacy_api = LegacyDialogApi()
    new_api = NewDialogApi()

    friend = make_user()
    friend.login()

    legacy_api.send(default_user, friend, "hello legacy")
    new_api.send(friend, default_user, "hello new")

    messages = legacy_api.list(friend, default_user)
    cmp_messages(messages, [
        legacy_api.DialogMessage(_from=friend.user_id, to=default_user.user_id, text="hello new"),
        legacy_api.DialogMessage(_from=default_user.user_id, to=friend.user_id, text="hello legacy")
    ])

    messages = new_api.list(friend, default_user)
    cmp_messages(messages, [
        new_api.DialogMessage(_from=friend.user_id, to=default_user.user_id, text="hello new"),
        new_api.DialogMessage(_from=default_user.user_id, to=friend.user_id, text="hello legacy")
    ])


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
            cmp_messages(messages, expected_messages)


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
