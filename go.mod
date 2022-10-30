module github.com/eleztian/gitlab-bot

go 1.18

require (
	github.com/davecgh/go-spew v1.1.1
	github.com/eleztian/gitlab-bot/common v0.0.0-00010101000000-000000000000
	github.com/fsnotify/fsnotify v1.6.0
	github.com/sirupsen/logrus v1.9.0
	github.com/xanzy/go-gitlab v0.74.0
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-retryablehttp v0.7.1 // indirect
	golang.org/x/net v0.0.0-20220805013720-a33c5aa5df48 // indirect
	golang.org/x/oauth2 v0.0.0-20220722155238-128564f6959c // indirect
	golang.org/x/sys v0.0.0-20220908164124-27713097b956 // indirect
	golang.org/x/time v0.0.0-20220722155302-e5dcc9cfc0b9 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
)

replace github.com/eleztian/gitlab-bot/common => ./common
