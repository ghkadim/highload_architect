from locust import run_single_user, constant_pacing, task
from openapi_client.model.dialog_user_id_send_post_request import DialogUserIdSendPostRequest
from user import OpenApiUser


class DialogUser(OpenApiUser):
    def __init__(self, env):
        super(DialogUser, self).__init__(env)
        self.login()
        self.num_friends = env.parsed_options.num_friends

    wait_time = constant_pacing(5)

    @task
    def dialog_send(self):
        friend = self.peek_friend(self.num_friends)
        if friend is None:
            return

        with self.client.rename_request("/dialog/[friend]/send"):
            self.openapi_client.dialog_user_id_send_post(
                friend,
                dialog_user_id_send_post_request=DialogUserIdSendPostRequest(text="hello"))

    @task
    def dialog_list(self):
        friend = self.peek_friend(self.num_friends)
        if friend is None:
            return

        with self.client.rename_request("/dialog/[friend]/list"):
            self.openapi_client.dialog_user_id_list_get(friend)


if __name__ == "__main__":
    run_single_user(DialogUser)
