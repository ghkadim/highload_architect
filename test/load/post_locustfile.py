import random

from locust import run_single_user, constant_pacing, task
from openapi_client.model.post_create_post_request import PostCreatePostRequest
from openapi_client.model.post_update_put_request import PostUpdatePutRequest
from user import OpenApiUser


class PostUser(OpenApiUser):
    def __init__(self, env):
        super(PostUser, self).__init__(env)
        self.login()
        self.num_friends = env.parsed_options.num_friends
        self.posts = list()

    wait_time = constant_pacing(30)

    @task(2)
    def post_create(self):
        post_id = self.openapi_client.post_create_post(
            post_create_post_request=PostCreatePostRequest(text="hello"))
        self.posts.append(post_id)

    @task
    def post_update(self):
        if len(self.posts) == 0:
            return

        post_id = random.choice(self.posts)
        self.openapi_client.post_update_put(
            post_update_put_request=PostUpdatePutRequest(id=post_id, text="hello world"))

    @task
    def post_delete(self):
        if len(self.posts) == 0:
            return
        post_id = random.choice(self.posts)
        with self.client.rename_request("/post/delete/[post]"):
            self.openapi_client.post_delete_id_put(id=post_id)
        self.posts.remove(post_id)


if __name__ == "__main__":
    run_single_user(PostUser)
