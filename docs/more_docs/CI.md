# CI Integration with tfacon

1. Store credentials like auth_token in tfacon-secrets(jenkins credentials)
2. Extract launch_name/launch_id from Test Manegement Platform(Report Portal) by your CI(__This part should be done as a part of your own automation__)
3. run ```tfacon run -v --auth-token ${tfacon-secrets} --project-name "project-name" --launch-name "launch-name"```