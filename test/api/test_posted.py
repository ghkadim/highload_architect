import time
import pytest

from user import User
from asyncapi import AsyncApi, WebSocketException, Post

from openapi_client.model.post_create_post_request import PostCreatePostRequest


def test_unauthorized(async_api: AsyncApi):
    with pytest.raises(WebSocketException) as e:
        with async_api.post_feed_posted():
            pass
    assert e.value.status_code == 401


def test_no_posts(async_api: AsyncApi, default_user: User, make_user):
    with async_api.post_feed_posted(default_user.access_token) as posted_feed:
        default_user.api.post_create_post(post_create_post_request=PostCreatePostRequest(text="hello"))
        other_user = make_user()
        other_user.login()
        other_user.api.post_create_post(post_create_post_request=PostCreatePostRequest(text="hello"))

        for _ in posted_feed:
            assert False, "posts not expected"


@pytest.mark.parametrize(
    "user_num,posts_per_user", [
        (1, 1),
        (2, 2),
        (4, 10),
        (11, 100)
    ])
def test_friends(async_api: AsyncApi, default_user: User, make_user, user_num, posts_per_user):
    users = []
    for i in range(user_num):
        user = make_user()
        user.login()
        users.append(user)
        default_user.api.friend_set_user_id_put(user.user_id)

    with async_api.post_feed_posted(default_user.access_token) as posted_feed:
        time.sleep(1)
        posts = []
        post_idx = 0
        for i in range(posts_per_user):
            for user in users:
                text = f"post {post_idx} from {user.user_id}"
                post_id = user.api.post_create_post(
                    post_create_post_request=PostCreatePostRequest(text=text))
                posts.append(Post(postId=post_id, postText=text, author_user_id=user.user_id))
                post_idx += 1

        assert list(posted_feed) == posts


def test_add_friends(async_api: AsyncApi, default_user: User, make_user):
    user = make_user()
    user.login()
    user.api.post_create_post(
        post_create_post_request=PostCreatePostRequest(text='should not be visible'))

    posts = []
    with async_api.post_feed_posted(default_user.access_token) as posted_feed:
        time.sleep(1)
        default_user.api.friend_set_user_id_put(user.user_id)
        time.sleep(2)

        text = f"post from {user.user_id}"
        post_id = user.api.post_create_post(
            post_create_post_request=PostCreatePostRequest(text=text))
        posts.append(Post(postId=post_id, postText=text, author_user_id=user.user_id))

        assert list(posted_feed) == posts


def test_delete_friends(async_api: AsyncApi, default_user: User, make_user):
    user = make_user()
    user.login()

    posts = []
    with async_api.post_feed_posted(default_user.access_token) as posted_feed:
        time.sleep(1)
        default_user.api.friend_set_user_id_put(user.user_id)
        time.sleep(2)

        text = f"post from {user.user_id}"
        post_id = user.api.post_create_post(
            post_create_post_request=PostCreatePostRequest(text=text))
        posts.append(Post(postId=post_id, postText=text, author_user_id=user.user_id))

        default_user.api.friend_delete_user_id_put(user.user_id)
        time.sleep(2)
        user.api.post_create_post(
            post_create_post_request=PostCreatePostRequest(text='should not be visible'))

        assert list(posted_feed) == posts


def test_multiple_connections(async_api: AsyncApi, default_user: User, make_user):
    user = make_user()
    user.login()
    user.api.post_create_post(
        post_create_post_request=PostCreatePostRequest(text='should not be visible'))

    posts = []
    with async_api.post_feed_posted(default_user.access_token) as posted_feed_1, \
            async_api.post_feed_posted(default_user.access_token) as posted_feed_2:

        time.sleep(1)
        default_user.api.friend_set_user_id_put(user.user_id)
        time.sleep(2)

        text = f"post from {user.user_id}"
        post_id = user.api.post_create_post(
            post_create_post_request=PostCreatePostRequest(text=text))
        posts.append(Post(postId=post_id, postText=text, author_user_id=user.user_id))

        default_user.api.friend_delete_user_id_put(user.user_id)
        time.sleep(2)
        user.api.post_create_post(
            post_create_post_request=PostCreatePostRequest(text='should not be visible'))

        assert list(posted_feed_1) == posts
        assert list(posted_feed_2) == posts
