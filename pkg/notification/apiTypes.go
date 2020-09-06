package notification

type User struct {
	Login             string `json:"login"`
	Id                int    `json:"id"`
	NodeId            string `json:"node_id"`
	AvatarUrl         string `json:"avatar_url"`
	GravatarId        string `json:"gravatar_id"`
	Url               string `json:"url"`
	HtmlUrl           string `json:"html_url"`
	FollowersUrl      string `json:"followers_url"`
	FollowingUrl      string `json:"following_url"`
	GistsUrl          string `json:"gists_url"`
	StarredUrl        string `json:"starred_url"`
	SubscriptionsUrl  string `json:"subscriptions_url"`
	OrganizationsUrl  string `json:"organizations_url"`
	ReposUrl          string `json:"repos_url"`
	EventsUrl         string `json:"events_url"`
	ReceivedEventsUrl string `json:"received_events_url"`
	Type              string `json:"type"`
	SiteAdmin         bool   `json:"site_admin"`
}

type Subject struct {
	Title            string `json:"title"`
	Url              string `json:"url"`
	LatestCommentUrl string `json:"latest_comment_url"`
	Type             string `json:"type"`
}

type Repository struct {
	Id               int    `json:"id"`
	NodeId           string `json:"node_id"`
	Name             string `json:"name"`
	FullName         string `json:"full_name"`
	Owner            User   `json:"owner"`
	Private          bool   `json:"private"`
	HtmlUrl          string `json:"html_url"`
	Description      string `json:"description"`
	Fork             bool   `json:"fork"`
	Url              string `json:"url"`
	ArchiveUrl       string `json:"archive_url"`
	AssigneesUrl     string `json:"assignees_url"`
	BlobsUrl         string `json:"blobs_url"`
	BranchesUrl      string `json:"branches_url"`
	CollaboratorsUrl string `json:"collaborators_url"`
	CommentsUrl      string `json:"comments_url"`
	CommitsUrl       string `json:"commits_url"`
	CompareUrl       string `json:"compare_url"`
	ContentsUrl      string `json:"contents_url"`
	ContributorsUrl  string `json:"contributors_url"`
	DeploymentsUrl   string `json:"deployments_url"`
	DownloadsUrl     string `json:"downloads_url"`
	EventsUrl        string `json:"events_url"`
	ForksUrl         string `json:"forks_url"`
	GitCommitsUrl    string `json:"git_commits_url"`
	GitRefsUrl       string `json:"git_refs_url"`
	GitTagsUrl       string `json:"git_tags_url"`
	GitUrl           string `json:"git_url"`
	IssueCommentUrl  string `json:"issue_comment_url"`
	IssueEventsUrl   string `json:"issue_events_url"`
	KeysUrl          string `json:"keys_url"`
	LabelsUrl        string `json:"labels_url"`
	LanguagesUrl     string `json:"languages_url"`
	MergesUrl        string `json:"merges_url"`
	MilestonesUrl    string `json:"milestones_url"`
	NotificationsUrl string `json:"notifications_url"`
	PullsUrl         string `json:"pulls_url"`
	ReleasesUrl      string `json:"releases_url"`
	SshUrl           string `json:"ssh_url"`
	StargazersUrl    string `json:"stargazers_url"`
	StatusesUrl      string `json:"statuses_url"`
	SubscribersUrl   string `json:"subscribers_url"`
	SubscriptionUrl  string `json:"subscription_url"`
	TagsUrl          string `json:"tags_url"`
	TeamsUrl         string `json:"teams_url"`
	TreesUrl         string `json:"trees_url"`
}

type Notification struct {
	Id         string     `json:"id"`
	Repository Repository `json:"repository"`
	Subject    Subject    `json:"subject"`
	Reason     string     `json:"reason"`
	Unread     bool       `json:"unread"`
	UpdatedAt  string     `json:"updated_at"`
	LastReadAt string     `json:"last_read_at"`
	Url        string     `json:"url"`
}

type Label struct {
	Id          int    `json:"id"`
	NodeId      string `json:"node_id"`
	Url         string `json:"url"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Color       string `json:"color"`
	Default     bool   `json:"default"`
}

type Milestone struct {
	Url          string `json:"url"`
	HtmlUrl      string `json:"html_url"`
	LabelsUrl    string `json:"labels_url"`
	Id           int    `json:"id"`
	NodeId       string `json:"node_id"`
	Number       int    `json:"number"`
	State        string `json:"state"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	Creator      User   `json:"creator"`
	OpenIssues   int    `json:"open_issues"`
	ClosedIssues int    `json:"closed_issues"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
	ClosedAt     string `json:"closed_at"`
	DueOn        string `json:"due_on"`
}

type PullRequest struct {
	Url      string `json:"url"`
	HtmlUrl  string `json:"html_url"`
	PatchUrl string `json:"patch_url"`
	DiffUrl  string `json:"diff_url"`
}

type Issue struct {
	Id               int         `json:"id"`
	NodeId           string      `json:"node_id"`
	Url              string      `json:"url"`
	RepositoryUrl    string      `json:"repository_url"`
	LabelsUrl        string      `json:"labels_url"`
	CommentsUrl      string      `json:"comments_url"`
	EventsUrl        string      `json:"events_url"`
	HtmlUrl          string      `json:"html_url"`
	Number           int         `json:"number"`
	State            string      `json:"state"`
	Title            string      `json:"title"`
	Body             string      `json:"body"`
	User             User        `json:"user"`
	Labels           []Label     `json:"labels"`
	Assignee         User        `json:"assignee"`
	Assignees        []User      `json:"assignees"`
	Milestone        Milestone   `json:"milestone"`
	Locked           bool        `json:"locked"`
	ActiveLockReason string      `json:"active_lock_reason"`
	Comments         int         `json:"comments"`
	PullRequest      PullRequest `json:"pull_request"`
	ClosedAt         string      `json:"closed_at"`
	CreatedAt        string      `json:"created_at"`
	UpdatedAt        string      `json:"updated_at"`
	ClosedBy         User        `json:"closed_by"`
}

type NotificationResponse []Notification
