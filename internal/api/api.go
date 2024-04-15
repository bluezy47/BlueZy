package api
import (
	"fmt"
	"log"
	"encoding/json"
	"net/http"
)

func GetGoogleUserDetails(accessTk string) (map[string]interface{}, error) {
	// get the user details from the google api
	userInfoReq, err := http.NewRequest(
		"GET",
		"https://www.googleapis.com/oauth2/v3/userinfo",
		nil,
	)
	if err != nil {
		log.Println("::[API][GetGoogleUserDetails] Error creating the request ::");
		return nil, err
	}
	// set the Header...
	userInfoReq.Header.Set(
		"Authorization", fmt.Sprintf("Bearer %s", accessTk),
	);

	// make the request...
	userInfoResp, err := http.DefaultClient.Do(userInfoReq)
	if err != nil {
		log.Println("::[API][GetGoogleUserDetails] Error making the request ::");
		return nil, err
	}
	defer userInfoResp.Body.Close(); // close the response body after the function returns

	var userInfo map[string]interface{};
	if err := json.NewDecoder(userInfoResp.Body).Decode(&userInfo); err != nil {
		log.Println("::[API][GetGoogleUserDetails] Error decoding the response ::");
		return nil, err
	}
	// success...
	return userInfo, nil
}
