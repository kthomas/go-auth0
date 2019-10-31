package auth0

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// func syncAuth0Users() ([]map[string]interface{}, error) {
// 	accessToken := getAuth0AccessToken()
// 	params := map[string]interface{}{
// 		"format": "json",
// 	}
// 	payload, _ := json.Marshal(params)
// 	url := fmt.Sprintf("%s/api/%s/jobs/users-exports", common.Auth0Domain, common.auth0APINamespace)
// 	req, _ := http.NewRequest("POST", url, bytes.NewReader(payload))
// 	req.Header.Add("authorization", fmt.Sprintf("bearer %s", *accessToken))
// 	req.Header.Add("content-type", "application/json")
// 	res, err := http.DefaultClient.Do(req)
// 	if err != nil {
// 		log.Warningf("failed to vend Auth0 API token; %s", err.Error())
// 		return nil, err
// 	}
// 	defer res.Body.Close()
// 	body, _ := ioutil.ReadAll(res.Body)
// 	log.Debugf("%s", string(body))
// 	resp := make([]map[string]interface{}, 0)
// 	err = json.Unmarshal(body, &resp)
// 	if err != nil {
// 		log.Warningf("failed to unmarshal Auth0 list users API response; %s", err.Error())
// 		return nil, err
// 	}
// 	return resp, nil
// }

// GetUser returns an auth0 user by id
func GetUser(auth0UserID string) (interface{}, error) {
	client, err := NewAuth0APIClient()
	if err != nil {
		log.Warningf("failed to fetch auth0 user; %s", err.Error())
		return nil, err
	}

	status, resp, err := client.Get(fmt.Sprintf("users/%s", auth0UserID), nil)
	if err != nil {
		log.Warningf("failed to fetch auth0 user; %s", err.Error())
		return nil, err
	}
	if status != 200 {
		log.Warningf("failed to fetch auth0 user; %s", err.Error())
		return nil, err
	}

	return resp, nil
}

// ExportUsers returns an export of all auth0 users
func ExportUsers() ([]interface{}, error) {
	client, err := NewAuth0APIClient()
	if err != nil {
		log.Warningf("failed to export auth0 users; %s", err.Error())
		return nil, err
	}

	status, resp, err := client.Post("jobs/users-exports", map[string]interface{}{
		"fields": []map[string]string{
			map[string]string{"name": "user_id"},
			map[string]string{"name": "name"},
			map[string]string{"name": "email"},
			map[string]string{"name": "app_metadata"},
		},
		"format": "json",
	})
	if err != nil {
		log.Warningf("failed to export auth0 users; %s", err.Error())
		return nil, err
	}
	if status != 200 {
		log.Warningf("failed to export auth0 users; status: %d; response: %s", status, resp)
		return nil, err
	}

	users := make([]interface{}, 0)
	auth0JobID := resp.(map[string]interface{})["id"].(string)
	for {
		job, err := GetJob(auth0JobID)
		if err != nil {
			log.Warningf("failed to fetch auth0 export users job; %s", err.Error())
		}
		if job != nil {
			if status, statusOk := job.(map[string]interface{})["status"].(string); statusOk {
				if status == "completed" {
					usersExportURL := job.(map[string]interface{})["location"].(string)
					status, resp, err := client.sendRequest("GET", usersExportURL, "", nil)
					if err != nil {
						log.Warningf("failed to fetch compressed auth0 export users artifact from location: %s; %s", usersExportURL, err.Error())
					} else if status == 200 {
						if decompressedUsersExport, decompressedUsersExportOk := resp.([]byte); decompressedUsersExportOk {
							log.Debugf("exported users:\n\n%s\n\n", string(decompressedUsersExport))
							lines := strings.Split(string(decompressedUsersExport), "\n")
							for i := range lines {
								if len(lines[i]) == 0 {
									continue
								}
								var usr map[string]interface{}
								err := json.Unmarshal([]byte(lines[i]), &usr)
								if err != nil {
									log.Warningf("failed to unmarshal auth0 user from exported users artifact on line %d; %s", i, err.Error())
								} else {
									users = append(users, usr)
								}
							}
						}
						break
					}
				}
			}
		}

		time.Sleep(time.Millisecond * 2500)
	}

	return users, nil
}

// GetJob returns an auth0 job by id
func GetJob(auth0JobID string) (interface{}, error) {
	client, err := NewAuth0APIClient()
	if err != nil {
		log.Warningf("failed to fetch auth0 user; %s", err.Error())
		return nil, err
	}

	status, resp, err := client.Get(fmt.Sprintf("jobs/%s", auth0JobID), nil)
	if err != nil {
		log.Warningf("failed to fetch auth0 job; %s", err.Error())
		return nil, err
	}
	if status != 200 {
		log.Warningf("failed to fetch auth0 job; %s", err.Error())
		return nil, err
	}

	return resp, nil
}
